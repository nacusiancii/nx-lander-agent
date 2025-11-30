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

// ═══════════════════════════════════════════════════════════════════════════
// 🔍 SEARCH TERM PROMPTS - Centralized Prompt Management
// ═══════════════════════════════════════════════════════════════════════════

const (
	// Initial generation system prompt
	INITIAL_GENERATION_SYSTEM_PROMPT = `
You are a SEO search term specialist. 
Your ONLY job is generating highly specific, 
conversion-focused search terms for book discovery and audiobook services.

You are an EXPERT at crafting search queries that real users type when looking for content.`

	// Initial generation user prompt template
	INITIAL_GENERATION_USER_PROMPT_TEMPLATE = `
Generate EXACTLY %d specific, must-target search terms for a Nextory landing page.

Theme: "%s"
Base Keywords: %s

REQUIREMENTS - You MUST include diverse search patterns:
✓ Comparison terms (e.g., "X vs Y", "X alternative")
✓ Question-based (e.g., "where to find X", "how to get X")
✓ Best/Top lists (e.g., "best X for Y", "top X in 2025")
✓ Value-focused (e.g., "unlimited X", "free X trial")
✓ Format combinations (e.g., "X audiobooks", "X ebooks")
✓ User intent (e.g., "X for beginners", "X for commute")
✓ Specific use cases (e.g., "X for family", "X for kids")

Make them SPECIFIC and CONVERSION-FOCUSED!
Use the submit_search_terms tool with EXACTLY %d terms.`

	// Refinement system prompt
	REFINEMENT_SYSTEM_PROMPT = `You are a SEO search term refinement specialist. You improve existing search terms by adding missing patterns and increasing diversity.`

	// Refinement user prompt template
	REFINEMENT_USER_PROMPT_TEMPLATE = `
Refine these %d search terms for theme "%s":

CURRENT TERMS:
%s

MISSING PATTERNS:
%s

Generate EXACTLY %d improved search terms that:
1. Keep the good ones from current terms
2. Add new terms covering missing patterns
3. Ensure high diversity and conversion focus

Use the submit_search_terms tool with EXACTLY %d terms.`
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
