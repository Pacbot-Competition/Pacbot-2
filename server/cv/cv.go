package cv

import (
	"image"
	"log"
	"sort"
	"sync"

	"gocv.io/x/gocv"
)

type CVModule struct {
	commands chan Command
	results  chan Result
	cam      *gocv.VideoCapture
	muCV     sync.RWMutex
	m        gocv.Mat
	maxDims  image.Point
}

type Coords image.Point

type Command int

const (
	f Command = iota
	c
	q
)

type Result struct {
	originalCommand Command
	result          []byte
}

func (c *Coords) sum() int {
	return c.X + c.Y
}

func (c *Coords) diff() int {
	return c.Y - c.X // TODO see if this right
}

func (c *Coords) toPoint() image.Point {
	return image.Point{c.X, c.Y}
}

func orderCorners(c1, c2, c3, c4 Coords) gocv.PointVector {
	ordered := gocv.NewPointVector()
	coords := []Coords{c1, c2, c3, c4}

	sort.Slice(coords, func(i, j int) bool {
		return coords[i].sum() < coords[j].sum()
	})

	ordered.Append(coords[0].toPoint()) // top left
	ordered.Append(coords[2].toPoint()) // bottom right

	sort.Slice(coords, func(i, j int) bool {
		return coords[i].diff() < coords[j].diff()
	})

	ordered.Append(coords[1].toPoint()) // top right
	ordered.Append(coords[3].toPoint()) // bottom left

	return gocv.NewPointVector() // TODO
}

func intMax(i, j int) int {
	if i < j {
		return j
	}
	return i
}

func dist(i, j image.Point) int {
	dx := i.X - j.X
	dy := i.Y - j.Y
	return dx*dx + dy*dy
}

func NewCVModule(commands chan Command, results chan Result, camId int) *CVModule {
	cam, err := gocv.VideoCaptureDevice(camId)
	if err != nil {
		log.Fatal(err)
	}
	cv := CVModule{
		commands: commands,
		results:  results,
		cam:      cam,
		muCV:     sync.RWMutex{},
		m:        gocv.NewMat(),
		maxDims:  image.Pt(0, 0),
	}
	return &cv
}

func (cv *CVModule) getDestRect(ordered gocv.PointVector) gocv.PointVector {
	dst := gocv.NewPointVector()

	tl, br, tr, bl := ordered.At(0), ordered.At(1), ordered.At(2), ordered.At(3)
	maxWidth := intMax(dist(br, bl), dist(tr, tl))
	maxHeight := intMax(dist(tr, br), dist(tl, bl))

	dst.Append(image.Pt(0, 0))
	dst.Append(image.Pt(maxWidth-1, 0))
	dst.Append(image.Pt(maxWidth-1, maxHeight-1))
	dst.Append(image.Pt(0, maxHeight-1))

	cv.muCV.Lock()
	cv.maxDims = image.Pt(maxWidth, maxHeight)
	cv.muCV.Unlock()

	return dst
}

func coordsToBytes(x, y float32) []byte {
	return []byte{byte(x), byte(y)}
}

func (cv *CVModule) SetCorners(c1, c2, c3, c4 Coords) {
	ordered := orderCorners(c1, c2, c3, c4)
	destRect := cv.getDestRect(ordered)
	cv.muCV.Lock()
	cv.m = gocv.GetPerspectiveTransform(ordered, destRect)
	cv.muCV.Unlock()
}

func (cv *CVModule) RunLoop() {
	for {
		newCommand := <-cv.commands

		img := gocv.NewMat()

		if ok := cv.cam.Read(&img); !ok {
			return // can include error handling here
		}

		switch newCommand {
		case q:
			return
		case f:
			cv.results <- Result{newCommand, img.ToBytes()}
		case c:
			warpedImg := gocv.NewMat()

			cv.muCV.RLock()
			gocv.WarpPerspective(img, &warpedImg, cv.m, cv.maxDims)
			cv.muCV.RUnlock()

			// TODO figure out flipping logic (if necessary)

			// may have to tweak median blur radius
			gocv.MedianBlur(warpedImg, &warpedImg, 5)

			mask := gocv.NewMat()
			gocv.CvtColor(warpedImg, &img, gocv.ColorBGRToHSV)     // TODO see if this is the right conversion
			gocv.InRange(img, gocv.NewMat(), gocv.NewMat(), &mask) // TODO add ranges
			gocv.BitwiseAndWithMask(img, img, &img, mask)
			grayscale := gocv.Split(img)[2]                                                      // TODO see if this is good conversion to grayscale
			contours := gocv.FindContours(grayscale, gocv.RetrievalList, gocv.ChainApproxSimple) // TODO check enums

			if contours.Size() == 0 {
				return
			}
			champ := contours.At(0)
			champArea := gocv.ContourArea(champ)
			for idx := 0; idx < contours.Size(); idx++ {
				contour := contours.At(idx)
				if area := gocv.ContourArea(contour); area > champArea {
					champ = contour
					champArea = area
				}
			}
			if champArea > 100 {
				rect := gocv.BoundingRect(champ)
				x, y := float32(rect.Min.X), float32(rect.Min.Y)
				w, h := float32(rect.Max.X)-x, float32(rect.Max.Y)-y
				x_c, y_c := x+w*0.5, y+h*0.5

				// TODO currently gives pixel coordinate
				// coordsToBytes will handle whatever coordinate scheme we decide
				cv.results <- Result{newCommand, coordsToBytes(x_c, y_c)}
			}
			contours.Close()
		}
	}
}

func (cv *CVModule) Quit() {
	cv.commands <- q
}
