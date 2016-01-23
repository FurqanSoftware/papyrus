package data

import (
	"time"

	"github.com/asaskevich/govalidator"

	"gopkg.in/mgo.v2"
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

func GetAccount(id bson.ObjectId) (*Account, error) {
	acc := Account{}

	err := sess.DB("").C(accountC).FindId(id).One(&acc)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func GetAccountEmail(addr string) (*Account, error) {
	addr, err := govalidator.NormalizeEmail(addr)
	if err != nil {
		return nil, nil
	}
	acc := Account{}
	err = sess.DB("").C(accountC).Find(bson.M{"emails": addr}).One(&acc)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
