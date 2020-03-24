package userssvc

import (
	"context"
	"fmt"

	users "github.com/microapis/users-api"
	"github.com/microapis/users-api/database"
	pb "github.com/microapis/users-api/proto"
	"github.com/microapis/users-api/service"
	"golang.org/x/crypto/bcrypt"
)

var _ pb.UserServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	usersSvc users.Service
}

// New ...
func New(store database.Store) *Service {
	return &Service{
		usersSvc: service.NewUsers(store),
	}
}

// Get Gets a user by ID.
func (us *Service) Get(ctx context.Context, gr *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	id := gr.GetUserId()
	fmt.Println(fmt.Sprintf("[gRPC][Users][Get][Request] id = %v", id))
	user, err := us.usersSvc.GetByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Get][Error] %v", err))
		return &pb.UserGetResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.UserGetResponse{
		Data:  user.ToProto(),
		Meta:  nil,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][Get][Response] %v", res))
	return res, nil
}

// GetByEmail get a user by Email
func (us *Service) GetByEmail(ctx context.Context, gr *pb.UserGetByEmailRequest) (*pb.UserGetByEmailResponse, error) {
	email := gr.GetEmail()
	fmt.Println(fmt.Sprintf("[gRPC][Users][GetByEmail][Request] email = %v", email))

	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][Users][GetByEmail][Error] %v", "email user params empty"))
		return &pb.UserGetByEmailResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    400,
				Message: "email user params empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][GetByEmail][Error] %v", "user not found"))
		return &pb.UserGetByEmailResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    404,
				Message: "user not found",
			},
		}, nil
	}

	res := &pb.UserGetByEmailResponse{
		Data:  user.ToProto(),
		Meta:  nil,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][GetByEmail][Response] %v", res))
	return res, nil
}

// Create creates a new user into database.
func (us *Service) Create(ctx context.Context, gr *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Request] data = %v", gr.GetData()))

	email := gr.GetData().GetEmail()
	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Error] %v", "email user param is empty"))
		return &pb.UserCreateResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    400,
				Message: "email user param is empty",
			},
		}, nil
	}

	_, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		name := gr.GetData().GetName()
		if name == "" {
			fmt.Println(fmt.Sprintf("[gRPC][TenpoUsers][Create][Error] %v", "name user param is empty"))
			return &pb.UserCreateResponse{
				Data: nil,
				Meta: nil,
				Error: &pb.UserError{
					Code:    400,
					Message: "name user param is empty",
				},
			}, nil
		}

		password := gr.GetData().GetPassword()
		if password == "" {
			fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Error] %v", "password user params is empty"))
			return &pb.UserCreateResponse{
				Data: nil,
				Meta: nil,
				Error: &pb.UserError{
					Code:    400,
					Message: "password user params is empty",
				},
			}, nil
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Error] %v", err))
			return &pb.UserCreateResponse{
				Data: nil,
				Meta: nil,
				Error: &pb.UserError{
					Code:    500,
					Message: "could not generate hashed password",
				},
			}, nil
		}

		user := &users.User{
			Email:    email,
			Name:     name,
			Password: string(hashedPassword),
		}

		if err := us.usersSvc.Create(user); err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Error] %v", err))
			return &pb.UserCreateResponse{
				Data: nil,
				Meta: nil,
				Error: &pb.UserError{
					Code:    500,
					Message: err.Error(),
				},
			}, nil
		}

		res := &pb.UserCreateResponse{
			Data:  user.ToProto(),
			Meta:  nil,
			Error: nil,
		}
		fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Response] %v", res))
		return res, nil
	}

	res := &pb.UserCreateResponse{
		Data: nil,
		Meta: nil,
		Error: &pb.UserError{
			Code:    400,
			Message: "user already exists",
		},
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][Create][Response] %v", res))
	return res, nil
}

// VerifyPassword ...
func (us *Service) VerifyPassword(ctx context.Context, gr *pb.UserVerifyPasswordRequest) (*pb.UserVerifyPasswordResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][Users][VerifyPassword][Request] email = %v password = %v", gr.GetEmail(), gr.GetPassword()))
	email := gr.GetEmail()
	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][Users][VerifyPassword][Error] %v", "email user param is empty"))
		return &pb.UserVerifyPasswordResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    400,
				Message: "email user param is empty",
			},
		}, nil
	}

	password := gr.GetPassword()
	if password == "" {
		fmt.Println(fmt.Sprintf("[gRPC][Users][VerifyPassword][Error] %v", "password user param is empty"))
		return &pb.UserVerifyPasswordResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    400,
				Message: "password user param is empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][VerifyPassword][Error] %v", "user not found"))
		return &pb.UserVerifyPasswordResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    404,
				Message: "user not found",
			},
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][VerifyPassword][Error] %v", "invalid password"))
		return &pb.UserVerifyPasswordResponse{
			Data: nil,
			Meta: nil,
			Error: &pb.UserError{
				Code:    400,
				Message: "invalid password",
			},
		}, nil
	}

	res := &pb.UserVerifyPasswordResponse{
		Data:  user.ToProto(),
		Meta:  nil,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][VerifyPassword][Response] %v", res))
	return res, nil
}

// List return a collection of users.
func (us *Service) List(ctx context.Context, gr *pb.UserListRequest) (*pb.UserListResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][Users][List][Request] empty = %v", ""))

	//TODO(ca): check bdd connection

	listedUsers, err := us.usersSvc.UserList()
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][List][Error] %v", err))
		return &pb.UserListResponse{
			Data: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	data := make([]*pb.User, 0)
	for _, user := range listedUsers {
		data = append(data, user.ToProto())
	}

	res := &pb.UserListResponse{
		Data:  data,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][List][Response] %v", res))
	return res, nil
}

// Update updates a user.
func (us *Service) Update(ctx context.Context, gr *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	id := gr.GetUserId()
	fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Request] user_id = %v", id))

	user, err := us.usersSvc.GetByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Error] %v", err))
		return &pb.UserUpdateResponse{
			Data: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	email := gr.GetData().GetEmail()
	name := gr.GetData().GetName()
	password := gr.GetData().GetPassword()

	if email == "" && name == "" && password == "" {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Error] %v", err))
		return &pb.UserUpdateResponse{
			Data: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: "not email, name or password to change",
			},
		}, nil
	}

	if email != "" {
		user.Email = email
	}
	if name != "" {
		user.Name = name
	}
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Error] %v", err))
			return &pb.UserUpdateResponse{
				Data: nil,
				Meta: nil,
				Error: &pb.UserError{
					Code:    500,
					Message: "could not generate hashed password",
				},
			}, nil
		}
		user.Password = string(hashedPassword)
	}

	if err := us.usersSvc.Update(user); err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Error] %v", err))
		return &pb.UserUpdateResponse{
			Data: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.UserUpdateResponse{
		Data: user.ToProto(),
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Response] %v", res))
	return res, nil
}

// Delete delete a user.
func (us *Service) Delete(ctx context.Context, gr *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	id := gr.GetUserId()
	fmt.Println(fmt.Sprintf("[gRPC][Users][Delete][Request] user_id = %v", id))

	user, err := us.usersSvc.GetByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Delete][Error] %v", err))
		return &pb.UserDeleteResponse{
			Data: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	if err := us.usersSvc.Delete(user); err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][Users][Delete][Error] %v", err))
		return &pb.UserDeleteResponse{
			Data: nil,
			Error: &pb.UserError{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.UserDeleteResponse{
		Data: user.ToProto(),
	}

	fmt.Println(fmt.Sprintf("[gRPC][Users][Update][Response] %v", res))
	return res, nil
}
