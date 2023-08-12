package cv

import (
	"image"
	"sync"

	"gocv.io/x/gocv"
)

type CVModule struct {
	cam       gocv.VideoCapture
	window    gocv.Window
	muCV      sync.Mutex
	m         gocv.Mat
	showLocal bool
}

type Coords struct {
	x float32
	y float32
}

type Corners struct {
	c1 Coords
	c2 Coords
	c3 Coords
	c4 Coords
}

type Command int

const (
	f Command = iota
	c
)

type Result struct {
	originalCommand Command
	result          []byte
}

func orderCorners(c Corners) gocv.PointVector {
	return gocv.NewPointVector() // TODO
}

func getDestRect(gocv.PointVector) gocv.PointVector {
	return gocv.NewPointVector() // TODO
}

func coordsToBytes(x, y float32) []byte {
	return []byte{}
}

func (cv *CVModule) SetCorners(c Corners) {
	ordered := orderCorners(c)
	destRect := getDestRect(ordered)
	cv.muCV.Lock()
	cv.m = gocv.GetPerspectiveTransform(ordered, destRect)
	cv.muCV.Unlock()
}

func (cv *CVModule) Process(commands chan Command, results chan Result) {
	newCommand := <-commands

	img := gocv.NewMat()
	// res is boolean check to see if read worked
	if ok := cv.cam.Read(&img); !ok {
		return // can include error handling here
	}

	// if cv.showLocal {
	// 	cv.window.IMShow(img)
	// }

	if newCommand == f {
		results <- Result{newCommand, img.ToBytes()}
	} else if newCommand == c {
		warpedImg := gocv.NewMat()

		cv.muCV.Lock()
		// get max width and height
		gocv.WarpPerspective(img, &warpedImg, cv.m, image.Pt(0, 0))
		cv.muCV.Unlock()

		// may have to tweak median blur radius
		gocv.MedianBlur(warpedImg, &warpedImg, 5)

		mask := gocv.NewMat()
		gocv.CvtColor(warpedImg, &img, 0)                      // add mode
		gocv.InRange(img, gocv.NewMat(), gocv.NewMat(), &mask) // add ranges
		gocv.BitwiseAndWithMask(img, img, &img, mask)
		grayscale := gocv.Split(img)[2]                // memory leaks?
		contours := gocv.FindContours(grayscale, 0, 0) // add modes

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
			results <- Result{newCommand, coordsToBytes(x_c, y_c)}
		}

		contours.Close() // see if this is necessary
	}
}
