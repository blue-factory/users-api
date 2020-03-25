package client

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	users "github.com/microapis/users-api"
	usersClient "github.com/microapis/users-api/client"
)

// TODO(ca): should implements code for validate createdAt and updateAt attrs values.

func before() (string, string, error) {
	host := os.Getenv("HOST")
	if host == "" {
		err := fmt.Errorf(fmt.Sprintf("Create: missing env variable HOST, failed with %s value", host))
		return "", "", err
	}

	port := os.Getenv("PORT")
	if port == "" {
		err := fmt.Errorf(fmt.Sprintf("Create: missing env variable PORT, failed with %s value", port))
		return "", "", err
	}

	return host, port, nil
}

// TestCreate ...
func TestCreate(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create with empty values for new user
	user := &users.User{}
	_, err = us.Create(user)
	if err != nil && err.Error() != "invalid name" {
		t.Errorf("Create: us.Create() error msg: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: us.Create() invalid error: %s", err.Error())
	}

	// Test create with invalid name new user
	user = &users.User{
		Name: "",
	}
	_, err = us.Create(user)
	if err != nil && err.Error() != "invalid name" {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
	}

	// Test create invalid email new user
	user = &users.User{
		Name:  "fake_user",
		Email: "",
	}
	_, err = us.Create(user)
	if err != nil && err.Error() != "invalid email" {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
	}

	// Test create invalid password new user
	user = &users.User{
		Name:     "fake_user",
		Email:    "fake_email_" + randomUUID.String(),
		Password: "",
	}
	_, err = us.Create(user)
	if err != nil && err.Error() != "invalid password" {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
	}

	// Test create valid new user
	user = &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("Create: us.Create(user) failed: %s", err.Error())
		return
	}

	expected := newUser.ID
	if expected == "" {
		t.Errorf("Create: newUser.ID(\"\") failed, expected %v, got %v", expected, user.ID)
		return
	}

	expected = newUser.Email
	if user.Email != expected {
		t.Errorf("Create: newUser.Email(\"\") failed, expected %v, got %v", expected, user.Email)
		return
	}

	expected = newUser.Password
	if user.Password == expected {
		t.Errorf("Create: newUser.Password(\"\") failed, expected %v, got %v", expected, user.Password)
		return
	}

	expected = newUser.Name
	if user.Name != expected {
		t.Errorf("Create: newUser.Name(\"\") failed, expected %v, got %v", expected, user.Name)
		return
	}

	// Delete created user
	err = us.Delete(newUser.ID)
	if err != nil {
		t.Errorf("Create: failed at Delete created user err=%v", err)
	}
}

// TestGet ...
func TestGet(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	user := &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("TestGet: us.Create(user) failed: %s", err.Error())
		return
	}

	// Test Get with invalid param
	getUser, err := us.Get("")
	if err == nil {
		t.Errorf("TestGet: us.Get(id) failed: %s", err.Error())
	}

	// Test Get new user
	getUser, err = us.Get(newUser.ID)
	if err != nil {
		t.Errorf("TestGet: us.Get(id) failed: %s", err.Error())
		return
	}

	expected := getUser.ID
	if newUser.ID != expected {
		t.Errorf("TestGet: getUser.ID(\"\") failed, expected %v, got %v", expected, newUser.ID)
		return
	}

	expected = getUser.Email
	if newUser.Email != expected {
		t.Errorf("TestGet: getUser.Email(\"\") failed, expected %v, got %v", expected, newUser.Email)
		return
	}

	expected = getUser.Password
	if newUser.Password != expected {
		t.Errorf("TestGet: getUser.Password(\"\") failed, expected %v, got %v", expected, newUser.Password)
		return
	}

	expected = getUser.Name
	if newUser.Name != expected {
		t.Errorf("TestGet: getUser.Name(\"\") failed, expected %v, got %v", expected, newUser.Name)
		return
	}

	// Delete created user
	err = us.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestGet: failed at Delete creaded user err=%v", err)
	}
}

// TestGetByEmail ...
func TestGetByEmail(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create new user
	user := &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("TestGetByEmail: us.Create(user) failed: %s", err.Error())
		return
	}

	// Test TestGetByEmail with invalid param
	getUserByEmail, err := us.GetByEmail("")
	if err == nil {
		t.Errorf("TestGetByEmail: us.GetByEmail() failed: %s", err.Error())
	}

	// Test TestGetByEmail new user
	getUserByEmail, err = us.GetByEmail(newUser.Email)
	if err != nil {
		t.Errorf("TestGetByEmail: us.GetByEmail(email) failed: %s", err.Error())
		return
	}

	expected := getUserByEmail.ID
	if newUser.ID != expected {
		t.Errorf("TestGetByEmail: getnewUser.ID(\"\") failed, expected %v, got %v", expected, user.ID)
		return
	}

	expected = getUserByEmail.Email
	if newUser.Email != expected {
		t.Errorf("TestGetByEmail: getUserByEmail.Email(\"\") failed, expected %v, got %v", expected, newUser.Email)
		return
	}

	expected = getUserByEmail.Password
	if newUser.Password != expected {
		t.Errorf("TestGetByEmail: getUserByEmail.Password(\"\") failed, expected %v, got %v", expected, newUser.Password)
		return
	}

	expected = getUserByEmail.Name
	if newUser.Name != expected {
		t.Errorf("TestGetByEmail: getUserByEmail.Name(\"\") failed, expected %v, got %v", expected, newUser.Name)
		return
	}

	// Delete created user
	err = us.Delete(getUserByEmail.ID)
	if err != nil {
		t.Errorf("TestGetByEmail: failed at Delete creaded user err=%v", err)
	}
}

