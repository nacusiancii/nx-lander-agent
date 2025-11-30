package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	openrouter "github.com/revrost/go-openrouter"
)

// ═══════════════════════════════════════════════════════════════════════════
// 🔍 SEARCH TERM SPECIALIST - A Hyper-Focused Search Query Generation Machine
// ═══════════════════════════════════════════════════════════════════════════
//
// This agent has ONE JOB: Generate amazing search terms for SEO optimization.
// It's NOT a generic framework - it's a SPECIALIST with hardcoded domain knowledge!
//
// Constraints:
// - Max n API calls per run (1 initial + up to n-1 refinements)
// - Each call has independent context (no exponential message history)
// - Quality evaluation happens locally (no extra API calls)
//
// ═══════════════════════════════════════════════════════════════════════════

const (
	MAX_REFINEMENT_ITERATIONS = 4  // Total: 1 initial + 4 refinements = 5 calls max
	TARGET_SEARCH_TERM_COUNT  = 15 // We want exactly 15 search terms
)

// SearchTermAgent - The obsessed search term craftsman
type SearchTermAgent struct {
	// Core inputs
	theme        string
	baseKeywords []string

	// API config
	apiKey    string
	modelName string
	providers []string

	// Current state
	currentTerms []string
	iteration    int
}

// SearchTermQuality - HARDCODED quality metrics for search terms
type SearchTermQuality struct {
	// SEO Pattern Coverage (hardcoded search term knowledge!)
	HasComparisons bool // e.g., "X vs Y", "X alternative"
	HasQuestions   bool // e.g., "where to find X", "how to get X"
	HasBestLists   bool // e.g., "best X for Y", "top X in 2025"
	HasValueTerms  bool // e.g., "unlimited X", "free X trial"
	HasFormatMix   bool // e.g., "X audiobooks", "X ebooks"
	HasUserIntent  bool // e.g., "X for beginners", "X for commute"

	// Diversity
	DiversityScore float64 // How unique are the terms?

	// Coverage
	TermCount int
}

// ═══════════════════════════════════════════════════════════════════════════
// 🎯 MAIN ENTRY POINT
// ═══════════════════════════════════════════════════════════════════════════

// NewSearchTermAgent creates a new specialized search term generator
func NewSearchTermAgent(theme string, baseKeywords []string, apiKey string, modelName string, providers []string) *SearchTermAgent {
	return &SearchTermAgent{
		theme:        theme,
		baseKeywords: baseKeywords,
		apiKey:       apiKey,
		modelName:    modelName,
		providers:    providers,
		iteration:    0,
	}
}

// Generate - The main loop: 1 initial call + up to 9 refinement calls
func (a *SearchTermAgent) Generate(ctx context.Context) ([]string, error) {
	log.Printf("🔍 Search Term Specialist started for theme: %s", a.theme)

	// CALL 1: Generate initial search terms
	terms, err := a.generateInitialTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("initial generation failed: %w", err)
	}
	a.currentTerms = terms
	log.Printf("✨ Generated %d initial terms", len(terms))

	// CALLS 2-10: Refinement iterations (up to 9 more calls)
	for a.iteration < MAX_REFINEMENT_ITERATIONS {
		// Evaluate quality locally (NO API call here!)
		quality := a.evaluateSearchTermQuality()

		// Check if we're good enough
		if a.isGoodEnough(quality) {
			log.Printf("✅ Quality target reached after %d total calls", a.iteration+1)
			break
		}

		// Refine the terms (1 API call per iteration)
		log.Printf("🔄 Refinement iteration %d: improving coverage...", a.iteration+1)
		refined, err := a.refineTermsIteration(ctx, quality)
		if err != nil {
			log.Printf("⚠️  Refinement %d failed, keeping current terms: %v", a.iteration+1, err)
			break // Don't fail completely, just stop refining
		}

		a.currentTerms = refined
		a.iteration++
	}

	log.Printf("🎉 Final: %d terms after %d total API calls", len(a.currentTerms), a.iteration+1)
	return a.currentTerms, nil
}

// ═══════════════════════════════════════════════════════════════════════════
// 🎬 GENERATION PHASE - The Initial Creative Burst
// ═══════════════════════════════════════════════════════════════════════════

