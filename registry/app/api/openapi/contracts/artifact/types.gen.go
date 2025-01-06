// Package artifact provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package artifact

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oapi-codegen/runtime"
)

// Defines values for AuthType.
const (
	AuthTypeAccessKeySecretKey AuthType = "AccessKeySecretKey"
	AuthTypeAnonymous          AuthType = "Anonymous"
	AuthTypeUserPassword       AuthType = "UserPassword"
)

// Defines values for ClientSetupStepType.
const (
	ClientSetupStepTypeGenerateToken ClientSetupStepType = "GenerateToken"
	ClientSetupStepTypeStatic        ClientSetupStepType = "Static"
)

// Defines values for PackageType.
const (
	PackageTypeDOCKER  PackageType = "DOCKER"
	PackageTypeGENERIC PackageType = "GENERIC"
	PackageTypeHELM    PackageType = "HELM"
	PackageTypeMAVEN   PackageType = "MAVEN"
)

// Defines values for RegistryType.
const (
	RegistryTypeUPSTREAM RegistryType = "UPSTREAM"
	RegistryTypeVIRTUAL  RegistryType = "VIRTUAL"
)

// Defines values for Status.
const (
	StatusERROR   Status = "ERROR"
	StatusFAILURE Status = "FAILURE"
	StatusSUCCESS Status = "SUCCESS"
)

// Defines values for UpstreamConfigSource.
const (
	UpstreamConfigSourceAwsEcr    UpstreamConfigSource = "AwsEcr"
	UpstreamConfigSourceCustom    UpstreamConfigSource = "Custom"
	UpstreamConfigSourceDockerhub UpstreamConfigSource = "Dockerhub"
)

// Defines values for RegistryTypeParam.
const (
	UPSTREAM RegistryTypeParam = "UPSTREAM"
	VIRTUAL  RegistryTypeParam = "VIRTUAL"
)

// Defines values for GetAllRegistriesParamsType.
const (
	GetAllRegistriesParamsTypeUPSTREAM GetAllRegistriesParamsType = "UPSTREAM"
	GetAllRegistriesParamsTypeVIRTUAL  GetAllRegistriesParamsType = "VIRTUAL"
)

// AccessKeySecretKey defines model for AccessKeySecretKey.
type AccessKeySecretKey struct {
	AccessKey                 *string `json:"accessKey,omitempty"`
	AccessKeySecretIdentifier *string `json:"accessKeySecretIdentifier,omitempty"`
	AccessKeySecretSpaceId    *int    `json:"accessKeySecretSpaceId,omitempty"`
	AccessKeySecretSpacePath  *string `json:"accessKeySecretSpacePath,omitempty"`
	SecretKeyIdentifier       string  `json:"secretKeyIdentifier"`
	SecretKeySpaceId          *int    `json:"secretKeySpaceId,omitempty"`
	SecretKeySpacePath        *string `json:"secretKeySpacePath,omitempty"`
}

// Anonymous defines model for Anonymous.
type Anonymous interface{}

// ArtifactLabelRequest defines model for ArtifactLabelRequest.
type ArtifactLabelRequest struct {
	Labels []string `json:"labels"`
}

// ArtifactMetadata Artifact Metadata
type ArtifactMetadata struct {
	DownloadsCount *int64    `json:"downloadsCount,omitempty"`
	Labels         *[]string `json:"labels,omitempty"`
	LastModified   *string   `json:"lastModified,omitempty"`
	Name           string    `json:"name"`

	// PackageType refers to package
	PackageType        *PackageType `json:"packageType,omitempty"`
	PullCommand        *string      `json:"pullCommand,omitempty"`
	RegistryIdentifier string       `json:"registryIdentifier"`
	RegistryPath       string       `json:"registryPath"`
	Version            *string      `json:"version,omitempty"`
}

// ArtifactStats Harness Artifact Stats
type ArtifactStats struct {
	DownloadCount    *int64 `json:"downloadCount,omitempty"`
	DownloadSize     *int64 `json:"downloadSize,omitempty"`
	TotalStorageSize *int64 `json:"totalStorageSize,omitempty"`
	UploadSize       *int64 `json:"uploadSize,omitempty"`
}

