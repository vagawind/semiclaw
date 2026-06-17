package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vagawind/semiclaw/internal/models/utils"
	"github.com/google/uuid"
)

const semiClawCloudEmbedPath = "/api/v1/embeddings"

// SemiClawCloudEmbedder 实现 embedding.Embedder 接口，对接 SemiClawCloud /api/v1/embeddings
type SemiClawCloudEmbedder struct {
	modelName                 string
	remoteModelName           string
	modelID                   string
	appID                     string
	apiKey                    string
	baseURL                   string
	dimensions                int
	supportsDimensionOverride bool
	client                    *http.Client
	EmbedderPooler
}

// NewSemiClawCloudEmbedder 构造 SemiClawCloudEmbedder
func NewSemiClawCloudEmbedder(config Config) (*SemiClawCloudEmbedder, error) {
	if config.AppID == "" {
		return nil, fmt.Errorf("SemiClawCloud embedder: AppID is required")
	}
	if config.AppSecret == "" {
		return nil, fmt.Errorf("SemiClawCloud embedder: AppSecret is required")
	}
	remoteModelName := ""
	if config.ExtraConfig != nil {
		remoteModelName = strings.TrimSpace(config.ExtraConfig["remote_model_name"])
	}
	return &SemiClawCloudEmbedder{
		modelName:                 config.ModelName,
		remoteModelName:           remoteModelName,
		modelID:                   config.ModelID,
		appID:                     config.AppID,
		apiKey:                    config.AppSecret,
		baseURL:                   strings.TrimRight(config.BaseURL, "/"),
		dimensions:                config.Dimensions,
		supportsDimensionOverride: config.SupportsDimensionOverride,
		client:                    &http.Client{Timeout: 60 * time.Second},
	}, nil
}

type semiClawCloudEmbedRequest struct {
	Model                string   `json:"model"`
	Input                []string `json:"input"`
	Dimensions           int      `json:"dimensions,omitempty"`
	TruncatePromptTokens int      `json:"truncate_prompt_tokens,omitempty"`
}

type semiClawCloudEmbedResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (e *SemiClawCloudEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	results, err := e.BatchEmbed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("semiclawcloud embedder: empty response")
	}
	return results[0], nil
}

func (e *SemiClawCloudEmbedder) BatchEmbed(ctx context.Context, texts []string) ([][]float32, error) {
	reqBody := semiClawCloudEmbedRequest{Model: e.effectiveModelName(), Input: texts}
	if e.supportsDimensionOverride && e.dimensions > 0 {
		reqBody.Dimensions = e.dimensions
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud embedder: marshal: %w", err)
	}

	requestID := uuid.New().String()
	headers := utils.Sign(e.appID, e.apiKey, requestID, string(bodyBytes))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+semiClawCloudEmbedPath, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud embedder: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud embedder: do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud embedder: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("semiclawcloud embedder: status %d: %s", resp.StatusCode, string(respBytes))
	}

	var embedResp semiClawCloudEmbedResponse
	if err := json.Unmarshal(respBytes, &embedResp); err != nil {
		return nil, fmt.Errorf("semiclawcloud embedder: unmarshal: %w", err)
	}

	result := make([][]float32, len(texts))
	for _, item := range embedResp.Data {
		if item.Index < len(result) {
			result[item.Index] = item.Embedding
		}
	}
	return result, nil
}

func (e *SemiClawCloudEmbedder) BatchEmbedWithPool(ctx context.Context, model Embedder, texts []string) ([][]float32, error) {
	return e.BatchEmbed(ctx, texts)
}

func (e *SemiClawCloudEmbedder) SetSupportsDimensionOverride(supported bool) {
	e.supportsDimensionOverride = supported
}

func (e *SemiClawCloudEmbedder) effectiveModelName() string {
	if e.remoteModelName != "" {
		return e.remoteModelName
	}
	return e.modelName
}

func (e *SemiClawCloudEmbedder) GetModelName() string { return e.modelName }
func (e *SemiClawCloudEmbedder) GetModelID() string   { return e.modelID }
func (e *SemiClawCloudEmbedder) GetDimensions() int   { return e.dimensions }
