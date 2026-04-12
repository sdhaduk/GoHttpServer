package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
	"testing"

	"github.com/sdhaduk/GoHttpServer/internal/users"
)

func TestHandleRoot(t *testing.T) {
	w := httptest.NewRecorder()

	handleRoot(w, nil)

	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome to our Homepage!\n")

	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleGoodbye(t *testing.T) {
	w := httptest.NewRecorder()

	handleGoodbye(w, nil)

	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Goodbye, World!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleHelloParameterized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?user=TestMan", nil)

	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)
	
	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleHelloParameterizedNoParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello/", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)
	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}
func TestHandleHelloParameterizedWrongParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?foo=bar", nil)
	w := httptest.NewRecorder()

	handleHelloParameterized(w, req)
	desiredCode := http.StatusOK

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, User!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleResponsesHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/responses/TestMan/hello/", nil)
	req.SetPathValue("user", "TestMan")

	w := httptest.NewRecorder()

	handleUserResponsesHello(w, req)
	desiredCode := http.StatusOK	
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Hello, TestMan!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleHelloHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)

	req.Header.Set("user", "Test Man")

	w := httptest.NewRecorder()

	handleHelloHeader(w, req)
	
	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Hello, Test Man!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleHelloHeaderNoHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/hello/", nil)

	w := httptest.NewRecorder()

	handleHelloHeader(w, req)
	
	desiredCode := http.StatusBadRequest
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("invalid username provided\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestHandleJSON(t *testing.T) {
	testRequest := UserData{FirstName: "Test Man",}

	marshalledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}
	
	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewBuffer(marshalledRequestBody))

	w := httptest.NewRecorder()

	handleJSON(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected %v, but got %v\nbody: %v", desiredCode, w.Code, w.Body.String())
	}
	
	expectedMessage := []byte("Hello, Test Man!\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got %v, expected %v", w.Body.String(), string(expectedMessage))
	}
}

func TestAddUser(t *testing.T) {
	testUser := UserData{
		FirstName: "Test",
		LastName: "User",
		Email: "TestMan@example.com",
	}

	marshalledRequestBody, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("error marshalling test data")
	}

	req := httptest.NewRequest(http.MethodPost, "/add-user", bytes.NewBuffer(marshalledRequestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	testManager := users.NewManager()
	testServer := server{
		userManager: testManager,
	}

	testServer.addUser(w, req)

	desiredCode := http.StatusCreated

	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v but got: %v\nbody: %s\n", desiredCode, w.Code, w.Body.String())
	}  

	resultUser, err := testManager.GetUserByName(testUser.FirstName, testUser.LastName)
	if err != nil {
		t.Fatalf("error retreiving user from manager\n%v", err)
	}

	convertedUser := convertUserToUserData(resultUser)

	if !reflect.DeepEqual(convertedUser, &testUser) {
		t.Errorf("convertedUser does not match the testUser")
	}
}

func TestConvertUserToUserData(t *testing.T) {
	testFirstName := "Test"
	testLastName := "User"
	testEmail, err := mail.ParseAddress("testUser@example.com")

	if err != nil {
		t.Fatalf("error parsing email")
	}

	testUser := users.User{
		FirstName: testFirstName,
		LastName: testLastName,
		Email: *testEmail,
	}

	result := convertUserToUserData(&testUser)
	
	expectedUser := &UserData{
		FirstName: testFirstName,
		LastName: testLastName,
		Email:   testEmail.Address,
	}

	if !reflect.DeepEqual(expectedUser, result) {
		t.Errorf("the expectedUser and result do not match")
	}
}

