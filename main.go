package main

import (
	"L0/run"
	"log"
)

func main() {
	err := run.Run()
	if err != nil {
		log.Fatal(err)
	}
}
