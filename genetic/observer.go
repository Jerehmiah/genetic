package genetic

import(
	"fmt"
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
	Actor int
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
	switch event.Type{
	case Born:
		o.gophers[event.Actor - 1] = NewGomoeba()
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

