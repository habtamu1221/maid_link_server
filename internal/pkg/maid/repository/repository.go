package maid_repo

import (
	"context"
	"errors"

	"github.com/samuael/Project/MaidLink/internal/pkg/maid"
	"github.com/samuael/Project/MaidLink/internal/pkg/model"
	"github.com/samuael/Project/MaidLink/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (maidrep *MaidRepo) CreateMaid(ctx context.Context) (*model.Maid, error) {
	maid := ctx.Value("maid").(*model.Maid)
	if oid, err := primitive.ObjectIDFromHex(maid.ID); err == nil {
		maid.BsonID = oid
		id := pkg.ObjectIDFromString(pkg.RemoveObjectIDPrefix(maid.ID))
		maid.ID = id
		insertResult, er := maidrep.DB.Collection(model.SMAID).InsertOne(ctx, maid)
		if insertResult == nil {
			return nil, errors.New("Unsuccesful")
		}
		if val := pkg.ObjectIDFromInsertResult(insertResult); val != "" && maid.ID == pkg.ObjectIDFromString(val) && er == nil {
			maid.ID = pkg.RemoveObjectIDPrefix(maid.ID)
			return maid, nil
		} else {
			return nil, er
		}
	} else {
		return nil, err
	}
}

func (userrepo *MaidRepo) GetImageUrls(contexts context.Context) ([]string, error) {
	userID := contexts.Value("user_id").(string)
	maid := &model.Maid{}
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	err = userrepo.DB.Collection(model.SMAID).FindOne(contexts, bson.D{{"_id", oid}}).Decode(maid)
	if err == nil {
		if maid.ProfileImages == nil {
			return []string{}, nil
		}
		return maid.ProfileImages, nil
	}
	return nil, err //
}

func (userrepo *MaidRepo) AddImageUrl(contexts context.Context) ([]string, error) {
	userID := contexts.Value("user_id").(string)
	imageUrl := contexts.Value("image_url").(string)
	if oid, er := primitive.ObjectIDFromHex(userID); er == nil {
		if imgurls, er := userrepo.GetImageUrls(contexts); er == nil {
			imgurls = append(imgurls, imageUrl)
			if updates, er := userrepo.DB.Collection(model.SMAID).UpdateOne(contexts, bson.D{
				{"_id", oid},
			}, bson.D{{"$set", bson.D{
				{"profileimages", imgurls},
			}}}); updates.ModifiedCount > 0 && er == nil {
				return imgurls, nil
			} else {
				return nil, er
			}
		} else {
			return nil, er
		}
	} else {
		return nil, er //
	}
}

// RemoveProfileImage ...
func (maidRepo *MaidRepo) RemoveProfileImage(contexts context.Context) error {
	imageUrl := contexts.Value("image_url").(string)
	id := contexts.Value("session").(*model.Session).UserID
	contexts = context.WithValue(contexts, "user_id", id)
	images, er := maidRepo.GetImageUrls(contexts)
	if er != nil {
		return er
	}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		for a, im := range images {
			if imageUrl == im {
				if a == 0 {
					if len(images) == 1 {
						images = []string{}
					} else {
						images = images[1:]
					}
				} else if a == len(images)-1 {
					images = images[0:a]
				} else {
					images = append(images[0:a], images[a+1:]...)
				}
			}
		}
		if uc, err := maidRepo.DB.Collection(model.SMAID).UpdateOne(contexts, bson.D{{"_id", oid}}, bson.D{{"$set", bson.D{{"profileimages", images}}}}); uc.ModifiedCount > 0 && err == nil {
			return nil
		} else {
			return errors.New("Updated Succesfuly")
		}
	} else {
		return err
	}
}

func (maidRepo *MaidRepo) GetUser(conte context.Context) (*model.Maid, error) {
	userID := conte.Value("user_id").(string)
	maid := &model.Maid{}
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	err = maidRepo.DB.Collection(model.SMAID).FindOne(conte, bson.D{{"_id", oid}}).Decode(maid)
	if err == nil {
		if maid.ProfileImages == nil {
			return nil, nil
		}
		return maid, nil
	}
	return nil, err //
}

func (maidRepo *MaidRepo) CreateWork(conte context.Context) (*model.Work, error) {
	work := conte.Value("work").(*model.Work)
	maidID := conte.Value("user_id").(string)
	var err error
	var oid primitive.ObjectID
	var works []*model.Work
	if works, err = maidRepo.GetWorks(conte); err == nil {
		work.NO = uint(len(works)) + 1
		if oid, err = primitive.ObjectIDFromHex(maidID); err == nil {
			works = append(works, work)
			uc, err := maidRepo.DB.Collection(model.SMAID).UpdateOne(conte, bson.D{{"_id", oid}}, bson.D{{"$set", bson.D{{"works", works}}}})
			if err != nil {
				return nil, err
			} else if uc.MatchedCount == 0 {
				return nil, errors.New("Matched Count")
			}
			return work, nil
		}
	}
	return nil, err
}

func (maidRepo *MaidRepo) GetWorks(conte context.Context) ([]*model.Work, error) {
	maidID := conte.Value("user_id").(string)
	maid := &model.Maid{}
	var err error
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(maidID); err == nil {
		if err = maidRepo.DB.Collection(model.SMAID).FindOne(conte, bson.D{{"_id", oid}}).Decode(maid); err == nil {
			return maid.Works, nil //
		}
		return nil, err
	}
	return nil, err
}

