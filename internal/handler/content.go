package handler

import (
	"net/http"
	"strconv"
	"lighthouse/internal/service"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	contentService      *service.ContentService
	recommendationService *service.RecommendationService
}

func NewContentHandler(contentService *service.ContentService, recommendationService *service.RecommendationService) *ContentHandler {
	return &ContentHandler{
		contentService:      contentService,
		recommendationService: recommendationService,
	}
}

// GetMovies godoc
// @Summary      Get movies list
// @Description  Get paginated list of movies
// @Tags         content
// @Produce      json
// @Param        limit   query     int  false  "Limit"  default(20)
// @Param        offset  query     int  false  "Offset" default(0)
// @Success      200     {object}  map[string]interface{}
// @Router       /content/movies [get]
func (h *ContentHandler) GetMovies(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	movies, err := h.contentService.GetMovies(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

func (h *ContentHandler) GetTVShows(c *gin.Context) {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	tvShows, err := h.contentService.GetTVShows(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tv_shows": tvShows})
}

func (h *ContentHandler) GetContentByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	content, err := h.contentService.GetContentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": content})
}

// Search godoc
// @Summary      Search content
// @Description  Search movies and TV shows by query
// @Tags         content
// @Produce      json
// @Param        q       query     string  true   "Search query"
// @Param        limit   query     int     false  "Limit"  default(20)
// @Param        offset  query     int     false  "Offset" default(0)
// @Success      200     {object}  map[string]interface{}
// @Router       /content/search [get]
func (h *ContentHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	results, err := h.contentService.Search(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// GetTrending godoc
// @Summary      Get trending content
// @Description  Get trending movies and TV shows
// @Tags         content
// @Produce      json
// @Param        limit   query     int  false  "Limit"  default(10)
// @Success      200     {object}  map[string]interface{}
// @Router       /content/trending [get]
func (h *ContentHandler) GetTrending(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	trending, err := h.contentService.GetTrending(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trending": trending})
}

// GetRecommendations godoc
// @Summary      Get personalized recommendations
// @Description  Get content recommendations based on user's watch history
// @Tags         content
// @Security     BearerAuth
// @Produce      json
// @Param        limit   query     int  false  "Limit"  default(10)
// @Success      200     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]string
// @Router       /content/recommendations [get]
func (h *ContentHandler) GetRecommendations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	limit := 10

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	recommendations, err := h.recommendationService.GetRecommendations(userID.(string), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}
