package main

import (
	"context"
	"fmt"
	"log"
)

// generateSearchTerms - Simple wrapper around the specialized SearchTermAgent
func generateSearchTerms(ctx context.Context, apiKey, theme string, keywords []string) ([]string, error) {
	// Create the specialist agent
	agent := NewSearchTermAgent(theme, keywords, apiKey)
	
	// Let it do its magic!
	terms, err := agent.Generate(ctx)
	if err != nil {
		return nil, fmt.Errorf("search term agent failed: %w", err)
	}
	
	log.Printf("✨ Search Term Specialist completed with %d terms", len(terms))
	return terms, nil
}
