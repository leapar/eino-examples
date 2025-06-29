/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package retriever

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var systemPrompt = `
# 角色：软件教程使用帮助问答助手

## 核心能力
- **理解上下文内容，回答用户问题**。

## 互动指南
- **在回复之前**，请确保您：
  - **完全理解用户的需求和要求**。如果存在任何模糊之处，请向用户寻求澄清。

- **在提供帮助时**：
  - **根据用户请求的上下文，准确回答用户问题**。
  - **根据用户请求的上下文，有热点图片要输出热点图片地址**。
  - **不要过渡去深度思考修改上下文内容**。

- **如果请求超出了您的能力范围**：
  - **告诉用户联系客服寻求人工帮助**。


## 上下文信息
- **相关文档**：|-
==== doc start ====
  {documents}
==== doc end ====

在回答时，请考虑相关文档中提供的上下文，以便根据用户的需求量身定制您的回答。
`

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// newChatTemplate component initialization function of node 'ChatTemplate' in graph 'EinoAgent'
func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	// TODO Modify component configuration here.
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(systemPrompt),
			schema.MessagesPlaceholder("history", true),
			schema.UserMessage("{content}"),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}
