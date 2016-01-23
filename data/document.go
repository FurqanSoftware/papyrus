package data

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Document struct {
	ID          bson.ObjectId `bson:"_id"`
	ShortId     string        `bson:"short_id"`
	Title       string        `bson:"title"`
	Content     string        `bson:"content"`
	Tags        []string      `bson:"tags"`
	Published   bool          `bson:"publishd"`
	PublishedAt time.Time     `bson:"pushlished_at"`
	AccessToken string        `bson:"access_token"`
	CreatedAt   time.Time     `bson:"created_at"`
	ModifiedAt  time.Time     `bson:"modified_at"`
}

func GetDocument(id bson.ObjectId) (*Document, error) {
	doc := Document{}
	err := sess.DB("").C(documentC).FindId(id).One(&doc)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
