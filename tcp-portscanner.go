package main

import (
  "fmt"
  "net"
  "sync"
  "time"
  "flag"
  "github.com/y3a/tcp-portscanner/matchname"
)

func worker(pChannel chan int, wg *sync.WaitGroup, ip string, timeout int) {
  for port := range pChannel {
    addr := fmt.Sprintf("%s:%d", ip, port)
    d := net.Dialer{Timeout: time.Duration(timeout) * time.Millisecond}
    conn, err := d.Dial("tcp", addr)
    if err == nil {
      detail := matchname.Parse(port)
      fmt.Printf("[+]Open: %s -> %s\n", addr, detail)
      conn.Close()
    }
    wg.Done()
  }
}

func main() {
  ip := flag.String("i", "", "IP Address/Hostname")
  timeout := flag.Int("t", 1000, "Timeout in milliseconds.")
  flag.Parse()
  if *ip == "" {
    flag.PrintDefaults()
    return
  }
  fmt.Printf("[*]Portscan against %s started, with a timeout of %d milliseconds\n", *ip, *timeout)
  var wg sync.WaitGroup
  pChannel := make(chan int, 100)
  for spawned := 0; spawned <= cap(pChannel); spawned ++ {
    go worker(pChannel, &wg, *ip, *timeout)
  }
  for port := 0; port <= 65535; port ++ {
    wg.Add(1)
    pChannel <- port
  }
  wg.Wait()
  close(pChannel)
}
