package stepper

import (
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

type Direction int

const (
	Up Direction = iota
	Down
)

type StepsData struct {
	Direction Direction
	Steps     int64
	SleepTime time.Duration
}

type Controller struct {
	StepPin rpio.Pin
	DirPin  rpio.Pin

	StepsChan chan StepsData
}

func New(stepPin int, dirPin int) *Controller {
	c := &Controller{
		StepPin:   rpio.Pin(stepPin),
		DirPin:    rpio.Pin(dirPin),
		StepsChan: make(chan StepsData),
	}

	c.StepPin.Mode(rpio.Output)
	c.DirPin.Mode(rpio.Output)

	return c
}

func (c *Controller) Move(direction Direction, steps int64, duration time.Duration) {
	data := StepsData{
		Direction: direction,
		Steps:     steps,
		SleepTime: time.Duration(int64(duration)/steps) / 2,
	}

	c.StepsChan <- data
}

func (c *Controller) update() {
	for {
		steps := <-c.StepsChan

		if steps.Direction == Up {
			c.DirPin.Write(rpio.High)
		} else if steps.Direction == Down {
			c.DirPin.Write(rpio.Low)
		}

		for i := int64(0); i < steps.Steps; i++ {
			c.StepPin.Write(rpio.High)
			time.Sleep(steps.SleepTime)
			c.StepPin.Write(rpio.Low)
			time.Sleep(steps.SleepTime)
		}
	}
}
