package company

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-arch-template/internal/api/domain/company"
	mongodb "go-arch-template/internal/api/storage/mongodb"
)

type MongoDBRepository struct {
	collection *mongo.Collection
}

func NewMongoDBRepository(db *mongodb.Client) *MongoDBRepository {
	return &MongoDBRepository{
		collection: db.Collection("companies"),
	}
}

type companyDocument struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (r *MongoDBRepository) Create(ctx context.Context, c *company.Company) error {
	if c.ID == "" {
		c.ID = primitive.NewObjectID().Hex()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	doc := companyDocument{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *MongoDBRepository) FindByID(ctx context.Context, id string) (*company.Company, error) {
	var doc companyDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &company.Company{
		ID:        doc.ID,
		Name:      doc.Name,
		Email:     doc.Email,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}, nil
}

func (r *MongoDBRepository) FindAll(ctx context.Context) ([]*company.Company, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []*company.Company
	for cursor.Next(ctx) {
		var doc companyDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		companies = append(companies, &company.Company{
			ID:        doc.ID,
			Name:      doc.Name,
			Email:     doc.Email,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		})
	}

	return companies, cursor.Err()
}

func (r *MongoDBRepository) Update(ctx context.Context, c *company.Company) error {
	c.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"name":       c.Name,
			"email":      c.Email,
			"updated_at": c.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": c.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *MongoDBRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

