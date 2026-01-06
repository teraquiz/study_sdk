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

type Category struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Description    string             `bson:"description"`
	Type           string             `bson:"type"`
	Icon           string             `bson:"icon,omitempty"`
	Color          string             `bson:"color,omitempty"`
	ParentID       string             `bson:"parent_id,omitempty"`
	TotalQuestions int                `bson:"total_questions"`
	Enabled        bool               `bson:"enabled"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

type CategoryFilter struct {
	Type      *string
	Enabled   *bool
	ProductID *string
}

type ProductFilter struct {
	Type    *string
	Enabled *bool
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	AreaIDs     []int              `bson:"area_ids"`
	Metadata    ProductMetadata    `bson:"metadata"`
	Enabled     bool               `bson:"enabled"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type ProductMetadata struct {
	Type            string   `bson:"type"`
	TotalQuestions  int      `bson:"total_questions"`
	TotalCategories int      `bson:"total_categories"`
	Languages       []string `bson:"languages"`
}
