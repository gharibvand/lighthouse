package service

import (
	"fmt"
	"lighthouse/internal/domain"
	"os"
	"path/filepath"
)

type StreamingService struct {
	storagePath string
}

func NewStreamingService(storagePath string) *StreamingService {
	return &StreamingService{
		storagePath: storagePath,
	}
}

func (s *StreamingService) GetVideoPath(contentID string, episodeID *string, quality domain.StreamingQuality) (string, error) {
	var videoPath string
	if episodeID != nil {
		// TV show episode
		videoPath = filepath.Join(s.storagePath, "content", contentID, 
			"episodes", *episodeID, string(quality))
	} else {
		// Movie
		videoPath = filepath.Join(s.storagePath, "content", contentID, string(quality))
	}

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("video not found")
	}

	return videoPath, nil
}

func (s *StreamingService) GetPlaylistPath(contentID string, episodeID *string, quality domain.StreamingQuality) (string, error) {
	videoPath, err := s.GetVideoPath(contentID, episodeID, quality)
	if err != nil {
		return "", err
	}

	playlistPath := filepath.Join(videoPath, "playlist.m3u8")
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		return "", fmt.Errorf("playlist not found")
	}

	return playlistPath, nil
}

func (s *StreamingService) GetSegmentPath(contentID string, episodeID *string, quality domain.StreamingQuality, segment string) (string, error) {
	videoPath, err := s.GetVideoPath(contentID, episodeID, quality)
	if err != nil {
		return "", err
	}

	segmentPath := filepath.Join(videoPath, segment)
	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		return "", fmt.Errorf("segment not found")
	}

	return segmentPath, nil
}
