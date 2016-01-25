package data

import "gopkg.in/mgo.v2"

var sess *mgo.Session

const (
	accountC      = "accounts"
	organizationC = "organizations"
	projectC      = "projects"
	memberC       = "members"
	documentC     = "documents"
)

func OpenDBSession(url string) (err error) {
	sess, err = mgo.Dial(url)
	return err
}
