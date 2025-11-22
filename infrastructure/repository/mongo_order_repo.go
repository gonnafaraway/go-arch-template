package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-arch-template/domain/entity"
)

type MongoOrderRepository struct {
	collection *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Database) *MongoOrderRepository {
	return &MongoOrderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *MongoOrderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	var order entity.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &order, err
}

func (r *MongoOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *MongoOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order)
	return err
}
