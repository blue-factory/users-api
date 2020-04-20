package usersclient_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	users "github.com/microapis/users-api"
	usersclient "github.com/microapis/users-api/client"
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

	client, err := usersclient.New(host + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}

	randomUUID := uuid.New()

	// Test create with empty values for new user
	user := &users.User{}
	_, err = client.Create(user)
	if err != nil && err.Error() != "invalid name" {
		t.Errorf("Create: client.Create() error msg: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: client.Create() invalid error: %s", err.Error())
	}

	// Test create with invalid name new user
	user = &users.User{
		Name: "",
	}
	_, err = client.Create(user)
	if err != nil && err.Error() != "invalid name" {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
	}

	// Test create invalid email new user
	user = &users.User{
		Name:  "fake_user",
		Email: "",
	}
	_, err = client.Create(user)
	if err != nil && err.Error() != "invalid email" {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
	}

	// Test create invalid password new user
	user = &users.User{
		Name:     "fake_user",
		Email:    "fake_email_" + randomUUID.String(),
		Password: "",
	}
	_, err = client.Create(user)
	if err != nil && err.Error() != "invalid password" {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
	}

	// Test create valid new user
	user = &users.User{
		Email:    "fake_email_" + randomUUID.String(),
		Password: "fake_password",
		Name:     "fake_name",
	}

	// Create new user
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("Create: client.Create(user) failed: %s", err.Error())
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

	expectedBool := newUser.Verified
	if false != expectedBool {
		t.Errorf("Create: newUser.Verified(\"\") failed, expected %v, got %v", expectedBool, false)
		return
	}

	// Delete created user
	err = client.Delete(newUser.ID)
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

	client, err := usersclient.New(host + ":" + port)
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
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("TestGet: client.Create(user) failed: %s", err.Error())
		return
	}

	// Test Get with invalid param
	getUser, err := client.Get("")
	if err == nil {
		t.Errorf("TestGet: client.Get(id) failed: %s", err.Error())
	}

	// Test Get new user
	getUser, err = client.Get(newUser.ID)
	if err != nil {
		t.Errorf("TestGet: client.Get(id) failed: %s", err.Error())
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

	expectedBool := getUser.Verified
	if false != expectedBool {
		t.Errorf("TestGet: getUser.Verified(\"\") failed, expected %v, got %v", expectedBool, false)
		return
	}

	// Delete created user
	err = client.Delete(newUser.ID)
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

	client, err := usersclient.New(host + ":" + port)
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
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("TestGetByEmail: client.Create(user) failed: %s", err.Error())
		return
	}

	// Test TestGetByEmail with invalid param
	getUserByEmail, err := client.GetByEmail("")
	if err == nil {
		t.Errorf("TestGetByEmail: client.GetByEmail() failed: %s", err.Error())
	}

	// Test TestGetByEmail new user
	getUserByEmail, err = client.GetByEmail(newUser.Email)
	if err != nil {
		t.Errorf("TestGetByEmail: client.GetByEmail(email) failed: %s", err.Error())
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

	expectedBool := getUserByEmail.Verified
	if false != expectedBool {
		t.Errorf("TestGetByEmail: getUserByEmail.Verified(\"\") failed, expected %v, got %v", expectedBool, false)
		return
	}

	// Delete created user
	err = client.Delete(getUserByEmail.ID)
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

	client, err := usersclient.New(host + ":" + port)
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
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("TestVerifyPassword: client.Create(user) failed: %s", err.Error())
		return
	}

	// Test verifyPasswordUser with invalid email param
	err = client.VerifyPassword("", user.Password)
	if err != nil && err.Error() != "invalid email" {
		t.Errorf(`TestVerifyPassword: client.VerifyPassword("", user.Password) failed: %s`, err.Error())
	}
	if err == nil {
		t.Errorf(`TestVerifyPassword: client.VerifyPassword("", user.Password) failed: %s`, err.Error())
	}

	// Test verifyPasswordUser with invalid password param
	err = client.VerifyPassword(user.Email, "")
	if err != nil && err.Error() != "invalid password" {
		t.Errorf(`TestVerifyPassword: client.VerifyPassword(email, "") failed: %s`, err.Error())
	}
	if err == nil {
		t.Errorf(`TestVerifyPassword: client.VerifyPassword(email, "") failed: %s`, err.Error())
	}

	// Test verifyPasswordUser with new user
	err = client.VerifyPassword(user.Email, user.Password)
	if err != nil {
		t.Errorf("TestVerifyPassword: client.VerifyPassword(email, password) failed: %s", err.Error())
		return
	}

	// Delete created user
	err = client.Delete(newUser.ID)
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

	client, err := usersclient.New(host + ":" + port)
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
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("TestUpdate: client.Create(user) failed: %s", err.Error())
		return
	}

	// Update user with invalid userID
	updateUser := &users.User{}
	_, err = client.Update("", updateUser)
	if err != nil && err.Error() != "invalid ID" {
		t.Errorf("TestUpdate: client.Update(id, user) failed: %s", err.Error())
		return
	}
	if err == nil {
		t.Errorf("TestUpdate: client.Update(id, user) failed: %s", err.Error())
		return
	}

	// Update user with empty values
	updateUser = &users.User{}
	_, err = client.Update(newUser.ID, updateUser)
	if err == nil {
		t.Errorf("TestUpdate: client.Update(id, user) failed: %s", err.Error())
		return
	}

	// update user
	updateUser = &users.User{
		ID:       "update_" + newUser.ID,
		Name:     "update_" + newUser.Name,
		Email:    "update_" + newUser.Email,
		Password: "update_" + newUser.Password,
		Verified: true,
	}

	updatedUser, err := client.Update(newUser.ID, updateUser)
	if err != nil {
		t.Errorf("TestUpdate: client.Update(id, user) failed: %s", err.Error())
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

	expectedBool := updatedUser.Verified
	if true != expectedBool {
		t.Errorf("TestUpdate: updatedUser.Verified(\"\") failed, expected %v, got %v", expectedBool, true)
		return
	}

	// TODO(ca): check if user password has changed

	// Delete created user
	err = client.Delete(newUser.ID)
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

	client, err := usersclient.New(host + ":" + port)
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
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("TestDelete: client.Create(user) failed: %s", err.Error())
		return
	}

	// Delete created user
	err = client.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestDelete: failed at Delete creaded user err=%v", err)
	}

	// Test get deleted user
	_, err = client.Get(newUser.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		t.Errorf("TestDelete: client.Get(id) failed: %s", err.Error())
	}
	if err == nil {
		t.Errorf("TestDelete: client.Get(id) failed: %s", err.Error())
	}
}

// TestList ...
func TestList(t *testing.T) {
	host, port, err := before()
	if err != nil {
		t.Errorf(err.Error())
	}

	client, err := usersclient.New(host + ":" + port)
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
	newUser, err := client.Create(user)
	if err != nil {
		t.Errorf("TestList: client.Create(user) failed: %s", err.Error())
		return
	}

	// Get user list
	list, err := client.List()
	if err != nil {
		t.Errorf("TestList: client.List() failed: %s", err.Error())
		return
	}

	// Validate if repsonse is a slice
	rt := reflect.TypeOf(list).Kind()
	if rt.String() != "slice" {
		t.Errorf("TestList: client.List() failed: %s", "response is not slice")
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
		t.Errorf("TestList: failed: %s", "created user is not found inside client.List()")
		return
	}

	// Delete created user
	err = client.Delete(newUser.ID)
	if err != nil {
		t.Errorf("TestList: failed at Delete creaded user err=%v", err)
	}
}
