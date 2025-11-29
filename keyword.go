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

// generateKeywords generates SEO keywords using the modular agent framework
func generateKeywords(ctx context.Context, apiKey, theme string) ([]string, error) {
	// 🎯 Configure the agent externally (maximum flexibility!)
	config := AgentConfig{
		ModelName: KEYWORD_MODEL,
		Providers: KEYWORD_PROVIDERS,
		
		SystemPrompt: "You are a SEO expert specializing in book discovery and audiobook streaming services.",
		
		UserPromptFormat: `Generate 8 SEO keywords for a Nextory landing page about "%s".

Consider various angles based on theme, for example:
- Format variations: audiobooks, ebooks, magazines
- Intent signals: best, top, popular, trending, recommendations
- Value propositions: unlimited, family, streaming, free trial
- Use cases: for commute, for family, for kids

Mix broad discovery terms with long-tail conversion keywords. Use the submit_keywords tool.`,

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
		
		Temperature:   0.7,
		MaxIterations: 3,
		
		APIKey:      apiKey,
		HTTPReferer: "https://github.com/booktok-hype-hub",
		XTitle:      "BookTok Landing Page Agent",
	}
	
	// 🚀 Run the thinking loop!
	result := RunAgent(ctx, config, theme)
	
	if !result.Success {
		return nil, fmt.Errorf("agent failed: %w", result.Error)
	}
	
	// Extract keywords using the helper function
	keywords, err := ExtractKeywords(result)
	if err != nil {
		return nil, err
	}
	
	// Add the theme itself as a keyword (preserving original behavior)
	keywords = append(keywords, strings.ToLower(theme))
	
	log.Printf("✨ Generated %d keywords in %d iterations", len(keywords), result.Iterations)
	return keywords, nil
}

func boolPtr(b bool) *bool {
	return &b
}
