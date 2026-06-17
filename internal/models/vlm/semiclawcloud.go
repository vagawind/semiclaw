package vlm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vagawind/semiclaw/internal/logger"
	"github.com/vagawind/semiclaw/internal/models/utils"
	"github.com/google/uuid"
)

const semiClawCloudVLMPath = "/api/v1/chat/completions"

// SemiClawCloudVLM implements VLM via the SemiClawCloud API.
type SemiClawCloudVLM struct {
	modelName       string
	remoteModelName string
	modelID         string
	appID           string
	apiKey          string
	baseURL         string
	client          *http.Client
}

// NewSemiClawCloudVLM creates a SemiClawCloud-backed VLM instance.
func NewSemiClawCloudVLM(config *Config) (*SemiClawCloudVLM, error) {
	if config.AppID == "" {
		return nil, fmt.Errorf("SemiClawCloud VLM: AppID is required")
	}
	if config.AppSecret == "" {
		return nil, fmt.Errorf("SemiClawCloud VLM: AppSecret is required")
	}
	remoteModelName := ""
	if config.Extra != nil {
		if v, ok := config.Extra["remote_model_name"]; ok {
			if vs, ok := v.(string); ok {
				remoteModelName = strings.TrimSpace(vs)
			}
		}
	}
	return &SemiClawCloudVLM{
		modelName:       config.ModelName,
		remoteModelName: remoteModelName,
		modelID:         config.ModelID,
		appID:           config.AppID,
		apiKey:          config.AppSecret,
		baseURL:         strings.TrimRight(config.BaseURL, "/"),
		client:          &http.Client{Timeout: vlmHTTPTimeout()},
	}, nil
}

type semiClawCloudVLMContentPart struct {
	Type     string                      `json:"type"`
	Text     string                      `json:"text,omitempty"`
	ImageURL *semiClawCloudVLMImageURL    `json:"image_url,omitempty"`
}

type semiClawCloudVLMImageURL struct {
	URL string `json:"url"`
}

type semiClawCloudVLMMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type semiClawCloudVLMRequest struct {
	Model       string                   `json:"model"`
	Messages    []semiClawCloudVLMMessage `json:"messages"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Temperature float64                  `json:"temperature,omitempty"`
	Stream      bool                     `json:"stream"`
}

type semiClawCloudVLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Predict sends images with a text prompt to the SemiClawCloud API.
func (v *SemiClawCloudVLM) Predict(ctx context.Context, imgBytesList [][]byte, prompt string) (string, error) {
	var parts []semiClawCloudVLMContentPart

	parts = append(parts, semiClawCloudVLMContentPart{
		Type: "text",
		Text: prompt,
	})

	for _, imgBytes := range imgBytesList {
		if len(imgBytes) > 0 {
			mimeType := detectImageMIME(imgBytes)
			b64 := base64.StdEncoding.EncodeToString(imgBytes)
			dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, b64)
			parts = append(parts, semiClawCloudVLMContentPart{
				Type: "image_url",
				ImageURL: &semiClawCloudVLMImageURL{
					URL: dataURI,
				},
			})
		}
	}

	reqBody := semiClawCloudVLMRequest{
		Model: v.effectiveModelName(),
		Messages: []semiClawCloudVLMMessage{
			{
				Role:    "user",
				Content: parts,
			},
		},
		MaxTokens:   defaultMaxToks,
		Temperature: float64(defaultTemp),
		Stream:      false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("semiclawcloud VLM: marshal: %w", err)
	}

	requestID := uuid.New().String()
	headers := utils.Sign(v.appID, v.apiKey, requestID, string(bodyBytes))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.baseURL+semiClawCloudVLMPath, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("semiclawcloud VLM: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, hv := range headers {
		req.Header.Set(k, hv)
	}

	totalImageSize := 0
	for _, img := range imgBytesList {
		totalImageSize += len(img)
	}
	logger.Infof(ctx, "[VLM] Calling SemiClawCloud API, model=%s, baseURL=%s, numImages=%d, totalImageSize=%d",
		v.effectiveModelName(), v.baseURL, len(imgBytesList), totalImageSize)

	resp, err := v.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("semiclawcloud VLM: do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("semiclawcloud VLM: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("semiclawcloud VLM: status %d: %s", resp.StatusCode, string(respBytes))
	}

	var vlmResp semiClawCloudVLMResponse
	if err := json.Unmarshal(respBytes, &vlmResp); err != nil {
		return "", fmt.Errorf("semiclawcloud VLM: unmarshal: %w", err)
	}
	if len(vlmResp.Choices) == 0 {
		return "", fmt.Errorf("semiclawcloud VLM: no choices in response")
	}

	content := vlmResp.Choices[0].Message.Content
	logger.Infof(ctx, "[VLM] SemiClawCloud response received, len=%d", len(content))
	return content, nil
}

func (v *SemiClawCloudVLM) effectiveModelName() string {
	if v.remoteModelName != "" {
		return v.remoteModelName
	}
	return v.modelName
}

func (v *SemiClawCloudVLM) GetModelName() string { return v.modelName }
func (v *SemiClawCloudVLM) GetModelID() string   { return v.modelID }
