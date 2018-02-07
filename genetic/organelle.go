package genetic
import (
  "time"
)

type Organelle struct {
	OsmosisStream chan Protein
	ObservationStream chan Event
	hunger int
	Death chan bool
	Identity float64
	identityChain chan float64
}

func NewOrganelle(osmosis chan Protein, observation chan Event, identityChain chan float64)  {
	tick := time.Tick(time.Second)
	organelle := &Organelle{osmosis, observation, 0, make(chan bool), <-identityChain, identityChain}
	identityChain <- organelle.Identity + 0.1
	go func(organelle *Organelle){
		OrganelleLoop: for {
			select {
			case <- tick:
				organelle.Tick()
			case <- organelle.Death:
				//This anonymous method should be the sole reference to the organelle, so it will get GC'd 
				//upon completion
				break OrganelleLoop
			}	
		}
	}(organelle)
	observation <- Event{Born,  organelle.Identity}
}

func (c *Organelle) Tick () {
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

func ( c *Organelle) osmosis(protein Protein){
	c.ObservationStream <- Event{Eat, c.Identity}
	c.react()
}

func (c *Organelle) react(){
	//Do something that takes 200 ms
	time.Sleep(200* time.Millisecond)
}

func (c *Organelle) split(){
	c.ObservationStream <- Event{Split, c.Identity}

	//Splitting still consumes a protein, so react
	c.react()
	NewCell(c.OsmosisStream, c.ObservationStream, c.identityChain)
}
