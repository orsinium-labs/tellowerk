package controllers

type Controller interface {
	Name() string

	Start() error
	Stop() error

	TakeOff() error
	Land() error

	Clockwise() error
}
