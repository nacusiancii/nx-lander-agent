package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	GLOBAL_AI_MODEL     = KIMI_K2_THINKING.Name()
	GLOBAL_AI_PROVIDERS = []string{KIMI_K2_THINKING["Google"]}
)

func main() {
	fmt.Println("🤖 Landing Page Agent Started")

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️  OPENROUTER_API_KEY not set, using fallback keywords")
	}

	fmt.Print("💡 What landing page idea? (e.g., romance books, thriller audiobooks): ")

	reader := bufio.NewReader(os.Stdin)
	idea, _ := reader.ReadString('\n')
	idea = strings.TrimSpace(idea)

	if idea == "" {
		fmt.Println("❌ No idea provided")
		return
	}

	fmt.Printf("\n💡 Building landing page for: %s\n", idea)
	fmt.Println("🔄 Generating SEO keywords...")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	keywords, err := generateKeywords(ctx, apiKey, idea)
	if err != nil {
		fmt.Printf("❌ Error generating keywords: %v\n", err)
		return
	}

	fmt.Println("\n✨ Generated Keywords:")
	fmt.Println(strings.Repeat("─", 50))
	for i, kw := range keywords {
		fmt.Printf("  %2d. %s\n", i+1, kw)
	}
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("\n📊 Total: %d keywords\n", len(keywords))

	// Generate specific search terms
	fmt.Println("\n🔍 Generating must-target search terms...")
	searchTerms, err := generateSearchTerms(ctx, apiKey, idea, keywords)
	if err != nil {
		fmt.Printf("❌ Error generating search terms: %v\n", err)
		return
	}

	fmt.Println("\n🎯 Must-Target Search Terms:")
	fmt.Println(strings.Repeat("═", 60))
	for i, term := range searchTerms {
		fmt.Printf("  %2d. %s\n", i+1, term)
	}
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("\n🎯 Total: %d search terms\n", len(searchTerms))
}
