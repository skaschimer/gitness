//  Copyright 2023 Harness, Inc.
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

package nuget

import (
	"context"
	"io"

	"github.com/harness/gitness/app/services/refcache"
	"github.com/harness/gitness/registry/app/api/openapi/contracts/artifact"
	"github.com/harness/gitness/registry/app/remote/adapter"
	"github.com/harness/gitness/registry/app/remote/adapter/nuget"
	"github.com/harness/gitness/registry/app/remote/registry"
	"github.com/harness/gitness/registry/types"
	"github.com/harness/gitness/secret"

	"github.com/rs/zerolog/log"
)

type RemoteRegistryHelper interface {
	GetFile(ctx context.Context, pkg, version, proxyEndpoint, fileName string) (io.ReadCloser, error)

	GetPackageMetadata(ctx context.Context, pkg, proxyEndpoint string) (io.ReadCloser, error)

	GetPackageVersionMetadataV2(ctx context.Context, pkg, version string) (io.ReadCloser, error)

	ListPackageVersion(ctx context.Context, pkg string) (io.ReadCloser, error)

	ListPackageVersionV2(ctx context.Context, pkg string) (io.ReadCloser, error)

	SearchPackageV2(ctx context.Context, searchTerm string, limit, offset int) (io.ReadCloser, error)

	SearchPackage(ctx context.Context, searchTerm string, limit, offset int) (io.ReadCloser, error)

	CountPackageV2(ctx context.Context, searchTerm string) (int64, error)

	CountPackageVersionV2(ctx context.Context, pkg string) (int64, error)
}

type remoteRegistryHelper struct {
	adapter  registry.NugetRegistry
	registry types.UpstreamProxy
}

func NewRemoteRegistryHelper(
	ctx context.Context,
	spaceFinder refcache.SpaceFinder,
	registry types.UpstreamProxy,
	service secret.Service,
) (RemoteRegistryHelper, error) {
	r := &remoteRegistryHelper{
		registry: registry,
	}
	if err := r.init(ctx, spaceFinder, service); err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to init remote registry for remote: %s", registry.RepoKey)
		return nil, err
	}
	return r, nil
}

func (r *remoteRegistryHelper) init(
	ctx context.Context,
	spaceFinder refcache.SpaceFinder,
	service secret.Service,
) error {
	key := string(artifact.PackageTypeNUGET)
	if r.registry.Source == string(artifact.UpstreamConfigSourceNugetOrg) {
		r.registry.RepoURL = nuget.NugetOrgURL
	}

	factory, err := adapter.GetFactory(key)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to get factory " + key)
		return err
	}

	adpt, err := factory.Create(ctx, spaceFinder, r.registry, service)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create factory " + key)
		return err
	}

	nugetReg, ok := adpt.(registry.NugetRegistry)
	if !ok {
		log.Ctx(ctx).Error().Msg("failed to cast factory to nuget registry")
		return err
	}
	r.adapter = nugetReg
	return nil
}

func (r *remoteRegistryHelper) GetFile(ctx context.Context, pkg,
	version, proxyEndpoint, fileName string) (io.ReadCloser, error) {
	v2, err := r.adapter.GetPackage(ctx, pkg, version, proxyEndpoint, fileName)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to get pkg: %s, version: %s", pkg, version)
	}
	return v2, err
}

func (r *remoteRegistryHelper) GetPackageMetadata(ctx context.Context,
	pkg, proxyEndpoint string) (io.ReadCloser, error) {
	metadata, err := r.adapter.GetPackageMetadata(ctx, pkg, proxyEndpoint)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to get metadata for pkg: %s", pkg)
		return nil, err
	}
	return metadata, nil
}

func (r *remoteRegistryHelper) ListPackageVersion(ctx context.Context,
	pkg string) (io.ReadCloser, error) {
	packageVersions, err := r.adapter.ListPackageVersion(ctx, pkg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to get package version for pkg: %s", pkg)
		return nil, err
	}
	return packageVersions, nil
}

func (r *remoteRegistryHelper) ListPackageVersionV2(ctx context.Context,
	pkg string) (io.ReadCloser, error) {
	packageVersions, err := r.adapter.ListPackageVersionV2(ctx, pkg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to get package version for pkg: %s", pkg)
		return nil, err
	}
	return packageVersions, nil
}

func (r *remoteRegistryHelper) GetPackageVersionMetadataV2(ctx context.Context,
	pkg, version string) (io.ReadCloser, error) {
	metadata, err := r.adapter.GetPackageVersionMetadataV2(ctx, pkg, version)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to get metadata for pkg: %s", pkg)
		return nil, err
	}
	return metadata, nil
}

func (r *remoteRegistryHelper) SearchPackageV2(ctx context.Context,
	searchTerm string, limit, offset int) (io.ReadCloser, error) {
	searchResults, err := r.adapter.SearchPackageV2(ctx, searchTerm, limit, offset)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to search packages with term: %s", searchTerm)
		return nil, err
	}
	return searchResults, nil
}

func (r *remoteRegistryHelper) SearchPackage(ctx context.Context,
	searchTerm string, limit, offset int) (io.ReadCloser, error) {
	searchResults, err := r.adapter.SearchPackage(ctx, searchTerm, limit, offset)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to search packages (v3) with term: %s", searchTerm)
		return nil, err
	}
	return searchResults, nil
}

func (r *remoteRegistryHelper) CountPackageV2(ctx context.Context,
	searchTerm string) (int64, error) {
	count, err := r.adapter.CountPackageV2(ctx, searchTerm)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to count packages with term: %s", searchTerm)
		return 0, err
	}
	return count, nil
}

func (r *remoteRegistryHelper) CountPackageVersionV2(ctx context.Context,
	pkg string) (int64, error) {
	count, err := r.adapter.CountPackageVersionV2(ctx, pkg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to count package versions for pkg: %s", pkg)
		return 0, err
	}
	return count, nil
}
