package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "strconv"
)

func main() {
    start := time.Now()
    ch := make(chan string)
    var count int64 = 1
    times,_ := strconv.ParseInt(os.Args[1],10,64)
    url := os.Args[2]
    for count = 0; count <= times; count++ {
        go fetch(url, ch, count) // start a goroutine
    }

    for count = 0; count <= times; count++ {
        fmt.Println(<-ch) // receive from channel ch
    }
    fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string, count int64) {
    start := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        ch <- fmt.Sprint(err) // send to channel ch
        return
    }
    nbytes, err := io.Copy(ioutil.Discard, resp.Body)
    resp.Body.Close() // don't leak resources
    if err != nil {
        ch <- fmt.Sprintf("while reading %s: %v", url, err)
        return
    }
    secs := time.Since(start).Seconds()
    ch <- fmt.Sprintf("%7d %.2fs  %7d  %s",count, secs, nbytes, url)
 }
