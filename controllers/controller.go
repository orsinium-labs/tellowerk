package controllers

import "gobot.io/x/gobot/platforms/dji/tello"

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
	Rotate(int) error
	OY(int) error
	OX(int) error
	OZ(int) error
	Hover() error

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
	SetVideoBitRate(tello.VideoBitRate) error
}
