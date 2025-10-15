package study_sdk

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlashcardRepository struct {
	flashcardCollection          *mongo.Collection
	flashcardCategoryCollection  *mongo.Collection
}

func NewFlashcardRepository(db *mongo.Database) *FlashcardRepository {
	return &FlashcardRepository{
		flashcardCollection:         db.Collection("flashcards"),
		flashcardCategoryCollection: db.Collection("flashcard_categories"),
	}
}

func (r *FlashcardRepository) FindByID(ctx context.Context, id string) (*Flashcard, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var flashcard Flashcard
	err = r.flashcardCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&flashcard)
	if err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (r *FlashcardRepository) FindByCategory(ctx context.Context, categoryID string) ([]Flashcard, error) {
	relations, err := r.findRelationsByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if len(relations) == 0 {
		return []Flashcard{}, nil
	}

	flashcardIDs := make([]primitive.ObjectID, len(relations))
	for i, rel := range relations {
		objID, err := primitive.ObjectIDFromHex(rel.FlashcardID)
		if err != nil {
			continue
		}
		flashcardIDs[i] = objID
	}

	cursor, err := r.flashcardCollection.Find(ctx, bson.M{"_id": bson.M{"$in": flashcardIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flashcards []Flashcard
	if err = cursor.All(ctx, &flashcards); err != nil {
		return nil, err
	}

	return flashcards, nil
}

func (r *FlashcardRepository) FindWithFilters(ctx context.Context, filter FlashcardFilter) ([]Flashcard, error) {
	query := bson.M{}

	if filter.Difficulty != nil {
		query["difficulty"] = *filter.Difficulty
	}

	if filter.Language != nil {
		query["language"] = *filter.Language
	}

	if filter.Verified != nil {
		query["verified"] = *filter.Verified
	}

	if filter.Enabled != nil {
		query["enabled"] = *filter.Enabled
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	if filter.CategoryID != nil {
		relations, err := r.findRelationsByCategory(ctx, *filter.CategoryID)
		if err != nil {
			return nil, err
		}

		flashcardIDs := make([]primitive.ObjectID, 0, len(relations))
		for _, rel := range relations {
			objID, err := primitive.ObjectIDFromHex(rel.FlashcardID)
			if err != nil {
				continue
			}
			flashcardIDs = append(flashcardIDs, objID)
		}

		query["_id"] = bson.M{"$in": flashcardIDs}
	}

	cursor, err := r.flashcardCollection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flashcards []Flashcard
	if err = cursor.All(ctx, &flashcards); err != nil {
		return nil, err
	}

	return flashcards, nil
}

type flashcardCategoryRelation struct {
	FlashcardID string `bson:"flashcard_id"`
	CategoryID  string `bson:"category_id"`
}

func (r *FlashcardRepository) findRelationsByCategory(ctx context.Context, categoryID string) ([]flashcardCategoryRelation, error) {
	cursor, err := r.flashcardCategoryCollection.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var relations []flashcardCategoryRelation
	if err = cursor.All(ctx, &relations); err != nil {
		return nil, err
	}

	return relations, nil
}
