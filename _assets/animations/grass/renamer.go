package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const baseDirectory = "/Users/kevinchou/go/src/github.com/kkevinchou/kito/assets/animations"

func main() {
	directory := filepath.Join(baseDirectory, "grass")
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return
	}

	counter := 0

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			newName := fmt.Sprintf("%d.png", counter)
			err := os.Rename(filepath.Join(directory, file.Name()), filepath.Join(directory, newName))
			if err != nil {
				fmt.Println(err)
				return
			}
			counter += 1
		}
	}
}
