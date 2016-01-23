package data

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Memeber struct {
	ID        bson.ObjectId `bson:"_id"`
	ProjectID bson.ObjectId `bson:"project_id"`
	AccountID bson.ObjectId `bson:"account_id"`
	InviterID bson.ObjectId `bson:"inviter_id"`
	InvitedAt time.Time     `bson:"invited_at"`

	ModifiedAt time.Time `bson:"modified_at"`
	CreatedAt  time.Time `bson:"created_at"`
}

func GetMember(id bson.ObjectId) (*Memeber, error) {
	mem := Memeber{}
	err := sess.DB("").C(memberC).FindId(id).One(&mem)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &mem, nil
}

func ListMembersProject(projectID bson.ObjectId, skip, limit int) ([]Memeber, error) {
	mems := []Memeber{}
	err := sess.DB("").C(memberC).
		Find(bson.M{"project_id": projectID}).
		Skip(skip).
		Limit(limit).
		Sort("-created_at").
		All(&mems)
	if err != nil {
		return nil, err
	}
	return mems, nil
}

func (m *Memeber) Account() (*Account, error) {
	return GetAccount(m.AccountID)
}

func (m *Memeber) Inviter() (*Account, error) {
	return GetAccount(m.InviterID)
}

func (m *Memeber) Put() error {
	m.ModifiedAt = time.Now()

	if m.ID == "" {
		m.ID = bson.NewObjectId()
		m.CreatedAt = m.ModifiedAt
	}
	_, err := sess.DB("").C(memberC).UpsertId(m.ID, m)
	return err
}
