package maid_service

import (
	"context"

	"github.com/habte/Project/MaidLink/internal/pkg/maid"
	"github.com/habte/Project/MaidLink/internal/pkg/model"
	"github.com/habte/Project/MaidLink/internal/pkg/user"
)

// MaidService maid service instance struct
type MaidService struct {
	Repo  maid.IMaidRepository
	URepo user.IUserRepo
}

func NewMaidService(repo maid.IMaidRepository, urepo user.IUserRepo) maid.IMaidService {
	return &MaidService{
		Repo:  repo,
		URepo: urepo,
	}
}

func (maidser *MaidService) CreateMaid(cotx context.Context) *model.Maid {
	maid := cotx.Value("maid").(*model.Maid)
	user := maid.User
	cotx = context.WithValue(cotx, "user", user)
	if user, er := maidser.URepo.CreateUser(cotx); user != nil && er == nil && user.ID != "" {
		maid.ID = user.ID //
		maid.User = user  //
		cotx = context.WithValue(cotx, "maid", maid)
		if maid, er := maidser.Repo.CreateMaid(cotx); er == nil && maid != nil {
			return maid
		} else {
			cotx = context.WithValue(cotx, "user_id", user.ID)
			maidser.URepo.RemoveUser(cotx)
		}
	}
	return nil
}

func (maidser *MaidService) AddImageUrl(contexts context.Context) ([]string, bool) {
	vals, err := maidser.Repo.AddImageUrl(contexts)
	return vals, (err == nil)
}

func (maidser *MaidService) RemoveProfileImage(contexts context.Context) bool {
	return maidser.Repo.RemoveProfileImage(contexts) == nil
}

func (maidser *MaidService) GetImageUrls(contexts context.Context) []string {
	val, er := maidser.Repo.GetImageUrls(contexts)
	if er == nil {
		return val
	}
	return nil
}

func (maidser *MaidService) GetUser(conte context.Context) *model.Maid {
	if maid, err := maidser.Repo.GetUser(conte); err == nil {
		if user, err := maidser.URepo.GetUserByID(conte); err == nil {
			user.Password = ""
			maid.User = user //
			return maid
		}
	}
	return nil
}
func (maidser *MaidService) GetWorks(contex context.Context) []*model.Work {
	if works, err := maidser.Repo.GetWorks(contex); err == nil {
		return works
	}
	return nil
}

func (maidser *MaidService) CreateWork(contex context.Context) *model.Work {
	work, er := maidser.Repo.CreateWork(contex) //
	if er == nil {
		return work
	}
	return nil
}
func (maidser *MaidService) DeleteWork(conte context.Context) bool {
	return maidser.Repo.DeleteWork(conte) == nil
}
func (maidser *MaidService) UpdateWork(conte context.Context) *model.Work {
	if work, err := maidser.Repo.UpdateWork(conte); err == nil && work != nil {
		return work
	}
	return nil
}

// GetMyMaids
func (maidser *MaidService) GetMyMaids(conte context.Context) []*model.Maid {
	maidso := []*model.Maid{}
	if maids, err := maidser.Repo.GetMyMaids(conte); err == nil {
		for _, md := range maids {
			conte = context.WithValue(conte, "user_id", md.ID)
			if user, er := maidser.URepo.GetUserByID(conte); er == nil {
				user.Password = ""
				md.User = user
				maidso = append(maidso, md)
			}
		}
		return maidso
	}
	return nil
}

// GetMaid ...
func (maidser *MaidService) GetMaid(conte context.Context) *model.Maid {
	if maid, er := maidser.Repo.GetMaid(conte); er == nil {
		conte = context.WithValue(conte, "user_id", maid.ID)
		if maid.User, er = maidser.URepo.GetUserByID(conte); er == nil && maid.User != nil {
			return maid
		}
	}
	return nil
}

func (maidser *MaidService) UpdateMaid(contr context.Context) *model.Maid {
	if maid, er := maidser.Repo.UpdateMaid(contr); er == nil {
		return maid
	}
	return nil
}

func (maidser *MaidService) GetMaids(conte context.Context) []*model.Maid {
	newMaids := []*model.Maid{}
	if maids, er := maidser.Repo.GetMaids(conte); er == nil {
		for _, maid := range maids {
			conte = context.WithValue(conte, "user_id", maid.ID)
			if maid.User, er = maidser.URepo.GetUserByID(conte); er == nil && maid.User != nil {
				newMaids = append(newMaids, maid)
			}
		}
		return newMaids
	}
	return nil
}

func (maidser *MaidService) SearchMaids(conte context.Context) []*model.Maid {
	if maids, er := maidser.Repo.SearchMaids(conte); er == nil && maids != nil {
		return maids
	} else {
		return nil
	}
}

// SearchIt uses "q"
func (maidser *MaidService) SearchIt(conte context.Context) []*interface{} {
	if maids, er := maidser.Repo.SearchIt(conte); er == nil {
		return maids
	}
	return nil
}
