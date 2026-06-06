package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/evidra/evidra/orchestrator/domain"
)

type LocalLLMEngine struct {
	baseURL     string
	model       string
	temperature float64
	maxTokens   int
	client      *http.Client
	builder     *PromptBuilder
}

type ollamaRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func NewLocalLLMEngine(baseURL, model string, temperature float64, maxTokens int) *LocalLLMEngine {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &LocalLLMEngine{
		baseURL:     baseURL,
		model:       model,
		temperature: temperature,
		maxTokens:   maxTokens,
		client:      &http.Client{Timeout: 120 * time.Second},
		builder:     NewPromptBuilder(),
	}
}

func (e *LocalLLMEngine) Name() string {
	return fmt.Sprintf("local-%s", e.model)
}

func (e *LocalLLMEngine) GenerateAnswer(ctx context.Context, question domain.Question, evidence []domain.Evidence) (domain.Draft, error) {
	systemPrompt, userPrompt := e.builder.Build(question, evidence)
	fullPrompt := systemPrompt + "\n\n" + userPrompt

	req := ollamaRequest{
		Model:       e.model,
		Prompt:      fullPrompt,
		Temperature: e.temperature,
		MaxTokens:   e.maxTokens,
		Stream:      false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: marshal request: %w", domain.ErrLLMError, err)
	}

	endpoint := fmt.Sprintf("%s/api/generate", e.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: create request: %w", domain.ErrLLMError, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: api call: %w", domain.ErrLLMError, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: read response: %w", domain.ErrLLMError, err)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Draft{}, fmt.Errorf("%w: local llm returned %d: %s", domain.ErrLLMError, resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return domain.Draft{}, fmt.Errorf("%w: parse response: %w", domain.ErrLLMError, err)
	}

	parsed, err := e.builder.ParseStructuredAnswer(ollamaResp.Response)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: %w", domain.ErrValidationFailed, err)
	}

	evidenceIDs := make([]uuid.UUID, len(evidence))
	for i, ev := range evidence {
		evidenceIDs[i] = ev.ID
	}

	draft := domain.NewDraft(
		question.ID, question.Text,
		parsed.Answer, parsed.Confidence,
		e.Name(), evidenceIDs, parsed.Reasoning,
	)
	return *draft, nil
}
