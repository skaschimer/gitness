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
	"fmt"

	"github.com/harness/gitness/registry/app/pkg"
	"github.com/harness/gitness/registry/app/pkg/base"
	"github.com/harness/gitness/registry/app/pkg/nuget"
	"github.com/harness/gitness/registry/app/pkg/response"
	nugettype "github.com/harness/gitness/registry/app/pkg/types/nuget"
	registrytypes "github.com/harness/gitness/registry/types"
)

func (c *controller) DownloadPackage(
	ctx context.Context,
	info nugettype.ArtifactInfo,
) *GetArtifactResponse {
	f := func(registry registrytypes.Registry, a pkg.Artifact) response.Response {
		info.RegIdentifier = registry.Name
		info.RegistryID = registry.ID
		nugetRegistry, ok := a.(nuget.Registry)
		if !ok {
			return &GetArtifactResponse{
				BaseResponse{
					fmt.Errorf("invalid registry type: expected nuget.Registry"),
					nil,
				},
				"", nil, nil,
			}
		}
		headers, fileReader, redirectURL, err := nugetRegistry.DownloadPackage(ctx, info)
		return &GetArtifactResponse{
			BaseResponse{
				err,
				headers,
			},
			redirectURL, fileReader, nil,
		}
	}

	result, err := base.ProxyWrapper(ctx, c.registryDao, f, info)
	if err != nil {
		return &GetArtifactResponse{
			BaseResponse{
				err,
				nil,
			},
			"", nil, nil,
		}
	}
	getResponse, ok := result.(*GetArtifactResponse)
	if !ok {
		return &GetArtifactResponse{
			BaseResponse{
				fmt.Errorf("invalid registry type: expected nuget.Registry"),
				nil,
			}, "", nil, nil,
		}
	}
	return getResponse
}
