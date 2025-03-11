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

package events

import (
	"context"

	"github.com/harness/gitness/events"
	"github.com/harness/gitness/registry/app/api/openapi/contracts/artifact"

	"github.com/rs/zerolog/log"
)

const ArtifactsCategory = "artifacts"

const ArtifactCreatedEvent events.EventType = "artifact-created"
const ArtifactUpdatedEvent events.EventType = "artifact-updated"
const ArtifactDeletedEvent events.EventType = "artifact-deleted"

type ArtifactCreatedPayload struct {
	RegistryID   int64                `json:"registry_id"`
	PrincipalID  int64                `json:"principal_id"`
	ArtifactType artifact.PackageType `json:"artifact_type"`
	Artifact     Artifact             `json:"artifact"`
}

type Artifact interface {
	GetInfo() string
}

type BaseArtifact struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
}

type DockerArtifact struct {
	BaseArtifact
	URL    string `json:"url"`
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
}

func (a *DockerArtifact) GetInfo() string {
	return a.Ref
}

type HelmArtifact struct {
	BaseArtifact
	URL    string `json:"url"`
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
}

func (a *HelmArtifact) GetInfo() string {
	return a.Ref
}

type ArtifactInfo struct {
	Type     artifact.PackageType `json:"type"`
	Name     string               `json:"name"`
	Version  string               `json:"version"`
	Artifact interface{}          `json:"artifact"`
}

type ArtifactChangeInfo struct {
	Type           artifact.PackageType `json:"type"`
	Name           string               `json:"name"`
	ArtifactChange interface{}          `json:"artifact_change"`
}

func (r *Reporter) ArtifactCreated(ctx context.Context, payload *ArtifactCreatedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, ArtifactCreatedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send artifact-created created event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported artifact-created event with id '%s'", eventID)
}

func (r *Reader) RegisterArtifactCreated(
	fn events.HandlerFunc[*ArtifactCreatedPayload],
	opts ...events.HandlerOption,
) error {
	return events.ReaderRegisterEvent(r.innerReader, ArtifactCreatedEvent, fn, opts...)
}

type ArtifactUpdatedPayload struct {
	RegistryID     int64                `json:"registry_id"`
	PrincipalID    int64                `json:"principal_id"`
	ArtifactType   artifact.PackageType `json:"artifact_type"`
	ArtifactChange ArtifactChange       `json:"artifact_change"`
}

type ArtifactChange struct {
	Old Artifact
	New Artifact
}

type DockerArtifactChange struct {
	Old DockerArtifact
	New DockerArtifact
}

type HelmArtifactChange struct {
	Old HelmArtifact
	New HelmArtifact
}

func (r *Reporter) ArtifactUpdated(ctx context.Context, payload *ArtifactUpdatedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, ArtifactUpdatedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send artifact updated event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported artifact updated event with id '%s'", eventID)
}

func (r *Reader) RegisterArtifactUpdated(
	fn events.HandlerFunc[*ArtifactUpdatedPayload],
	opts ...events.HandlerOption,
) error {
	return events.ReaderRegisterEvent(r.innerReader, ArtifactUpdatedEvent, fn, opts...)
}

type ArtifactDeletedPayload struct {
	RegistryID   int64                `json:"registry_id"`
	PrincipalID  int64                `json:"principal_id"`
	ArtifactType artifact.PackageType `json:"artifact_type"`
	Artifact     Artifact             `json:"artifact"`
}

func (r *Reporter) ArtifactDeleted(ctx context.Context, payload *ArtifactDeletedPayload) {
	eventID, err := events.ReporterSendEvent(r.innerReporter, ctx, ArtifactDeletedEvent, payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to send artifact deleted event")
		return
	}

	log.Ctx(ctx).Debug().Msgf("reported artifact deleted event with id '%s'", eventID)
}

func (r *Reader) RegisterArtifactDeleted(
	fn events.HandlerFunc[*ArtifactDeletedPayload],
	opts ...events.HandlerOption,
) error {
	return events.ReaderRegisterEvent(r.innerReader, ArtifactDeletedEvent, fn, opts...)
}
