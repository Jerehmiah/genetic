package main
import (
  	"github.com/jerehmiah/genetic/genetic"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func Run(){
	
	cfg := pixelgl.WindowConfig{
		Title: "Life",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync: true,
	}
	win,err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	go func() {
		observation := genetic.GetObservationStream(win)	
		environment := genetic.CreateEnvironment()
		environment.Abiogenesis(observation)
	}()
	

	for !win.Closed() {
		win.Clear(colornames.Firebrick)
		win.Update()
	}
	
}

func main(){

	pixelgl.Run(Run)
}

