package main

import (
	"context"
	"fmt"
	"log"
)

var (
	SEARCH_TERMS_MODEL     = MINIMAX_M2.Name()
	SEARCH_TERMS_PROVIDERS = []string{MINIMAX_M2["Google"]}
)

// generateSearchTerms - Simple wrapper around the specialized SearchTermAgent
func generateSearchTerms(ctx context.Context, apiKey, theme string, keywords []string) ([]string, error) {
	// Create the specialist agent
	agent := NewSearchTermAgent(theme, keywords, apiKey, SEARCH_TERMS_MODEL, SEARCH_TERMS_PROVIDERS)

	// Let it do its magic!
	terms, err := agent.Generate(ctx)
	if err != nil {
		return nil, fmt.Errorf("search term agent failed: %w", err)
	}

	log.Printf("✨ Search Term Specialist completed with %d terms", len(terms))
	return terms, nil
}
