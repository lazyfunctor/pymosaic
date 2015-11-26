package mosaic

import (
    "fmt"
    "os"
    "image"
    "container/heap"
    _ "image/color"
    "math"
    "sync"
    "image/draw"
    "image/png"
    "encoding/json"
    )


type Distance struct {
    dist float64
    path string
}

type DistHeap []*Distance

func (h DistHeap) Len() int {
    return len(h)
}

func (h DistHeap) Less(i, j int) bool {
    return h[i].dist < h[j].dist
}

func (h DistHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *DistHeap) Push(x interface{}) {
    *h = append(*h, x.(*Distance))
}

func (h *DistHeap) Pop() interface{} {
    oldh := *h
    x := oldh[len(oldh)-1]
    newh := oldh[0:len(oldh)-1]
    *h = newh
    return x
}

func Test() {
    var dh DistHeap
    heap.Init(&dh)
    heap.Push(&dh, &Distance{dist: 0.9, path: "foo1"})
    heap.Push(&dh, &Distance{dist: 0.2, path: "foo2"})
    heap.Push(&dh, &Distance{dist: 0.3, path: "foo3"})
    fmt.Println(heap.Pop(&dh))
    fmt.Println(heap.Pop(&dh))
    fmt.Println(heap.Pop(&dh))
}

//var usedTiles = make(map[string]int)
//var mutex = &sync.Mutex{}
const resetReps = 30

type RepeatCheck struct {
    m *sync.Mutex
    usedTiles map[string]int
    reset int
}

func calcEuclidean(color1 [3]float64, color2 [3]float64) (dist float64) {
    r1, g1, b1 := color1[0], color1[1], color1[2]
    r2, g2, b2 := color2[0], color2[1], color2[2]
    return math.Sqrt((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2))
}

type Grid struct {
    Color [3]float64
    X int
    Y int
}

type Tile struct {
    X int
    Y int
    path string
}

func TileFinder(gridChan chan *Grid, tileData map[string][3]float64, outChan chan *Tile, r *RepeatCheck, wg *sync.WaitGroup) {
    defer wg.Done()
    for grid := range(gridChan) {
        var dh DistHeap
        heap.Init(&dh)
        for path, color := range tileData {
            dist := calcEuclidean(color, grid.Color)
            heap.Push(&dh, &Distance{dist: dist, path:path})
        }
        r.m.Lock()
        var bestTile string
        for {
            tile := heap.Pop(&dh).(*Distance)
            //fmt.Println(tile.path)
            used, _ := r.usedTiles[tile.path]
            if used == 0 {
                r.usedTiles[tile.path] = 1
                bestTile = tile.path
                break
            } else {
                r.usedTiles[tile.path] += 1
                if r.usedTiles[tile.path] == r.reset {
                    r.usedTiles[tile.path] = 0
                }
            }

        }
        r.m.Unlock()
        outChan <- &Tile{X: grid.X, Y: grid.Y, path: bestTile}
    }
    return
}

const gridSize = 10

func calcGridAvg(x int, y int, img image.Image) (color [3]float64){
    var reds, greens, blues, pixels uint32
    for i := x*gridSize + 0; i <= x*gridSize + 9; i++ {
        for j := y*gridSize + 0; j <= y*gridSize + 9; j++ {
            pixel := img.At(i, j)
            red, green, blue, _ := pixel.RGBA()
            reds += red
            greens += green
            blues += blue
            pixels += 1
        }
    }
    color = [3]float64{float64(reds)/float64(pixels), float64(blues)/float64(pixels), float64(greens)/float64(pixels)}
    return
}

func ImageAssembler(jigsaw chan *Tile, cellsX int, cellsY int, gridSize int, imgSync *sync.WaitGroup, output string) {
    defer imgSync.Done()
    tileSize := 75
    img := image.NewRGBA(image.Rect(0, 0, cellsX*tileSize, cellsX*tileSize))
    for tile := range jigsaw {
        reader, err := os.Open(tile.path)
        if err != nil {
            fmt.Println(err)
            return
        }
        src, _, err := image.Decode(reader)
        if err != nil {
            fmt.Println(err)
            return
        }
        sr := src.Bounds()
        dp := image.Pt(tile.X*tileSize, tile.Y*tileSize)
        r := sr.Sub(sr.Min).Add(dp)
        draw.Draw(img, r, src, sr.Min, draw.Src)
        reader.Close()
    }
    fmt.Println(output)
    out, err := os.Create(output)
    if err != nil {
        fmt.Println(err)
        return
    }
    err = png.Encode(out, img)
    if err != nil {
        fmt.Println(err)
        return
    }
    out.Close()
}

func AnalyzeTarget(inp, tileDataString, output string ) (err error) {
    var tileData map[string][3]float64
    err = json.Unmarshal([]byte(tileDataString), &tileData)
    if err != nil {
        return
    }
    reader, err := os.Open(inp)
    if err != nil {
        return
    }
    image, _, err := image.Decode(reader)
    if err != nil {
        return
    }
    bounds := image.Bounds()
    cellsX := (bounds.Max.X + 1)/gridSize
    cellsY := (bounds.Max.Y + 1)/gridSize
    gridChan := make(chan *Grid)
    outChan := make(chan *Tile)
    repCheck := &RepeatCheck{m: &sync.Mutex{}, usedTiles: make(map[string]int), reset: resetReps}
    var imgSync sync.WaitGroup
    go ImageAssembler(outChan, cellsX, cellsY, gridSize, &imgSync, output)
    var wg sync.WaitGroup
    imgSync.Add(1)
    workers := 10
    wg.Add(workers)
    for i := 0; i < workers; i++ {
        go TileFinder(gridChan, tileData, outChan, repCheck, &wg)
    }
    for x := 0; x < cellsX; x++ {
        for y := 0; y < cellsY; y++ {
            gridColor := calcGridAvg(x, y, image)
            gridChan <- &Grid{X:x, Y:y, Color: gridColor}
            //FindBestTile(gridColor, tileData)
        }
    }
    close(gridChan)
    wg.Wait()
    close(outChan)
    imgSync.Wait()
    return
}

