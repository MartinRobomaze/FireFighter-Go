package motors

import (
	"FireFighter/IO"
	"FireFighter/utils"
	"time"
)

type Controller struct {
	Handler    *IO.HwHandler
	BrakeDelay time.Duration
}

func (c *Controller) Forward(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: IO.Forward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: IO.Forward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: IO.Forward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: IO.Forward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Backward(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: IO.Backward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: IO.Backward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: IO.Backward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: IO.Backward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Left(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: IO.Backward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: IO.Backward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: IO.Forward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: IO.Forward,
			Speed:     speed,
		},
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}

func (c *Controller) Right(speed int) {
	motors := []IO.MotorData{
		{
			Motor:     "A",
			Direction: IO.Forward,
			Speed:     speed,
		},
		{
			Motor:     "B",
			Direction: IO.Forward,
			Speed:     speed,
		},
		{
			Motor:     "C",
			Direction: IO.Backward,
			Speed:     speed,
		},
		{
			Motor:     "D",
			Direction: IO.Backward,
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
				{Motor: "A", Direction: IO.Forward, Speed: speed},
				{Motor: "B", Direction: IO.Forward, Speed: BCSpeed},
				{Motor: "C", Direction: IO.Forward, Speed: BCSpeed},
				{Motor: "D", Direction: IO.Forward, Speed: speed},
			}
		} else if angle <= 90 {
			BCSpeed := utils.ValMap(angle, 45, 90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Forward, Speed: speed},
				{Motor: "B", Direction: IO.Backward, Speed: BCSpeed},
				{Motor: "C", Direction: IO.Backward, Speed: BCSpeed},
				{Motor: "D", Direction: IO.Forward, Speed: speed},
			}
		} else if angle <= 135 {
			ADSpeed := utils.ValMap(angle, 135, 90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Forward, Speed: ADSpeed},
				{Motor: "B", Direction: IO.Backward, Speed: speed},
				{Motor: "C", Direction: IO.Backward, Speed: speed},
				{Motor: "D", Direction: IO.Forward, Speed: ADSpeed},
			}
		} else if angle <= 180 {
			BCSpeed := utils.ValMap(angle, 135, 180, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Backward, Speed: speed},
				{Motor: "B", Direction: IO.Backward, Speed: BCSpeed},
				{Motor: "C", Direction: IO.Backward, Speed: BCSpeed},
				{Motor: "D", Direction: IO.Backward, Speed: speed},
			}
		}
	} else {
		if angle >= -45 {
			ADSpeed := utils.ValMap(angle, -45, 0, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Forward, Speed: ADSpeed},
				{Motor: "B", Direction: IO.Forward, Speed: speed},
				{Motor: "C", Direction: IO.Forward, Speed: speed},
				{Motor: "D", Direction: IO.Forward, Speed: ADSpeed},
			}
		} else if angle >= -90 {
			ADSpeed := utils.ValMap(angle, -45, -90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Backward, Speed: ADSpeed},
				{Motor: "B", Direction: IO.Forward, Speed: speed},
				{Motor: "C", Direction: IO.Forward, Speed: speed},
				{Motor: "D", Direction: IO.Backward, Speed: ADSpeed},
			}
		} else if angle >= -135 {
			BCSpeed := utils.ValMap(angle, -135, -90, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Backward, Speed: speed},
				{Motor: "B", Direction: IO.Forward, Speed: BCSpeed},
				{Motor: "C", Direction: IO.Forward, Speed: BCSpeed},
				{Motor: "D", Direction: IO.Backward, Speed: speed},
			}
		} else if angle >= -180 {
			ADSpeed := utils.ValMap(angle, -135, -180, 0, float64(speed))

			motors = []IO.MotorData{
				{Motor: "A", Direction: IO.Backward, Speed: ADSpeed},
				{Motor: "B", Direction: IO.Backward, Speed: speed},
				{Motor: "C", Direction: IO.Backward, Speed: speed},
				{Motor: "D", Direction: IO.Backward, Speed: ADSpeed},
			}
		}
	}

	c.Handler.MotorsHandler.SetMotors(motors)
}
