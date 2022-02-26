package id

// UserID is the user id
type UserID string

func (u UserID) String() string {
	return string(u)
}
