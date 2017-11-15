package genetic
import(
	"time"
	"fmt"
)

var environmentProteins = [][]string{
	{"",""}, 
	{""}, 
	{"","","","",""},
	{"","","","","","","","",""},
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""},
	{""},
	{""},
	{""},
	{"","","","",""},
	{"","","","","","","","",""},
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""}, 
	{"","",""},
	{""},
    {},
    {},
    {},
    {}}

type Environment struct {
	environmentStream chan *Protein
}

func CreateEnvironment() *Environment{
	inputStream := make(chan *Protein, 12)
	return &Environment{inputStream}
}

func (e *Environment) Abiogenesis(observation chan Event){
	_ = NewCell(e.environmentStream, e.environmentStream, observation)	
	tick := time.Tick(time.Second)
	idx := 0
	AbioLoop: for {
		select {
		case <- tick:
			if idx >= len(environmentProteins){
				break AbioLoop
			}
			observation <- Event{Feed, len(environmentProteins[idx]) }

			for _, p := range environmentProteins[idx]{
				e.InputProtein(p)
			}
			idx = idx + 1
		}
	}
	observation <- Event{Terminate, 0}
	fmt.Println("Simulation concluded")
}

func (e *Environment) InputProteins(proteins []Protein){
	for _, p := range proteins {
		e.InputProtein(&p)
	}
}

func (e *Environment) InputProtein(protein Protein){
	e.environmentStream <- &protein
}

