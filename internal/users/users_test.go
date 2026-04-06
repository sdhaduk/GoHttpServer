package users

import (
	"errors"
	"net/mail"
	"reflect"
	"testing"
)

func TestAddUser(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}
	
	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}
	
	if len(testManager.users) != 1 {
		t.Errorf("bad test manager user count, should be 1")
		if len(testManager.users) < 1 {
			t.Fatal()
		}
	}

	expectedUser := User{
		FirstName: testFirstName,
		LastName: testLastName,
		Email: *testEmail,
	}

	foundUser := testManager.users[0]

	if !reflect.DeepEqual(expectedUser, foundUser) {
		t.Errorf("added user data is not correct\nwanted:%+v\ngot:%+v\n", expectedUser, foundUser)
	}
}

func TestAddUserEmptyFirstName(t *testing.T) {
	testManager := NewManager()

	testFirstName := ""
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Errorf("did not get error for empty first name")
	} else {
		if err.Error() == "invalid first name" {
			return
		} else {
			t.Errorf("received incorrect error: %s", err)
		}
	}
}

func TestAddUserEmptyLastName(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := ""
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Errorf("did not get error for empty first name")
	} else {
		if err.Error() == "invalid last name" {
			return
		} else {
			t.Errorf("received incorrect error: %s", err)
		}
	}
}

func TestGetUserByName(t *testing.T) {
	testManager := NewManager()

	err := testManager.AddUser("foo", "bar", "foo@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v\n", err)
	}

	err = testManager.AddUser("bar", "baz", "bbaz@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v\n", err)
	}

	err = testManager.AddUser("foo", "baz", "foobaz@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v\n", err)
	}
	
	err = testManager.AddUser("bar", "foo", "barfoo@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v\n", err)
	}

	tests := map[string]struct{
		first string
		last string
		expected *User
		expectedErr error
	}{
		"simple lookup": {
			first: "foo",
			last: "bar",
			expected: &testManager.users[0],
			expectedErr: nil,
		},
		"last element lookup": {
			first: "bar",
			last: "foo",
			expected: &testManager.users[3],
			expectedErr: nil,
		},
		"no match lookup": {
			first: "john",
			last: "smith",
			expected: nil,
			expectedErr: ErrNoResultsFound,
		},
		"partial match lookup": {
			first: "foo",
			last: "foo",
			expected: nil,
			expectedErr: ErrNoResultsFound,
		},
		"empty first name": {
			first: "",
			last: "baz",
			expected: nil,
			expectedErr: ErrNoResultsFound,
		},
		"empty last name": {
			first: "foo",
			last: "",
			expected: nil,
			expectedErr: ErrNoResultsFound,
		},
	}

	for name, test := range tests {
		result, err := testManager.GetUserByName(test.first, test.last)

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("%s: invalid result\ngot: %+v\nwanted: %+v\n", name, result, test.expected)
		}
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: invalid error\ngot: %+v\nwanted: %+v\n", name, err, test.expectedErr)
		}
	}
}

func TestGetUserByNameDuplicate(t *testing.T) {
	testingManager := NewManager()

	err := testingManager.AddUser("foo", "bar", "foobar@example.com")

	if err != nil {
		t.Fatalf("error adding user")
	}

	err = testingManager.AddUser("foo", "bar", "foobar@example.com")
	if err.Error() != "user already exists" {
		t.Errorf("incorrect error received")
	}
}
