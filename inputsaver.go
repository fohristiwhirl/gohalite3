package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var input_chan = make(chan string, 100)

func input_reader() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		ok := scanner.Scan()
		input_chan <- scanner.Text()
		if ok == false {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}
}

func main() {

	fmt.Printf("InputSaver")

	outfile, _ := os.Create("inputsaver.txt")
	defer outfile.Close()

	go input_reader()

	last_comms := time.Now()

	for {

		select {

		case line := <- input_chan:

			last_comms = time.Now()
			fmt.Fprintf(outfile, line)
			fmt.Fprintf(outfile, "\n")

		default:

			// Basically the idea is to send \n if there's been no comms in either direction for 20 ms

			time.Sleep(10 * time.Millisecond)

			if time.Now().Sub(last_comms) > (20 * time.Millisecond) {
				fmt.Printf("\n")
				last_comms = time.Now()
			}
		}
	}
}
