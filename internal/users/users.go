package users

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type User struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Balance  int    `yaml:"balance"`
}

type Users struct {
	Users []User `yaml:"users"`
}

func LoadUsers() (*Users, error) {
	data, err := os.ReadFile("docs/users.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading users file: %v", err)
	}

	var users Users
	err = yaml.Unmarshal(data, &users)
	if err != nil {
		return nil, fmt.Errorf("error parsing users yaml: %v", err)
	}

	return &users, nil
}

func (u *Users) ShowPlayers() {
	fmt.Println("\nAvailable Players:")
	for _, user := range u.Users {
		fmt.Printf("- %s\n", user.Name)
	}
}

func (u *Users) Authenticate(name, password string) (*User, error) {
	attempts := 3
	var currentUser *User

	for attempts > 0 {
		for i := range u.Users {
			if u.Users[i].Name == name {
				currentUser = &u.Users[i]
				break
			}
		}

		if currentUser == nil {
			return nil, fmt.Errorf("user %s not found", name)
		}

		if currentUser.Balance <= 0 {
			fmt.Println("Sorry, you're broke 🥀")
			fmt.Print("Please enter another username: ")
			var newName string
			fmt.Scanln(&newName)
			return u.Authenticate(newName, password)
		}

		if currentUser.Password == password {
			return currentUser, nil
		}

		attempts--
		if attempts > 0 {
			fmt.Printf("Wrong password! %d attempts remaining\n", attempts)
			fmt.Print("Enter password again: ")
			fmt.Scanln(&password)
		} else {
			return nil, fmt.Errorf("authentication failed after 3 attempts")
		}
	}

	return nil, fmt.Errorf("authentication failed")
}
