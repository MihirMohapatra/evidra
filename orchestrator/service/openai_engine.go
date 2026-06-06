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

const openaiEndpoint = "https://api.openai.com/v1/chat/completions"

type OpenAIEngine struct {
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
	client      *http.Client
	builder     *PromptBuilder
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type openAIChoice struct {
	Message openAIMessage `json:"message"`
}

type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
}

func NewOpenAIEngine(apiKey, model string, temperature float64, maxTokens int) *OpenAIEngine {
	return &OpenAIEngine{
		apiKey:      apiKey,
		model:       model,
		temperature: temperature,
		maxTokens:   maxTokens,
		client:      &http.Client{Timeout: 60 * time.Second},
		builder:     NewPromptBuilder(),
	}
}

func (e *OpenAIEngine) Name() string {
	return fmt.Sprintf("openai-%s", e.model)
}

func (e *OpenAIEngine) GenerateAnswer(ctx context.Context, question domain.Question, evidence []domain.Evidence) (domain.Draft, error) {
	systemPrompt, userPrompt := e.builder.Build(question, evidence)

	req := openAIRequest{
		Model: e.model,
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: e.temperature,
		MaxTokens:   e.maxTokens,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: marshal request: %w", domain.ErrLLMError, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openaiEndpoint, bytes.NewReader(body))
	if err != nil {
		return domain.Draft{}, fmt.Errorf("%w: create request: %w", domain.ErrLLMError, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+e.apiKey)

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
		return domain.Draft{}, fmt.Errorf("%w: openai returned %d: %s", domain.ErrLLMError, resp.StatusCode, string(respBody))
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return domain.Draft{}, fmt.Errorf("%w: parse response: %w", domain.ErrLLMError, err)
	}

	if len(openAIResp.Choices) == 0 {
		return domain.Draft{}, fmt.Errorf("%w: no choices returned", domain.ErrLLMError)
	}

	parsed, err := e.builder.ParseStructuredAnswer(openAIResp.Choices[0].Message.Content)
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
