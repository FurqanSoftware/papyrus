package data

import (
	"math"
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
	err = sess.DB("").C(accountC).Find(bson.M{"emails.address_norm": addr}).One(&acc)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (a *Account) NOrganizations() (int, error) {
	n, err := sess.DB("").C(organizationC).Find(bson.M{"owner_id": a.ID}).Count()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (a *Account) Organizations() ([]Organization, error) {
	orgs, err := ListOraganizationsOwner(a.ID, 0, math.MaxInt32)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (a *Account) PrimaryEmail() AccountEmail {
	for _, email := range a.Emails {
		if email.Primary {
			return email
		}
	}
	return AccountEmail{}
}

func (a *Account) Put() error {
	a.ModifiedAt = time.Now()

	if a.ID == "" {
		a.ID = bson.NewObjectId()
		a.CreatedAt = a.ModifiedAt
	}
	_, err := sess.DB("").C(accountC).UpsertId(a.ID, a)
	return err
}