// ArtifactSummary Harness Artifact Summary
type ArtifactSummary struct {
	CreatedAt      *string   `json:"createdAt,omitempty"`
	DownloadsCount *int64    `json:"downloadsCount,omitempty"`
	ImageName      string    `json:"imageName"`
	Labels         *[]string `json:"labels,omitempty"`
	ModifiedAt     *string   `json:"modifiedAt,omitempty"`

	// PackageType refers to package
	PackageType PackageType `json:"packageType"`
}

// ArtifactVersionMetadata Artifact Version Metadata
type ArtifactVersionMetadata struct {
	DigestCount     *int    `json:"digestCount,omitempty"`
	DownloadsCount  *int64  `json:"downloadsCount,omitempty"`
	IslatestVersion *bool   `json:"islatestVersion,omitempty"`
	LastModified    *string `json:"lastModified,omitempty"`
	Name            string  `json:"name"`

	// PackageType refers to package
	PackageType        *PackageType `json:"packageType,omitempty"`
	PullCommand        *string      `json:"pullCommand,omitempty"`
	RegistryIdentifier string       `json:"registryIdentifier"`
	RegistryPath       string       `json:"registryPath"`
	Size               *string      `json:"size,omitempty"`
}

// ArtifactVersionSummary Docker Artifact Version Summary
type ArtifactVersionSummary struct {
	ImageName       string `json:"imageName"`
	IsLatestVersion *bool  `json:"isLatestVersion,omitempty"`

	// PackageType refers to package
	PackageType PackageType `json:"packageType"`
	Version     string      `json:"version"`
}

// AuthType Authentication type
type AuthType string

// CleanupPolicy Cleanup Policy for Harness Artifact Registries
type CleanupPolicy struct {
	ExpireDays    *int      `json:"expireDays,omitempty"`
	Name          *string   `json:"name,omitempty"`
	PackagePrefix *[]string `json:"packagePrefix,omitempty"`
	VersionPrefix *[]string `json:"versionPrefix,omitempty"`
}

// ClientSetupDetails Client Setup Details
type ClientSetupDetails struct {
	MainHeader string               `json:"mainHeader"`
	SecHeader  string               `json:"secHeader"`
	Sections   []ClientSetupSection `json:"sections"`
}

// ClientSetupSection Client Setup Section
type ClientSetupSection struct {
	Header *string            `json:"header,omitempty"`
	Steps  *[]ClientSetupStep `json:"steps,omitempty"`
}

// ClientSetupStep Client Setup Step
type ClientSetupStep struct {
	Commands *[]ClientSetupStepCommand `json:"commands,omitempty"`
	Header   *string                   `json:"header,omitempty"`

	// Type ClientSetupStepType type
	Type *ClientSetupStepType `json:"type,omitempty"`
}

// ClientSetupStepCommand Client Setup Step Command
type ClientSetupStepCommand struct {
	Label *string `json:"label,omitempty"`
	Value *string `json:"value,omitempty"`
}

// ClientSetupStepType ClientSetupStepType type
type ClientSetupStepType string

// DockerArtifactDetail Docker Artifact Detail
type DockerArtifactDetail struct {
	CreatedAt       *string `json:"createdAt,omitempty"`
	DownloadsCount  *int64  `json:"downloadsCount,omitempty"`
	ImageName       string  `json:"imageName"`
	IsLatestVersion *bool   `json:"isLatestVersion,omitempty"`
	ModifiedAt      *string `json:"modifiedAt,omitempty"`

	// PackageType refers to package
	PackageType  PackageType `json:"packageType"`
	PullCommand  *string     `json:"pullCommand,omitempty"`
	RegistryPath string      `json:"registryPath"`
	Size         *string     `json:"size,omitempty"`
	Url          string      `json:"url"`
	Version      string      `json:"version"`
}

// DockerArtifactManifest Docker Artifact Manifest
type DockerArtifactManifest struct {
	Manifest string `json:"manifest"`
}

// DockerLayerEntry Harness Artifact Layers
type DockerLayerEntry struct {
	Command string  `json:"command"`
	Size    *string `json:"size,omitempty"`
}

// DockerLayersSummary Harness Layers Summary
type DockerLayersSummary struct {
	Digest string              `json:"digest"`
	Layers *[]DockerLayerEntry `json:"layers,omitempty"`
	OsArch *string             `json:"osArch,omitempty"`
}

