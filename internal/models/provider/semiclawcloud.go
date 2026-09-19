package provider

import "github.com/vagawind/semiclaw/internal/types"

const (
	ProviderSemiClawCloud ProviderName = "weknoracloud"

	// SemiClawCloudBaseURL SemiClawCloud 服务硬编码 Base URL（统一入口，路径由各实现拼接）
	SemiClawCloudBaseURL = "https://weknora.weixin.qq.com"
)

type SemiClawCloudProvider struct{}

func init() {
	Register(&SemiClawCloudProvider{})
}

func (p *SemiClawCloudProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderSemiClawCloud,
		DisplayName: "SemiClawCloud",
		Description: "SemiClaw云服务，模型：chat, embedding, rerank, vlm",
		DefaultURLs: map[types.ModelType]string{
			types.ModelTypeKnowledgeQA: SemiClawCloudBaseURL,
			types.ModelTypeEmbedding:   SemiClawCloudBaseURL,
			types.ModelTypeRerank:      SemiClawCloudBaseURL,
			types.ModelTypeVLLM:        SemiClawCloudBaseURL,
		},
		ModelTypes: []types.ModelType{
			types.ModelTypeKnowledgeQA,
			types.ModelTypeEmbedding,
			types.ModelTypeRerank,
			types.ModelTypeVLLM,
		},
		RequiresAuth: true,
	}
}

func (p *SemiClawCloudProvider) ValidateConfig(config *Config) error {
	// AppID/AppSecret 通过专用初始化接口写入，此处仅做结构校验。
	// 其中 AppSecret 字段当前实际承载上游 API Key。
	return nil
}
