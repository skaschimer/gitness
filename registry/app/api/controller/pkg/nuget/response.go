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
	"io"

	"github.com/harness/gitness/registry/app/pkg/commons"
	"github.com/harness/gitness/registry/app/pkg/response"
	"github.com/harness/gitness/registry/app/pkg/types/nuget"
	"github.com/harness/gitness/registry/app/storage"
)

var _ response.Response = (*GetServiceEndpointResponse)(nil)
var _ response.Response = (*GetArtifactResponse)(nil)
var _ response.Response = (*PutArtifactResponse)(nil)

type BaseResponse struct {
	Error           error
	ResponseHeaders *commons.ResponseHeaders
}

func (r BaseResponse) GetError() error {
	return r.Error
}

type GetServiceEndpointResponse struct {
	BaseResponse
	ServiceEndpoint *nuget.ServiceEndpoint
}

type GetServiceEndpointV2Response struct {
	BaseResponse
	ServiceEndpoint *nuget.ServiceEndpointV2
}

type GetServiceMetadataV2Response struct {
	BaseResponse
	ServiceMetadata *nuget.ServiceMetadataV2
}

type GetArtifactResponse struct {
	BaseResponse
	RedirectURL string
	Body        *storage.FileReader
	ReadCloser  io.ReadCloser
}
type PutArtifactResponse struct {
	BaseResponse
}
type DeleteArtifactResponse struct {
	BaseResponse
}

type ListPackageVersionResponse struct {
	BaseResponse
	PackageVersion *nuget.PackageVersion
}

type ListPackageVersionV2Response struct {
	BaseResponse
	FeedResponse *nuget.FeedResponse
}

type SearchPackageV2Response struct {
	BaseResponse
	FeedResponse *nuget.FeedResponse
}

type SearchPackageResponse struct {
	BaseResponse
	SearchResponse *nuget.SearchResultResponse
}

type EntityCountResponse struct {
	BaseResponse
	Count int64
}

type GetPackageMetadataResponse struct {
	BaseResponse
	RegistrationResponse nuget.RegistrationResponse
}

type GetPackageVersionMetadataV2Response struct {
	BaseResponse
	FeedEntryResponse *nuget.FeedEntryResponse
}

type RegistrationResponse interface {
	isRegistrationResponse() // marker method
}

type GetPackageVersionMetadataResponse struct {
	BaseResponse
	RegistrationLeafResponse *nuget.RegistrationLeafResponse
}