func (a *SearchTermAgent) generateInitialTerms(ctx context.Context) ([]string, error) {
	client := openrouter.NewClient(
		a.apiKey,
		openrouter.WithHTTPReferer("https://github.com/booktok-hype-hub"),
		openrouter.WithXTitle("Search Term Specialist"),
	)

	keywordList := strings.Join(a.baseKeywords, ", ")

	// HARDCODED search term generation prompt - SPECIALIZED!
	systemPrompt := `You are a SEO search term specialist. Your ONLY job is generating highly specific, conversion-focused search terms for book discovery and audiobook services.

You are an EXPERT at crafting search queries that real users type when looking for content.`

	userPrompt := fmt.Sprintf(`Generate EXACTLY %d specific, must-target search terms for a Nextory landing page.

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
Use the submit_search_terms tool with EXACTLY %d terms.`,
		TARGET_SEARCH_TERM_COUNT, a.theme, keywordList, TARGET_SEARCH_TERM_COUNT)

	resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: a.modelName,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleSystem,
				Content: openrouter.Content{Text: systemPrompt},
			},
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: userPrompt},
			},
		},
		Tools:       a.getSearchTermTool(),
		Temperature: 0.8, // Creative but focused
		Provider: &openrouter.ChatProvider{
			Order:          a.providers,
			AllowFallbacks: boolPtr(false),
		},
	})

	if err != nil {
		return nil, err
	}

	return a.extractSearchTerms(resp)
}

// ═══════════════════════════════════════════════════════════════════════════
// ✨ REFINEMENT PHASE - Stateless Iteration (Fresh Context Each Time!)
// ═══════════════════════════════════════════════════════════════════════════

func (a *SearchTermAgent) refineTermsIteration(ctx context.Context, quality SearchTermQuality) ([]string, error) {
	client := openrouter.NewClient(
		a.apiKey,
		openrouter.WithHTTPReferer("https://github.com/booktok-hype-hub"),
		openrouter.WithXTitle("Search Term Specialist"),
	)

	// Build FOCUSED refinement prompt - NO MESSAGE HISTORY!
	// Just current terms + what's missing = STATELESS!
	missingPatterns := a.identifyMissingPatterns(quality)

	systemPrompt := `You are a SEO search term refinement specialist. You improve existing search terms by adding missing patterns and increasing diversity.`

	userPrompt := fmt.Sprintf(`Refine these %d search terms for theme "%s":

CURRENT TERMS:
%s

MISSING PATTERNS:
%s

Generate EXACTLY %d improved search terms that:
1. Keep the good ones from current terms
2. Add new terms covering missing patterns
3. Ensure high diversity and conversion focus

Use the submit_search_terms tool with EXACTLY %d terms.`,
		len(a.currentTerms),
		a.theme,
		a.formatTermsForPrompt(a.currentTerms),
		missingPatterns,
		TARGET_SEARCH_TERM_COUNT,
		TARGET_SEARCH_TERM_COUNT)

	resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: a.modelName,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleSystem,
				Content: openrouter.Content{Text: systemPrompt},
			},
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: userPrompt},
			},
		},
		Tools:       a.getSearchTermTool(),
		Temperature: 0.7, // Slightly more deterministic for refinement
		Provider: &openrouter.ChatProvider{
			Order:          a.providers,
			AllowFallbacks: boolPtr(false),
		},
	})

	if err != nil {
		return nil, err
	}

	return a.extractSearchTerms(resp)
}

// ═══════════════════════════════════════════════════════════════════════════
// 📊 QUALITY EVALUATION - Local Logic (NO API CALLS!)
// ═══════════════════════════════════════════════════════════════════════════

// evaluateSearchTermQuality - HARDCODED search term pattern detection
func (a *SearchTermAgent) evaluateSearchTermQuality() SearchTermQuality {
	quality := SearchTermQuality{
		TermCount: len(a.currentTerms),
	}

	termsLower := make([]string, len(a.currentTerms))
	for i, term := range a.currentTerms {
		termsLower[i] = strings.ToLower(term)
	}

	// HARDCODED pattern detection - SEARCH TERM SPECIFIC!
	for _, term := range termsLower {
		// Comparison patterns
		if strings.Contains(term, " vs ") || strings.Contains(term, " versus ") ||
			strings.Contains(term, "alternative") || strings.Contains(term, "comparison") {
			quality.HasComparisons = true
		}

		// Question patterns
		if strings.HasPrefix(term, "where ") || strings.HasPrefix(term, "how ") ||
			strings.HasPrefix(term, "what ") || strings.HasPrefix(term, "which ") {
			quality.HasQuestions = true
		}

		// Best/Top lists
		if strings.Contains(term, "best ") || strings.Contains(term, "top ") ||
			strings.Contains(term, "most popular") {
			quality.HasBestLists = true
		}

		// Value terms
		if strings.Contains(term, "unlimited") || strings.Contains(term, "free") ||
			strings.Contains(term, "trial") || strings.Contains(term, "affordable") {
			quality.HasValueTerms = true
		}

		// Format mix
		if strings.Contains(term, "audiobook") || strings.Contains(term, "ebook") ||
			strings.Contains(term, "book") || strings.Contains(term, "magazine") {
			quality.HasFormatMix = true
		}

		// User intent
		if strings.Contains(term, " for ") {
			quality.HasUserIntent = true
		}
	}

	// Calculate diversity (simple: unique word count ratio)
	quality.DiversityScore = a.calculateDiversity(termsLower)

	return quality
}

