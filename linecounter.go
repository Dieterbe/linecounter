package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

var freq int
var field int

func init() {
	flag.IntVar(&freq, "freq", 1000, "report frequency in ms (default:1000)")
	flag.IntVar(&field, "f", -1, "seggregate by given whitespace separated field value. fields are numbered from 1 like cut and awk. (default -1 to disable)")
}

func main() {
	flag.Parse()
	if field == 0 {
		fmt.Fprintln(os.Stderr, "fields are numbered from 1")
		os.Exit(2)
	}
	if field < 0 {
		lineTracker()
	} else {
		seggregatedLineTracker(field - 1)
	}
}

func lineTracker() {
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

func seggregatedLineTracker(index int) {
	data := make(chan string)
	exit := make(chan int)
	reader := bufio.NewReader(os.Stdin)
	counter := make(map[string]int)

	processLine := func(line string) {
		fields := strings.Fields(line)
		if index < len(fields) {
			data <- fields[index]
		} else {
			data <- "null"
		}
	}

	printCounters := func(header string) {
		fmt.Println(header)
		keys := make([]string, 0, len(counter))
		for field := range counter {
			keys = append(keys, field)
		}
		sort.Strings(keys)
		for _, field := range keys {
			fmt.Println(field, counter[field])
		}
	}

	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				if len(line) > 0 {
					processLine(line)
				}
				exit <- 0
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				exit <- 2
			}
			processLine(line)
		}
	}()
	tick := time.NewTicker(time.Duration(freq) * time.Millisecond)
	for {
		select {
		case <-tick.C:
			printCounters("===============================")
			newCounter := make(map[string]int)
			for key := range counter {
				newCounter[key] = 0
			}
			counter = newCounter
		case ret := <-exit:
			printCounters("==== (incomplete interval) ====")
			os.Exit(ret)
		case str := <-data:
			counter[str] += 1
		}
	}
}
