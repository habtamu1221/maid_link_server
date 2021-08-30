package client

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IClientRepository interface {
	// CreateClient .. to create a clinet instance
	// value name 'clinet'
	// Only the client specific information is to be registerd with
	// not the User Information ...
	CreateClient(context.Context) (*model.Client, error)
	GetClient(context.Context) (*model.Client, error)
	// UpdateClient a method to update a client instance using only
	// "client" *model.Client as a parameter
	UpdateMyMaids(context.Context) ([]string, error)
}
