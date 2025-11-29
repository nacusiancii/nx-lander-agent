package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	openrouter "github.com/revrost/go-openrouter"
)

// ═══════════════════════════════════════════════════════════════════════════
// 🧠 AI AGENT - Pure Thinking Engine with THINK → ACT → OBSERVE → REFINE Loop
// ═══════════════════════════════════════════════════════════════════════════

// AgentConfig defines the complete external configuration for an agent instance.
// Zero hardcoded values - everything is injected!
type AgentConfig struct {
	// Model configuration
	ModelName string   // e.g., "minimax/minimax-m2"
	Providers []string // e.g., ["minimax/fp8"]
	
	// Prompts (injected externally for maximum flexibility)
	SystemPrompt     string // Agent's role and expertise
	UserPromptFormat string // Task description with %s placeholder for theme/input
	
	// Tools (injected externally - agent is tool-agnostic)
	Tools []openrouter.Tool
	
	// Behavior tuning
	Temperature   float64 // Creativity vs determinism
	MaxIterations int     // Safety limit for the thinking loop
	
	// API credentials
	APIKey     string
	HTTPReferer string
	XTitle      string
}

// AgentState tracks the agent's internal state throughout its execution.
type AgentState struct {
	iteration int                                   // Current iteration number
	messages  []openrouter.ChatCompletionMessage   // Conversation history
	phase     string                                // Current phase: THINK/ACT/OBSERVE/REFINE
	result    interface{}                           // Final result (if completed)
	completed bool                                  // Whether task is done
}

// AgentResult represents the final output from the agent.
type AgentResult struct {
	Success    bool        // Whether the task completed successfully
	Result     interface{} // The actual result data
	Iterations int         // Number of iterations taken
	Error      error       // Error if failed
}

// ═══════════════════════════════════════════════════════════════════════════
// 🎯 CORE AGENT: The Beautiful Thinking Loop
// ═══════════════════════════════════════════════════════════════════════════

// RunAgent executes the AI agent with THINK → ACT → OBSERVE → REFINE loop.
// This is the main entry point - pure, modular, and magnificent!
func RunAgent(ctx context.Context, config AgentConfig, input string) AgentResult {
	log.Printf("🚀 Agent started with model: %s", config.ModelName)
	
	// Initialize the agent's state
	state := &AgentState{
		iteration: 0,
		messages:  initializeMessages(config, input),
		phase:     "THINK",
		completed: false,
	}
	
	// Create OpenRouter client
	client := openrouter.NewClient(
		config.APIKey,
		openrouter.WithHTTPReferer(config.HTTPReferer),
		openrouter.WithXTitle(config.XTitle),
	)
	
	// 🔄 THE MAGNIFICENT LOOP - Where the magic happens!
	for state.iteration < config.MaxIterations && !state.completed {
		state.iteration++
		log.Printf("\n═══ Iteration %d/%d ═══", state.iteration, config.MaxIterations)
		
		// 💭 THINK: Agent reasons about current state
		if err := think(ctx, client, config, state); err != nil {
			return AgentResult{Success: false, Error: err, Iterations: state.iteration}
		}
		
		// 🎬 ACT: Agent executes action (tool call or answer)
		response, err := act(ctx, client, config, state)
		if err != nil {
			return AgentResult{Success: false, Error: err, Iterations: state.iteration}
		}
		
		// 👁️ OBSERVE: Agent processes the action result
		if err := observe(state, response); err != nil {
			return AgentResult{Success: false, Error: err, Iterations: state.iteration}
		}
		
		// ✨ REFINE: Agent decides whether to continue or conclude
		if refine(state) {
			log.Printf("✅ Task completed in %d iterations", state.iteration)
			return AgentResult{
				Success:    true,
				Result:     state.result,
				Iterations: state.iteration,
			}
		}
	}
	
	// Max iterations reached without completion
	if !state.completed {
		log.Printf("⚠️  Max iterations (%d) reached", config.MaxIterations)
		return AgentResult{
			Success:    false,
			Error:      fmt.Errorf("max iterations reached without completion"),
			Iterations: state.iteration,
		}
	}
	
	return AgentResult{Success: true, Result: state.result, Iterations: state.iteration}
}

// ═══════════════════════════════════════════════════════════════════════════
// 🧩 PHASE IMPLEMENTATIONS: Each phase with single responsibility
// ═══════════════════════════════════════════════════════════════════════════

