module github.com/cloudwego/eino-examples

go 1.24.0

toolchain go1.24.2

require (
	github.com/bytedance/sonic v1.13.2
	github.com/cloudwego/eino v0.3.20
	github.com/cloudwego/eino-ext/components/document/loader/file v0.0.0-20250415073426-726b929afbc2
	github.com/cloudwego/eino-ext/components/document/parser/html v0.0.0-20250117061805-cd80d1780d76
	github.com/cloudwego/eino-ext/components/document/parser/pdf v0.0.0-20250117061805-cd80d1780d76
	github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown v0.0.0-20250415073426-726b929afbc2
	github.com/cloudwego/eino-ext/components/embedding/ollama v0.0.0-00010101000000-000000000000
	github.com/cloudwego/eino-ext/components/indexer/chromem v0.0.0-00010101000000-000000000000
	github.com/cloudwego/eino-ext/components/model/ark v0.1.0
	github.com/cloudwego/eino-ext/components/model/deepseek v0.0.0-20250221090944-e8ef7aabbe10
	github.com/cloudwego/eino-ext/components/model/openai v0.0.0-20250221090944-e8ef7aabbe10
	github.com/cloudwego/eino-ext/components/retriever/chromem v0.0.0-00010101000000-000000000000
	github.com/cloudwego/eino-ext/components/retriever/volc_vikingdb v0.0.0-20250319082935-6219ec437e56
	github.com/cloudwego/eino-ext/components/tool/duckduckgo v0.0.0-20250221090944-e8ef7aabbe10
	github.com/cloudwego/eino-ext/components/tool/mcp v0.0.0-20250415073426-726b929afbc2
	github.com/cloudwego/eino-ext/devops v0.1.7
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/getkin/kin-openapi v0.118.0
	github.com/mark3labs/mcp-go v0.20.1
	github.com/ollama/ollama v0.6.5
	github.com/philippgille/chromem-go v0.7.0
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/eino-ext/components/model/ollama v0.0.0-20250417123744-154d7ca4d3cd
	github.com/cloudwego/eino-ext/libs/acl/openai v0.0.0-20250221090944-e8ef7aabbe10 // indirect
	github.com/cohesion-org/deepseek-go v1.2.3 // indirect
	github.com/dslipak/pdf v0.0.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/invopop/yaml v0.3.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/matoous/go-nanoid v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/nikolalohinski/gonja v1.5.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sashabaranov/go-openai v1.37.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/slongfield/pyfmt v0.0.0-20220222012616-ea85ff4c361f // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/volcengine/volc-sdk-golang v1.0.196 // indirect
	github.com/volcengine/volcengine-go-sdk v1.0.185 // indirect
	github.com/yargevad/filepathx v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/crypto v0.34.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require github.com/cloudwego/eino-ext/components/tool/bingsearch v0.0.0-20250417123744-154d7ca4d3cd

replace (
	github.com/cloudwego/eino-ext/components/embedding/ollama => D:\gocodes\ai\szp-eino-ext\components\embedding\ollama
	github.com/cloudwego/eino-ext/components/indexer/chromem => D:\gocodes\ai\szp-eino-ext\components\indexer\chromem
	github.com/cloudwego/eino-ext/components/model/ollama => D:\gocodes\ai\szp-eino-ext\components\model\ollama
	github.com/cloudwego/eino-ext/components/retriever => D:\gocodes\ai\szp-eino-ext\components\retriever
	github.com/cloudwego/eino-ext/components/retriever/chromem => D:\gocodes\ai\szp-eino-ext\components\retriever\chromem
	github.com/cloudwego/eino => D:\gocodes\ai\eino

)
