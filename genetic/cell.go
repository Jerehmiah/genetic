package genetic
import (
  "time"
)

type Cell struct {
	OsmosisStream chan Protein
	ObservationStream chan Event
	hunger int
	Death chan bool
	Identity float64
	identityChain chan float64
}
var starvation = 3

func NewCell(osmosis chan Protein, observation chan Event, identityChain chan float64)  {
	tick := time.Tick(time.Second)
	cell := &Cell{osmosis, observation, 0, make(chan bool), <-identityChain, identityChain}
	identityChain <- cell.Identity + 1
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
	observation <- Event{Born,  cell.Identity}
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
		  c.ObservationStream <- Event{Hunger, c.Identity}
		}
		if c.hunger == starvation{
			c.ObservationStream <- Event{Starve, c.Identity}
		}

		if c.hunger > starvation{
			c.ObservationStream <- Event{Die, c.Identity}
			c.Death <- true
		}
	}
}

func ( c *Cell) osmosis(protein Protein){
	c.ObservationStream <- Event{Eat, c.Identity}
	c.react()
}

func (c *Cell) react(){
	//Do something that takes 200 ms
	time.Sleep(200* time.Millisecond)
}

func (c *Cell) split(){
	c.ObservationStream <- Event{Split, c.Identity}

	//Splitting still consumes a protein, so react
	c.react()
	NewCell(c.OsmosisStream, c.ObservationStream, c.identityChain)
}
