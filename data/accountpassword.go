package data

type AccountPassword struct {
	Salt       []byte `bson:"salt"`
	Iteration  int    `bson:"iteration"`
	KeyLength  int    `bson:"keyLength"`
	DerivedKey []byte `bson:"derivedKey"`
}
