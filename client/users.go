package client

import (
	"context"
	"errors"

	users "github.com/microapis/users-api"
	pb "github.com/microapis/users-api/proto"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	Client pb.UserServiceClient
}

// New ...
func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := pb.NewUserServiceClient(conn)

	return &Client{
		Client: c,
	}, nil
}

// Get ...
func (c *Client) Get(userID string) (*users.User, error) {
	if userID == "" {
		return nil, errors.New("invalid userID")
	}

	gr, err := c.Client.Get(context.Background(), &pb.UserGetRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return nil, errors.New(msg)
	}

	u := &users.User{}
	user := u.FromProto(gr.GetData())
	return user, nil
}

// GetByEmail ...
func (c *Client) GetByEmail(email string) (*users.User, error) {
	if email == "" {
		return nil, errors.New("invalid email")
	}

	gr, err := c.Client.GetByEmail(context.Background(), &pb.UserGetByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return nil, errors.New(msg)
	}

	u := new(users.User)
	user := u.FromProto(gr.GetData())
	return user, nil
}

// Create ...
func (c *Client) Create(user *users.User) (*users.User, error) {
	if user.Name == "" {
		return nil, errors.New("invalid name")
	}

	if user.Email == "" {
		return nil, errors.New("invalid email")
	}

	if user.Password == "" {
		return nil, errors.New("invalid password")
	}

	gr, err := c.Client.Create(context.Background(), &pb.UserCreateRequest{
		Data: user.ToProto(),
	})
	if err != nil {
		return nil, err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return nil, errors.New(msg)
	}

	u := new(users.User)
	uu := u.FromProto(gr.GetData())
	return uu, nil
}

// VerifyPassword ...
func (c *Client) VerifyPassword(email, password string) error {
	if email == "" {
		return errors.New("invalid email")
	}

	if password == "" {
		return errors.New("invalid password")
	}

	gr, err := c.Client.VerifyPassword(context.Background(), &pb.UserVerifyPasswordRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return errors.New(msg)
	}

	return nil
}

// List ...
func (c *Client) List() ([]*users.User, error) {
	gr, err := c.Client.List(context.Background(), &pb.UserListRequest{})
	if err != nil {
		return nil, err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return nil, errors.New(msg)
	}

	protoUsers := gr.GetData()

	data := make([]*users.User, 0)
	for _, user := range protoUsers {
		u := new(users.User)
		data = append(data, u.FromProto(user))
	}

	return data, nil
}

// Update ...
func (c *Client) Update(ID string, u *users.User) (*users.User, error) {
	if ID == "" {
		return nil, errors.New("invalid ID")
	}

	gr, err := c.Client.Update(context.Background(), &pb.UserUpdateRequest{
		UserId: ID,
		Data:   u.ToProto(),
	})
	if err != nil {
		return nil, err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return nil, errors.New(msg)
	}

	uu := &users.User{}
	user := uu.FromProto(gr.GetData())
	return user, nil
}

// Delete ...
func (c *Client) Delete(ID string) error {
	if ID == "" {
		return errors.New("invalid ID")
	}

	gr, err := c.Client.Delete(context.Background(), &pb.UserDeleteRequest{
		UserId: ID,
	})
	if err != nil {
		return err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return errors.New(msg)
	}

	return nil
}
