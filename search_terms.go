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
	SEARCH_TERMS_MODEL     = MINIMAX_M2.Name()
	SEARCH_TERMS_PROVIDERS = []string{MINIMAX_M2["Minimax"]}
)

// generateSearchTerms generates 15 specific must-target search terms
// using the modular agent framework
func generateSearchTerms(ctx context.Context, apiKey, theme string, keywords []string) ([]string, error) {
	keywordList := strings.Join(keywords, ", ")
	
	// 🎯 Configure the agent externally (no hardcoded prompts in agent!)
	config := AgentConfig{
		ModelName: SEARCH_TERMS_MODEL,
		Providers: SEARCH_TERMS_PROVIDERS,
		
		SystemPrompt: "You are a SEO expert specializing in search term optimization for book discovery and audiobook streaming services. Focus on high-intent, conversion-oriented search terms.",
		
		UserPromptFormat: `Generate EXACTLY 15 specific, must-target search terms for a Nextory landing page.

Theme: "%s"
Base Keywords: ` + keywordList + `

Requirements:
- Generate EXACTLY 15 search terms (no more, no less)
- Each term should be highly specific and conversion-focused
- Combine the theme with various modifiers and intents
- Include different search patterns based on theme:
  * Comparison terms (e.g., "X vs Y", "X alternative")
  * Question-based (e.g., "where to find X", "how to get X")
  * Best/Top lists (e.g., "best X for Y", "top X in 2025")
  * Value-focused (e.g., "unlimited X", "free X trial")
  * Format combinations (e.g., "X audiobooks", "X ebooks")
  * User intent (e.g., "X for beginners", "X for commute")
  * Specific use cases (e.g., "X for family", "X for kids")
  
Use the submit_search_terms tool with EXACTLY 15 terms.`,

		Tools: []openrouter.Tool{
			{
				Type: openrouter.ToolTypeFunction,
				Function: &openrouter.FunctionDefinition{
					Name:        "submit_search_terms",
					Description: "Submit exactly 15 specific must-target search terms",
					Parameters: json.RawMessage(`{
						"type": "object",
						"properties": {
							"search_terms": {
								"type": "array",
								"items": {"type": "string"},
								"description": "Array of exactly 15 specific search terms",
								"minItems": 15,
								"maxItems": 15
							}
						},
						"required": ["search_terms"]
					}`),
				},
			},
		},
		
		Temperature:   0.8,
		MaxIterations: 3,
		
		APIKey:      apiKey,
		HTTPReferer: "https://github.com/booktok-hype-hub",
		XTitle:      "BookTok Landing Page Agent",
	}
	
	// 🚀 Run the beautiful thinking loop!
	result := RunAgent(ctx, config, theme)
	
	if !result.Success {
		return nil, fmt.Errorf("agent failed: %w", result.Error)
	}
	
	// Extract search terms using the helper function
	searchTerms, err := ExtractSearchTerms(result)
	if err != nil {
		return nil, err
	}
	
	log.Printf("✨ Generated %d search terms in %d iterations", len(searchTerms), result.Iterations)
	return searchTerms, nil
}
