package data

import (
	"time"

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
