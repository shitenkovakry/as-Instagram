package models

type UserRegistration struct {
	ID    int    `bson:"user_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

type UsersRegistration []*UserRegistration
