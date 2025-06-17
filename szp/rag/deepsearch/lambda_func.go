package deepsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func newLambda(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	return input.Query, nil
}

func newLambda2(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content": input.Query,
		"history": input.History,
		"date":    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

const SUB_QUERY_PROMPT = `To answer this question more comprehensively, please break down the original question into up to four sub-questions. Return as list of str.
If this is a very simple question and no decomposition is necessary, then keep the only one original question in the python code list.

Original Question: {original_query}


<EXAMPLE>
Example input:
"Explain deep learning"

Example output:
[
    "What is deep learning?",
    "What is the difference between deep learning and machine learning?",
    "What is the history of deep learning?"
]
</EXAMPLE>

请用json格式字符串数组返回数据内容:
`

func newRetriverLambda(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	// TODO Modify component configuration here.
	tpConfig := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.UserMessage(SUB_QUERY_PROMPT),
		},
	}
	ctp := prompt.FromMessages(tpConfig.FormatType, tpConfig.Templates...)

	msg, err := ctp.Format(ctx, map[string]any{
		"original_query": input.Query,
	})
	if err != nil {
		fmt.Println(err)
	}
	model, err := newChatModel(ctx)
	if err != nil {
		fmt.Println(err)
	}
	res, err := model.Generate(ctx, msg)
	if err != nil {
		fmt.Println(err)
	}
	subQueries, err := literalEval(res.Content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(subQueries)

	retriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}

	retDocs := make([]*schema.Document, 0)
	for _, subQuery := range subQueries {
		docs, err := retriever.Retrieve(ctx, subQuery)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, doc := range docs {
				skip := false
				for i := 0; i < len(retDocs); i++ {
					preDoc := retDocs[i]
					if preDoc.ID == doc.ID {
						skip = true
						break
					}
				}
				if !skip {
					retDocs = append(retDocs, doc)
				}
			}
		}
	}

	result := ""
	for i := 0; i < len(retDocs); i++ {
		preDoc := retDocs[i]
		result += fmt.Sprintf("<Document %d>\n%s\n<\\Document %d>\n", i, preDoc.Content, i)
	}

	questions := ""
	for i := 0; i < len(subQueries); i++ {
		txt := subQueries[i]
		questions += fmt.Sprintf("%d. %s\n", i, txt)
	}

	return map[string]any{
		"history":        input.History,
		"date":           time.Now().Format("2006-01-02 15:04:05"),
		"question":       input.Query,
		"mini_questions": questions,
		"mini_chunk_str": result,
	}, nil
}
