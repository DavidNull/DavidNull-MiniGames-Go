package persistence

import (
	"fmt"
	"os"

	"davgames/internal/users"

	"gopkg.in/yaml.v3"
)

func SaveUserBalance(user *users.User) error {
	data, err := os.ReadFile("docs/users.yaml")
	if err != nil {
		return fmt.Errorf("error leyendo archivo users.yaml: %v", err)
	}

	var usersData users.Users
	err = yaml.Unmarshal(data, &usersData)
	if err != nil {
		return fmt.Errorf("error parseando YAML: %v", err)
	}

	for i := range usersData.Users {
		if usersData.Users[i].Name == user.Name {
			usersData.Users[i].Balance = user.Balance
			break
		}
	}

	yamlData, err := yaml.Marshal(&usersData)
	if err != nil {
		return fmt.Errorf("error convirtiendo a YAML: %v", err)
	}

	err = os.WriteFile("docs/users.yaml", yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error guardando archivo YAML: %v", err)
	}

	return nil
}
