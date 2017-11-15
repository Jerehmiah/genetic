package genetic
import (
  "time"
  "fmt"
)

type Cell struct {
	OsmosisStream chan Protein
	EmissionStream chan Protein
	ObservationStream chan Event
	hunger int
	Death chan bool
	ProteinName int
}
var ProteinCount = 0
var starvation = 3

func NewCell(osmosis chan Protein, emission chan Protein, observation chan Event)  {
	ProteinCount = ProteinCount + 1
	tick := time.Tick(time.Second)
	cell := &Cell{osmosis, emission, observation, 0, make(chan bool), ProteinCount }
	go func(cell *Cell){
		CellLoop: for {
			select {
			case <- tick:
				cell.Tick()
			case <- cell.Death:
				//This anonymous method should be the sole reference to the cell, so it will get GC'd 
				//upon completion
				break CellLoop
			}	
		}
	}(cell)
	observation <- Event{Born,  cell.ProteinName}
}

func (c *Cell) Tick () {
	select {
	case protein := <- c.OsmosisStream:
		
		c.osmosis(protein)
		
		//If I wasn't hungry, see if there's enough food to split
		if(c.hunger == 0){
			select {
			case <-c.OsmosisStream:
				c.split()
			default:
			  //do nothing	
			}
		}
		
		c.hunger = 0
	default:
		c.hunger = c.hunger + 1
		if c.hunger == 1{
		  c.ObservationStream <- Event{Hunger, c.ProteinName}
		}
		if c.hunger == starvation{
			c.ObservationStream <- Event{Starve, c.ProteinName}
		}

		if c.hunger > starvation{
			c.ObservationStream <- Event{Die, c.ProteinName}
			c.Death <- true
		}
	}
}

func ( c *Cell) osmosis(protein Protein){
	c.ObservationStream <- Event{Eat, c.ProteinName}
	c.react()
}

func (c *Cell) react(){
	//Do something that takes 200 ms
	time.Sleep(200* time.Millisecond)
}

func (c *Cell) emit(protein Protein){
	fmt.Printf("Cell number %d emitted a protein\n", c.ProteinName)
	c.EmissionStream <- protein
} 

func (c *Cell) split(){
	c.ObservationStream <- Event{Split, c.ProteinName}

	//Splitting still consumes a protein, so react
	c.react()
	NewCell(c.OsmosisStream, c.EmissionStream, c.ObservationStream)
}
