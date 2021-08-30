package maid

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

// IMaidRepository represents the repository object
type IMaidRepository interface {
	CreateMaid(context.Context) (*model.Maid, error)
	// GetImageUrls user_id
	GetImageUrls(contexts context.Context) ([]string, error)
	// AddImageUrl add image to the maid list of profile images
	AddImageUrl(context.Context) ([]string, error)
	// RemoveProfileImage ...
	RemoveProfileImage(context.Context) error
	// GetUser user_id
	GetUser(context.Context) (*model.Maid, error)
	// GetWorks ...
	GetWorks(context.Context) ([]*model.Work, error)
	// CreateWork
	CreateWork(context.Context) (*model.Work, error)
	// DeleteWork  DELETE
	DeleteWork(context.Context) error
	// UpdateWork user_id , work
	UpdateWork(conte context.Context) (*model.Work, error)
	// GetMyMaids  user_id user is ADMIN
	GetMyMaids(conte context.Context) ([]*model.Maid, error)
	// GetMaid returns a maid instance using 'user_id' in the context
	GetMaid(conte context.Context) (*model.Maid, error)
	// UpdateMaid ...
	UpdateMaid(context.Context) (*model.Maid, error)
	// GetMaids takes "offset"  "limit" to return a list of Maids
	GetMaids(conte context.Context) ([]*model.Maid, error)
	// MyMaidsWhichIPayedFor "user_id" returns []*string
	MyMaidsWhichIPayedFor(context.Context) ([]string, error)
}