// DeleteWork
func (maidRepo *MaidRepo) DeleteWork(conte context.Context) error {
	userID := conte.Value("user_id").(string)
	var err error
	var works []*model.Work
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		workNo := conte.Value("work_no").(int)
		if works, err = maidRepo.GetWorks(conte); err == nil {
			for i, wk := range works {
				if wk.NO == uint(workNo) {
					println(string(pkg.GetJson(wk)))
					if i == 0 {
						if len(works) == 1 {
							works = []*model.Work{}
						} else {
							works = works[1:]
						}
					} else {
						if len(works)-1 == i {
							works = works[:i]
						} else {
							works = append(works[:i], works[i+1:]...)
						}
					}
				}
			}
			if uc, er := maidRepo.DB.Collection(model.SMAID).UpdateOne(conte, bson.D{{"_id", oid}}, bson.D{{"$set", bson.D{{"works", works}}}}); er == nil && uc.ModifiedCount > 0 {
				return nil
			} else {
				return errors.New("Update Was Not Succesful ")
			}
		} else {
			return err
		}
	}
	return err
}

func (maidRepo *MaidRepo) UpdateWork(conte context.Context) (*model.Work, error) {
	work := conte.Value("work").(*model.Work)
	var works []*model.Work
	var er error
	if oid, er := primitive.ObjectIDFromHex(conte.Value("user_id").(string)); er == nil {
		works, er = maidRepo.GetWorks(conte)
		// search for the work and Replace it with the new One
		for a, wk := range works {
			if wk.NO == work.NO {
				works[a] = work
			}
		}
		if uc, er := maidRepo.DB.Collection(model.SMAID).UpdateOne(conte, bson.D{{"_id", oid}}, bson.D{{"$set",
			bson.D{{"works", works}}}}); er == nil && uc.ModifiedCount > 0 {
			return work, nil
		}
		return nil, errors.New("not modified")
	}
	return nil, er
}

func (maidRepo *MaidRepo) GetMyMaids(conte context.Context) ([]*model.Maid, error) {
	maids := []*model.Maid{}
	if cursor, er := maidRepo.DB.Collection(model.SMAID).Find(conte, bson.D{{"createdby", conte.Value("user_id").(string)}}); er == nil {
		for cursor.Next(conte) {
			maidin := &model.Maid{}
			err := cursor.Decode(maidin)
			if err == nil {
				maidin.ID = pkg.RemoveObjectIDPrefix(maidin.BsonID.String())
				maids = append(maids, maidin)
			}
		}
		return maids, nil
	} else {
		println(er.Error())
		return maids, er
	}
}

func (maidRepo *MaidRepo) GetMaid(conte context.Context) (*model.Maid, error) {
	userID := conte.Value("maid_id").(string)
	println(userID)
	maid := &model.Maid{}
	if oid, er := primitive.ObjectIDFromHex(userID); er == nil {
		if er = maidRepo.DB.Collection(model.SMAID).FindOne(conte, bson.D{{"_id", oid}}).Decode(maid); er == nil {
			maid.ID = pkg.RemoveObjectIDPrefix(maid.BsonID.String())
			println(string(pkg.GetJson(maid)))
		} else {
			println(er.Error())
		}
		return maid, er
	} else {
		println(er)
		return nil, er
	}
}

func (maidRepo *MaidRepo) UpdateMaid(conte context.Context) (*model.Maid, error) {
	if oid, er := primitive.ObjectIDFromHex(conte.Value("maid_id").(string)); er == nil {
		maid := conte.Value("maid").(*model.Maid)
		if uc, er := maidRepo.DB.Collection(model.SMAID).UpdateOne(conte, bson.D{{"_id", oid}}, bson.D{
			{"$set", bson.D{
				{"phone", maid.Phone},
				{"address", maid.Address},
				{"rated_by", maid.RatedBy},
				{"bio", maid.Bio},
				{"rates", maid.Rates},
				{"ratecount", maid.RateCount},
				{"carrers", maid.Carrers},
				{"works", maid.Works},
			}}}); er == nil && uc.ModifiedCount > 0 {
			return maid, er
		} else {
			return nil, er
		}
	} else {
		return nil, er
	}
}

// GetMaids
func (maidRepo *MaidRepo) GetMaids(conte context.Context) ([]*model.Maid, error) {
	offset := conte.Value("offset").(int)
	limit := conte.Value("limit").(int)
	create := func(x int64) *int64 {
		return &x
	}
	insta := struct {
		OFFSET *int64
		LIMIT  *int64
	}{
		OFFSET: create(int64(offset)),
		LIMIT:  create(int64(limit)),
	}
	maids := []*model.Maid{}
	if cursor, err := maidRepo.DB.Collection(model.SMAID).Find(conte, bson.D{}, &options.FindOptions{Limit: insta.LIMIT}, &options.FindOptions{Skip: insta.OFFSET}); err == nil {
		for cursor.Next(conte) {
			maid := &model.Maid{}
			er := cursor.Decode(maid)
			if er == nil {
				maids = append(maids, maid)
			}
		}
		return maids, nil
	} else {
		println(err.Error())
		return maids, err
	}
}

// MyMaidsWhichIPayedFor "user_id" returns []*string
func (maidRepo *MaidRepo) MyMaidsWhichIPayedFor(conte context.Context) ([]string, error) {
	maid := &model.Maid{}
	if oid, er := primitive.ObjectIDFromHex(conte.Value("user_id").(string)); er == nil {
		if er = maidRepo.DB.Collection(model.SMAID).FindOne(conte, bson.D{{"_id", oid}}).Decode(maid); er == nil {
			// maid.
		} //, &options.FindOneOptions{Projection: })
	}
	return nil, nil
}
