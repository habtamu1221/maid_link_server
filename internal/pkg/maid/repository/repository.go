package maid_repo

import (
	"github.com/samuael/Project/MaidLink/internal/pkg/maid"
	"go.mongodb.org/mongo-driver/mongo"
)

type MaidRepo struct {
	DB *mongo.Database
}

// NewMaidRepo is a function returning the link to the maid repository.
func NewMadiRepo(db *mongo.Database) maid.IMaidRepository {
	return &MaidRepo{
		DB: db,
	}
}
