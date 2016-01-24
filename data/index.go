package data

import "gopkg.in/mgo.v2"

var indexes = map[string][]mgo.Index{
	accountC: {
		{
			Key:    []string{"emails.address_norm"},
			Unique: true,
		},
	},
	documentC: {
		{
			Key:    []string{"short_id"},
			Unique: true,
		},
		{
			Key: []string{"project_id"},
		},
	},
	memberC: {
		{
			Key:    []string{"project_id", "account_id"},
			Unique: true,
		},
		{
			Key: []string{"account_id", "organization_id"},
		},
	},
	organizationC: {
		{
			Key: []string{"owner_id"},
		},
	},
	projectC: {
		{
			Key: []string{"organization_id"},
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
