package handler

import (
	"net/http"
	"strconv"
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
	"lighthouse/internal/service"

	"github.com/gin-gonic/gin"
)

type StreamingHandler struct {
	streamingService  *service.StreamingService
	watchHistoryRepo *postgres.WatchHistoryRepository
}

func NewStreamingHandler(streamingService *service.StreamingService, watchHistoryRepo *postgres.WatchHistoryRepository) *StreamingHandler {
	return &StreamingHandler{
		streamingService:  streamingService,
		watchHistoryRepo: watchHistoryRepo,
	}
}

func (h *StreamingHandler) GetPlaylist(c *gin.Context) {
	contentID := c.Param("contentId")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	quality := domain.StreamingQuality(c.Query("quality"))
	if quality == "" {
		quality = domain.Quality720p
	}

	var episodeID *string
	if epID := c.Query("episode_id"); epID != "" {
		episodeID = &epID
	}

	playlistPath, err := h.streamingService.GetPlaylistPath(contentID, episodeID, quality)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	c.File(playlistPath)
}

func (h *StreamingHandler) GetSegment(c *gin.Context) {
	contentID := c.Param("contentId")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	quality := domain.StreamingQuality(c.Param("quality"))
	segment := c.Param("segment")

	var episodeID *string
	if epID := c.Query("episode_id"); epID != "" {
		episodeID = &epID
	}

	segmentPath, err := h.streamingService.GetSegmentPath(contentID, episodeID, quality, segment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
		return
	}

	c.File(segmentPath)
}

func (h *StreamingHandler) UpdateWatchProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	contentID := c.Param("contentId")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	var req struct {
		ProfileID string  `json:"profile_id" binding:"required"`
		EpisodeID *string `json:"episode_id"`
		Progress  int     `json:"progress" binding:"required"`
		Duration  int     `json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if history exists
	history, _ := h.watchHistoryRepo.GetByContentID(userID.(string), contentID)
	if history != nil {
		history.Progress = req.Progress
		history.Duration = req.Duration
		history.EpisodeID = req.EpisodeID
		if err := h.watchHistoryRepo.Update(history); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		history = &domain.WatchHistory{
			UserID:    userID.(string),
			ProfileID: req.ProfileID,
			ContentID: contentID,
			EpisodeID: req.EpisodeID,
			Progress:  req.Progress,
			Duration:  req.Duration,
		}
		if err := h.watchHistoryRepo.Create(history); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
