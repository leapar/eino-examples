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

package deepsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var systemPrompt = `
You are a AI content analysis expert, good at summarizing content. Please summarize a specific and detailed answer or report based on the previous queries and the retrieved document chunks.

Original Query: {question}

Previous Sub Queries: {mini_questions}

Related Chunks: 
{mini_chunk_str}
`

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

/*
Parse a string response into a Python object using ast.literal_eval.

This method attempts to extract and parse JSON or Python literals from the response content,
handling various formats like code blocks and special tags.

Args:
	response_content: The string content to parse.

Returns:
	The parsed Python object.

Raises:
	ValueError: If the response content cannot be parsed.
*/

func literalEval(responseContent string) ([]string, error) {
	// Trim the content
	responseContent = strings.TrimSpace(responseContent)

	// Remove content between <think> and </think>
	if strings.Contains(responseContent, "<think>") && strings.Contains(responseContent, "</think>") {
		endOfThink := strings.Index(responseContent, "</think>") + len("</think>")
		responseContent = responseContent[endOfThink:]
	}

	// Handle code blocks
	if strings.HasPrefix(responseContent, "```") && strings.HasSuffix(responseContent, "```") {
		if strings.HasPrefix(responseContent, "```python") {
			responseContent = responseContent[9 : len(responseContent)-3]
		} else if strings.HasPrefix(responseContent, "```json") {
			responseContent = responseContent[7 : len(responseContent)-3]
		} else if strings.HasPrefix(responseContent, "```str") {
			responseContent = responseContent[6 : len(responseContent)-3]
		} else if strings.HasPrefix(responseContent, "```\n") {
			responseContent = responseContent[4 : len(responseContent)-3]
		} else {
			return nil, fmt.Errorf("invalid code block format")
		}
	}

	// Try to parse as JSON
	result := make([]string, 0)
	if err := json.Unmarshal([]byte(responseContent), &result); err == nil {
		return result, nil
	}

	// Try to find JSON/List patterns
	re := regexp.MustCompile(`(\[.*?\]|\{.*?\})`)
	matches := re.FindAllString(responseContent, -1)

	if len(matches) != 1 {
		return nil, fmt.Errorf("invalid JSON/List format for response content:\n%s", responseContent)
	}

	jsonPart := matches[0]
	if err := json.Unmarshal([]byte(jsonPart), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON part: %v", err)
	}

	return result, nil
}

// newChatTemplate component initialization function of node 'ChatTemplate' in graph 'EinoAgent'
func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	// TODO Modify component configuration here.
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(systemPrompt),
			schema.UserMessage("{question}"),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}
