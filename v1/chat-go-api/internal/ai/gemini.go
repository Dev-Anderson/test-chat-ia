package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Intent string

const (
	Greeting  Intent = "greeting"
	About     Intent = "about"
	Schedule  Intent = "schedule"
	Plans     Intent = "plans"
	BookClass Intent = "book_class"
	Unknown   Intent = "unknown"
)

func buildPrompt(userMessage string) string {
	return fmt.Sprintf(`
Você é um classificador de intenção para um chatbot de atendimento de um Box de CrossFit.
Retorne APENAS uma das intenções abaixo (sem texto adicional):

- greeting
- about
- schedule
- plans
- book_class
- unknown

Mensagem do usuário: "%s"
`, userMessage)
}

func ClassifyIntent(ctx context.Context, message string) (Intent, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return Unknown, fmt.Errorf("GEMINI_API_KEY não configurada")
	}

	model := "gemini-2.5-flash"
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model, key,
	)

	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{"text": buildPrompt(message)},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Unknown, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Unknown, fmt.Errorf("gemini http status: %d", resp.StatusCode)
	}

	var r struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Unknown, err
	}

	if len(r.Candidates) == 0 || len(r.Candidates[0].Content.Parts) == 0 {
		return Unknown, nil
	}

	intentStr := strings.TrimSpace(strings.ToLower(r.Candidates[0].Content.Parts[0].Text))
	switch intentStr {
	case "greeting":
		return Greeting, nil
	case "about":
		return About, nil
	case "schedule":
		return Schedule, nil
	case "plans":
		return Plans, nil
	case "book_class":
		return BookClass, nil
	default:
		return Unknown, nil
	}
}
