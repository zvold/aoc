package util

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
)

var inputFlag = flag.String("input", "", "Filename to read the input from. When set to the special value '-', read from the standard input.")

type CloserFunc func()

func OpenInputFile(f embed.FS) (fs.File, CloserFunc) {
	var file fs.File
	var err error
	switch *inputFlag {
	case "":
		fmt.Println("Reading from embedded input.")
		file, err = f.Open("input-1.txt")
	case "-":
		fmt.Println("Reading from standard input.")
		file = os.Stdin
	default:
		fmt.Printf("Reading from %s.\n", *inputFlag)
		file, err = os.Open(*inputFlag)
	}
	if err != nil {
		log.Fatal(err)
	}

	return file, func() {
		if *inputFlag != "-" {
			err := file.Close()
			if err != nil {
				return
			}
		}
	}
}
