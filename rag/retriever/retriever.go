package retriever

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/retriever/chromem"
	"github.com/cloudwego/eino/components/retriever"
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

// newRetriever component initialization function of node 'RedisRetriever' in graph 'EinoAgent'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	config := &chromem.RetrieverConfig{
		Client:         db,
		TopK:           8,
		ScoreThreshold: 0.7,
	}
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	rtr, err = chromem.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}
