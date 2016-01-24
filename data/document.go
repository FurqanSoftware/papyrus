package data

import (
	"crypto/rand"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Document struct {
	ID          bson.ObjectId `bson:"_id"`
	ProjectID   bson.ObjectId `bson:"project_id"`
	ShortID     string        `bson:"short_id,omitempty"`
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

func GetDocumentShortID(shortID string) (*Document, error) {
	doc := Document{}
	err := sess.DB("").C(documentC).Find(bson.M{"short_id": shortID}).One(&doc)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

const shortIDAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := []byte{}
	for _, v := range b {
		s = append(s, shortIDAlpha[v%62])
	}
	return string(s), nil
}

func ListDocumentsProject(projectID bson.ObjectId, skip, limit int) ([]Document, error) {
	docs := []Document{}
	err := sess.DB("").C(documentC).
		Find(bson.M{"project_id": projectID}).
		Skip(skip).
		Limit(limit).
		Sort("-created_at").
		All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func (d *Document) Project() (*Project, error) {
	return GetProject(d.ProjectID)
}

func (d *Document) Put() error {
	d.ModifiedAt = time.Now()

	if d.ID == "" {
		d.ID = bson.NewObjectId()
		d.CreatedAt = d.ModifiedAt
	}
	_, err := sess.DB("").C(documentC).UpsertId(d.ID, d)
	return err
}
