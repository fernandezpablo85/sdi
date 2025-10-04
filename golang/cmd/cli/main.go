package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello,from cli!")
	if len(os.Args) > 1 {
		fmt.Println("Args:", os.Args[1:])
	}
}
