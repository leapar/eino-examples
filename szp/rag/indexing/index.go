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
	"fmt"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func newLambda(ctx context.Context, docs []*schema.Document, opts ...any) (output []*schema.Document, err error) {

	headers := []string{
		"title1",
		"title2",
		"title3",
		"title4",
		"title5",
	}

	for _, doc := range docs {
		title := `这段内容目录结构是[%s]，用"___"分割。\r\n`
		lanmu := ""
		for _, header := range headers {
			val := doc.MetaData[header]
			if val != nil {
				if len(lanmu) > 0 {
					lanmu = fmt.Sprintf("%s___%s", lanmu, val.(string))
				} else {
					lanmu = val.(string)
				}
			}
		}

		title = fmt.Sprintf(title, lanmu)

		doc.Content = title + doc.Content
	}

	return docs, nil
}

func BuildKnowledgeIndexing(ctx context.Context) (r compose.Runnable[document.Source, []string], err error) {
	const (
		FileLoader       = "FileLoader"
		MarkdownSplitter = "MarkdownSplitter"
		RedisIndexer     = "RedisIndexer"
		HeaderAppender   = "HeaderAppender"
	)
	g := compose.NewGraph[document.Source, []string]()
	fileLoaderKeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(FileLoader, fileLoaderKeyOfLoader)
	markdownSplitterKeyOfDocumentTransformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}

	_ = g.AddLambdaNode(HeaderAppender, compose.InvokableLambdaWithOption(newLambda), compose.WithNodeName("DocsAppendHeaders"))

	_ = g.AddDocumentTransformerNode(MarkdownSplitter, markdownSplitterKeyOfDocumentTransformer)
	redisIndexerKeyOfIndexer, err := newIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(RedisIndexer, redisIndexerKeyOfIndexer)
	_ = g.AddEdge(compose.START, FileLoader)
	_ = g.AddEdge(RedisIndexer, compose.END)
	_ = g.AddEdge(FileLoader, MarkdownSplitter)
	_ = g.AddEdge(MarkdownSplitter, HeaderAppender)
	_ = g.AddEdge(HeaderAppender, RedisIndexer)

	r, err = g.Compile(ctx, compose.WithGraphName("KnowledgeIndexing"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
