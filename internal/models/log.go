package models

import "time"

type TempLog struct {
	DeviceID       string
	DeviceTime     time.Time
	Latitude       float64
	Longitude      float64
	Altitude       float64
	Course         float64
	Satellites     int
	SpeedOTG       float32
	AccelerationX1 float32
	AccelerationY1 float32
	Signal         int
	PowerSupply    int
}

type HasWarnings struct {
	DeviceID    string
	WarningTime time.Time
	WarningType int
}
