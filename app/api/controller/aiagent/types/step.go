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

package types

type PipelineStepData struct {
	Yaml string `json:"yaml_step"`
}

// create.
type GeneratePipelineStepInput struct {
	Prompt       string            `json:"prompt"`
	Metadata     map[string]string `json:"metadata"`
	Conversation []Conversation    `json:"conversation"`
}

type GeneratePipelineStepOutput struct {
	Error string           `json:"error"`
	Data  PipelineStepData `json:"data"`
}

// update.
type UpdatePipelineStepInput struct {
	Prompt       string           `json:"prompt"`
	Data         PipelineStepData `json:"data"`
	Conversation []Conversation   `json:"conversation"`
}

type UpdatePipelineStepOutput struct {
	Error string           `json:"error"`
	Data  PipelineStepData `json:"data"`
}

func (in *GeneratePipelineStepInput) GetConversation() []Conversation {
	return in.Conversation
}

func (in *GeneratePipelineStepInput) GetPrompt() string {
	return in.Prompt
}

func (in *GeneratePipelineStepInput) GetValidationPrompt() string {
	return "Create a step-yaml with the following query: " + in.GetPrompt()
}
