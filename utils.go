package main

import (
	"errors"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

func Min(a ...int) int {
	min := int(^uint(0) >> 1)
	for _, i := range a {
		if i < min {
			min = i
		}
	}
	return min
}

func Range(start, end, step int) []int {
	var out []int

	for i := start; i <= end; i += step {
		out = append(out, i)
	}

	return out
}

func FRange(start, end, step float64) []float64 {
	var out []float64

	for i := start; i <= end; i += step {
		out = append(out, i)
	}

	return out
}

func Walk(searchDir string) []string {
	fileList := []string{}

	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	return fileList
}

func MakeDirs(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

const PI float64 = 3.141592653589793

func Radians(deg float64) float64 {
	return deg * (PI / 180)
}

func Degrees(rad float64) float64 {
	return rad * (180 / PI)
}

func DegToNum(lat, lon float64, zoom int) (int, int) {
	latrad := Radians(lat)
	n := math.Pow(2.0, float64(zoom))
	xtile := int((lon + 180.0) / 360.0 * n)
	ytile := int((1.0 - math.Log(math.Tan(latrad)+(1/math.Cos(latrad)))/PI) / 2.0 * n)

	return xtile, ytile
}

func LatLonToTile(lat, lon float64, zoom int) (int, int) {
	latrad := Radians(lat)
	n := 1 << zoom
	x := int(float64(n) * ((lon + 180.0) / 360.0))
	y := int(float64(n) * (1 - (math.Log(math.Tan(latrad)+1/math.Cos(latrad)) / PI)) / 2.0)

	return x, y
}

func TileToLatLon(x, y, zoom int) (float64, float64) {
	n := 1 << zoom
	latrad := math.Atan(math.Sinh(PI * (1.0 - 2.0*float64(y)/float64(n))))
	lat := latrad * 180 / PI
	lon := float64(360*x/n - 180.0)

	return lat, lon
}

func TileBounds(x, y, zoom int) (float64, float64, float64, float64) {
	n := math.Pow(2.0, float64(zoom))
	lonmin := float64(x)/float64(n)*360.0 - 180.0
	latminRad := math.Atan(math.Sinh(PI * (1 - 2*float64(y)/float64(n))))
	latmin := Degrees(latminRad)
	lonmax := (float64(x)+1)/float64(n)*360.0 - 180.0
	latmaxRad := math.Atan(math.Sinh(PI * (1 - 2*(float64(y)+1)/float64(n))))
	latmax := Degrees(latmaxRad)

	return latmin, latmax, lonmin, lonmax
}

func DownloadFile(filePath, URL string) (*os.File, error) {
	response, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	file, err := os.Create(filePath)
	_, err = io.Copy(file, response.Body)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, nil
}
