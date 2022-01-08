package main

import (
	"FireFighter/IO"
	"FireFighter/comm"
	"FireFighter/firefinder"
	"FireFighter/lepton"
	"FireFighter/motors"
	"FireFighter/stepper"
	"fmt"
	gouvc "github.com/MartinRobomaze/go-uvc"
	"github.com/MartinRobomaze/gocv"
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"image"
	"log"
	"time"
)

const (
	baseSpeed = 120

	fanPin = 17

	stepPin = 27
	dirPin  = 18
)

func initSensorInterface() (commHandler *comm.Handler, handler *IO.HwHandler, controller *motors.Controller) {
	commHandler, err := comm.New("/dev/ttyUSB0")
	if err != nil {
		logrus.WithError(err).Log(logrus.PanicLevel, "ERROR initializing commHandler")
	}

	handler = IO.New(commHandler)

	controller = &motors.Controller{
		Handler:    handler,
		BrakeDelay: 50 * time.Millisecond,
	}

	return commHandler, handler, controller
}

func initRPiGPIOServices() (controller *stepper.Controller, fan rpio.Pin) {
	err := rpio.Open()
	if err != nil {
		logrus.WithError(err).Log(logrus.PanicLevel, "ERROR Initializing RPi GPIO")
	}

	fan = rpio.Pin(fanPin)
	fan.Output()

	controller = stepper.New(stepPin, dirPin)

	return controller, fan
}

func initFireFinder() (cam *lepton.Driver, frameChan <-chan *gouvc.Frame, fireFinder *firefinder.FireFinder) {
	cam, err := lepton.New()
	if err != nil {
		logrus.WithError(err).Log(logrus.PanicLevel, "ERROR connecting to thermal camera")
	}

	frameChan, err = cam.StartStreaming()
	if err != nil {
		logrus.WithError(err).Log(logrus.PanicLevel, "ERROR starting stream from thermal camera")
	}

	fireFinder = firefinder.New(
		40,
		1.5,
		0.3,
		image.Point{X: 5, Y: 5},
		[]int{160, 120},
		57,
		42.5,
	)

	return cam, frameChan, fireFinder
}

func closeServices(commHandler *comm.Handler, cam *lepton.Driver) {
	err := cam.Close()
	if err != nil {
		logrus.WithError(err).Log(logrus.PanicLevel, "ERROR closing thermal camera connection")
	}

	err = commHandler.Close()
	if err != nil {
		logrus.WithError(err).Log(logrus.PanicLevel, "ERROR closing serial connection")
	}
}

func main() {
	commHandler, _, motorsController := initSensorInterface()

	stepperController, _ := initRPiGPIOServices()

	cam, frameChan, fireFinder := initFireFinder()

	defer closeServices(commHandler, cam)

	for {
		thermalDataRaw := <-frameChan
		thermalData := gocv.NewMatFromInt16Arr(120, 160, gocv.MatTypeCV16U, thermalDataRaw.Data)

		fire, fireData := fireFinder.FindFire(thermalData)
		if !fire {
			continue
		}

		fireAngleHor, fireAngleVer := fireFinder.FireCoordsToAngles(fireData.FireLocation)

		if fireAngleHor < -20 {
			motorsController.Left(baseSpeed)
		} else if fireAngleHor < 20 {
			motorsController.Slide(fireAngleHor, baseSpeed)
			log.Println("mopslik")
		} else {
			motorsController.Right(baseSpeed)
		}

		if fireAngleVer < -10 {
			stepperController.Move(stepper.Up, 200, 1000/9*time.Millisecond)
		} else if fireAngleVer < -5 {
			stepperController.Move(stepper.Up, 50, 1000/9*time.Millisecond)
		} else if 5 < fireAngleVer && fireAngleVer < 10 {
			stepperController.Move(stepper.Down, 50, 1000/9*time.Millisecond)
		} else if fireAngleVer > 10 {
			stepperController.Move(stepper.Down, 200, 1000/9*time.Millisecond)
		}

		logrus.WithFields(logrus.Fields{
			"fire temperature":      fireData.FireTemperature,
			"fire location":         fmt.Sprintf("%d \t %d", fireData.FireLocation.X, fireData.FireLocation.Y),
			"fire angle horizontal": fireAngleHor,
			"fire angle vertical":   fireAngleVer,
		}).Log(logrus.InfoLevel, "FIRE")
	}
}
