package study_sdk

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flashcard struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Language   string             `bson:"language"`
	Front      string             `bson:"front"`
	Back       string             `bson:"back"`
	Hint       string             `bson:"hint,omitempty"`
	Difficulty string             `bson:"difficulty"`
	Images     []FlashcardImage   `bson:"images,omitempty"`
	Tags       []string           `bson:"tags,omitempty"`
	Enabled    bool               `bson:"enabled"`
	Verified   bool               `bson:"verified"`
	CreatedBy  string             `bson:"created_by"`
	VerifiedBy string             `bson:"verified_by,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

type FlashcardImage struct {
	URL     string `bson:"url"`
	Caption string `bson:"caption,omitempty"`
}

type FlashcardFilter struct {
	CategoryID *string
	Difficulty *string
	Language   *string
	Verified   *bool
	Enabled    *bool
	Tags       []string
}