// DockerManifestDetails Harness Artifact Layers
type DockerManifestDetails struct {
	CreatedAt      *string `json:"createdAt,omitempty"`
	Digest         string  `json:"digest"`
	DownloadsCount *int64  `json:"downloadsCount,omitempty"`
	OsArch         string  `json:"osArch"`
	Size           *string `json:"size,omitempty"`
}

// DockerManifests Harness Manifests
type DockerManifests struct {
	ImageName       string                   `json:"imageName"`
	IsLatestVersion *bool                    `json:"isLatestVersion,omitempty"`
	Manifests       *[]DockerManifestDetails `json:"manifests,omitempty"`
	Version         string                   `json:"version"`
}

// Error defines model for Error.
type Error struct {
	// Code The http error code
	Code string `json:"code"`

	// Details Additional details about the error
	Details *map[string]interface{} `json:"details,omitempty"`

	// Message The reason the request failed
	Message string `json:"message"`
}

// HelmArtifactDetail Helm Artifact Detail
type HelmArtifactDetail struct {
	Artifact        *string `json:"artifact,omitempty"`
	CreatedAt       *string `json:"createdAt,omitempty"`
	DownloadsCount  *int64  `json:"downloadsCount,omitempty"`
	IsLatestVersion *bool   `json:"isLatestVersion,omitempty"`
	ModifiedAt      *string `json:"modifiedAt,omitempty"`

	// PackageType refers to package
	PackageType  PackageType `json:"packageType"`
	PullCommand  *string     `json:"pullCommand,omitempty"`
	RegistryPath string      `json:"registryPath"`
	Size         *string     `json:"size,omitempty"`
	Url          string      `json:"url"`
	Version      string      `json:"version"`
}

// HelmArtifactManifest Helm Artifact Manifest
type HelmArtifactManifest struct {
	Manifest string `json:"manifest"`
}

// ListArtifact A list of Artifacts
type ListArtifact struct {
	// Artifacts A list of Artifact
	Artifacts []ArtifactMetadata `json:"artifacts"`

	// ItemCount The total number of items
	ItemCount *int64 `json:"itemCount,omitempty"`

	// PageCount The total number of pages
	PageCount *int64 `json:"pageCount,omitempty"`

	// PageIndex The current page
	PageIndex *int64 `json:"pageIndex,omitempty"`

	// PageSize The number of items per page
	PageSize *int `json:"pageSize,omitempty"`
}

// ListArtifactLabel A list of Harness Artifact Labels
type ListArtifactLabel struct {
	// ItemCount The total number of items
	ItemCount *int64   `json:"itemCount,omitempty"`
	Labels    []string `json:"labels"`

	// PageCount The total number of pages
	PageCount *int64 `json:"pageCount,omitempty"`

	// PageIndex The current page
	PageIndex *int64 `json:"pageIndex,omitempty"`

	// PageSize The number of items per page
	PageSize *int `json:"pageSize,omitempty"`
}

// ListArtifactVersion A list of Artifact versions
type ListArtifactVersion struct {
	// ArtifactVersions A list of Artifact versions
	ArtifactVersions *[]ArtifactVersionMetadata `json:"artifactVersions,omitempty"`

	// ItemCount The total number of items
	ItemCount *int64 `json:"itemCount,omitempty"`

	// PageCount The total number of pages
	PageCount *int64 `json:"pageCount,omitempty"`

	// PageIndex The current page
	PageIndex *int64 `json:"pageIndex,omitempty"`

	// PageSize The number of items per page
	PageSize *int `json:"pageSize,omitempty"`
}

// ListRegistry A list of Harness Artifact Registries
type ListRegistry struct {
	// ItemCount The total number of items
	ItemCount *int64 `json:"itemCount,omitempty"`

	// PageCount The total number of pages
	PageCount *int64 `json:"pageCount,omitempty"`

	// PageIndex The current page
	PageIndex *int64 `json:"pageIndex,omitempty"`

	// PageSize The number of items per page
	PageSize *int `json:"pageSize,omitempty"`

	// Registries A list of Harness Artifact Registries
	Registries []RegistryMetadata `json:"registries"`
}

