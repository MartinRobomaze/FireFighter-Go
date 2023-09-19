# FireFighter-Go
FireFighter robot RPi program in Go.  
Uses a modified version of [gocv](https://github.com/MartinRobomaze/gocv) and [go-uvc](https://github.com/MartinRobomaze/go-uvc) to communicate with thermal camera and to detect fire. It communicates with [FireFighterSensorInterface](https://github.com/MartinRobomaze/FireFighterSensorInterface) to get data from light sensors and to control motors.
## Installation
1. OpenCV has to be installed either by compiling it from source or by installing precompiled libraries. Precompiled binaries for Raspberry Pi can be found [here](https://lindevs.com/install-precompiled-opencv-on-raspberry-pi/).  
2. Install golang. For Ubuntu/Debian run `apt install golang` as root. 
3. Clone this repository.
4. Install dependencies by running `go get` inside the project dirctory.
5. Build the project using with `go build`.
6. Run the program with `./FireFighter`, root permissions or setting permissions using [udev](https://stackoverflow.com/questions/22713834/libusb-cannot-open-usb-device-permission-isse-netbeans-ubuntu) might be needed for accessing the camera.
