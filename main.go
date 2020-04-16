package main

import (
	"fmt"
	"strconv"
	"sync"
)

/**
Odessa   -> 30.0246 46.1599 31.6341 46.7747 (46.4682, 30.8293)
Kiev     -> 29.6081 49.9300 31.4813 50.8857 (50.4179, 30.6867)
Nikolaev -> 30.9839 46.7646 32.5934 47.3726 (47.0694, 31.7886)
DNR      -> 36.969 46.811 39.568 48.691     (47.9524, 38.0884)
*/

/**
https://services.arcgisonline.com/arcgis/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile
https://services.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile
*/
const (
	minZoom    = 10
	maxZoom    = 15
	tileScheme = "zyx" // zxy or zyx
	targetUrl  = "https://services.arcgisonline.com/arcgis/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile"
	outDir     = "Kiev_Adm"
)

type Bounds struct {
	west  float64
	south float64
	east  float64
	north float64
}

func main() {
	var urange [][2]string
	bounds := Bounds{
		west:  29.6081,
		south: 49.9300,
		east:  31.4813,
		north: 50.8857,
	}

	zoomRange := Range(minZoom, maxZoom+1, 1)
	for _, zoom := range zoomRange {
		minX, minY := DegToNum(bounds.north, bounds.west, zoom)
		maxX, maxY := DegToNum(bounds.south, bounds.east, zoom)

		yRange := Range(minY, maxY+1, 1)
		for _, yIndex := range yRange {
			xRange := Range(minX, maxX+1, 1)
			filePath := outDir + "\\" + strconv.Itoa(zoom) + "\\" + strconv.Itoa(yIndex)
			MakeDirs(filePath)

			for _, xIndex := range xRange {
				tileUrl := targetUrl + "/" + strconv.Itoa(zoom) + "/" + strconv.Itoa(yIndex) + "/" + strconv.Itoa(xIndex)
				fileName := filePath + "\\" + strconv.Itoa(xIndex) + ".png"

				if !FileExists(fileName) {
					urange = append(urange, [2]string{fileName, tileUrl})
				} else {
					continue
				}
			}
		}
	}

	wg := sync.WaitGroup{}
	for _, uu := range urange {
		filepath := uu[0]
		tileurl := uu[1]
		wg.Add(1)
		DownloadFile(filepath, tileurl)
		defer wg.Done()
		fmt.Println(tileurl)
	}

	defer wg.Wait()
	return
}
