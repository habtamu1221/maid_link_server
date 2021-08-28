package clientService

import (
	"context"

	"github.com/samuael/Project/MaidLink/internal/pkg/client"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/internal/pkg/user"
)

type ClientService struct {
	Repo  client.IClientRepository
	URepo user.IUserRepo
}

func NewClientService(repo client.IClientRepository, urepo user.IUserRepo) client.IClientService {
	return &ClientService{
		Repo:  repo,
		URepo: urepo,
	}
}

// CreateClient ...
func (clientser *ClientService) CreateClient(contex context.Context) *model.Client {
	client := contex.Value("client").(*model.Client)
	if client == nil {
		return client
	}
	user := client.User
	contex = context.WithValue(contex, "user", user)
	user, err := clientser.URepo.CreateUser(contex)
	if err == nil && user != nil {
		client.ID = user.ID
		contex = context.WithValue(contex, "client", client)
		if client, err = clientser.Repo.CreateClient(contex); err == nil && client != nil {
			return client
		}
		if err != nil {
			println(err.Error())
		}
	}
	// remove the
	contex = context.WithValue(contex, "user_id", user.ID)
	clientser.URepo.RemoveUser(contex)
	return nil
}

func (clientser *ClientService) GetClient(conte context.Context) *model.Client {
	client, er := clientser.Repo.GetClient(conte)
	if er == nil {
		if client.User, er = clientser.URepo.GetUserByID(conte); client.User != nil {
			return client
		}
	}
	return nil
}
