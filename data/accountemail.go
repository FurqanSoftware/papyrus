package data

import "time"

type AccountEmail struct {
	Address     string `bson:"address"`
	AddressNorm string `bson:"address_norm"`

	Primary bool `bson:"primary"`

	Verified   bool      `bson:"verified"`
	VerifiedAt time.Time `bson:"verified_at"`
	Token      string    `bson:"token"`
}
