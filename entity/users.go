package entity

/** UserID **/
type UserID interface {
	ID
}

/** User **/
type User interface {
	UserID() UserID
	Name() string
	Email() Email
}

func NewUser(userID UserID, name string, email Email) User {
	return &user{
		userID: userID,
		name:   name,
		email:  email,
	}
}

type user struct {
	userID UserID
	name   string
	email  Email
}

func (u *user) UserID() UserID {
	return u.userID
}
func (u *user) Name() string {
	return u.name
}
func (u *user) Email() Email {
	return u.email
}

/** Repository **/
type UsersReader interface {
	FindAll() ([]User, error)
	FindByID(id UserID) (User, error)
}

type UsersWriter interface {
	Insert(user User) (User, error)
	Delete(id UserID) error
}

type UsersReadWriter interface {
	UsersReader
	UsersWriter
}
