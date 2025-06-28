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

// 在doc中拼接目录结构
func appendHeadersLambda(ctx context.Context, docs []*schema.Document, opts ...any) (output []*schema.Document, err error) {

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
		ChromemIndexer   = "ChromemIndexer"
		HeaderAppender   = "HeaderAppender"
	)

	g := compose.NewGraph[document.Source, []string]()
	fileLoaderKeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}

	//问答加载节点
	_ = g.AddLoaderNode(FileLoader, fileLoaderKeyOfLoader)

	//切割器节点
	markdownSplitterTransformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(MarkdownSplitter, markdownSplitterTransformer)
	_ = g.AddLambdaNode(HeaderAppender, compose.InvokableLambdaWithOption(appendHeadersLambda), compose.WithNodeName("DocsAppendHeaders"))

	//向量化节点
	chromemIndexer, err := newChromemIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(ChromemIndexer, chromemIndexer)

	//构建图
	_ = g.AddEdge(compose.START, FileLoader)
	_ = g.AddEdge(FileLoader, MarkdownSplitter)
	_ = g.AddEdge(MarkdownSplitter, HeaderAppender)
	_ = g.AddEdge(HeaderAppender, ChromemIndexer)
	_ = g.AddEdge(ChromemIndexer, compose.END)

	r, err = g.Compile(ctx, compose.WithGraphName("KnowledgeIndexing"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