// ListRegistryArtifact A list of Artifacts
type ListRegistryArtifact struct {
	// Artifacts A list of Artifact
	Artifacts []RegistryArtifactMetadata `json:"artifacts"`

	// ItemCount The total number of items
	ItemCount *int64 `json:"itemCount,omitempty"`

	// PageCount The total number of pages
	PageCount *int64 `json:"pageCount,omitempty"`

	// PageIndex The current page
	PageIndex *int64 `json:"pageIndex,omitempty"`

	// PageSize The number of items per page
	PageSize *int `json:"pageSize,omitempty"`
}

// PackageType refers to package
type PackageType string

// Registry Harness Artifact Registry
type Registry struct {
	AllowedPattern *[]string        `json:"allowedPattern,omitempty"`
	BlockedPattern *[]string        `json:"blockedPattern,omitempty"`
	CleanupPolicy  *[]CleanupPolicy `json:"cleanupPolicy,omitempty"`

	// Config SubConfig specific for Virtual or Upstream Registry
	Config      *RegistryConfig `json:"config,omitempty"`
	CreatedAt   *string         `json:"createdAt,omitempty"`
	Description *string         `json:"description,omitempty"`
	Identifier  string          `json:"identifier"`
	Labels      *[]string       `json:"labels,omitempty"`
	ModifiedAt  *string         `json:"modifiedAt,omitempty"`

	// PackageType refers to package
	PackageType PackageType `json:"packageType"`
	Url         string      `json:"url"`
}

// RegistryArtifactMetadata Artifact Metadata
type RegistryArtifactMetadata struct {
	DownloadsCount *int64    `json:"downloadsCount,omitempty"`
	Labels         *[]string `json:"labels,omitempty"`
	LastModified   *string   `json:"lastModified,omitempty"`
	LatestVersion  string    `json:"latestVersion"`
	Name           string    `json:"name"`

	// PackageType refers to package
	PackageType        *PackageType `json:"packageType,omitempty"`
	RegistryIdentifier string       `json:"registryIdentifier"`
	RegistryPath       string       `json:"registryPath"`
}

// RegistryConfig SubConfig specific for Virtual or Upstream Registry
type RegistryConfig struct {
	// Type refers to type of registry i.e virtual or upstream
	Type  RegistryType `json:"type"`
	union json.RawMessage
}

// RegistryMetadata Harness Artifact Registry Metadata
type RegistryMetadata struct {
	ArtifactsCount *int64    `json:"artifactsCount,omitempty"`
	Description    *string   `json:"description,omitempty"`
	DownloadsCount *int64    `json:"downloadsCount,omitempty"`
	Identifier     string    `json:"identifier"`
	Labels         *[]string `json:"labels,omitempty"`
	LastModified   *string   `json:"lastModified,omitempty"`

	// PackageType refers to package
	PackageType  PackageType `json:"packageType"`
	Path         *string     `json:"path,omitempty"`
	RegistrySize *string     `json:"registrySize,omitempty"`

	// Type refers to type of registry i.e virtual or upstream
	Type RegistryType `json:"type"`
	Url  string       `json:"url"`
}

// RegistryRequest defines model for RegistryRequest.
type RegistryRequest struct {
	AllowedPattern *[]string        `json:"allowedPattern,omitempty"`
	BlockedPattern *[]string        `json:"blockedPattern,omitempty"`
	CleanupPolicy  *[]CleanupPolicy `json:"cleanupPolicy,omitempty"`

	// Config SubConfig specific for Virtual or Upstream Registry
	Config      *RegistryConfig `json:"config,omitempty"`
	Description *string         `json:"description,omitempty"`
	Identifier  string          `json:"identifier"`
	Labels      *[]string       `json:"labels,omitempty"`

	// PackageType refers to package
	PackageType PackageType `json:"packageType"`
	ParentRef   *string     `json:"parentRef,omitempty"`
}

// RegistryType refers to type of registry i.e virtual or upstream
type RegistryType string

// Status Indicates if the request was successful or not
type Status string

// UpstreamConfig Configuration for Harness Artifact UpstreamProxies
type UpstreamConfig struct {
	Auth *UpstreamConfig_Auth `json:"auth,omitempty"`

	// AuthType Authentication type
	AuthType AuthType              `json:"authType"`
	Source   *UpstreamConfigSource `json:"source,omitempty"`
	Url      *string               `json:"url,omitempty"`
}

// UpstreamConfig_Auth defines model for UpstreamConfig.Auth.
type UpstreamConfig_Auth struct {
	union json.RawMessage
}

