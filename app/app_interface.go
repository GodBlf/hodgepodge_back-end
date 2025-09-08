package app

type App interface {
	Login() error
	RollCall() error
}
