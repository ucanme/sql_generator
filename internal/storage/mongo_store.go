package storage

import (
	"context"
	"fmt"
	"time"

	"awesomeProject2/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store defines the interface for data storage
type Store interface {
	// Table operations
	CreateTable(table *models.Table) error
	GetTableByName(name string) (*models.Table, error)
	SearchTables(keyword string, limit, offset int) ([]*models.Table, error)
	ListTables(limit, offset int) ([]*models.Table, error)
	UpdateTable(name string, table *models.Table) error
	DeleteTable(name string) error

	// Query operations
	CreateQuery(query *models.Query) error
	GetQueryByID(id string) (*models.Query, error)
	ListQueries(limit, offset int) ([]*models.Query, error)
}

// MongoStore implements Store interface with MongoDB
type MongoStore struct {
	client   *mongo.Client
	db       *mongo.Database
	tables   *mongo.Collection
	queries  *mongo.Collection
	context  context.Context
}

// NewMongoStore creates a new MongoDB storage
func NewMongoStore(mongoURI, dbName string) (*MongoStore, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(dbName)
	tables := database.Collection("tables")
	queries := database.Collection("queries")

	// Create indexes
	_, err = tables.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"name", 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("name_unique_index"),
		},
		{
			Keys: bson.D{{"description", "text"}, {"columns.description", "text"}, {"columns.name", "text"}},
			Options: options.Index().
				SetName("text_search_index"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &MongoStore{
		client:  client,
		db:      database,
		tables:  tables,
		queries: queries,
		context: ctx,
	}, nil
}

// CreateTable saves a table definition
func (s *MongoStore) CreateTable(table *models.Table) error {
	now := time.Now()
	table.CreatedAt = now
	table.UpdatedAt = now

	_, err := s.tables.InsertOne(s.context, table)
	if err != nil {
		return fmt.Errorf("failed to insert table: %w", err)
	}
	return nil
}

// GetTableByName retrieves a table by name
func (s *MongoStore) GetTableByName(name string) (*models.Table, error) {
	var table models.Table
	err := s.tables.FindOne(s.context, bson.M{"name": name}).Decode(&table)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("table not found: %s", name)
		}
		return nil, fmt.Errorf("failed to find table: %w", err)
	}
	return &table, nil
}

// SearchTables searches tables by keywords
func (s *MongoStore) SearchTables(keyword string, limit, offset int) ([]*models.Table, error) {
	if limit <= 0 {
		limit = 10
	}
	
	if limit > 100 {
		limit = 100 // Cap at 100 results
	}

	filter := bson.M{"$text": bson.M{"$search": keyword}}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	cursor, err := s.tables.Find(s.context, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search tables: %w", err)
	}
	defer cursor.Close(s.context)

	var tables []*models.Table
	if err = cursor.All(s.context, &tables); err != nil {
		return nil, fmt.Errorf("failed to decode tables: %w", err)
	}

	return tables, nil
}

// ListTables returns tables with pagination
func (s *MongoStore) ListTables(limit, offset int) ([]*models.Table, error) {
	if limit <= 0 {
		limit = 10
	}
	
	if limit > 100 {
		limit = 100 // Cap at 100 results
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := s.tables.Find(s.context, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer cursor.Close(s.context)

	var tables []*models.Table
	if err = cursor.All(s.context, &tables); err != nil {
		return nil, fmt.Errorf("failed to decode tables: %w", err)
	}

	return tables, nil
}

// UpdateTable updates a table by name
func (s *MongoStore) UpdateTable(name string, table *models.Table) error {
	table.UpdatedAt = time.Now()
	
	result, err := s.tables.UpdateOne(
		s.context,
		bson.M{"name": name},
		bson.M{"$set": table},
	)
	if err != nil {
		return fmt.Errorf("failed to update table: %w", err)
	}
	
	if result.MatchedCount == 0 {
		return fmt.Errorf("table not found: %s", name)
	}
	
	return nil
}

// DeleteTable removes a table by name
func (s *MongoStore) DeleteTable(name string) error {
	result, err := s.tables.DeleteOne(s.context, bson.M{"name": name})
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}
	
	if result.DeletedCount == 0 {
		return fmt.Errorf("table not found: %s", name)
	}
	
	return nil
}

// CreateQuery saves a generated query
func (s *MongoStore) CreateQuery(query *models.Query) error {
	query.CreatedAt = time.Now()
	
	_, err := s.queries.InsertOne(s.context, query)
	if err != nil {
		return fmt.Errorf("failed to insert query: %w", err)
	}
	return nil
}

// GetQueryByID retrieves a query by ID
func (s *MongoStore) GetQueryByID(id string) (*models.Query, error) {
	var query models.Query
	err := s.queries.FindOne(s.context, bson.M{"_id": id}).Decode(&query)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("query not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find query: %w", err)
	}
	return &query, nil
}

// ListQueries returns all queries with pagination
func (s *MongoStore) ListQueries(limit, offset int) ([]*models.Query, error) {
	if limit <= 0 {
		limit = 10
	}
	
	if limit > 100 {
		limit = 100 // Cap at 100 results
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := s.queries.Find(s.context, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list queries: %w", err)
	}
	defer cursor.Close(s.context)

	var queries []*models.Query
	if err = cursor.All(s.context, &queries); err != nil {
		return nil, fmt.Errorf("failed to decode queries: %w", err)
	}

	return queries, nil
}