package login

type Login interface {
	GetLoginPage() (salt, execution, lt string, err error)
	Login(username, password string) (string, bool, error)
}
