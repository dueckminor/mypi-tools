package users

import (
	"errors"
	"os"

	"github.com/dueckminor/mypi-api/go/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound = errors.New("User not found")
)

type User struct {
	Name     string
	Password string
}

type UserCfg struct {
	config.Config
}

func (cfg *UserCfg) GetUser(username string) (*User, error) {
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
			userEntry.SetString("password", user.Password)
			return nil
		}
	}
	return ErrNotFound
}

func ReadUserCfg() (result *UserCfg, err error) {
	result = &UserCfg{}
	result.Config, err = config.ReadConfigFile("etc/users.yml")
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		result.Config = config.New("etc/users.yml", make([]interface{}, 0))
	}
	return result, nil
}

func AddUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
		userCfg.UpdateUser(user)
	} else {
		userCfg.AddArrayElement(User{
			Name:     username,
			Password: string(hash),
		})
	}

	return userCfg.Write()
}

func CheckPasswd(username, password string) bool {
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
