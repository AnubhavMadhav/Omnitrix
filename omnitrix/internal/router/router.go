package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/AnubhavMadhav/Omnitrix/omnitrix/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/llm"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/provider"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/reflex"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/resolver"
	"github.com/AnubhavMadhav/Omnitrix/omnitrix/internal/utils"
)

type ChatRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type ChatResponse struct {
	Source        string `json:"source" example:"llm_generation"`
	DetectedLabel string `json:"detected_label" example:"coding"`
	TargetModel   string `json:"target_model" example:"DeepSeek-Coder-V2-Lite"`
	Response      string `json:"response" example:"def sort(l): return sorted(l)"`
}

func New(reflexEngine *reflex.Engine, llmClient llm.Client, factory *provider.Factory) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/chat", chatHandler(reflexEngine, llmClient, factory))
	}

	return r
}

// @Summary      Process a User Prompt
// @Description  Routes the prompt to the correct LLM based on intent and user tier.
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        X-User-Tier header string false "User Tier (free/premium)" default(free)
// @Param        request body ChatRequest true "User Prompt"
// @Success      200  {object}  ChatResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chat [post]

func chatHandler(reflexEngine *reflex.Engine, llmClient llm.Client, factory *provider.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Reflex Logic (For Greetings or Blocked words)
		intent, response := reflexEngine.Search(req.Prompt)
		intentType := "none"
		if intent != reflex.IntentNone {
			if intent == 1 {
				intentType = "greeting"
			} else {
				intentType = "blocked"
			}
			c.JSON(http.StatusOK, ChatResponse{
				Source:        "reflex",
				DetectedLabel: intentType,
				Response:      response,
			})
			return
		}

		fmt.Println("intent : ", intent)

		// Classification Logic (decide the "label")
		classification, err := llmClient.Classify(c.Request.Context(), req.Prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Brain freeze: " + err.Error()})
			return
		}

		// Resolver Logic (Decide which model to use)
		userTier := c.GetHeader("X-User-Tier")
		targetModel := resolver.ResolveModel(userTier, classification.Label)

		fmt.Println("targetModel : ", targetModel)

		// Execution Logic
		generator, err := factory.GetGenerator(targetModel)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Provider unavailable: " + err.Error()})
			return
		}

		finalResponse, err := generator.Generate(c.Request.Context(), targetModel, req.Prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed: " + err.Error()})
			return
		}

		utils.PrintCleanLog(targetModel, finalResponse)

		c.JSON(http.StatusOK, ChatResponse{
			Source:        "llm_generation",
			DetectedLabel: string(classification.Label),
			TargetModel:   targetModel,
			Response:      finalResponse,
		})
	}
}
