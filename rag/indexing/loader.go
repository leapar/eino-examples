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

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/parser/csv"
	"github.com/cloudwego/eino-ext/components/document/parser/doc"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino-ext/components/document/parser/xlsx"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
)

// newLoader component initialization function of node 'FileLoader' in graph 'KnowledgeIndexing'
func newLoader(ctx context.Context) (ldr document.Loader, err error) {
	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{ToPages: true})
	if err != nil {
		return nil, err
	}

	docParser, err := doc.NewDocParser()
	if err != nil {
		return nil, err
	}

	csvParser, err := csv.NewCsvParser()
	if err != nil {
		return nil, err
	}

	xlsxParser, err := xlsx.NewXlsxParser()
	if err != nil {
		return nil, err
	}
	txtParser := &parser.TextParser{}

	extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		FallbackParser: parser.TextParser{},
		Parsers: map[string]parser.Parser{
			".pdf":  pdfParser,
			".txt":  txtParser,
			".md":   txtParser,
			".doc":  docParser,
			".docx": docParser,
			".xlsx": xlsxParser,
			".csv":  csvParser,
		},
	})
	if err != nil {
		return nil, err
	}

	// TODO Modify component configuration here.
	config := &file.FileLoaderConfig{
		Parser: extParser,
	}
	ldr, err = file.NewFileLoader(ctx, config)
	if err != nil {
		return nil, err
	}

	return ldr, nil
}
