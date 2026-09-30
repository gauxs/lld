package code

type User struct {
	userID uint
}

func NewUser(userID uint) *User {
	return &User{
		userID: userID,
	}
}

func (u *User) GetUserID() uint {
	if u == nil {
		// zero is invalid userID
		return 0
	}

	return u.userID
}
