package genetic

import(
	"image"
	"os"
	_ "image/png"
	"math/rand"
	"time"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

const (
	Born = 0
	Hunger = 1
	Starve = 2
	Split = 3
	Die = 4
	Eat = 5
	Feed = 6
	Terminate = 7
)

type Event struct{
	Type int
	Actor int
}

type Observer struct{
	observationStream chan Event
	window *pixelgl.Window
	Done bool
	gophers []chan int
}
var GopherPic pixel.Picture
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetObservationStream(window *pixelgl.Window) chan Event{
	observationStream := make(chan Event, 2000)
	observer := &Observer{observationStream, window, false, make([]chan int, 100)}
	
	pic, err := loadPicture("gopher.png")
	if err != nil {
		panic(err)
	}
	GopherPic = pic

	go func(observer *Observer){
		for !observer.Done{
			select{
			case event := <-observationStream:
				observer.ProcessEvent(event)
			}
		}
		fmt.Println("Observer was done")
	}(observer)
	return observationStream
}

func (o *Observer) ProcessEvent(event Event){
	switch event.Type{
	case Born:
		gopherChan := make(chan int)
		o.gophers[event.Actor - 1] = gopherChan
		go func(o *Observer){
			angle := 0.0
			sprite := pixel.NewSprite(GopherPic, GopherPic.Bounds())
			last := time.Now()
			x := float64(r.Intn(924)) + 50
			y := float64(r.Intn(668)) + 50
			GopherLoop: for {
				select {
				case n := <- gopherChan:
					switch n{
					case Die:
						break GopherLoop
					}
				default:
					dt := time.Since(last).Seconds()
					last = time.Now()
					angle += 0.3 * dt
					mat := pixel.IM
					
					mat = mat.ScaledXY(pixel.ZV, pixel.V(0.05, 0.05))
					mat = mat.Moved(pixel.V(x, y))
					mat = mat.Rotated(pixel.V(x, y), angle)
					sprite.Draw(o.window, mat)
				}
				
			}
		}(o)
		
		fmt.Printf(Blue("Cell number %d was born\n"), event.Actor)
	case Hunger:
		fmt.Printf("Cell number %d is hungry\n", event.Actor)
	case Starve:
	  	fmt.Printf(Yellow("Cell number %d is starving\n"), event.Actor)
	case Split:
		fmt.Printf("Cell number %d consumed an additional protein, and is splitting\n", event.Actor	)
	case Die:
		o.gophers[event.Actor - 1] <- Die
		fmt.Printf(Red("Cell number %d died\n"), event.Actor)
	case Eat:
		fmt.Printf("Cell number %d consumed a protein\n", event.Actor)
	case Feed:
		fmt.Printf(White("***********  Adding %d proteins\n"), event.Actor)
	case Terminate:
		fmt.Println("Received Terminate signal")
		o.Done = true
		o.window.SetClosed(true)
	default:
		fmt.Println("Got some wacky event")
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}