// UpstreamConfigSource defines model for UpstreamConfig.Source.
type UpstreamConfigSource string

// UserPassword defines model for UserPassword.
type UserPassword struct {
	SecretIdentifier *string `json:"secretIdentifier,omitempty"`
	SecretSpaceId    *int    `json:"secretSpaceId,omitempty"`
	SecretSpacePath  *string `json:"secretSpacePath,omitempty"`
	UserName         string  `json:"userName"`
}

// VirtualConfig Configuration for Harness Virtual Artifact Registries
type VirtualConfig struct {
	UpstreamProxies *[]string `json:"upstreamProxies,omitempty"`
}

// LabelsParam defines model for LabelsParam.
type LabelsParam []string

// RegistryIdentifierParam defines model for RegistryIdentifierParam.
type RegistryIdentifierParam []string

// RegistryTypeParam defines model for RegistryTypeParam.
type RegistryTypeParam string

// ArtifactParam defines model for artifactParam.
type ArtifactParam string

// ArtifactPathParam defines model for artifactPathParam.
type ArtifactPathParam string

// DigestParam defines model for digestParam.
type DigestParam string

// FromDateParam defines model for fromDateParam.
type FromDateParam string

// LatestVersion defines model for latestVersion.
type LatestVersion bool

// PackageTypeParam defines model for packageTypeParam.
type PackageTypeParam []string

// PageNumber defines model for pageNumber.
type PageNumber int64

// PageSize defines model for pageSize.
type PageSize int64

// RecursiveParam defines model for recursiveParam.
type RecursiveParam bool

// RegistryRefPathParam defines model for registryRefPathParam.
type RegistryRefPathParam string

// SearchTerm defines model for searchTerm.
type SearchTerm string

// SortField defines model for sortField.
type SortField string

// SortOrder defines model for sortOrder.
type SortOrder string

// SpaceRefPathParam defines model for spaceRefPathParam.
type SpaceRefPathParam string

// SpaceRefQueryParam defines model for spaceRefQueryParam.
type SpaceRefQueryParam string

// ToDateParam defines model for toDateParam.
type ToDateParam string

// VersionParam defines model for versionParam.
type VersionParam string

// VersionPathParam defines model for versionPathParam.
type VersionPathParam string

