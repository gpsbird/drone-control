package drone

import (
	"log"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"image"
	"fmt"
	"time"
	"github.com/socketbind/drone-control/decoder"
)

func DroneControl(videoChannel chan *image.Image, commandChannel chan interface{}) {
	drone := tello.NewDriver("8890")

	imageHandler := func(im *image.Image) {
		videoChannel <- im
	}

	work := func() {
		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(4)
			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			decoder.Decode(pkt, imageHandler)
		})

		for {
			select {
			case cmd := <-commandChannel:
				switch cmd := cmd.(type) {
				case TakeOffCommand:
					log.Printf("Going to take off, %q", cmd)
					drone.TakeOff()
				case LandCommand:
					log.Printf("Going to land, %q", cmd)
					drone.Land()
				case RotateCounterClockwiseCommand:
					log.Printf("Rotating counter-clockwise %d", cmd.Value)
					drone.CounterClockwise(cmd.Value)
				case RotateClockwiseCommand:
					log.Printf("Rotating clockwise %d", cmd.Value)
					drone.Clockwise(cmd.Value)
				case UpCommand:
					log.Printf("Going up %d", cmd.Value)
					drone.Up(cmd.Value)
				case DownCommand:
					log.Printf("Going down %d", cmd.Value)
					drone.Down(cmd.Value)
				case LeftCommand:
					log.Printf("Going left %d", cmd.Value)
					drone.Left(cmd.Value)
				case RightCommand:
					log.Printf("Going right %d", cmd.Value)
					drone.Right(cmd.Value)
				case ForwardCommand:
					log.Printf("Going forward %d", cmd.Value)
					drone.Forward(cmd.Value)
				case BackwardCommand:
					log.Printf("Going backward %d", cmd.Value)
					drone.Backward(cmd.Value)
				}
			}
		}
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
