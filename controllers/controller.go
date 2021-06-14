package controllers

type Controller interface {
	Name() string

	Start() error
	Stop() error

	TakeOff() error
	Land() error

	Clockwise(int) error
	CounterClockwise(int) error
	Forward(int) error
	Backward(int) error
	Left(int) error
	Right(int) error
	Down(int) error
	Up(int) error
}
