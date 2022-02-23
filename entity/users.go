package entity

/** UserID **/
type UserID interface {
	ID
}

/** User **/
type User interface {
	UserID() UserID
	Name() string
}

func NewUser(userID UserID, name string) User {
	return &user{
		userID: userID,
		name:   name,
	}
}

type user struct {
	userID UserID
	name   string
}

func (u *user) UserID() UserID {
	return u.userID
}
func (u *user) Name() string {
	return u.name
}

/** Repository **/
type UsersReader interface {
	FindAll() ([]User, error)
	FindByID(id UserID) (User, error)
}
