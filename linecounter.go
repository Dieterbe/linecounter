package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var freq int
var data chan int

func init() {
	flag.IntVar(&freq, "freq", 1000, "report frequency in ms")
	data = make(chan int)
}

func main() {
    flag.Parse()
	go tracker()
	reader()
}

func reader() {
	reader := bufio.NewReader(os.Stdin)
	for {
		_, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
		data <- 1
	}
}

func tracker() {
	counter := 0
	tick := time.NewTicker(time.Duration(freq) * time.Millisecond)
    for {
        select {
        case <-tick.C:
            fmt.Println(counter)
            counter = 0
        case diff := <-data:
            counter += diff
        }
    }
}
