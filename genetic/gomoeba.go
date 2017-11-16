package genetic
import (
	"os"
	"image"
	_ "image/png"
	"math/rand"
	"time"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)

type Gomoeba struct {
	GopherChan chan int
}
var Window *pixelgl.Window
var GopherPic pixel.Picture
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewGomoeba() chan int{
	gomoeba := &Gomoeba{make(chan int, 5)}
	direction := (r.Float64() * 2.0) - 1.0
	go func(gomoeba *Gomoeba){
		angle := 0.0
		sprite := pixel.NewSprite(GopherPic, GopherPic.Bounds())
		last := time.Now()
		x := float64(r.Intn(824)) + 100
		y := float64(r.Intn(568)) + 100
		
		GopherLoop: for {
			select {
			case n := <- gomoeba.GopherChan:
				switch n{
				case Die:
					//close(o.gophers[event.Actor - 1])
					break GopherLoop
				}
			default:
				dt := time.Since(last).Seconds()
				last = time.Now()
				angle += 0.3 * dt * direction
				mat := pixel.IM
				mat = mat.ScaledXY(pixel.ZV, pixel.V(0.05, 0.05))
				mat = mat.Moved(pixel.V(x, y))
				mat = mat.Rotated(pixel.V(x, y), angle)
				sprite.Draw(Window, mat)
			}
			
		}
	}(gomoeba)
	return gomoeba.GopherChan
}

func GetReadyForGomoebas(window *pixelgl.Window){
	pic, err := loadPicture("gopher.png")
	if err != nil {
		panic(err)
	}
	GopherPic = pic
	Window = window
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