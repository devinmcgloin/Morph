package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devinmcgloin/morph/src/api/image"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File provided")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(image.GetExif(f))
}