// initializeMessages sets up the initial conversation context.
func initializeMessages(config AgentConfig, input string) []openrouter.ChatCompletionMessage {
	return []openrouter.ChatCompletionMessage{
		{
			Role: openrouter.ChatMessageRoleSystem,
			Content: openrouter.Content{
				Text: config.SystemPrompt,
			},
		},
		{
			Role: openrouter.ChatMessageRoleUser,
			Content: openrouter.Content{
				Text: fmt.Sprintf(config.UserPromptFormat, input),
			},
		},
	}
}

// think prepares the agent for the next action.
// In this implementation, thinking happens implicitly in the LLM call.
func think(ctx context.Context, client *openrouter.Client, config AgentConfig, state *AgentState) error {
	state.phase = "THINK"
	log.Printf("💭 THINK: Planning next action...")
	// Thinking is embedded in the model's reasoning process
	return nil
}

// act executes the agent's chosen action by calling the LLM.
func act(ctx context.Context, client *openrouter.Client, config AgentConfig, state *AgentState) (*openrouter.ChatCompletionResponse, error) {
	state.phase = "ACT"
	log.Printf("🎬 ACT: Executing action...")
	
	// Build the request
	req := openrouter.ChatCompletionRequest{
		Model:       config.ModelName,
		Messages:    state.messages,
		Tools:       config.Tools,
		Temperature: float32(config.Temperature),
	}
	
	// Add provider configuration if specified
	if len(config.Providers) > 0 {
		req.Provider = &openrouter.ChatProvider{
			Order:          config.Providers,
			AllowFallbacks: boolPtr(false),
		}
	}
	
	// Execute the action (LLM call)
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("action failed: %w", err)
	}
	
	return &resp, nil
}

// observe processes the result of the action.
func observe(state *AgentState, response *openrouter.ChatCompletionResponse) error {
	state.phase = "OBSERVE"
	log.Printf("👁️  OBSERVE: Processing response...")
	
	if len(response.Choices) == 0 {
		return fmt.Errorf("no response choices received")
	}
	
	choice := response.Choices[0]
	
	// Check if agent made a tool call
	if len(choice.Message.ToolCalls) > 0 {
		toolCall := choice.Message.ToolCalls[0]
		log.Printf("   → Tool called: %s", toolCall.Function.Name)
		
		// Parse tool arguments to get the result
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &result); err != nil {
			return fmt.Errorf("failed to parse tool arguments: %w", err)
		}
		
		state.result = result
		state.completed = true
		return nil
	}
	
	// Check if agent provided a text response (no tool call)
	if choice.Message.Content.Text != "" {
		log.Printf("   → Text response received")
		state.result = choice.Message.Content.Text
		// Continue iterating - agent might need more thinking
		return nil
	}
	
	return fmt.Errorf("unexpected response format")
}

// refine determines whether the agent should continue or conclude.
func refine(state *AgentState) bool {
	state.phase = "REFINE"
	log.Printf("✨ REFINE: Evaluating completion status...")
	
	if state.completed {
		log.Printf("   → Task complete!")
		return true
	}
	
	log.Printf("   → Continuing to next iteration...")
	return false
}

// ═══════════════════════════════════════════════════════════════════════════
// 🛠️ HELPER UTILITIES
// ═══════════════════════════════════════════════════════════════════════════

// ExtractSearchTerms is a convenience function for extracting search terms from agent results.
func ExtractSearchTerms(result AgentResult) ([]string, error) {
	if !result.Success {
		return nil, result.Error
	}
	
	// Extract from tool call result
	if resultMap, ok := result.Result.(map[string]interface{}); ok {
		if terms, ok := resultMap["search_terms"].([]interface{}); ok {
			searchTerms := make([]string, len(terms))
			for i, term := range terms {
				searchTerms[i] = term.(string)
			}
			return searchTerms, nil
		}
	}
	
	return nil, fmt.Errorf("result does not contain search_terms")
}

// ExtractKeywords is a convenience function for extracting keywords from agent results.
func ExtractKeywords(result AgentResult) ([]string, error) {
	if !result.Success {
		return nil, result.Error
	}
	
	// Extract from tool call result
	if resultMap, ok := result.Result.(map[string]interface{}); ok {
		if kws, ok := resultMap["keywords"].([]interface{}); ok {
			keywords := make([]string, len(kws))
			for i, kw := range kws {
				keywords[i] = kw.(string)
			}
			return keywords, nil
		}
	}
	
	return nil, fmt.Errorf("result does not contain keywords")
}
