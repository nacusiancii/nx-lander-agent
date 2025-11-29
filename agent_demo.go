package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	openrouter "github.com/revrost/go-openrouter"
)

// ═══════════════════════════════════════════════════════════════════════════
// 🎪 AGENT DEMO - Showcasing the modular agent with different configurations
// ═══════════════════════════════════════════════════════════════════════════

// This demo shows how the SAME agent framework can handle completely different
// tasks just by changing the external configuration!

func runAgentDemo() {
	apiKey := "your-api-key-here" // Replace with actual API key
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("\n🎪 ═══════════════════════════════════════════════════════")
	fmt.Println("🎪 AI AGENT DEMO - Pure Thinking Engine")
	fmt.Println("🎪 ═══════════════════════════════════════════════════════\n")

	// ═══════════════════════════════════════════════════════════════════════
	// Demo 1: Math Problem Solver
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n📊 DEMO 1: Math Problem Solver")
	fmt.Println("─────────────────────────────────────────────────────────")

	mathConfig := AgentConfig{
		ModelName: KIMI_K2_THINKING.Name(),
		Providers: []string{KIMI_K2_THINKING["Google"]},

		SystemPrompt: "You are a brilliant mathematician. Solve problems step-by-step with clear reasoning.",

		UserPromptFormat: "Solve this math problem: %s\n\nUse the submit_answer tool to provide your final answer.",

		Tools: []openrouter.Tool{
			{
				Type: openrouter.ToolTypeFunction,
				Function: &openrouter.FunctionDefinition{
					Name:        "submit_answer",
					Description: "Submit the final answer to the math problem",
					Parameters: json.RawMessage(`{
						"type": "object",
						"properties": {
							"answer": {
								"type": "string",
								"description": "The final answer with explanation"
							}
						},
						"required": ["answer"]
					}`),
				},
			},
		},

		Temperature:   0.3, // Low for deterministic reasoning
		MaxIterations: 5,

		APIKey:      apiKey,
		HTTPReferer: "https://github.com/demo-agent",
		XTitle:      "Math Problem Solver Agent",
	}

	result := RunAgent(ctx, mathConfig, "What is the sum of all prime numbers less than 20?")
	fmt.Printf("✅ Result: %+v\n", result)

	// ═══════════════════════════════════════════════════════════════════════
	// Demo 2: Creative Writing Assistant
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n✍️  DEMO 2: Creative Writing Assistant")
	fmt.Println("─────────────────────────────────────────────────────────")

	writingConfig := AgentConfig{
		ModelName: MINIMAX_M2.Name(),
		Providers: []string{MINIMAX_M2["Minimax"]},

		SystemPrompt: "You are a creative writing assistant who helps brainstorm story ideas.",

		UserPromptFormat: "Generate 5 unique story ideas based on this theme: %s\n\nUse the submit_stories tool.",

		Tools: []openrouter.Tool{
			{
				Type: openrouter.ToolTypeFunction,
				Function: &openrouter.FunctionDefinition{
					Name:        "submit_stories",
					Description: "Submit 5 creative story ideas",
					Parameters: json.RawMessage(`{
						"type": "object",
						"properties": {
							"stories": {
								"type": "array",
								"items": {"type": "string"},
								"description": "Array of 5 story ideas",
								"minItems": 5,
								"maxItems": 5
							}
						},
						"required": ["stories"]
					}`),
				},
			},
		},

		Temperature:   0.9, // High for creativity!
		MaxIterations: 5,

		APIKey:      apiKey,
		HTTPReferer: "https://github.com/demo-agent",
		XTitle:      "Creative Writing Agent",
	}

	result = RunAgent(ctx, writingConfig, "time travel paradoxes")
	fmt.Printf("✅ Result: %+v\n", result)

	// ═══════════════════════════════════════════════════════════════════════
	// Demo 3: Code Reviewer
	// ═══════════════════════════════════════════════════════════════════════
	fmt.Println("\n🔍 DEMO 3: Code Reviewer")
	fmt.Println("─────────────────────────────────────────────────────────")

	codeReviewConfig := AgentConfig{
		ModelName: KIMI_K2_THINKING.Name(),
		Providers: []string{KIMI_K2_THINKING["Google"]},

		SystemPrompt: "You are an expert code reviewer. Analyze code for bugs, best practices, and improvements.",

		UserPromptFormat: `Review this code and provide feedback:\n\n%s\n\nUse the submit_review tool.`,

		Tools: []openrouter.Tool{
			{
				Type: openrouter.ToolTypeFunction,
				Function: &openrouter.FunctionDefinition{
					Name:        "submit_review",
					Description: "Submit code review feedback",
					Parameters: json.RawMessage(`{
						"type": "object",
						"properties": {
							"issues": {
								"type": "array",
								"items": {"type": "string"},
								"description": "List of issues found"
							},
							"suggestions": {
								"type": "array",
								"items": {"type": "string"},
								"description": "List of improvement suggestions"
							},
							"rating": {
								"type": "string",
								"description": "Overall code quality rating"
							}
						},
						"required": ["issues", "suggestions", "rating"]
					}`),
				},
			},
		},

		Temperature:   0.5, // Balanced
		MaxIterations: 5,

		APIKey:      apiKey,
		HTTPReferer: "https://github.com/demo-agent",
		XTitle:      "Code Review Agent",
	}

	sampleCode := `
def calculate_average(numbers):
    sum = 0
    for i in numbers:
        sum += i
    return sum / len(numbers)
`

	result = RunAgent(ctx, codeReviewConfig, sampleCode)
	fmt.Printf("✅ Result: %+v\n", result)

	fmt.Println("\n🎪 ═══════════════════════════════════════════════════════")
	fmt.Println("🎪 Demo Complete! Same agent, different tasks!")
	fmt.Println("🎪 ═══════════════════════════════════════════════════════\n")
}

// ═══════════════════════════════════════════════════════════════════════════
// 🎯 KEY INSIGHT: The agent is a PURE FRAMEWORK
// ═══════════════════════════════════════════════════════════════════════════
//
// Notice how the SAME RunAgent() function handles:
// - Math problems (deterministic, temp 0.3)
// - Creative writing (high creativity, temp 0.9)
// - Code review (balanced, temp 0.5)
//
// Each with completely different:
// - System prompts
// - User prompts
// - Tools
// - Temperature settings
// - Models
//
// This is the power of MODULARITY! 🚀
