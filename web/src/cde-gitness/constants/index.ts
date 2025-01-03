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
import vsCodeWebIcon from 'cde-gitness/assests/vsCodeWeb.svg?url'
import vsCodeIcon from 'cde-gitness/assests/VSCode.svg?url'
import intellijIcon from 'cde-gitness/assests/intellij.svg?url'
import type { StringsMap } from 'framework/strings/stringTypes'

export const docLink = 'https://developer.harness.io/docs/cloud-development-environments'

export enum IDEType {
  VSCODE = 'vs_code',
  VSCODEWEB = 'vs_code_web',
  INTELLIJ = 'intellij'
}

export interface ideType {
  label: keyof StringsMap
  value: string
  icon: any
}

export const getIDETypeOptions = (getString: any) => [
  {
    label: getString('cde.ide.desktop'),
    value: IDEType.VSCODE,
    icon: vsCodeIcon
  },
  {
    label: getString('cde.ide.browser'),
    value: IDEType.VSCODEWEB,
    icon: vsCodeWebIcon
  },
  {
    label: getString('cde.ide.intellij'),
    value: IDEType.INTELLIJ,
    icon: intellijIcon
  }
]

export enum EnumGitspaceCodeRepoType {
  GITHUB = 'github',
  GITLAB = 'gitlab',
  HARNESS_CODE = 'harness_code',
  BITBUCKET = 'bitbucket',
  UNKNOWN = 'unknown',
  GITNESS = 'gitness'
}

export enum GitspaceStatus {
  RUNNING = 'running',
  STOPPED = 'stopped',
  UNKNOWN = 'unknown',
  ERROR = 'error',
  STARTING = 'starting',
  STOPPING = 'stopping',
  UNINITIALIZED = 'uninitialized'
}

export interface GitspaceStatusTypesListItem {
  label: string
  value: GitspaceStatus
}

export const GitspaceStatusTypes = (getString: any) => [
  {
    label: getString('cde.gitspaceStatus.active'),
    value: GitspaceStatus.RUNNING
  },
  {
    label: getString('cde.gitspaceStatus.stopped'),
    value: GitspaceStatus.STOPPED
  },
  {
    label: getString('cde.gitspaceStatus.error'),
    value: GitspaceStatus.ERROR
  }
]

export const getIDEOption: any = (type = '', getString = null) => {
  let ideItem = null
  if (type && getString) {
    ideItem = getIDETypeOptions(getString).find((ide: ideType) => ide?.value === type)
  }
  return ideItem
}

export enum GitspaceOwnerType {
  SELF = 'self',
  ALL = 'all'
}

export interface GitspaceOwnerTypeListItem {
  label: string
  value: GitspaceOwnerType
}

export const GitspaceOwnerTypes = (getString: any) => [
  {
    label: getString('cde.gitspaceOwners.allGitspaces'),
    value: GitspaceOwnerType.ALL
  },
  {
    label: getString('cde.gitspaceOwners.myGitspaces'),
    value: GitspaceOwnerType.SELF
  }
]

export enum GitspaceActionType {
  START = 'start',
  STOP = 'stop'
}

export enum GitspaceRegion {
  USEast = 'us-east',
  USWest = 'us-west',
  Europe = 'Europe',
  Australia = 'Australia'
}

export enum SortByType {
  CREATED = 'created',
  LAST_USED = 'last_used',
  LAST_ACTIVATED = 'last_activated'
}

export interface SortByTypeListItem {
  label: string
  value: SortByType
}

export const SortByTypes = (getString: any) => [
  {
    label: getString('cde.created'),
    value: SortByType.CREATED
  },
  {
    label: getString('cde.lastUsed'),
    value: SortByType.LAST_USED
  },
  {
    label: getString('cde.lastStarted'),
    value: SortByType.LAST_ACTIVATED
  }
]
