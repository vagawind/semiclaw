package rerank

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

const semiClawCloudRerankPath = "/api/v1/rerank"

// SemiClawCloudReranker 实现 rerank.Reranker 接口，对接 SemiClawCloud /api/v1/rerank
type SemiClawCloudReranker struct {
	modelName       string
	remoteModelName string
	modelID         string
	appID           string
	apiKey          string
	baseURL         string
	client          *http.Client
}

// NewSemiClawCloudReranker 构造 SemiClawCloudReranker
func NewSemiClawCloudReranker(config *RerankerConfig) (*SemiClawCloudReranker, error) {
	if config.AppID == "" {
		return nil, fmt.Errorf("SemiClawCloud reranker: AppID is required")
	}
	if config.AppSecret == "" {
		return nil, fmt.Errorf("SemiClawCloud reranker: AppSecret is required")
	}
	remoteModelName := ""
	if config.ExtraConfig != nil {
		remoteModelName = strings.TrimSpace(config.ExtraConfig["remote_model_name"])
	}
	return &SemiClawCloudReranker{
		modelName:       config.ModelName,
		remoteModelName: remoteModelName,
		modelID:         config.ModelID,
		appID:           config.AppID,
		apiKey:          config.AppSecret,
		baseURL:         strings.TrimRight(config.BaseURL, "/"),
		client:          &http.Client{Timeout: 60 * time.Second},
	}, nil
}

type semiClawCloudRerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

type semiClawCloudRerankResponse struct {
	Results []struct {
		Index          int     `json:"index"`
		RelevanceScore float64 `json:"relevance_score"`
		Document       struct {
			Text string `json:"text"`
		} `json:"document"`
	} `json:"results"`
}

func (r *SemiClawCloudReranker) Rerank(ctx context.Context, query string, documents []string) ([]RankResult, error) {
	reqBody := semiClawCloudRerankRequest{
		Model:     r.effectiveModelName(),
		Query:     query,
		Documents: documents,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud reranker: marshal: %w", err)
	}

	requestID := uuid.New().String()
	headers := utils.Sign(r.appID, r.apiKey, requestID, string(bodyBytes))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+semiClawCloudRerankPath, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud reranker: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud reranker: do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("semiclawcloud reranker: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("semiclawcloud reranker: status %d: %s", resp.StatusCode, string(respBytes))
	}

	var rerankResp semiClawCloudRerankResponse
	if err := json.Unmarshal(respBytes, &rerankResp); err != nil {
		return nil, fmt.Errorf("semiclawcloud reranker: unmarshal: %w", err)
	}

	results := make([]RankResult, 0, len(rerankResp.Results))
	for _, item := range rerankResp.Results {
		results = append(results, RankResult{
			Index:          item.Index,
			RelevanceScore: item.RelevanceScore,
			Document:       DocumentInfo{Text: item.Document.Text},
		})
	}
	return results, nil
}

func (r *SemiClawCloudReranker) effectiveModelName() string {
	if r.remoteModelName != "" {
		return r.remoteModelName
	}
	return r.modelName
}

func (r *SemiClawCloudReranker) GetModelName() string { return r.modelName }
func (r *SemiClawCloudReranker) GetModelID() string   { return r.modelID }
