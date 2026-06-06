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

const claudeEndpoint = "https://api.anthropic.com/v1/messages"

type ClaudeEngine struct {
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
	client      *http.Client
	builder     *PromptBuilder
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model       string          `json:"model"`
	Messages    []claudeMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
}

type claudeContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type claudeResponse struct {
	Content []claudeContent `json:"content"`
}

func NewClaudeEngine(apiKey, model string, temperature float64, maxTokens int) *ClaudeEngine {
	return &ClaudeEngine{
		apiKey:      apiKey,
		model:       model,
		temperature: temperature,
		maxTokens:   maxTokens,
		client:      &http.Client{Timeout: 60 * time.Second},
		builder:     NewPromptBuilder(),
	}
}

func (e *ClaudeEngine) Name() string {
	return fmt.Sprintf("claude-%s", e.model)
}

func (e *ClaudeEngine) GenerateAnswer(ctx context.Context, question domain.Question, evidence []domain.Evidence) (domain.Draft, error) {
	systemPrompt, userPrompt := e.builder.Build(question, evidence)

	req := claudeRequest{
		Model: e.model,
		Messages: []claudeMessage{
			{Role: "user", Content: systemPrompt + "\n\n" + userPrompt},
		},
		MaxTokens:   e.maxTokens,
		Temperature: e.temperature,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: marshal request: %w", domain.ErrLLMError, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, claudeEndpoint, bytes.NewReader(body))
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: create request: %w", domain.ErrLLMError, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", e.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: api call: %w", domain.ErrLLMError, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: read response: %w", domain.ErrLLMError, err)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Draft{}, fmt.Errorf("%w: claude returned %d: %s", domain.ErrLLMError, resp.StatusCode, string(respBody))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return domain.Draft{}, fmt.Errorf("%w: parse response: %w", domain.ErrLLMError, err)
	}

	if len(claudeResp.Content) == 0 {
		return domain.Draft{}, fmt.Errorf("%w: no content returned", domain.ErrLLMError)
	}

	parsed, err := e.builder.ParseStructuredAnswer(claudeResp.Content[0].Text)
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
