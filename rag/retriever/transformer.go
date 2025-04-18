package retriever

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

type ReRankransformer struct {
	config *ReRankransformerConfig
}

type ReRankransformerConfig struct {
	Model           string
	ReturnDocuments bool
	TopK            int
	Query           string
	ScoreThreshold  float64
	ApiKey          string
}

func newReRankransformer(ctx context.Context, opt *ReRankransformerConfig, query string) (tfr document.Transformer, err error) {
	config := &ReRankransformerConfig{
		Model:           "gte-rerank",
		ReturnDocuments: false,
		TopK:            5,
		Query:           query,
		ScoreThreshold:  0.7,
	}
	if opt != nil {
		if opt.TopK > 0 {
			config.TopK = opt.TopK
		}
		config.ScoreThreshold = opt.ScoreThreshold
		config.ReturnDocuments = opt.ReturnDocuments
		config.Model = opt.Model
		config.ApiKey = opt.ApiKey
	}
	tfr = &ReRankransformer{config: config}
	return tfr, nil
}

type ReRankConfigInput struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

type ReRankConfigParams struct {
	ReturnDocuments bool `json:"return_documents"`
	TopK            int  `json:"top_n"`
}

type ReRankConfig struct {
	Model      string              `json:"model"`
	Input      *ReRankConfigInput  `json:"input"`
	Parameters *ReRankConfigParams `json:"parameters"`
}

type ReRankDataUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type ReRankDataOutputResult struct {
	Index int     `json:"index"`
	Score float64 `json:"relevance_score"`
}

type ReRankDataOutput struct {
	Results []*ReRankDataOutputResult `json:"results"`
}

type ReRankData struct {
	Output    *ReRankDataOutput `json:"output"`
	Usage     *ReRankDataUsage  `json:"usage"`
	RequestId string            `json:"request_id"`
}

func (impl *ReRankransformer) rerankAli(config *ReRankConfig) (*ReRankData, error) {
	param, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	var jsonData io.Reader
	//把[]byte 转成实现了read接口的Reader结构体
	if param != nil {
		jsonData = bytes.NewReader(param)
	}
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank", jsonData)
	if err != nil {
		err = fmt.Errorf("网络故障")
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", impl.config.ApiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reRankData := &ReRankData{}
	err = json.Unmarshal(resp_body, reRankData)
	if err != nil {
		fmt.Println(string(resp_body))
	}
	if reRankData.Output == nil {
		return nil, fmt.Errorf("no data")
	}
	return reRankData, err
}

func (impl *ReRankransformer) Transform(ctx context.Context, src []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	config := &ReRankConfig{
		Model: impl.config.Model,
		Input: &ReRankConfigInput{
			Query:     impl.config.Query,
			Documents: make([]string, 0),
		},
		Parameters: &ReRankConfigParams{
			ReturnDocuments: impl.config.ReturnDocuments,
			TopK:            impl.config.TopK,
		},
	}
	for _, v := range src {
		config.Input.Documents = append(config.Input.Documents, v.Content)
	}
	reRankData, err := impl.rerankAli(config)
	if err != nil {
		return src, nil
	}
	dst := make([]*schema.Document, 0)

	for i := 0; i < len(reRankData.Output.Results); i++ {
		res := reRankData.Output.Results[i]
		if res.Score >= impl.config.ScoreThreshold {
			dst = append(dst, src[res.Index])
		}

	}
	return dst, nil
}
