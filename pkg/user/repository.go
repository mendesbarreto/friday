package user

type Repository interface {
    CreateUser(user *User) (*User, error)
    FindUsersById(id string) (*User, error)
    DeleteUserById(id string) (error)
}
