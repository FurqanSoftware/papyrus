package data

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	ID             bson.ObjectId   `bson:"_id"`
	Name           string          `bson:"name"`
	OwnerID        bson.ObjectId   `bson:"owner_id"`
	OrganizationID bson.ObjectId   `bson:"organization_id"`
	MemberIDs      []bson.ObjectId `bson:"member_ids"`
	CreatedAt      time.Time       `bson:"created_at"`
	ModifiedAt     time.Time       `bson:"modified_at"`
}

func ListProjectsOrganization(orgID bson.ObjectId, skip, limit int) ([]Project, error) {
	prjs := []Project{}

	err := sess.DB("").C(projectC).
		Find(bson.M{"organization_id": orgID}).
		Skip(skip).
		Limit(limit).
		Sort("-created_at").
		All(&prjs)
	if err != nil {
		return nil, err
	}

	return prjs, nil
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

func (p *Project) Members() ([]Member, error) {
	mems := []Member{}
	err := sess.DB("").C(memberC).
		Find(bson.M{"project_id": p.ID}).
		All(&mems)

	if err != nil {
		return nil, err
	}
	return mems, nil
}

func (p *Project) Organization() (*Organization, error) {
	return GetOraganization(p.OrganizationID)
}

func (p *Project) Owner() (*Account, error) {
	return GetAccount(p.OwnerID)
}

func (p *Project) Put() error {
	p.ModifiedAt = time.Now()

	if p.ID == "" {
		p.ID = bson.NewObjectId()
		p.CreatedAt = p.ModifiedAt
	}
	_, err := sess.DB("").C(projectC).UpsertId(p.ID, p)
	return err
}
