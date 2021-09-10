package client

import (
	"context"

	"github.com/habte/Project/MaidLink/internal/pkg/model"
)

type IClientService interface {
	// CreateClient .. to create a clinet instance
	// value name 'clinet'
	// Only the client specific information is to be registerd with
	// not the User Information ...
	CreateClient(context.Context) *model.Client
	// GetClient uses "user_id" which is to be passed in the context and returns a client instance.
	GetClient(context.Context) *model.Client
	// UpdateClient a method which uses the context as a parameter and updates
	// the variable "my_maids" holds the real instance.
	UpdateMyMaids(context.Context) []string
	// MyMaidsWhichIPayedFor "user_id" returns []*string
	// NOT FINISHED YET
	MyMaidsWhichIPayedForString(context.Context) []string
}
