package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/evidra/evidra/questionnaire/domain"
)

type DetectedQuestion struct {
	Text    string
	Type    domain.QuestionType
	Order   int
	Options []string
}

var (
	numberedPattern = regexp.MustCompile(`^(\d+)[.)]\s+(.+)`)
	letterPattern   = regexp.MustCompile(`^([a-zA-Z])[.)]\s+(.+)`)
	optionPattern   = regexp.MustCompile(`^([a-zA-Z])[.)]\s+(.+)`)
	trueFalsePat    = regexp.MustCompile(`(?i)^(true|false)\s*[.)]?\s*$`)
	fillBlankPat    = regexp.MustCompile(`_{3,}|\(\)|\[\]`)
)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) Find(text string) []DetectedQuestion {
	lines := strings.Split(text, "\n")
	var questions []DetectedQuestion
	var currentAnswerLines []string
	order := 0
	inOptions := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if m := numberedPattern.FindStringSubmatch(line); m != nil {
			if len(currentAnswerLines) > 0 && len(questions) > 0 {
				questions[len(questions)-1].Options = append(questions[len(questions)-1].Options, currentAnswerLines...)
				currentAnswerLines = nil
			}
			inOptions = false
			order++
			qtype := detectType(m[2], line)
			questions = append(questions, DetectedQuestion{
				Text:  m[2],
				Type:  qtype,
				Order: order,
			})
			continue
		}

		if len(questions) > 0 {
			if m := optionPattern.FindStringSubmatch(line); m != nil && inOptions {
				currentAnswerLines = append(currentAnswerLines, line)
				continue
			}

			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				currentAnswerLines = append(currentAnswerLines, strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "))
				inOptions = true
				continue
			}

			if strings.Contains(line, "?") && len(line) < 200 {
				if len(currentAnswerLines) > 0 {
					questions[len(questions)-1].Options = append(questions[len(questions)-1].Options, currentAnswerLines...)
					currentAnswerLines = nil
				}
				order++
				qtype := detectType(line, line)
				questions = append(questions, DetectedQuestion{
					Text:  line,
					Type:  qtype,
					Order: order,
				})
				inOptions = false
				continue
			}

			currentAnswerLines = append(currentAnswerLines, line)
			inOptions = true
		}
	}

	if len(currentAnswerLines) > 0 && len(questions) > 0 {
		questions[len(questions)-1].Options = append(questions[len(questions)-1].Options, currentAnswerLines...)
	}

	return questions
}

func detectType(line, fullLine string) domain.QuestionType {
	if fillBlankPat.MatchString(line) {
		return domain.QuestionTypeFillBlank
	}
	if trueFalsePat.MatchString(line) || containsAny(strings.ToLower(line), "true or false", "true/false") {
		return domain.QuestionTypeTrueFalse
	}
	if containsAny(line, "select all", "choose all", "multiple") {
		return domain.QuestionTypeMultipleChoice
	}
	if containsAny(line, "select one", "choose one", "single") || hasOptions(fullLine) {
		return domain.QuestionTypeSingleChoice
	}
	if strings.Contains(line, "?") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "what") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "why") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "how") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "describe") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "explain") || strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "list") {
		return domain.QuestionTypeOpenEnded
	}
	return domain.QuestionTypeOpenEnded
}

func hasOptions(text string) bool {
	lines := strings.Split(text, "\n")
	if len(lines) < 3 {
		return false
	}
	optionCount := 0
	for _, l := range lines[1:] {
		l = strings.TrimSpace(l)
		if letterPattern.MatchString(l) || strings.HasPrefix(l, "- ") || strings.HasPrefix(l, "* ") {
			optionCount++
		}
	}
	return optionCount >= 2
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
