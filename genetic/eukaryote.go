package genetic
import (
  "time"
)

type Eukaryote struct {
	OsmosisStream chan Protein
	NucleusStream chan Protein
	ObservationStream chan Event
	hunger int
	Death chan bool
	Identity float64
	identityChain chan float64
}

func NewEukaryote(osmosis chan Protein, observation chan Event, identityChain chan float64)  {
	tick := time.Tick(time.Second)
	organelleChan := make(chan Protein,30)
	eukaryote := &Eukaryote{osmosis, organelleChan, observation, 0, make(chan bool), <-identityChain, identityChain}
	
	//Make organelles
	organelleIdentity := make(chan float64, 1)
	organelleIdentity <- float64(eukaryote.Identity) + 0.1
	NewOrganelle(organelleChan, observation, organelleIdentity)
	NewOrganelle(organelleChan, observation, organelleIdentity)

	identityChain <- eukaryote.Identity + 1
	go func(eukaryote *Eukaryote){
		EukaryoteLoop: for {
			select {
			case <- tick:
				eukaryote.Tick()
			case <- eukaryote.Death:
				//This anonymous method should be the sole reference to the eukaryote, so it will get GC'd 
				//upon completion
				break EukaryoteLoop
			}	
		}
	}(eukaryote)
	observation <- Event{Born,  eukaryote.Identity}
}

func (c *Eukaryote) Tick () {
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

func ( c *Eukaryote) osmosis(protein Protein){
	c.ObservationStream <- Event{Eat, c.Identity}
	c.NucleusStream <- protein
	c.react()
}

func (c *Eukaryote) react(){
	//Do something that takes 200 ms
	time.Sleep(200* time.Millisecond)
}

func (c *Eukaryote) split(){
	c.ObservationStream <- Event{Split, c.Identity}

	//Splitting still consumes a protein, so react
	c.react()
	NewEukaryote(c.OsmosisStream, c.ObservationStream, c.identityChain)
}
