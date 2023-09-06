// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/harness/gitness/build/file"
	"github.com/harness/gitness/build/scheduler"
	"github.com/harness/gitness/internal/store"
	urlprovider "github.com/harness/gitness/internal/url"
	"github.com/harness/gitness/livelog"
	gitness_store "github.com/harness/gitness/store"
	"github.com/harness/gitness/types"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

var noContext = context.Background()

var _ ExecutionManager = (*Manager)(nil)

type (
	// Request provides filters when requesting a pending
	// build from the queue. This allows an agent, for example,
	// to request a build that matches its architecture and kernel.
	Request struct {
		Kind    string            `json:"kind"`
		Type    string            `json:"type"`
		OS      string            `json:"os"`
		Arch    string            `json:"arch"`
		Variant string            `json:"variant"`
		Kernel  string            `json:"kernel"`
		Labels  map[string]string `json:"labels,omitempty"`
	}

	// Config represents a pipeline config file.
	Config struct {
		Data string `json:"data"`
		Kind string `json:"kind"`
	}

	// Context represents the minimum amount of information
	// required by the runner to execute a build.
	Context struct {
		Repo    *types.Repository `json:"repository"`
		Build   *types.Execution  `json:"build"`
		Stage   *types.Stage      `json:"stage"`
		Secrets []*types.Secret   `json:"secrets"`
		Config  *file.File        `json:"config"`
	}

	// ExecutionManager encapsulates complex build operations and provides
	// a simplified interface for build runners.
	ExecutionManager interface {
		// Request requests the next available build stage for execution.
		Request(ctx context.Context, args *Request) (*types.Stage, error)

		// Accept accepts the build stage for execution.
		Accept(ctx context.Context, stage int64, machine string) (*types.Stage, error)

		// Write writes a line to the build logs.
		Write(ctx context.Context, step int64, line *livelog.Line) error

		// Details returns details about stage.
		Details(ctx context.Context, stageID int64) (*Context, error)

		// Upload uploads the full logs.
		Upload(ctx context.Context, step int64, r io.Reader) error

		// UploadBytes uploads the full logs.
		UploadBytes(ctx context.Context, step int64, b []byte) error

		// Before signals the build step is about to start.
		Before(ctx context.Context, step *types.Step) error

		// After signals the build step is complete.
		After(ctx context.Context, step *types.Step) error

		// BeforeAll signals the build stage is about to start.
		BeforeAll(ctx context.Context, stage *types.Stage) error

		// AfterAll signals the build stage is complete.
		AfterAll(ctx context.Context, stage *types.Stage) error
	}
)

// Manager provides a simplified interface to the build runner so that it
// can more easily interact with the server.
type Manager struct {
	Executions  store.ExecutionStore
	Config      *types.Config
	FileService file.FileService
	Pipelines   store.PipelineStore
	urlProvider *urlprovider.Provider
	// Converter  store.ConvertService
	// Events     store.Pubsub
	// Globals    store.GlobalSecretStore
	Logs store.LogStore
	Logz livelog.LogStream
	// Netrcs     store.NetrcService
	Repos     store.RepoStore
	Scheduler scheduler.Scheduler
	Secrets   store.SecretStore
	// Status  store.StatusService
	Stages store.StageStore
	Steps  store.StepStore
	// System  *store.System
	Users store.PrincipalStore
	// Webhook store.WebhookSender
}

func New(
	config *types.Config,
	executionStore store.ExecutionStore,
	pipelineStore store.PipelineStore,
	urlProvider *urlprovider.Provider,
	fileService file.FileService,
	logStore store.LogStore,
	logStream livelog.LogStream,
	repoStore store.RepoStore,
	scheduler scheduler.Scheduler,
	secretStore store.SecretStore,
	stageStore store.StageStore,
	stepStore store.StepStore,
	userStore store.PrincipalStore,
) *Manager {
	return &Manager{
		Config:      config,
		Executions:  executionStore,
		Pipelines:   pipelineStore,
		urlProvider: urlProvider,
		FileService: fileService,
		Logs:        logStore,
		Logz:        logStream,
		Repos:       repoStore,
		Scheduler:   scheduler,
		Secrets:     secretStore,
		Stages:      stageStore,
		Steps:       stepStore,
		Users:       userStore,
	}
}

