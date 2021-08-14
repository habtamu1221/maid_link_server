package client_service

import (
	"github.com/samuael/Project/MaidLink/internal/pkg/client"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
)

type ClientService struct {
	Repo client.IClientRepository
}

func NewClientService(repo client.IClientRepository) client.IClientService {
	return &ClientService{
		Repo: repo,
	}
}

func CreateClient(client *model.Client) *model.Client {

	return nil
}
