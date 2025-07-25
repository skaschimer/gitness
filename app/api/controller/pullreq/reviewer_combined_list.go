// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pullreq

import (
	"context"
	"errors"
	"fmt"

	"github.com/harness/gitness/app/auth"
	"github.com/harness/gitness/store"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type CombinedListResponse struct {
	Reviewers          []*types.PullReqReviewer   `json:"reviewers,omitempty"`
	UserGroupReviewers []*types.UserGroupReviewer `json:"user_group_reviewers,omitempty"`
}

// ReviewersListCombined returns the combined reviewer list for the pull request.
func (c *Controller) ReviewersListCombined(
	ctx context.Context,
	session *auth.Session,
	repoRef string,
	prNum int64,
) (*CombinedListResponse, error) {
	repo, err := c.getRepoCheckAccess(ctx, session, repoRef, enum.PermissionRepoView)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire access to repo: %w", err)
	}

	pr, err := c.pullreqStore.FindByNumber(ctx, repo.ID, prNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find pull request by number: %w", err)
	}

	reviewers, err := c.reviewerStore.List(ctx, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviewers: %w", err)
	}

	userGroupReviewers, err := c.userGroupReviewerStore.List(ctx, pr.ID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, fmt.Errorf("failed to list user group reviewers: %w", err)
	}
	if errors.Is(err, store.ErrResourceNotFound) || len(userGroupReviewers) == 0 {
		return &CombinedListResponse{
			Reviewers: reviewers,
		}, nil
	}
	userGroupReviewersMap := make(map[int64]*types.UserGroupReviewer, len(userGroupReviewers))
	for _, userGroupReviewer := range userGroupReviewers {
		userGroupReviewersMap[userGroupReviewer.UserGroupID] = userGroupReviewer
	}

	addedByIDs := make([]int64, len(userGroupReviewers))
	userGroupIDs := make([]int64, len(userGroupReviewers))
	for i, v := range userGroupReviewers {
		addedByIDs[i] = v.CreatedBy
		userGroupIDs[i] = v.UserGroupID
	}
	userGroupsMap, err := c.userGroupStore.Map(ctx, userGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to map usergroups: %w", err)
	}
	principalInfoCacheMap, err := c.principalInfoCache.Map(ctx, addedByIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load PR principal infos: %w", err)
	}

	groupPrincipalsMap, err := c.userGroupService.MapGroupIDsToPrincipals(ctx, userGroupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to map group IDs to principals: %w", err)
	}

	reviewersMap := reviewersMap(reviewers)

	for groupID, principals := range groupPrincipalsMap {
		// userGroupReviewersMap, userGroupsMap and groupPrincipalsMap depend on userGroupReviewers
		// if a group doesn't exist in userGroupReviewers it won't exist in any of these
		// and if it exists in userGroupReviewers it will exist in all of these
		userGroupReviewer := userGroupReviewersMap[groupID]
		userGroupReviewer.UserGroup = *userGroupsMap[groupID].ToUserGroupInfo()

		// principal could be deleted/removed without group being, so we check for its existence
		if addedBy, ok := principalInfoCacheMap[userGroupReviewer.CreatedBy]; ok {
			userGroupReviewer.AddedBy = *addedBy
		}

		userGroupReviewer.Decision = enum.PullReqReviewDecisionPending
		var userGroupReviewerDecisions []types.UserGroupReviewerDecision
		for _, principal := range principals {
			reviewer, ok := reviewersMap[principal.ID]
			if !ok {
				continue
			}

			userGroupReviewerDecisions = append(
				userGroupReviewerDecisions,
				types.UserGroupReviewerDecision{
					Decision: reviewer.ReviewDecision,
					SHA:      reviewer.SHA,
					Reviewer: reviewer.Reviewer,
				},
			)

			userGroupReviewer.Decision = getHighestOrderDecision(
				userGroupReviewer.Decision, reviewer.ReviewDecision,
			)
		}
		userGroupReviewer.UserDecisions = userGroupReviewerDecisions
	}

	return &CombinedListResponse{
		Reviewers:          reviewers,
		UserGroupReviewers: userGroupReviewers,
	}, nil
}

func getHighestOrderDecision(
	d1 enum.PullReqReviewDecision,
	d2 enum.PullReqReviewDecision,
) enum.PullReqReviewDecision {
	if d1 == enum.PullReqReviewDecisionChangeReq || d2 == enum.PullReqReviewDecisionChangeReq {
		return enum.PullReqReviewDecisionChangeReq
	}
	if d1 == enum.PullReqReviewDecisionApproved || d2 == enum.PullReqReviewDecisionApproved {
		return enum.PullReqReviewDecisionApproved
	}
	if d1 == enum.PullReqReviewDecisionReviewed || d2 == enum.PullReqReviewDecisionReviewed {
		return enum.PullReqReviewDecisionReviewed
	}
	return enum.PullReqReviewDecisionPending
}