// Request requests the next available build stage for execution.
func (m *Manager) Request(ctx context.Context, args *Request) (*types.Stage, error) {
	log := log.With().
		Str("kind", args.Kind).
		Str("type", args.Type).
		Str("os", args.OS).
		Str("arch", args.Arch).
		Str("kernel", args.Kernel).
		Str("variant", args.Variant).
		Logger()
	log.Debug().Msg("manager: request queue item")

	stage, err := m.Scheduler.Request(ctx, scheduler.Filter{
		Kind:    args.Kind,
		Type:    args.Type,
		OS:      args.OS,
		Arch:    args.Arch,
		Kernel:  args.Kernel,
		Variant: args.Variant,
		Labels:  args.Labels,
	})
	if err != nil && ctx.Err() != nil {
		log.Debug().Err(err).Msg("manager: context canceled")
		return nil, err
	}
	if err != nil {
		log.Warn().Err(err).Msg("manager: request queue item error")
		return nil, err
	}
	return stage, nil
}

// Accept accepts the build stage for execution. It is possible for multiple
// agents to pull the same stage from the queue.
func (m *Manager) Accept(ctx context.Context, id int64, machine string) (*types.Stage, error) {
	log := log.With().
		Int64("stage-id", id).
		Str("machine", machine).
		Logger()
	log.Debug().Msg("manager: accept stage")

	stage, err := m.Stages.Find(noContext, id)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot find stage")
		return nil, err
	}
	if stage.Machine != "" {
		log.Debug().Msg("manager: stage already assigned. abort.")
		return nil, fmt.Errorf("stage already assigned, abort")
	}

	stage.Machine = machine
	stage.Status = types.StatusPending
	stage.Updated = time.Now().Unix()
	err = m.Stages.Update(noContext, stage)
	if errors.Is(err, gitness_store.ErrVersionConflict) {
		log.Debug().Err(err).Msg("manager: stage processed by another agent")
	} else if err != nil {
		log.Debug().Err(err).Msg("manager: cannot update stage")
	} else {
		log.Debug().Err(err).Msg("manager: stage accepted")
	}
	return stage, err
}

// Write writes a line to the build logs.
func (m *Manager) Write(ctx context.Context, step int64, line *livelog.Line) error {
	fmt.Println("line is: ", line)
	err := m.Logz.Write(ctx, step, line)
	if err != nil {
		log.Warn().Int64("step-id", step).Err(err).Msg("manager: cannot write to log stream")
		return err
	}
	return nil
}

// Upload uploads the full logs.
func (m *Manager) Upload(ctx context.Context, step int64, r io.Reader) error {
	err := m.Logs.Create(ctx, step, r)
	if err != nil {
		log.Error().Err(err).Int64("step-id", step).Msg("manager: cannot upload complete logs")
		return err
	}
	return nil
}

// UploadBytes uploads the full logs.
func (m *Manager) UploadBytes(ctx context.Context, step int64, data []byte) error {
	buf := bytes.NewBuffer(data)
	err := m.Logs.Create(ctx, step, buf)
	if err != nil {
		log.Error().Err(err).Int64("step-id", step).Msg("manager: cannot upload complete logs")
		return err
	}
	return nil
}

