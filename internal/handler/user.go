package handler

import (
	"net/http"
	"strconv"
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
	"lighthouse/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService      *service.UserService
	watchHistoryRepo *postgres.WatchHistoryRepository
	watchlistRepo    *postgres.WatchlistRepository
	ratingRepo       *postgres.RatingRepository
}

func NewUserHandler(userService *service.UserService, watchHistoryRepo *postgres.WatchHistoryRepository,
	watchlistRepo *postgres.WatchlistRepository, ratingRepo *postgres.RatingRepository) *UserHandler {
	return &UserHandler{
		userService:      userService,
		watchHistoryRepo: watchHistoryRepo,
		watchlistRepo:    watchlistRepo,
		ratingRepo:       ratingRepo,
	}
}

func (h *UserHandler) GetProfiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profiles, err := h.userService.GetProfiles(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

func (h *UserHandler) CreateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var profile domain.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.UserID = userID.(string)
	if err := h.userService.CreateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"profile": profile})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var profile domain.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ID = id
	profile.UserID = userID.(string)

	if err := h.userService.UpdateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (h *UserHandler) DeleteProfile(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteProfile(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}

func (h *UserHandler) GetWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	limit := 20
	offset := 0

	watchlist, err := h.watchlistRepo.GetByUserID(userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"watchlist": watchlist})
}

func (h *UserHandler) AddToWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		ProfileID string `json:"profile_id" binding:"required"`
		ContentID string `json:"content_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	watchlist := &domain.Watchlist{
		UserID:    userID.(string),
		ProfileID: req.ProfileID,
		ContentID: req.ContentID,
	}

	if err := h.watchlistRepo.Create(watchlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"watchlist": watchlist})
}

func (h *UserHandler) RemoveFromWatchlist(c *gin.Context) {
	userID, _ := c.Get("user_id")
	contentID := c.Param("contentId")

	if err := h.watchlistRepo.Delete(userID.(string), contentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from watchlist"})
}

func (h *UserHandler) GetWatchHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	limit := 20
	offset := 0

	history, err := h.watchHistoryRepo.GetByUserID(userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": history})
}

func (h *UserHandler) AddRating(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var rating domain.Rating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating.UserID = userID.(string)

	// Check if rating exists
	existing, _ := h.ratingRepo.GetUserRating(userID.(string), rating.ContentID)
	if existing != nil {
		rating.ID = existing.ID
		if err := h.ratingRepo.Update(&rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := h.ratingRepo.Create(&rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"rating": rating})
}
