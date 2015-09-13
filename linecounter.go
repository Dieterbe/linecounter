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

func init() {
	flag.IntVar(&freq, "freq", 1000, "report frequency in ms")
}

func main() {
	flag.Parse()
	tracker()
}

func tracker() {
	data := make(chan int)
	exit := make(chan int)
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			l, err := reader.ReadString('\n')
			if err == io.EOF {
				if len(l) > 0 {
					data <- 1
				}
				exit <- 0
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				exit <- 2
			}
			data <- 1
		}
	}()
	counter := 0
	tick := time.NewTicker(time.Duration(freq) * time.Millisecond)
	for {
		select {
		case <-tick.C:
			fmt.Println(counter)
			counter = 0
		case ret := <-exit:
			fmt.Println(counter, "(incomplete interval)")
			os.Exit(ret)
		case diff := <-data:
			counter += diff
		}
	}
}