// Details provides details about the stage.
func (m *Manager) Details(ctx context.Context, stageID int64) (*Context, error) {
	log := log.With().
		Int64("stage-id", stageID).
		Logger()
	log.Debug().Msg("manager: fetching stage details")

	stage, err := m.Stages.Find(noContext, stageID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot find stage")
		return nil, err
	}
	execution, err := m.Executions.Find(noContext, stage.ExecutionID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot find build")
		return nil, err
	}
	pipeline, err := m.Pipelines.Find(noContext, execution.PipelineID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot find pipeline")
		return nil, err
	}
	repo, err := m.Repos.Find(noContext, execution.RepoID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot find repo")
		return nil, err
	}
	// Backfill clone URL
	repo.GitURL, err = m.createCustomCloneURL(repo.Path)
	if err != nil {
		log.Warn().Err(err).Msg("manager: could not create custom clone url")
		return nil, err
	}
	stages, err := m.Stages.List(ctx, stage.ExecutionID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot list stages")
		return nil, err
	}
	execution.Stages = stages
	log = log.With().
		Int64("build", execution.Number).
		Str("repo", repo.GetGitUID()).
		Logger()

	// TODO: Currently we fetch all the secrets from the same space.
	// This logic can be updated when needed.
	secrets, err := m.Secrets.ListAll(noContext, repo.ParentID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot list secrets")
		return nil, err
	}

	// Fetch contents of YAML from the execution ref at the pipeline config path.
	file, err := m.FileService.Get(ctx, repo, pipeline.ConfigPath, execution.After)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot fetch file")
		return nil, err
	}

	return &Context{
		Repo:    repo,
		Build:   execution,
		Stage:   stage,
		Secrets: secrets,
		Config:  file,
	}, nil
}

// Before signals the build step is about to start.
func (m *Manager) Before(ctx context.Context, step *types.Step) error {
	log := log.With().
		Str("step.status", step.Status).
		Str("step.name", step.Name).
		Int64("step.id", step.ID).
		Logger()

	log.Debug().Msg("manager: updating step status")

	err := m.Logz.Create(noContext, step.ID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: cannot create log stream")
		return err
	}
	updater := &updater{
		Executions: m.Executions,
		Repos:      m.Repos,
		Steps:      m.Steps,
		Stages:     m.Stages,
	}
	return updater.do(ctx, step)
}

// After signals the build step is complete.
func (m *Manager) After(ctx context.Context, step *types.Step) error {
	log := log.With().
		Str("step.status", step.Status).
		Str("step.name", step.Name).
		Int64("step.id", step.ID).
		Logger()
	log.Debug().Msg("manager: updating step status")

	var errs error
	updater := &updater{
		Executions: m.Executions,
		Repos:      m.Repos,
		Steps:      m.Steps,
		Stages:     m.Stages,
	}

	if err := updater.do(ctx, step); err != nil {
		errs = multierror.Append(errs, err)
		log.Warn().Err(errs).Msg("manager: cannot update step")
	}

	if err := m.Logz.Delete(noContext, step.ID); err != nil {
		log.Warn().Err(err).Msg("manager: cannot teardown log stream")
	}
	return errs
}

// BeforeAll signals the build stage is about to start.
func (m *Manager) BeforeAll(ctx context.Context, stage *types.Stage) error {
	s := &setup{
		Executions: m.Executions,
		Repos:      m.Repos,
		Steps:      m.Steps,
		Stages:     m.Stages,
		Users:      m.Users,
	}

	err := s.do(ctx, stage)
	return err
}

// AfterAll signals the build stage is complete.
func (m *Manager) AfterAll(ctx context.Context, stage *types.Stage) error {
	t := &teardown{
		Executions: m.Executions,
		Logs:       m.Logz,
		Repos:      m.Repos,
		Scheduler:  m.Scheduler,
		Steps:      m.Steps,
		Stages:     m.Stages,
	}
	return t.do(ctx, stage)
}

// createCustomCloneURL creates an endpoint to interact with gitness
// using the network (if provided) that gitness is running on.
func (m *Manager) createCustomCloneURL(repoPath string) (string, error) {
	// We use http to interact with gitness from the build containers.
	endpoint := "http://" + m.Config.Server.HTTP.Network + m.Config.Server.HTTP.Bind
	url, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	return m.urlProvider.GenerateCustomRepoCloneURL(url, repoPath), nil
}
