package sharemongo

import (
	"fmt"
	"gService/share/id"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToObjID convert Stringer to ObjID.
func ToObjID(sid fmt.Stringer) (primitive.ObjectID, error) {
	// 只要sid实现了String方法，即可调用
	return primitive.ObjectIDFromHex(sid.String())
}

// ToUserID convert ObjectID to UserID.
func ToUserID(oid primitive.ObjectID) id.UserID {
	return id.UserID(oid.Hex())
}
