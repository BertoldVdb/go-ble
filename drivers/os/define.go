package oshci

import "errors"

var (
	// ErrorClosed is returned when the driver is closed using Close().
	ErrorClosed = errors.New("Device was closed")

	// ErrorHUP is returned when the connection to the device is no longer valid.
	ErrorHUP = errors.New("Device was disconnected")
)
