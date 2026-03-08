package uuid

import (
	"github.com/google/uuid"
)

// GenerateV7 generates a UUID v7 (time-ordered UUID)
// UUID v7 is time-ordered and provides better index performance than v4
// Benefits:
// - 24% smaller indexes
// - 11x faster index creation
// - Better cache utilization
// - Time-ordered (newest records at the end of B-tree)
func GenerateV7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// MustGenerateV7 generates a UUID v7 and panics on error
// Use this when you're certain the generation will succeed
func MustGenerateV7() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}
