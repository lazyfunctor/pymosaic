package client

import (
    "net/http"
    "net/url"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "path/filepath"
    "io"
    "os"
    "strings"
    "sync"
    )


const (
    searchURL = "https://api.flickr.com/services/rest/?method=flickr.photos.search&license=4&extras=usage&format=json&nojsoncallback=1&tags=%s&api_key=%s&page=%d"
    workers = 15
)

type Result struct {
    Details struct {
        Page int `json:"page"`
        Pages int `json:"pages"`
        Total string `json:"total"`
        PerPage int `json: "perpage"`
        Photos []struct {
            ID string `json:"id"`
            Owner string `json:"owner"`
            Secret string `json:"secret"`
            Title string `json:"title"`
            CanDownload int `json:"can_download"`
        } `json:"photo"`
    } `json:"photos"`
}

type SizeResult struct {
    Details struct {
        CanDownload int `json:"can_download"`
        Sizes []struct {
            Label string `json:"label"`
            // Width string `json:"width"`
            // Height string `json:"height"`
            Source string `json:"source"`
        } `json:"size"`
    } `json:"sizes"`
}


func downloadImage(photoID, outputDir, apiKey string) (err error) {
        sizeURL := "https://api.flickr.com/services/rest/?method=flickr.photos.getSizes&api_key=%s&format=json&nojsoncallback=1&photo_id=%s"
        downURL := fmt.Sprintf(sizeURL, apiKey, photoID)
        resp, err := http.Get(downURL)
        if err != nil {
            return
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        // fmt.Println(string(body))
        var sr SizeResult
        err = json.Unmarshal(body, &sr)
        if err != nil {
            return
        }
        //fmt.Println(sr)
        var source string
        for _, size := range sr.Details.Sizes {
            if size.Label == "Square" {
                source = size.Source
                break
            }
        }
        fileURL, err := url.Parse(source)
        if err != nil {
            return
        }
        path := fileURL.Path
        segments := strings.Split(path, "/")
        fmt.Println(path)
        fmt.Println(segments[len(segments) - 1])
        fileName := filepath.Join(outputDir, segments[len(segments) - 1])
        file, err := os.Create(fileName)
        if err != nil {
            return
        }
        defer file.Close()


        fileResp, err := http.Get(source)
        if err != nil {
            return
        }
        defer fileResp.Body.Close()
        _, err = io.Copy(file, fileResp.Body)
        return
}

func worker(downloads chan string, outputDir string, wg *sync.WaitGroup, errChan chan error, apiKey string) {
    defer wg.Done()
    for photoID := range downloads {
        err := downloadImage(photoID, outputDir, apiKey)
        if err != nil {
            //fmt.Println(err)
            errChan <- err
        }

    }
}

func monitor(errChan chan error) {
    for err := range errChan {
        fmt.Println(err)
    }
}

func Download(tag string, reqdCount int, outputDir string, apiKey string) {
    downloads := make(chan string)
    // outputDir := "/home/abhinav/dev/mosaic/tiles"
    errChan := make(chan error)
    go monitor(errChan)
    var wg sync.WaitGroup
    wg.Add(workers)
    for i:=0; i < workers; i++ {
        go worker(downloads, outputDir, &wg, errChan, apiKey)
    }

    pager := 1
    count := 0
    finish := false
    for ! finish {
        apiUrl := fmt.Sprintf(searchURL, tag, apiKey, pager)
        fmt.Println(apiUrl)
        resp, err := http.Get(apiUrl)
        if err != nil {
            panic("get failed")
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        var r Result
        //fmt.Println(string(body))
        err = json.Unmarshal(body, &r)
        if err != nil {
            fmt.Println(err)
            panic("Error")
        }
        for _, photo := range r.Details.Photos {
            if photo.CanDownload == 1 {
                count += 1
                //fmt.Println(photo)
                downloads <- photo.ID
            }
            if count >= reqdCount {
                finish = true
                break
            }
        }
        pager += 1
    }
    close(downloads)
    wg.Wait()
    close(errChan)
}

