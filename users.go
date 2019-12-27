package users

import (
	"time"

	pb "github.com/microapis/lib/proto"
)

// User ...
type User struct {
	ID string `json:"id" db:"id"`

	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// Service ...
type Service interface {
	Getter
	Create(*User) error
	UserList() ([]*User, error)

	// TODO(ca): below methods are not implemented
	Update(*User) error
	Delete(*User) error
}

// Getter ...
type Getter interface {
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}

// Query ...
type Query struct {
	ID    string
	Email string
}

// ToProto ...
func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt.UnixNano(),
		UpdatedAt: u.UpdatedAt.UnixNano(),
	}
}

// FromProto ...
func (u *User) FromProto(uu *pb.User) *User {
	u.ID = uu.Id
	u.Email = uu.Email
	u.Password = uu.Password
	u.CreatedAt = time.Unix(uu.CreatedAt, 0)
	u.UpdatedAt = time.Unix(uu.UpdatedAt, 0)

	return u
}
