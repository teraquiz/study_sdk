package study_sdk

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db         *mongo.Database
	repository *FlashcardRepository
}

type Config struct {
	MongoURI     string
	DatabaseName string
	Timeout      time.Duration
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := mongoClient.Database(cfg.DatabaseName)
	repository := NewFlashcardRepository(db)

	return &Client{
		db:         db,
		repository: repository,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.db.Client().Disconnect(ctx)
}

func (c *Client) GetFlashcardsByCategory(ctx context.Context, categoryID string) ([]Flashcard, error) {
	return c.repository.FindByCategory(ctx, categoryID)
}

func (c *Client) GetFlashcardsByCategories(ctx context.Context, categoryIDs []string) ([]Flashcard, error) {
	return c.repository.FindByCategories(ctx, categoryIDs)
}

func (c *Client) GetFlashcardByID(ctx context.Context, id string) (*Flashcard, error) {
	return c.repository.FindByID(ctx, id)
}

func (c *Client) ListFlashcards(ctx context.Context, filter FlashcardFilter) ([]Flashcard, error) {
	return c.repository.FindWithFilters(ctx, filter)
}

func (c *Client) GetFlashcardsByIDs(ctx context.Context, ids []string) ([]Flashcard, error) {
	return c.repository.FindByIDs(ctx, ids)
}

func (c *Client) GetCategoriesByFlashcard(ctx context.Context, flashcardID string) ([]Category, error) {
	return c.repository.FindCategoriesByFlashcard(ctx, flashcardID)
}

func (c *Client) GetCategoryByID(ctx context.Context, id string) (*Category, error) {
	return c.repository.FindCategoryByID(ctx, id)
}

func (c *Client) ListCategories(ctx context.Context, filter CategoryFilter) ([]Category, error) {
	return c.repository.FindCategoriesWithFilters(ctx, filter)
}

func (c *Client) GetProductByID(ctx context.Context, id string) (*Product, error) {
	return c.repository.FindProductByID(ctx, id)
}

func (c *Client) GetProductsByCategory(ctx context.Context, categoryID string) ([]Product, error) {
	return c.repository.FindProductsByCategory(ctx, categoryID)
}

func (c *Client) ListProducts(ctx context.Context, filter ProductFilter) ([]Product, error) {
	return c.repository.FindAllProducts(ctx, filter)
}
