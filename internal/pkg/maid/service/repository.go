package maid_service

import (
	"github.com/samuael/Project/MaidLink/internal/pkg/maid"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

// MaidService maid service instance struct
type MaidService struct {
	Repo maid.IMaidRepository
}

func NewMaidService(repo *maid.IMaidRepository) maid.IMaidService {
	return &MaidService{
		Repo: repo,
	}
}

func CreateMaid(maid *model.Maid) {

}
