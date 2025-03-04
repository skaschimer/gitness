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

import React from 'react'
import { Layout } from '@harnessio/uicore'
import type { ArtifactVersionSummary } from '@harnessio/react-har-service-client'

import { String } from '@ar/frameworks/strings'
import { RepositoryPackageType } from '@ar/common/types'
import { VersionListColumnEnum } from '@ar/pages/version-list/components/VersionListTable/types'
import ArtifactActions from '@ar/pages/artifact-details/components/ArtifactActions/ArtifactActions'
import VersionListTable, {
  type CommonVersionListTableProps
} from '@ar/pages/version-list/components/VersionListTable/VersionListTable'
import {
  type ArtifactActionProps,
  type VersionDetailsHeaderProps,
  type VersionDetailsTabProps,
  type VersionListTableProps,
  VersionStep
} from '@ar/frameworks/Version/Version'

import NpmVersionOverviewPage from './pages/overview/NpmVersionOverviewPage'
import { VersionDetailsTab } from '../components/VersionDetailsTabs/constants'
import VersionDetailsHeaderContent from '../components/VersionDetailsHeaderContent/VersionDetailsHeaderContent'
import NpmVersionArtifactDetailsPage from './pages/artifact-dertails/NpmVersionArtifactDetailsPage'

export class NpmVersionType extends VersionStep<ArtifactVersionSummary> {
  protected packageType = RepositoryPackageType.NPM
  protected allowedVersionDetailsTabs: VersionDetailsTab[] = [
    VersionDetailsTab.OVERVIEW,
    VersionDetailsTab.ARTIFACT_DETAILS,
    VersionDetailsTab.CODE
  ]

  versionListTableColumnConfig: CommonVersionListTableProps['columnConfigs'] = {
    [VersionListColumnEnum.Name]: { width: '100%' },
    [VersionListColumnEnum.Size]: { width: '100%' },
    [VersionListColumnEnum.DownloadCount]: { width: '100%' },
    [VersionListColumnEnum.PullCommand]: { width: '100%' },
    [VersionListColumnEnum.LastModified]: { width: '100%' }
  }

  renderVersionListTable(props: VersionListTableProps): JSX.Element {
    return <VersionListTable {...props} columnConfigs={this.versionListTableColumnConfig} />
  }

  renderVersionDetailsHeader(props: VersionDetailsHeaderProps<ArtifactVersionSummary>): JSX.Element {
    return <VersionDetailsHeaderContent {...props} />
  }

  renderVersionDetailsTab(props: VersionDetailsTabProps): JSX.Element {
    switch (props.tab) {
      case VersionDetailsTab.OVERVIEW:
        return <NpmVersionOverviewPage />
      case VersionDetailsTab.ARTIFACT_DETAILS:
        return <NpmVersionArtifactDetailsPage />
      case VersionDetailsTab.OSS:
        return (
          <Layout.Vertical spacing="xlarge">
            <NpmVersionOverviewPage />
            <NpmVersionArtifactDetailsPage />
          </Layout.Vertical>
        )
      default:
        return <String stringID="tabNotFound" />
    }
  }

  renderArtifactActions(props: ArtifactActionProps): JSX.Element {
    return <ArtifactActions {...props} />
  }
}
