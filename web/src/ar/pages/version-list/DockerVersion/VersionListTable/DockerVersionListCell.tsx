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
import type { Cell, CellValue, ColumnInstance, Renderer, Row, TableInstance } from 'react-table'
import { Layout } from '@harnessio/uicore'
import { Icon } from '@harnessio/icons'
import type { ArtifactVersionMetadata } from '@harnessio/react-har-service-client'

import TableCells from '@ar/components/TableCells/TableCells'

import css from './DockerVersionListTable.module.scss'

type CellTypeWithActions<D extends Record<string, any>, V = any> = TableInstance<D> & {
  column: ColumnInstance<D>
  row: Row<D>
  cell: Cell<D, V>
  value: CellValue<V>
}

type CellType = Renderer<CellTypeWithActions<ArtifactVersionMetadata>>

export const DockerVersionNameCell: CellType = ({ value }) => {
  return (
    <Layout.Horizontal className={css.nameCellContainer} spacing="small">
      <Icon name="store-artifact-bundle" size={24} />
      <TableCells.TextCell value={value} />
    </Layout.Horizontal>
  )
}
