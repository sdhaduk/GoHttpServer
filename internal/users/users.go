package users

import (
	"errors"
	"fmt"
	"net/mail"
)

var ErrNoResultsFound = errors.New("no results found")

type User struct {
	FirstName string
	LastName string
	Email mail.Address
}

type Manager struct {
	users []User
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) AddUser(firstName, lastName, email string) error {
	parsedAddress, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email %s", email)
	}

	if firstName == "" {
		return fmt.Errorf("invalid first name")
	}

	if lastName == "" {
		return fmt.Errorf("invalid last name")
	}

	newUser := User{
		FirstName: firstName,
		LastName: lastName,
		Email: *parsedAddress,
	} 

	existingUser, err := m.GetUserByName(firstName, lastName)

	if err != nil && !errors.Is(err, ErrNoResultsFound) {
		return fmt.Errorf("error searching through users")
	} 

	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	m.users = append(m.users, newUser)
	return nil
}

func (m *Manager) GetUserByName(firstName, lastName string) (*User, error) {
	for i, user := range m.users {
		if user.FirstName == firstName && user.LastName == lastName {
			result := m.users[i]
			return &result, nil
		}
	}

	return nil, ErrNoResultsFound
}