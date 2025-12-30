package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/iannsp/shiftopt/internal/models"
	"google.golang.org/api/option"
)

// ParseConstraint attempts to use Gemini AI, falling back to Mock if no key is found.
func ParseConstraint(input string) models.Unavailability {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("[AI] No API Key found. Using Mock Parser.")
		return MockParse(input)
	}

	result, err := callGemini(context.Background(), apiKey, input)
	if err != nil {
		log.Printf("[AI] Gemini Error: %v. Falling back to Mock.\n", err)
		return MockParse(input)
	}

	return result
}

func callGemini(ctx context.Context, apiKey, input string) (models.Unavailability, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return models.Unavailability{}, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	
	// 1. Configure for JSON Mode (Structured Output)
	model.ResponseMIMEType = "application/json"
	
	// 2. The System Prompt (The "Brain" Instructions)
	prompt := fmt.Sprintf(`
	You are a scheduling assistant. Extract availability constraints from the user's text.
	Return a SINGLE JSON object with this exact schema:
	{
		"EmployeeName": "string (Matches one of: Alice (Vet), Bob (Vet), Carol (Vet), Dave (Jun), Eve (Jun), Frank (Jun), Grace (Grinder), Hank (Grinder))",
		"StartHour": int (0-23),
		"EndHour": int (0-23, exclusive),
		"Reason": "string (short summary)"
	}
	
	Rules:
	- "Morning" = 08:00 to 12:00
	- "Afternoon" = 13:00 to 17:00
	- "All day" = 08:00 to 20:00
	- Use fuzzy matching to map names like "Alice" to "Alice (Vet)".
	
	User Input: "%s"
	`, input)

	// 3. Generate
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return models.Unavailability{}, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return models.Unavailability{}, fmt.Errorf("empty response from AI")
	}

	// 4. Parse JSON
	// Gemini returns the JSON string in the first part
	rawJSON, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return models.Unavailability{}, fmt.Errorf("unexpected response format")
	}

	var parsed models.Unavailability
	if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
		return models.Unavailability{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return parsed, nil
}

// MockParse (Preserved for Fallback/Testing)
func MockParse(input string) models.Unavailability {
	input = strings.ToLower(input)
	result := models.Unavailability{}

	// Simple Heuristics (Same as before)
	if strings.Contains(input, "alice") {
		result.EmployeeName = "Alice (Vet)"
	} else if strings.Contains(input, "bob") {
		result.EmployeeName = "Bob (Vet)"
	} else {
		result.EmployeeName = "Dave (Jun)" // Default fallback
	}

	if strings.Contains(input, "morning") {
		result.StartHour = 8; result.EndHour = 12
	} else if strings.Contains(input, "afternoon") {
		result.StartHour = 13; result.EndHour = 17
	} else {
		result.StartHour = 9; result.EndHour = 10
	}

	if strings.Contains(input, "dentist") {
		result.Reason = "Medical (Dentist)"
	} else {
		result.Reason = "Personal"
	}
	return result
}
