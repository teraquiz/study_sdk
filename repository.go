package study_sdk

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlashcardRepository struct {
	flashcardCollection         *mongo.Collection
	flashcardCategoryCollection *mongo.Collection
	categoryCollection          *mongo.Collection
	categoryProductCollection   *mongo.Collection
	productCollection           *mongo.Collection
}

func NewFlashcardRepository(db *mongo.Database) *FlashcardRepository {
	return &FlashcardRepository{
		flashcardCollection:         db.Collection("flashcards"),
		flashcardCategoryCollection: db.Collection("flashcard_categories"),
		categoryCollection:          db.Collection("categories"),
		categoryProductCollection:   db.Collection("category_products"),
		productCollection:           db.Collection("products"),
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

func (r *FlashcardRepository) FindByIDs(ctx context.Context, ids []string) ([]Flashcard, error) {
	if len(ids) == 0 {
		return []Flashcard{}, nil
	}

	objIDs := stringsToObjectIDs(ids)
	if len(objIDs) == 0 {
		return []Flashcard{}, nil
	}

	return findAll[Flashcard](ctx, r.flashcardCollection, bson.M{"_id": bson.M{"$in": objIDs}})
}

func (r *FlashcardRepository) FindByCategory(ctx context.Context, categoryID string) ([]Flashcard, error) {
	return r.FindByCategories(ctx, []string{categoryID})
}

func (r *FlashcardRepository) FindByCategories(ctx context.Context, categoryIDs []string) ([]Flashcard, error) {
	if len(categoryIDs) == 0 {
		return []Flashcard{}, nil
	}

	relations, err := findAll[flashcardCategoryRelation](ctx, r.flashcardCategoryCollection, bson.M{"category_id": bson.M{"$in": categoryIDs}})
	if err != nil || len(relations) == 0 {
		return []Flashcard{}, err
	}

	flashcardIDsMap := make(map[string]bool)
	for _, rel := range relations {
		flashcardIDsMap[rel.FlashcardID] = true
	}

	flashcardIDs := make([]string, 0, len(flashcardIDsMap))
	for id := range flashcardIDsMap {
		flashcardIDs = append(flashcardIDs, id)
	}

	return findAll[Flashcard](ctx, r.flashcardCollection, bson.M{"_id": bson.M{"$in": stringsToObjectIDs(flashcardIDs)}})
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
		relations, err := findAll[flashcardCategoryRelation](ctx, r.flashcardCategoryCollection, bson.M{"category_id": *filter.CategoryID})
		if err != nil {
			return nil, err
		}

		ids := make([]string, len(relations))
		for i, rel := range relations {
			ids[i] = rel.FlashcardID
		}
		query["_id"] = bson.M{"$in": stringsToObjectIDs(ids)}
	}

	return findAll[Flashcard](ctx, r.flashcardCollection, query)
}

type flashcardCategoryRelation struct {
	FlashcardID string `bson:"flashcard_id"`
	CategoryID  string `bson:"category_id"`
}

type categoryProductRelation struct {
	CategoryID string `bson:"category_id"`
	ProductID  string `bson:"product_id"`
}

func stringsToObjectIDs(ids []string) []primitive.ObjectID {
	result := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			result = append(result, objID)
		}
	}
	return result
}

func findAll[T any](ctx context.Context, collection *mongo.Collection, filter bson.M) ([]T, error) {
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *FlashcardRepository) FindCategoriesByFlashcard(ctx context.Context, flashcardID string) ([]Category, error) {
	relations, err := findAll[flashcardCategoryRelation](ctx, r.flashcardCategoryCollection, bson.M{"flashcard_id": flashcardID})
	if err != nil || len(relations) == 0 {
		return []Category{}, err
	}

	ids := make([]string, len(relations))
	for i, rel := range relations {
		ids[i] = rel.CategoryID
	}

	return findAll[Category](ctx, r.categoryCollection, bson.M{"_id": bson.M{"$in": stringsToObjectIDs(ids)}})
}

func (r *FlashcardRepository) FindCategoryByID(ctx context.Context, id string) (*Category, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var category Category
	err = r.categoryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *FlashcardRepository) FindCategoriesWithFilters(ctx context.Context, filter CategoryFilter) ([]Category, error) {
	query := bson.M{}

	if filter.Type != nil {
		query["type"] = *filter.Type
	}
	if filter.Enabled != nil {
		query["enabled"] = *filter.Enabled
	}

	if filter.ProductID != nil {
		relations, err := findAll[categoryProductRelation](ctx, r.categoryProductCollection, bson.M{"product_id": *filter.ProductID})
		if err != nil {
			return nil, err
		}
		if len(relations) == 0 {
			return []Category{}, nil
		}

		ids := make([]string, len(relations))
		for i, rel := range relations {
			ids[i] = rel.CategoryID
		}
		query["_id"] = bson.M{"$in": stringsToObjectIDs(ids)}
	}

	return findAll[Category](ctx, r.categoryCollection, query)
}

func (r *FlashcardRepository) FindProductByID(ctx context.Context, id string) (*Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product Product
	err = r.productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *FlashcardRepository) FindProductsByCategory(ctx context.Context, categoryID string) ([]Product, error) {
	relations, err := findAll[categoryProductRelation](ctx, r.categoryProductCollection, bson.M{"category_id": categoryID})
	if err != nil || len(relations) == 0 {
		return []Product{}, err
	}

	ids := make([]string, len(relations))
	for i, rel := range relations {
		ids[i] = rel.ProductID
	}

	return findAll[Product](ctx, r.productCollection, bson.M{"_id": bson.M{"$in": stringsToObjectIDs(ids)}})
}

func (r *FlashcardRepository) FindAllProducts(ctx context.Context, filter ProductFilter) ([]Product, error) {
	query := bson.M{}

	if filter.Type != nil {
		query["metadata.type"] = *filter.Type
	}
	if filter.Enabled != nil {
		query["enabled"] = *filter.Enabled
	}

	return findAll[Product](ctx, r.productCollection, query)
}
