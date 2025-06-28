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

	"github.com/cloudwego/eino-ext/components/indexer/chromem"
	"github.com/cloudwego/eino/components/indexer"
	cg "github.com/philippgille/chromem-go"
)

var db *cg.DB

func init() {
	var err error
	db, err = cg.NewPersistentDB("./db", true)
	if err != nil {
		fmt.Println(err)
	}
}

// newIndexer component initialization function of node 'RedisIndexer' in graph 'KnowledgeIndexing'
func newChromemIndexer(ctx context.Context) (idr indexer.Indexer, err error) {

	config := &chromem.IndexerConfig{
		Client: db,
	}

	embedding, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embedding
	idr, err = chromem.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
