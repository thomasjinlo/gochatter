package models

type User struct {
	username string
}

func NewUser() *User {
	return &User{}
}

func (u *User) Username() string {
	return u.username
}

type localRepository struct {
	users map[string]*User
}

func NewLocalUserRepository() *localRepository {
	return &localRepository{
		users: make(map[string]*User),
	}
}

func (r *localRepository) Login(username string) error {
	_, ok := r.users[username]
	if !ok {
		r.users[username] = &User{username: username}
	}
	return nil
}

func (r *localRepository) GetAll() []*User {
	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}
