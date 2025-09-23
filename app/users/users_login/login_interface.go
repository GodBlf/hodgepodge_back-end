package users_login

import "github.com/go-resty/resty/v2"

type UsersLoginInterface interface {
	Login(username string, password string) (*resty.Client, error)
}
