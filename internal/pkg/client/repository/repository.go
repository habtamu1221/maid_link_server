package client_repository

import (
	"github.com/samuael/Project/MaidLink/internal/pkg/client"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClientRepository struct representing the client repository of mongo dataase
type ClientRepository struct {
	DB *mongo.Database //
}

// NewClientRepository returns a new client repository class performing on a mongo database
func NewClientRepository(db *mongo.Database) client.IClientRepository {
	return &ClientRepository{
		DB: db,
	}
}