// ArtifactLabelResponse defines model for ArtifactLabelResponse.
type ArtifactLabelResponse struct {
	// Data Harness Artifact Summary
	Data ArtifactSummary `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ArtifactStatsResponse defines model for ArtifactStatsResponse.
type ArtifactStatsResponse struct {
	// Data Harness Artifact Stats
	Data ArtifactStats `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ArtifactSummaryResponse defines model for ArtifactSummaryResponse.
type ArtifactSummaryResponse struct {
	// Data Harness Artifact Summary
	Data ArtifactSummary `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ArtifactVersionSummaryResponse defines model for ArtifactVersionSummaryResponse.
type ArtifactVersionSummaryResponse struct {
	// Data Docker Artifact Version Summary
	Data ArtifactVersionSummary `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// BadRequest defines model for BadRequest.
type BadRequest Error

// ClientSetupDetailsResponse defines model for ClientSetupDetailsResponse.
type ClientSetupDetailsResponse struct {
	// Data Client Setup Details
	Data ClientSetupDetails `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// DockerArtifactDetailResponse defines model for DockerArtifactDetailResponse.
type DockerArtifactDetailResponse struct {
	// Data Docker Artifact Detail
	Data DockerArtifactDetail `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// DockerArtifactManifestResponse defines model for DockerArtifactManifestResponse.
type DockerArtifactManifestResponse struct {
	// Data Docker Artifact Manifest
	Data DockerArtifactManifest `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// DockerLayersResponse defines model for DockerLayersResponse.
type DockerLayersResponse struct {
	// Data Harness Layers Summary
	Data DockerLayersSummary `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// DockerManifestsResponse defines model for DockerManifestsResponse.
type DockerManifestsResponse struct {
	// Data Harness Manifests
	Data DockerManifests `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// HelmArtifactDetailResponse defines model for HelmArtifactDetailResponse.
type HelmArtifactDetailResponse struct {
	// Data Helm Artifact Detail
	Data HelmArtifactDetail `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// HelmArtifactManifestResponse defines model for HelmArtifactManifestResponse.
type HelmArtifactManifestResponse struct {
	// Data Helm Artifact Manifest
	Data HelmArtifactManifest `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// InternalServerError defines model for InternalServerError.
type InternalServerError Error

// ListArtifactLabelResponse defines model for ListArtifactLabelResponse.
type ListArtifactLabelResponse struct {
	// Data A list of Harness Artifact Labels
	Data ListArtifactLabel `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ListArtifactResponse defines model for ListArtifactResponse.
type ListArtifactResponse struct {
	// Data A list of Artifacts
	Data ListArtifact `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ListArtifactVersionResponse defines model for ListArtifactVersionResponse.
type ListArtifactVersionResponse struct {
	// Data A list of Artifact versions
	Data ListArtifactVersion `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ListRegistryArtifactResponse defines model for ListRegistryArtifactResponse.
type ListRegistryArtifactResponse struct {
	// Data A list of Artifacts
	Data ListRegistryArtifact `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// ListRegistryResponse defines model for ListRegistryResponse.
type ListRegistryResponse struct {
	// Data A list of Harness Artifact Registries
	Data ListRegistry `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// NotFound defines model for NotFound.
type NotFound Error

// RegistryResponse defines model for RegistryResponse.
type RegistryResponse struct {
	// Data Harness Artifact Registry
	Data Registry `json:"data"`

	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// Success defines model for Success.
type Success struct {
	// Status Indicates if the request was successful or not
	Status Status `json:"status"`
}

// Unauthenticated defines model for Unauthenticated.
type Unauthenticated Error

// Unauthorized defines model for Unauthorized.
type Unauthorized Error

// CreateRegistryParams defines parameters for CreateRegistry.
type CreateRegistryParams struct {
	// SpaceRef Unique space path
	SpaceRef *SpaceRefQueryParam `form:"space_ref,omitempty" json:"space_ref,omitempty"`
}

// ListArtifactLabelsParams defines parameters for ListArtifactLabels.
type ListArtifactLabelsParams struct {
	// Page Current page number
	Page *PageNumber `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of items per page
	Size *PageSize `form:"size,omitempty" json:"size,omitempty"`

	// SearchTerm search Term.
	SearchTerm *SearchTerm `form:"search_term,omitempty" json:"search_term,omitempty"`
}

// GetArtifactStatsForRegistryParams defines parameters for GetArtifactStatsForRegistry.
type GetArtifactStatsForRegistryParams struct {
	// From Date. Format - MM/DD/YYYY
	From *FromDateParam `form:"from,omitempty" json:"from,omitempty"`

	// To Date. Format - MM/DD/YYYY
	To *ToDateParam `form:"to,omitempty" json:"to,omitempty"`
}

// GetArtifactStatsParams defines parameters for GetArtifactStats.
type GetArtifactStatsParams struct {
	// From Date. Format - MM/DD/YYYY
	From *FromDateParam `form:"from,omitempty" json:"from,omitempty"`

	// To Date. Format - MM/DD/YYYY
	To *ToDateParam `form:"to,omitempty" json:"to,omitempty"`
}

// GetDockerArtifactDetailsParams defines parameters for GetDockerArtifactDetails.
type GetDockerArtifactDetailsParams struct {
	// Digest Digest.
	Digest DigestParam `form:"digest" json:"digest"`
}

// GetDockerArtifactLayersParams defines parameters for GetDockerArtifactLayers.
type GetDockerArtifactLayersParams struct {
	// Digest Digest.
	Digest DigestParam `form:"digest" json:"digest"`
}

// GetDockerArtifactManifestParams defines parameters for GetDockerArtifactManifest.
type GetDockerArtifactManifestParams struct {
	// Digest Digest.
	Digest DigestParam `form:"digest" json:"digest"`
}

// GetAllArtifactVersionsParams defines parameters for GetAllArtifactVersions.
type GetAllArtifactVersionsParams struct {
	// Page Current page number
	Page *PageNumber `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of items per page
	Size *PageSize `form:"size,omitempty" json:"size,omitempty"`

	// SortOrder sortOrder
	SortOrder *SortOrder `form:"sort_order,omitempty" json:"sort_order,omitempty"`

	// SortField sortField
	SortField *SortField `form:"sort_field,omitempty" json:"sort_field,omitempty"`

	// SearchTerm search Term.
	SearchTerm *SearchTerm `form:"search_term,omitempty" json:"search_term,omitempty"`
}

// GetAllArtifactsByRegistryParams defines parameters for GetAllArtifactsByRegistry.
type GetAllArtifactsByRegistryParams struct {
	// Label Label.
	Label *LabelsParam `form:"label,omitempty" json:"label,omitempty"`

	// Page Current page number
	Page *PageNumber `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of items per page
	Size *PageSize `form:"size,omitempty" json:"size,omitempty"`

	// SortOrder sortOrder
	SortOrder *SortOrder `form:"sort_order,omitempty" json:"sort_order,omitempty"`

	// SortField sortField
	SortField *SortField `form:"sort_field,omitempty" json:"sort_field,omitempty"`

	// SearchTerm search Term.
	SearchTerm *SearchTerm `form:"search_term,omitempty" json:"search_term,omitempty"`
}

// GetClientSetupDetailsParams defines parameters for GetClientSetupDetails.
type GetClientSetupDetailsParams struct {
	// Artifact Artifat
	Artifact *ArtifactParam `form:"artifact,omitempty" json:"artifact,omitempty"`

	// Version Version
	Version *VersionParam `form:"version,omitempty" json:"version,omitempty"`
}

// GetArtifactStatsForSpaceParams defines parameters for GetArtifactStatsForSpace.
type GetArtifactStatsForSpaceParams struct {
	// From Date. Format - MM/DD/YYYY
	From *FromDateParam `form:"from,omitempty" json:"from,omitempty"`

	// To Date. Format - MM/DD/YYYY
	To *ToDateParam `form:"to,omitempty" json:"to,omitempty"`
}

// GetAllArtifactsParams defines parameters for GetAllArtifacts.
type GetAllArtifactsParams struct {
	// RegIdentifier Registry Identifier
	RegIdentifier *RegistryIdentifierParam `form:"reg_identifier,omitempty" json:"reg_identifier,omitempty"`

	// Page Current page number
	Page *PageNumber `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of items per page
	Size *PageSize `form:"size,omitempty" json:"size,omitempty"`

	// SortOrder sortOrder
	SortOrder *SortOrder `form:"sort_order,omitempty" json:"sort_order,omitempty"`

	// SortField sortField
	SortField *SortField `form:"sort_field,omitempty" json:"sort_field,omitempty"`

	// SearchTerm search Term.
	SearchTerm *SearchTerm `form:"search_term,omitempty" json:"search_term,omitempty"`

	// LatestVersion Latest Version Filter.
	LatestVersion *LatestVersion `form:"latest_version,omitempty" json:"latest_version,omitempty"`

	// PackageType Registry Package Type
	PackageType *PackageTypeParam `form:"package_type,omitempty" json:"package_type,omitempty"`
}

// GetAllRegistriesParams defines parameters for GetAllRegistries.
type GetAllRegistriesParams struct {
	// PackageType Registry Package Type
	PackageType *PackageTypeParam `form:"package_type,omitempty" json:"package_type,omitempty"`

	// Type Registry Type
	Type *GetAllRegistriesParamsType `form:"type,omitempty" json:"type,omitempty"`

	// Page Current page number
	Page *PageNumber `form:"page,omitempty" json:"page,omitempty"`

	// Size Number of items per page
	Size *PageSize `form:"size,omitempty" json:"size,omitempty"`

	// SortOrder sortOrder
	SortOrder *SortOrder `form:"sort_order,omitempty" json:"sort_order,omitempty"`

	// SortField sortField
	SortField *SortField `form:"sort_field,omitempty" json:"sort_field,omitempty"`

	// SearchTerm search Term.
	SearchTerm *SearchTerm `form:"search_term,omitempty" json:"search_term,omitempty"`

	// Recursive Whether to list registries recursively.
	Recursive *RecursiveParam `form:"recursive,omitempty" json:"recursive,omitempty"`
}

// GetAllRegistriesParamsType defines parameters for GetAllRegistries.
type GetAllRegistriesParamsType string

// CreateRegistryJSONRequestBody defines body for CreateRegistry for application/json ContentType.
type CreateRegistryJSONRequestBody RegistryRequest

// ModifyRegistryJSONRequestBody defines body for ModifyRegistry for application/json ContentType.
type ModifyRegistryJSONRequestBody RegistryRequest

// UpdateArtifactLabelsJSONRequestBody defines body for UpdateArtifactLabels for application/json ContentType.
type UpdateArtifactLabelsJSONRequestBody ArtifactLabelRequest

// AsVirtualConfig returns the union data inside the RegistryConfig as a VirtualConfig
func (t RegistryConfig) AsVirtualConfig() (VirtualConfig, error) {
	var body VirtualConfig
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromVirtualConfig overwrites any union data inside the RegistryConfig as the provided VirtualConfig
func (t *RegistryConfig) FromVirtualConfig(v VirtualConfig) error {
	t.Type = "VIRTUAL"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeVirtualConfig performs a merge with any union data inside the RegistryConfig, using the provided VirtualConfig
func (t *RegistryConfig) MergeVirtualConfig(v VirtualConfig) error {
	t.Type = "VIRTUAL"

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsUpstreamConfig returns the union data inside the RegistryConfig as a UpstreamConfig
func (t RegistryConfig) AsUpstreamConfig() (UpstreamConfig, error) {
	var body UpstreamConfig
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromUpstreamConfig overwrites any union data inside the RegistryConfig as the provided UpstreamConfig
func (t *RegistryConfig) FromUpstreamConfig(v UpstreamConfig) error {
	t.Type = "UPSTREAM"

	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeUpstreamConfig performs a merge with any union data inside the RegistryConfig, using the provided UpstreamConfig
func (t *RegistryConfig) MergeUpstreamConfig(v UpstreamConfig) error {
	t.Type = "UPSTREAM"

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t RegistryConfig) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"type"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t RegistryConfig) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "UPSTREAM":
		return t.AsUpstreamConfig()
	case "VIRTUAL":
		return t.AsVirtualConfig()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t RegistryConfig) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	if err != nil {
		return nil, err
	}
	object := make(map[string]json.RawMessage)
	if t.union != nil {
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	}

	object["type"], err = json.Marshal(t.Type)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'type': %w", err)
	}

	b, err = json.Marshal(object)
	return b, err
}

func (t *RegistryConfig) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	object := make(map[string]json.RawMessage)
	err = json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["type"]; found {
		err = json.Unmarshal(raw, &t.Type)
		if err != nil {
			return fmt.Errorf("error reading 'type': %w", err)
		}
	}

	return err
}

// AsUserPassword returns the union data inside the UpstreamConfig_Auth as a UserPassword
func (t UpstreamConfig_Auth) AsUserPassword() (UserPassword, error) {
	var body UserPassword
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromUserPassword overwrites any union data inside the UpstreamConfig_Auth as the provided UserPassword
func (t *UpstreamConfig_Auth) FromUserPassword(v UserPassword) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeUserPassword performs a merge with any union data inside the UpstreamConfig_Auth, using the provided UserPassword
func (t *UpstreamConfig_Auth) MergeUserPassword(v UserPassword) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsAnonymous returns the union data inside the UpstreamConfig_Auth as a Anonymous
func (t UpstreamConfig_Auth) AsAnonymous() (Anonymous, error) {
	var body Anonymous
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAnonymous overwrites any union data inside the UpstreamConfig_Auth as the provided Anonymous
func (t *UpstreamConfig_Auth) FromAnonymous(v Anonymous) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAnonymous performs a merge with any union data inside the UpstreamConfig_Auth, using the provided Anonymous
func (t *UpstreamConfig_Auth) MergeAnonymous(v Anonymous) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsAccessKeySecretKey returns the union data inside the UpstreamConfig_Auth as a AccessKeySecretKey
func (t UpstreamConfig_Auth) AsAccessKeySecretKey() (AccessKeySecretKey, error) {
	var body AccessKeySecretKey
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAccessKeySecretKey overwrites any union data inside the UpstreamConfig_Auth as the provided AccessKeySecretKey
func (t *UpstreamConfig_Auth) FromAccessKeySecretKey(v AccessKeySecretKey) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAccessKeySecretKey performs a merge with any union data inside the UpstreamConfig_Auth, using the provided AccessKeySecretKey
func (t *UpstreamConfig_Auth) MergeAccessKeySecretKey(v AccessKeySecretKey) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t UpstreamConfig_Auth) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *UpstreamConfig_Auth) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
