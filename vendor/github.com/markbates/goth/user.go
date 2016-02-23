package goth

import "encoding/gob"

func init() {
	gob.Register(User{})
}

// User contains the information common amongst most OAuth and OAuth2 providers.
// All of the "raw" datafrom the provider can be found in the `RawData` field.
type User struct {
	RawData           map[string]interface{}
	Provider          string
	Email             string
	Name              string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
}
