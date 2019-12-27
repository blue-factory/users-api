package userssvc

import (
	"context"
	"fmt"

	pb "github.com/microapis/lib/proto"
	users "github.com/microapis/users-api"
	"github.com/microapis/users-api/database"
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

// UserGet Gets a user by ID.
func (us *Service) UserGet(ctx context.Context, gr *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	id := gr.GetUserId()
	fmt.Println(fmt.Sprintf("[gRPC][UsersService][Get][Request] id = %v", id))
	user, err := us.usersSvc.GetByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][Get][Error] %v", err))
		return &pb.UserGetResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.UserGetResponse{
		Meta:  nil,
		Data:  user.ToProto(),
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][UsersService][Get][Response] %v", res))
	return res, nil
}

// UserGetByEmail get a user by Email
func (us *Service) UserGetByEmail(ctx context.Context, gr *pb.UserGetByEmailRequest) (*pb.UserGetByEmailResponse, error) {
	email := gr.GetEmail()
	fmt.Println(fmt.Sprintf("[gRPC][UsersService][GetByEmail][Request] email = %v", email))

	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][GetByEmail][Error] %v", "email user params empty"))
		return &pb.UserGetByEmailResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "email user params empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][GetByEmail][Error] %v", "user not found"))
		return &pb.UserGetByEmailResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    404,
				Message: "user not found",
			},
		}, nil
	}

	res := &pb.UserGetByEmailResponse{
		Meta:  nil,
		Data:  user.ToProto(),
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][UsersService][GetByEmail][Response] %v", res))
	return res, nil
}

// UserCreate creates a new user into database.
func (us *Service) UserCreate(ctx context.Context, gr *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Request] data = %v", gr.GetData()))

	email := gr.GetData().GetEmail()
	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Error] %v", "email user param is empty"))
		return &pb.UserCreateResponse{
			Meta: nil,
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "email user param is empty",
			},
		}, nil
	}

	_, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		password := gr.GetData().GetPassword()
		if password == "" {
			fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Error] %v", "password user params is empty"))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    400,
					Message: "password user params is empty",
				},
			}, nil
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Error] %v", err))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    500,
					Message: "could not generate hashed password",
				},
			}, nil
		}

		user := &users.User{
			Email:    email,
			Password: string(hashedPassword),
		}

		if err := us.usersSvc.Create(user); err != nil {
			fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Error] %v", err))
			return &pb.UserCreateResponse{
				Meta: nil,
				Data: nil,
				Error: &pb.Error{
					Code:    500,
					Message: err.Error(),
				},
			}, nil
		}

		res := &pb.UserCreateResponse{
			Meta:  nil,
			Data:  user.ToProto(),
			Error: nil,
		}
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Response] %v", res))
		return res, nil
	}

	res := &pb.UserCreateResponse{
		Meta: nil,
		Data: nil,
		Error: &pb.Error{
			Code:    400,
			Message: "user already exists",
		},
	}

	fmt.Println(fmt.Sprintf("[gRPC][UsersService][Create][Response] %v", res))
	return res, nil
}

// UserVerifyPassword ...
func (us *Service) UserVerifyPassword(ctx context.Context, gr *pb.UserVerifyPasswordRequest) (*pb.UserVerifyPasswordResponse, error) {
	fmt.Println(fmt.Sprintf("[gRPC][UsersService][VerifyPassword][Request] email = %v password = %v", gr.GetEmail(), gr.GetPassword()))
	email := gr.GetEmail()
	if email == "" {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][VerifyPassword][Error] %v", "email user param is empty"))
		return &pb.UserVerifyPasswordResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    400,
				Message: "email user param is empty",
			},
		}, nil
	}

	password := gr.GetPassword()
	if password == "" {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][VerifyPassword][Error] %v", "password user param is empty"))
		return &pb.UserVerifyPasswordResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    400,
				Message: "password user param is empty",
			},
		}, nil
	}

	user, err := us.usersSvc.GetByEmail(email)
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][VerifyPassword][Error] %v", "user not found"))
		return &pb.UserVerifyPasswordResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    404,
				Message: "user not found",
			},
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(fmt.Sprintf("[gRPC][UsersService][VerifyPassword][Error] %v", "invalid password"))
		return &pb.UserVerifyPasswordResponse{
			Valid: false,
			Error: &pb.Error{
				Code:    400,
				Message: "invalid password",
			},
		}, nil
	}

	res := &pb.UserVerifyPasswordResponse{
		Valid: true,
		Error: nil,
	}

	fmt.Println(fmt.Sprintf("[gRPC][UsersService][VerifyPassword][Response] %v", res))
	return res, nil
}

// UserList return a collection of users.
func (us *Service) UserList(ctx context.Context, gr *pb.UserListRequest) (*pb.UserListResponse, error) {
	fmt.Println(fmt.Sprintf("[GRPC][UsersService][List][Request] empty = %v", ""))
	//TODO(ca): check bdd connection
	listedUsers, err := us.usersSvc.UserList()
	if err != nil {
		fmt.Println(fmt.Sprintf("[GRPC][UsersService][List][Error] %v", err))
		return &pb.UserListResponse{
			Data: nil,
			Error: &pb.Error{
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

	fmt.Println(fmt.Sprintf("[GRPC][UsersService][List][Response] %v", res))
	return res, nil
}
