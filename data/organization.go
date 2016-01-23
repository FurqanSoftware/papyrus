package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Organization struct {
	Name       string        `bson:"name"`
	OwnerID    bson.ObjectId `bson:"owner_id"`
	CreatorID  bson.ObjectId `bson:"creator_id"`
	CreatedAt  time.Time     `bson:"created_at"`
	ModifiedAt time.Time     `bson:"modified_at"`
}