func (a *SearchTermAgent) calculateDiversity(terms []string) float64 {
	wordSet := make(map[string]bool)
	totalWords := 0

	for _, term := range terms {
		words := strings.Fields(term)
		totalWords += len(words)
		for _, word := range words {
			wordSet[word] = true
		}
	}

	if totalWords == 0 {
		return 0
	}

	return float64(len(wordSet)) / float64(totalWords)
}

// isGoodEnough - HARDCODED quality thresholds for search terms
func (a *SearchTermAgent) isGoodEnough(quality SearchTermQuality) bool {
	// Must have correct count
	if quality.TermCount != TARGET_SEARCH_TERM_COUNT {
		return false
	}

	// Must cover at least 4 out of 6 patterns
	patternCount := 0
	if quality.HasComparisons {
		patternCount++
	}
	if quality.HasQuestions {
		patternCount++
	}
	if quality.HasBestLists {
		patternCount++
	}
	if quality.HasValueTerms {
		patternCount++
	}
	if quality.HasFormatMix {
		patternCount++
	}
	if quality.HasUserIntent {
		patternCount++
	}

	// Must have good diversity
	return patternCount >= 4 && quality.DiversityScore >= 0.6
}

// identifyMissingPatterns - HARDCODED search term pattern knowledge
func (a *SearchTermAgent) identifyMissingPatterns(quality SearchTermQuality) string {
	var missing []string

	if !quality.HasComparisons {
		missing = append(missing, "- Comparison terms (e.g., 'X vs Y', 'X alternative')")
	}
	if !quality.HasQuestions {
		missing = append(missing, "- Question-based (e.g., 'where to find X', 'how to get X')")
	}
	if !quality.HasBestLists {
		missing = append(missing, "- Best/Top lists (e.g., 'best X for Y', 'top X in 2025')")
	}
	if !quality.HasValueTerms {
		missing = append(missing, "- Value-focused (e.g., 'unlimited X', 'free X trial')")
	}
	if !quality.HasFormatMix {
		missing = append(missing, "- Format combinations (e.g., 'X audiobooks', 'X ebooks')")
	}
	if !quality.HasUserIntent {
		missing = append(missing, "- User intent (e.g., 'X for beginners', 'X for commute')")
	}

	if len(missing) == 0 {
		return "None - improve diversity and specificity!"
	}

	return strings.Join(missing, "\n")
}

// ═══════════════════════════════════════════════════════════════════════════
// 🛠️ HELPERS - Search Term Specific Utilities
// ═══════════════════════════════════════════════════════════════════════════

func (a *SearchTermAgent) getSearchTermTool() []openrouter.Tool {
	return []openrouter.Tool{
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
	}
}

func (a *SearchTermAgent) extractSearchTerms(resp openrouter.ChatCompletionResponse) ([]string, error) {
	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		return nil, fmt.Errorf("no tool call in response")
	}

	var result struct {
		SearchTerms []string `json:"search_terms"`
	}

	args := resp.Choices[0].Message.ToolCalls[0].Function.Arguments
	if err := json.Unmarshal([]byte(args), &result); err != nil {
		return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
	}

	if len(result.SearchTerms) != TARGET_SEARCH_TERM_COUNT {
		return nil, fmt.Errorf("expected %d terms, got %d", TARGET_SEARCH_TERM_COUNT, len(result.SearchTerms))
	}

	return result.SearchTerms, nil
}

func (a *SearchTermAgent) formatTermsForPrompt(terms []string) string {
	var formatted []string
	for i, term := range terms {
		formatted = append(formatted, fmt.Sprintf("%d. %s", i+1, term))
	}
	return strings.Join(formatted, "\n")
}
