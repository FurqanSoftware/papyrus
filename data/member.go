package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Memeber struct {
	ID        bson.ObjectId "bson:`_id`"
	ProjectID bson.ObjectId "bson:`project_id`"
	AccountID bson.ObjectId "bson:`account_id`"
	InviterID bson.ObjectId "bson:`inviter_id`"
	InvitedAt time.Time     "bson:`invited_at`"
}
