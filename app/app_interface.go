package app

type App interface {
	AppLogin() error
	AppImplRollCall() error
}
