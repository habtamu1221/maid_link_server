package maid

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type IMaidService interface {
	// CreateMaid |  <<maid>>
	CreateMaid(context.Context) *model.Maid
	// GetImageUrls user_id
	GetImageUrls(contexts context.Context) []string
	// AddImageUrl | <<image_url>>
	AddImageUrl(context.Context) ([]string, bool)
	// RemoveProfileImage .. | <<image_url>> ...
	RemoveProfileImage(context.Context) bool
	// GetUser user_id
	GetUser(context.Context) *model.Maid
	// GetWorks ...
	GetWorks(context.Context) []*model.Work
	// CreateWork
	CreateWork(context.Context) *model.Work
	// DeleteWork  DELETE
	DeleteWork(context.Context) bool
	// UpdateWork  user_id and work value
	UpdateWork(conte context.Context) *model.Work
	// GetMyMaids user_id user is ADMIN
	GetMyMaids(conte context.Context) []*model.Maid
	// GetMaid returns a maid instance using 'maid_id' in the context
	GetMaid(conte context.Context) *model.Maid
	// UpdateMaid ...this service context variables 'user_id' and 'maid'
	// then returns an updated maid instance
	// Remember the instance maid is *model.Maid and also the return value of this
	// method too...
	UpdateMaid(conte context.Context) *model.Maid
	// GetMaids takes "offset"  "limit" to return a list of Maids
	GetMaids(conte context.Context) []*model.Maid
	// MyMaidsWhichIPayedFor "user_id" returns []*string
	// NOT FINISHED YET
	MyMaidsWhichIPayedFor(context.Context) []string
}
