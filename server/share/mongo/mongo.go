package sharemongo

import "go.mongodb.org/mongo-driver/bson"

// SetOnInsert is option $setOnInsert.
func SetOnInsert(v interface{}) bson.M {
	return bson.M{"$setOnInsert": v}
}
