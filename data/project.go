package data

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	ID         bson.ObjectId   `bson:"_id"`
	Name       string          `bson:"name"`
	OwnerId    bson.ObjectId   `bson:"owner_id"`
	MemberIds  []bson.ObjectId `bson:"member_ids"`
	CreatedAt  time.Time       `bson:"created_at"`
	ModifiedAt time.Time       `bson:"modified_at"`
}

func GetProject(id bson.ObjectId) (*Project, error) {
	pro := Project{}
	err := sess.DB("").C(projectC).FindId(id).One(&pro)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pro, nil
}
