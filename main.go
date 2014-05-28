package main

import (
    "code.google.com/p/go.net/websocket"
    "net/http"
    "io"
    "io/ioutil"
    "fmt"
)

const Uri = "127.0.0.1:8080"

type WaitedClient struct {
    C *websocket.Conn
    W chan int
}

var Clients = make(chan WaitedClient)

func ChatHandler(client *websocket.Conn) {
    wait := make(chan int)
    Clients <- WaitedClient{client, wait}
    <-wait
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("index.html")
    fmt.Fprintf(w, string(body), "ws://" + Uri + "/ws")
}

func WaitedCopy(first, second WaitedClient) {
    io.Copy(first.C, second.C)
    fmt.Fprintf(first.C, "[!] Disconnected")
    first.W <- 0
}

func Communicate(first, second WaitedClient) {
    go WaitedCopy(first, second)
    go WaitedCopy(second, first)
}

func ChooseTwo() {
    for {
        first := <-Clients
        fmt.Fprintf(first.C, "[!] Wait...")
        second := <-Clients
        fmt.Fprintf(first.C, "[!] Go!")
        fmt.Fprintf(second.C, "[!] Go!")
        go Communicate(first, second)
    }
}

func main() {
    http.Handle("/ws", websocket.Handler(ChatHandler))
    http.HandleFunc("/", IndexHandler)
    go ChooseTwo()
    http.ListenAndServe(Uri, nil)
}
