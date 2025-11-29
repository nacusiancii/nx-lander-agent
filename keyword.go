package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	openrouter "github.com/revrost/go-openrouter"
)

var (
	KEYWORD_MODEL     = GLOBAL_AI_MODEL
	KEYWORD_PROVIDERS = GLOBAL_AI_PROVIDERS
)

func generateKeywords(ctx context.Context, apiKey, theme string) ([]string, error) {
	client := openrouter.NewClient(
		apiKey,
		openrouter.WithXTitle("BookTok Landing Page Agent"),
		openrouter.WithHTTPReferer("https://github.com/booktok-hype-hub"),
	)

	resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: KEYWORD_MODEL,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role: openrouter.ChatMessageRoleSystem,
				Content: openrouter.Content{
					Text: "You are a SEO expert specializing in book discovery and audiobook streaming services.",
				},
			},
			{
				Role: openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{
					Text: fmt.Sprintf(`Generate 8 SEO keywords for a Nextory landing page about "%s".

Consider various angles based on theme, for example:
- Format variations: audiobooks, ebooks, magazines
- Intent signals: best, top, popular, trending, recommendations
- Value propositions: unlimited, family, streaming, free trial
- Use cases: for commute, for family, for kids

Mix broad discovery terms with long-tail conversion keywords. Use the submit_keywords tool.`, theme),
				},
			},
		},
		Tools: []openrouter.Tool{
			{
				Type: openrouter.ToolTypeFunction,
				Function: &openrouter.FunctionDefinition{
					Name:        "submit_keywords",
					Description: "Submit the generated SEO keywords",
					Parameters: json.RawMessage(`{
						"type": "object",
						"properties": {
							"keywords": {
								"type": "array",
								"items": {"type": "string"},
								"description": "Array of 8 SEO keywords"
							}
						},
						"required": ["keywords"]
					}`),
				},
			},
		},
		Temperature: 0.7,
		Provider: &openrouter.ChatProvider{
			Order:          KEYWORD_PROVIDERS,
			AllowFallbacks: boolPtr(false),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("AI failed: %w", err)
	}

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		var keywordResult struct {
			Keywords []string `json:"keywords"`
		}
		args := resp.Choices[0].Message.ToolCalls[0].Function.Arguments
		if json.Unmarshal([]byte(args), &keywordResult) == nil {
			keywordResult.Keywords = append(keywordResult.Keywords, strings.ToLower(theme))
			log.Printf("✨ Generated %d keywords", len(keywordResult.Keywords))
			return keywordResult.Keywords, nil
		}
	}

	log.Println("⚠️  No valid tool call, returning error")
	return nil, fmt.Errorf("no valid tool call")
}

func boolPtr(b bool) *bool {
	return &b
}
