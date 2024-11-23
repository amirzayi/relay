package main

import (
	"fmt"

	"github.com/AmirMirzayi/relay/cmd"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("error occured: %v", r)
		}
	}()
	cmd.Execute()
}
