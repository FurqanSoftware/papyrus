package data

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Organization struct {
	ID         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"name"`
	OwnerID    bson.ObjectId `bson:"owner_id"`
	CreatorID  bson.ObjectId `bson:"creator_id"`
	CreatedAt  time.Time     `bson:"created_at"`
	ModifiedAt time.Time     `bson:"modified_at"`
}

func GetOraganization(id bson.ObjectId) (*Organization, error) {
	org := Organization{}
	err := sess.DB("").C(organizationC).FindId(id).One(&org)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}
