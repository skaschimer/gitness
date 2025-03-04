/*
 * Copyright 2024 Harness, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { createContext, type PropsWithChildren } from 'react'
import { PageError, PageSpinner } from '@harnessio/uicore'
import { type FileDetailResponseResponse, useGetArtifactFilesQuery } from '@harnessio/react-har-service-client'

import { encodeRef } from '@ar/hooks/useGetSpaceRef'
import type { VersionDetailsPathParams } from '@ar/routes/types'
import { useDecodedParams, useGetSpaceRef, useParentHooks } from '@ar/hooks'
import type { UseUpdateQueryParamsReturn } from '@ar/__mocks__/hooks/useUpdateQueryParams'

import {
  ArtifactFileListPageQueryParams,
  useArtifactFileListQueryParamOptions
} from '../components/ArtifactFileListTable/utils'

interface VersionFilesProviderProps {
  data: FileDetailResponseResponse
  updateQueryParams: UseUpdateQueryParamsReturn<Partial<ArtifactFileListPageQueryParams>>['updateQueryParams']
  queryParams: ArtifactFileListPageQueryParams
  refetch: () => void
  sort: string[]
}

export const VersionFilesContext = createContext<VersionFilesProviderProps>({} as VersionFilesProviderProps)

const VersionFilesProvider = (props: PropsWithChildren<unknown>) => {
  const pathParams = useDecodedParams<VersionDetailsPathParams>()
  const spaceRef = useGetSpaceRef()

  const { useQueryParams, useUpdateQueryParams } = useParentHooks()
  const { updateQueryParams } = useUpdateQueryParams<Partial<ArtifactFileListPageQueryParams>>()

  const queryParamOptions = useArtifactFileListQueryParamOptions()
  const queryParams = useQueryParams<ArtifactFileListPageQueryParams>(queryParamOptions)
  const { page, size, sort } = queryParams

  const [sortField, sortOrder] = sort || []

  const {
    data,
    isFetching: loading,
    error,
    refetch
  } = useGetArtifactFilesQuery({
    registry_ref: spaceRef,
    artifact: encodeRef(pathParams.artifactIdentifier),
    version: pathParams.versionIdentifier,
    queryParams: {
      page,
      size,
      sort_field: sortField,
      sort_order: sortOrder
    }
  })
  const responseData = data?.content

  return (
    <>
      {loading ? <PageSpinner /> : null}
      {error && !loading ? <PageError message={error.message} onClick={() => refetch()} /> : null}
      {!error && !loading && responseData ? (
        <VersionFilesContext.Provider value={{ data: responseData, refetch, updateQueryParams, queryParams, sort }}>
          {props.children}
        </VersionFilesContext.Provider>
      ) : null}
    </>
  )
}

export default VersionFilesProvider
