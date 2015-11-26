package mosaic

import (
    "runtime"
    "strings"
    "sync"
    _ "fmt"
    "path/filepath"
    "os"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "encoding/json"
)

const workers = 10

func AnalyzeImage(inp string) (color [3]float64, err error) {
    reader, err := os.Open(inp)
    if err != nil {
        return
    }
    defer reader.Close()
    image, _, err := image.Decode(reader)
    if err != nil {
        return
    }
    bounds := image.Bounds()
    var reds uint32
    var greens uint32
    var blues uint32
    var pixels uint32
    for i := 0; i <= bounds.Max.X; i++ {
        for j := 0; j <= bounds.Max.Y; j++ {
            pixel := image.At(i, j)
            red, green, blue, _ := pixel.RGBA()
            reds += red
            blues += blue
            greens += green
            pixels += 1
        }
    }
    color = [3]float64{float64(reds)/float64(pixels), float64(blues)/float64(pixels), float64(greens)/float64(pixels)}
    return
}


func analyzer(inpChan chan string, resultChan chan [3]float64) {
    for inpPath := range inpChan {
        color, _ := AnalyzeImage(inpPath)
        resultChan <- color
    }
    return
}

type AnalysisResult struct {
    Path string
    Color [3]float64
}


func AnalyzeAll(inputDir string) string {
    runtime.GOMAXPROCS(4)
    analysis := make(map[string][3]float64)
    inpChan := make(chan string)
    resultChan := make(chan *AnalysisResult)
    var wg sync.WaitGroup
    wg.Add(workers)
    for i:=0; i < workers; i++ {
        go func() {
            defer wg.Done()
            for inpPath := range inpChan {
                color, _ := AnalyzeImage(inpPath)
                resultChan <- &AnalysisResult{Color: color, Path: inpPath}
            }
            return
        }()
    }
    go func() {
        for result := range resultChan {
            analysis[result.Path] = result.Color
        }
    }()
    var walker = func(path string, info os.FileInfo, err error) (outErr error) {
        if ! info.IsDir() && info.Size() > 0 && ! strings.HasPrefix(info.Name(), ".") {
            inpChan <- path
            //fmt.Println(path)
        }
        return
    }
    filepath.Walk(inputDir, walker)
    close(inpChan)
    wg.Wait()
    close(resultChan)
    analysisJson, err := json.Marshal(analysis)
    if err != nil {
        panic("Conversion to json failed")
    }
    return string(analysisJson)
}

