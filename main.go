package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		if err := Serve(); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(time.Second)
	if err := Request(); err != nil {
		fmt.Println(err)
	}
}
