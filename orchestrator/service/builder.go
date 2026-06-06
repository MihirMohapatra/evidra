package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/evidra/evidra/orchestrator/domain"
)

type PromptBuilder struct{}

type StructuredAnswer struct {
	Answer     string  `json:"answer"`
	Confidence float64 `json:"confidence"`
	Reasoning  string  `json:"reasoning"`
}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (b *PromptBuilder) Build(question domain.Question, evidence []domain.Evidence) (systemPrompt, userPrompt string) {
	systemPrompt = `You are a compliance evidence analyst for Evidra. Your role is to analyze compliance questions using provided evidence items and generate structured answers.

For each question:
1. Review all relevant evidence items carefully
2. Determine if the evidence supports, contradicts, or is insufficient for the question
3. Provide a clear, factual answer based on the evidence
4. Rate your confidence in the answer (0.0 to 1.0)
5. Explain your reasoning step by step

Respond in JSON format with the following fields:
- answer: string (your concise answer)
- confidence: number (0.0 to 1.0)
- reasoning: string (detailed step-by-step reasoning)`

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Question\n%s\n\n", question.Text))
	if question.Context != "" {
		sb.WriteString(fmt.Sprintf("## Context\n%s\n\n", question.Context))
	}

	if len(evidence) > 0 {
		sb.WriteString("## Relevant Evidence\n\n")
		for i, ev := range evidence {
			sb.WriteString(fmt.Sprintf("### Evidence %d: %s\n", i+1, ev.Title))
			sb.WriteString(fmt.Sprintf("- Category: %s\n", ev.Category))
			sb.WriteString(fmt.Sprintf("- Relevance Score: %.4f\n", ev.Score))
			sb.WriteString(fmt.Sprintf("- Content: %s\n\n", ev.Content))
		}
	} else {
		sb.WriteString("## Relevant Evidence\nNo relevant evidence found. Answer based on general knowledge and indicate low confidence.\n")
	}

	userPrompt = sb.String()
	return
}

func (b *PromptBuilder) ParseStructuredAnswer(raw string) (StructuredAnswer, error) {
	cleaned := cleanJSON(raw)

	var result StructuredAnswer
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return result, fmt.Errorf("parse structured answer: %w", err)
	}

	if result.Confidence < 0 || result.Confidence > 1 {
		result.Confidence = 0
	}
	if result.Answer == "" {
		return result, fmt.Errorf("answer is empty")
	}

	return result, nil
}

func cleanJSON(raw string) string {
	start := strings.Index(raw, "{")
	if start == -1 {
		return raw
	}
	end := strings.LastIndex(raw, "}")
	if end == -1 || end < start {
		return raw
	}
	return raw[start : end+1]
}
