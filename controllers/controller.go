package controllers

type Controller interface {
	Name() string

	Start() error
	Stop() error

	// take off and land
	TakeOff() error
	ThrowTakeOff() error
	Land() error
	PalmLand() error
	StopLanding() error

	// movement
	Clockwise(int) error
	CounterClockwise(int) error
	Forward(int) error
	Backward(int) error
	Left(int) error
	Right(int) error
	Down(int) error
	Up(int) error

	// tricks
	LeftFlip() error
	RightFlip() error
	FrontFlip() error
	BackFlip() error
	Bounce() error

	// settings
	SetFastMode() error
	SetSlowMode() error
	SetExposure(int) error
}
