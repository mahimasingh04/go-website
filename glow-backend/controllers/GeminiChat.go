package controllers

import ("github.com/joho/godotenv"
 "github.com/google/generative-ai-go/genai"
"os"
"context"
"log"
"fmt"
"google.golang.org/api/option"
"github.com/mahimasingh04/go-website/glow-backend/models"
"github.com/gofiber/fiber/v2"

"gorm.io/gorm" // Import gorm
	

)

type ChatRequest struct {
	
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
	Budget  int    `json:"budget,omitempty"` // Optional budget for product suggestions
}

type GeminiChatController struct { // Create a struct to hold dependencies
	DB *gorm.DB
}

func NewGeminiChatController (db *gorm.DB) *GeminiChatController{ // Constructor function
	return &GeminiChatController{DB: db}
}



func (ch *GeminiChatController) ChatWithAI(c *fiber.Ctx) error  {
	req := new(ChatRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	if err := ch.DB.First(&user, req.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Determine the prompt based on the user's message and known information
	prompt, err := ch.generatePrompt(user, req.Message, req.Budget)
	if err != nil {
		log.Printf("Error generating prompt: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate AI prompt"})
	}

	geminiResponse, err := queryGemini(user.ConversationHistory, prompt)
	if err != nil {
		log.Printf("Gemini API error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get response from AI"})
	}

	updatedConversationHistory := user.ConversationHistory + "\nUser: " + req.Message + "\nAI: " + geminiResponse

	user.ConversationHistory = updatedConversationHistory
	if err := ch.DB.Save(&user).Error; err != nil {
    log.Printf("Error saving user: %v", err)
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save conversation"})
}

	return c.JSON(fiber.Map{"response": geminiResponse})
}


func (ch *GeminiChatController) generatePrompt(user models.User, userMessage string, budget int) (string, error) {
	// Add more sophisticated logic here to categorize the query and generate the appropriate prompt
	// Example (very basic):

	if user.SkinType == "" {
		return fmt.Sprintf("Identify the user's skin type (oily, dry, combination, sensitive, normal) based on the following symptoms: %s. Provide a concise skin type classification.", userMessage), nil
	} else if budget > 0 {
		return fmt.Sprintf("The user has %s skin and a budget of $%d. Suggest 2-3 specific skincare products (cleanser, moisturizer, sunscreen) that are suitable for their skin type and within their budget.  Provide links to purchase the products.", user.SkinType, budget), nil
	} else if containsKeywords(userMessage, []string{"acne", "hormonal", "food"}) {
		
		return fmt.Sprintf("The user is experiencing hormonal acne and asking about food recommendations. Recommend 3-5 foods to include in their diet and 3-5 foods to avoid to help manage hormonal acne. Explain the reasoning behind each recommendation."), nil
	} else if containsKeywords(userMessage, []string{"disease", "symptoms"}) {
		return fmt.Sprintf("The user is describing the following skin symptoms: %s. Based on these symptoms, what are the most likely skin diseases or conditions? List 2-3 possibilities and briefly describe each. Indicate if more information is needed for a more accurate diagnosis and remind that the AI does not provide medical advice and a dermatologist should be contacted.", userMessage), nil
	} else {
		return fmt.Sprintf("Answer the following question about skincare: %s. Provide a clear and concise response.", userMessage), nil
	}
}

// Helper function to check if the user message contains certain keywords
func containsKeywords(message string, keywords []string) bool {
	messageLower := string(message) //strings.ToLower(message)
	for _, keyword := range keywords {
		keywordLower := string(keyword) //strings.ToLower(keyword)
		if string(messageLower) == string(keywordLower) {
			return true
		}
	}
	return false
}



func queryGemini(conversationHistory, prompt string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
		return "", fmt.Errorf("error loading .env file: %w", err)  // Return error, don't panic
	}

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY") 
	 fmt.Println("Gemini API Key:", apiKey)

	if apiKey == "" {
		log.Println("GEMINI_API_KEY not set in .env")
		return "", fmt.Errorf("GEMINI_API_KEY not set in .env")  //Handle the error
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to create Gemini client: %w", err)  // Return error
	}
	defer client.Close()

  model := client.GenerativeModel("gemini-pro")
	cs := model.StartChat()

	// Prepend conversation history
	fullPrompt := conversationHistory + "\nUser: " + prompt // Use the generated prompt
	resp, err := cs.SendMessage(ctx, genai.Text(fullPrompt))

	if err != nil {
		log.Println("Error sending message:", err)
		return "", err
	}
responseString := ""
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 { // Added safety check

        for _, part := range resp.Candidates[0].Content.Parts {
			responseString += fmt.Sprintf("%v", part)
		}
	} else {
		log.Println("No candidates or parts in the response!") // Added Error logging
		return "", fmt.Errorf("no content in Gemini response")
	}

return responseString,nil

}
