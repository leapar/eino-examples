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

package indexing

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

func newDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	// TODO Modify component configuration here.
	config := &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":     "title1",
			"##":    "title2",
			"###":   "title3",
			"####":  "title4",
			"#####": "title5",
		},
		TrimHeaders: false}
	tfr, err = markdown.NewHeaderSplitter(ctx, config)
	if err != nil {
		return nil, err
	}

	// tfr, err = recursive.NewSplitter(ctx, &recursive.Config{
	// 	ChunkSize:   1500,
	// 	OverlapSize: 300,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	return tfr, nil
}
