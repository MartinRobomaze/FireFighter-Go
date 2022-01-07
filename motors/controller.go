package motors

import (
	"FireFighter/IO"
	"FireFighter/utils"
	"time"
)

type MotorDirection string

const (
	Forward  MotorDirection = "F"
	Backward                = "B"
)

type Controller struct {
	Handler    IO.HwHandler
	BrakeDelay time.Duration
}

func (c *Controller) Forward(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: Forward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: Forward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: Forward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: Forward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Backward(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: Backward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: Backward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: Backward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: Backward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Left(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: Backward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: Backward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: Forward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: Forward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Right(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: Forward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: Forward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: Backward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: Backward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Slide(angle float64, speed int) {
	var motors []IO.MotorData

	if angle >= 0 {
		if angle <= 45 {
			BCSpeed := utils.ValMap(angle, 45, 0, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Forward, Speed: speed},
				{Motor: "B", Direction: Forward, Speed: BCSpeed},
				{Motor: "C", Direction: Forward, Speed: BCSpeed},
				{Motor: "D", Direction: Forward, Speed: speed},
			}
		} else if angle <= 90 {
			BCSpeed := utils.ValMap(angle, 45, 90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Forward, Speed: speed},
				{Motor: "B", Direction: Backward, Speed: BCSpeed},
				{Motor: "C", Direction: Backward, Speed: BCSpeed},
				{Motor: "D", Direction: Forward, Speed: speed},
			}
		} else if angle <= 135 {
			ADSpeed := utils.ValMap(angle, 135, 90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Forward, Speed: ADSpeed},
				{Motor: "B", Direction: Backward, Speed: speed},
				{Motor: "C", Direction: Backward, Speed: speed},
				{Motor: "D", Direction: Forward, Speed: ADSpeed},
			}
		} else if angle <= 180 {
			BCSpeed := utils.ValMap(angle, 135, 180, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Backward, Speed: speed},
				{Motor: "B", Direction: Backward, Speed: BCSpeed},
				{Motor: "C", Direction: Backward, Speed: BCSpeed},
				{Motor: "D", Direction: Backward, Speed: speed},
			}
		}
	} else {
		if angle >= -45 {
			ADSpeed := utils.ValMap(angle, -45, 0, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Forward, Speed: ADSpeed},
				{Motor: "B", Direction: Forward, Speed: speed},
				{Motor: "C", Direction: Forward, Speed: speed},
				{Motor: "D", Direction: Forward, Speed: ADSpeed},
			}
		} else if angle >= -90 {
			ADSpeed := utils.ValMap(angle, -45, -90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Backward, Speed: ADSpeed},
				{Motor: "B", Direction: Forward, Speed: speed},
				{Motor: "C", Direction: Forward, Speed: speed},
				{Motor: "D", Direction: Backward, Speed: ADSpeed},
			}
		} else if angle >= -135 {
			BCSpeed := utils.ValMap(angle, -135, -90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Backward, Speed: speed},
				{Motor: "B", Direction: Forward, Speed: BCSpeed},
				{Motor: "C", Direction: Forward, Speed: BCSpeed},
				{Motor: "D", Direction: Backward, Speed: speed},
			}
		} else if angle >= -180 {
			ADSpeed := utils.ValMap(angle, -135, -180, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: Backward, Speed: ADSpeed},
				{Motor: "B", Direction: Backward, Speed: speed},
				{Motor: "C", Direction: Backward, Speed: speed},
				{Motor: "D", Direction: Backward, Speed: ADSpeed},
			}
		}
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}
