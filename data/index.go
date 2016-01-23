package data

import "gopkg.in/mgo.v2"

var indexes = map[string][]mgo.Index{
	accountC: {
		{
			Key:    []string{"emails.address_norm"},
			Unique: true,
		},
	},
}

func MakeIndexes() error {
	for col, indexes := range indexes {
		for _, index := range indexes {
			err := sess.DB("").C(col).EnsureIndex(index)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