// TestVerifyPassword ...
func TestVerifyPassword(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create new user
	user := &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("TestVerifyPassword: us.Create(user) failed: %s", err.Error())
		return
	}

	// Test verifyPasswordUser with invalid email param
	err = us.VerifyPassword("", user.Password)
	if err != nil && err.Error() != "invalid email" {
		t.Errorf(`TestVerifyPassword: us.VerifyPassword("", user.Password) failed: %s`, err.Error())
	}
	if err == nil {
		t.Errorf(`TestVerifyPassword: us.VerifyPassword("", user.Password) failed: %s`, err.Error())
	}

	// Test verifyPasswordUser with invalid password param
	err = us.VerifyPassword(user.Email, "")
	if err != nil && err.Error() != "invalid password" {
		t.Errorf(`TestVerifyPassword: us.VerifyPassword(email, "") failed: %s`, err.Error())
	}
	if err == nil {
		t.Errorf(`TestVerifyPassword: us.VerifyPassword(email, "") failed: %s`, err.Error())
	}

	// Test verifyPasswordUser with new user
	err = us.VerifyPassword(user.Email, user.Password)
	if err != nil {
		t.Errorf("TestVerifyPassword: us.VerifyPassword(email, password) failed: %s", err.Error())
		return
	}

	// Delete created user
	err = us.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestVerifyPassword: failed at Delete creaded user err=%v", err)
	}
}

// TestUpdate ...
func TestUpdate(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create new user
	user := &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("TestUpdate: us.Create(user) failed: %s", err.Error())
		return
	}

	// Update user with invalid userID
	updateUser := &users.User{}
	_, err = us.Update("", updateUser)
	if err != nil && err.Error() != "invalid ID" {
		t.Errorf("TestUpdate: us.Update(id, user) failed: %s", err.Error())
		return
	}
	if err == nil {
		t.Errorf("TestUpdate: us.Update(id, user) failed: %s", err.Error())
		return
	}

	// Update user with empty values
	updateUser = &users.User{}
	_, err = us.Update(newUser.ID, updateUser)
	if err == nil {
		t.Errorf("TestUpdate: us.Update(id, user) failed: %s", err.Error())
		return
	}

	// update user
	updateUser = &users.User{
		ID:       "update_" + newUser.ID,
		Name:     "update_" + newUser.Name,
		Email:    "update_" + newUser.Email,
		Password: "update_" + newUser.Password,
	}

	updatedUser, err := us.Update(newUser.ID, updateUser)
	if err != nil {
		t.Errorf("TestUpdate: us.Update(id, user) failed: %s", err.Error())
		return
	}

	expected := updatedUser.ID
	if newUser.ID != expected {
		t.Errorf("TestUpdate: updatedUser.ID(\"\") failed, expected %v, got %v", expected, newUser.ID)
		return
	}

	expected = updatedUser.Name
	if updateUser.Name != expected {
		t.Errorf("TestUpdate: updatedUser.Name(\"\") failed, expected %v, got %v", expected, updateUser.Name)
		return
	}

	expected = updatedUser.Email
	if updateUser.Email != expected {
		t.Errorf("TestUpdate: updatedUser.Email(\"\") failed, expected %v, got %v", expected, updateUser.Email)
		return
	}

	// TODO(ca): check if user password has changed

	// Delete created user
	err = us.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestUpdate: failed at Delete creaded user err=%v", err)
	}
}

// TestDelete ...
func TestDelete(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create new user
	user := &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("TestDelete: us.Create(user) failed: %s", err.Error())
		return
	}

	// Delete created user
	err = us.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestDelete: failed at Delete creaded user err=%v", err)
	}

	// Test get deleted user
	_, err = us.Get(newUser.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		t.Errorf("TestDelete: us.Get(id) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("TestDelete: us.Get(id) failed: %s", err.Error())
	}
}

// TestList ...
func TestList(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	us, err := usersClient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create new user
	user := &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := us.Create(user)
	if err != nil {
		t.Errorf("TestList: us.Create(user) failed: %s", err.Error())
		return
	}

	// Get user list
	list, err := us.List()
	if err != nil {
		t.Errorf("TestList: us.List() failed: %s", err.Error())
		return
	}

	// Validate if repsonse is a slice
	rt := reflect.TypeOf(list).Kind()
	if rt.String() != "slice" {
		t.Errorf("TestList: us.List() failed: %s", "response is not slice")
		return
	}

	// Find created user into list slice response
	var u *users.User
	for _, v := range list {
		if v.ID == newUser.ID {
			u = v
		}
	}

	// Check if finded new user is not null
	if u == nil {
		t.Errorf("TestList: failed: %s", "created user is not found inside us.List()")
		return
	}

	// Delete created user
	err = us.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestList: failed at Delete creaded user err=%v", err)
	}
}
