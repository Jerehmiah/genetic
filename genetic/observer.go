package genetic

import(
	
	"fmt"
	"math"
	"github.com/faiface/pixel/pixelgl"
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
	Actor float64
}

type Observer struct{
	observationStream chan Event
	window *pixelgl.Window
	Done bool
	gophers []chan int
}

func GetObservationStream(window *pixelgl.Window) chan Event{
	observationStream := make(chan Event, 2000)
	observer := &Observer{observationStream, window, false, make([]chan int, 100)}
	
	GetReadyForGomoebas(window)

	go func(observer *Observer){
		for !observer.Done{
			observer.ProcessEvent(<-observationStream)
		}
	}(observer)
	return observationStream
}

func (o *Observer) ProcessEvent(event Event){
	isOrganelle := !(math.Trunc(event.Actor) == event.Actor)
	display :=  "Eukaryote"
	if isOrganelle {
		display = "Organelle"
	}
	switch event.Type{
	case Born:
		fmt.Printf(Blue("%s number %f was born\n"), display, event.Actor)
		if !isOrganelle{
			o.gophers[int(event.Actor - 1)] = NewGomoeba()
		}
	case Hunger:
		fmt.Printf("%s number %f is hungry\n", display, event.Actor)
	case Starve:
	  	fmt.Printf(Yellow("%s number %f is starving\n"), display, event.Actor)
	case Split:
		fmt.Printf("%s number %f consumed an additional protein, and is splitting\n", display, event.Actor	)
	case Die:
		if !isOrganelle {
			o.gophers[int(event.Actor - 1)] <- Die	
		}
		fmt.Printf(Red("%s number %f died\n"), display, event.Actor)
	case Eat:
		fmt.Printf("%s number %f consumed a protein\n", display, event.Actor)
	case Feed:
		fmt.Printf(White("***********  Adding %f proteins\n"), event.Actor)
	case Terminate:
		fmt.Println("Received Terminate signal")
		o.Done = true
		o.window.SetClosed(true)
	default:
		fmt.Println("Got some wacky event")
	}
}

