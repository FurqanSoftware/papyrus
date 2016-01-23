package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID              bson.ObjectId   `bson:"_id"`
	Emails          []AccountEmail  `bson:"emails"`
	Password        AccountPassword `bson:"password"`
	OrganizationIDs []bson.ObjectId `bson:"organization_ids"`

	CreatedAt  time.Time `bson:"created_at"`
	ModifiedAt time.Time `bson:"modified_at"`
}
