package users

import (
	"errors"
	"os"
	"strings"

	"github.com/dueckminor/mypi-tools/go/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound = errors.New("User not found")
)

type User struct {
	Name     string
	Password string
	Mail     string
	Groups   []string
}

type UserCfg struct {
	config.Config
}

func parseUser(userEntry config.Config) *User {
	return &User{
		Name:     userEntry.GetString("name"),
		Password: userEntry.GetString("password"),
	}
}

func (cfg *UserCfg) GetUsers() ([]*User, error) {
	userEntries := cfg.GetArray()

	result := make([]*User, 0, len(userEntries))

	for _, userEntry := range userEntries {
		result = append(result, parseUser(userEntry))
	}
	return result, nil
}

func (cfg *UserCfg) GetUser(username string) (*User, error) {
	username = strings.ToLower(username)
	for _, userEntry := range cfg.GetArray() {
		name := userEntry.GetString("name")
		if name == username {
			return &User{
				Name:     name,
				Password: userEntry.GetString("password"),
			}, nil
		}
	}
	return nil, ErrNotFound
}

func (cfg *UserCfg) UpdateUser(user *User) error {
	for _, userEntry := range cfg.GetArray() {
		name := userEntry.GetString("name")
		if name == user.Name {
			return userEntry.SetString("password", user.Password)
		}
	}
	return ErrNotFound
}

func ReadUserCfg() (result *UserCfg, err error) {
	result = &UserCfg{}
	result.Config, err = config.ReadConfigFile("etc/mypi-auth/users.yml")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		result.Config = config.New("etc/mypi-auth/users.yml", make([]interface{}, 0))
	}
	return result, nil
}

func AddUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost-2, // DefaultCost takes to long on a Raspberry-Pi
	)
	if err != nil {
		return err
	}

	userCfg, err := ReadUserCfg()
	if err != nil {
		return err
	}

	user, _ := userCfg.GetUser(username)
	if user != nil {
		user.Password = string(hash)
		err = userCfg.UpdateUser(user)
		if err != nil {
			return err
		}
	} else {
		err = userCfg.AddArrayElement(User{
			Name:     username,
			Password: string(hash),
		})
		if err != nil {
			return err
		}
	}

	return userCfg.Write()
}

func CheckPassword(username, password string) bool {
	userCfg, err := ReadUserCfg()
	if err != nil {
		return false
	}

	user, _ := userCfg.GetUser(username)
	if user == nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
