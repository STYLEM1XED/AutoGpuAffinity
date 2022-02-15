package main

import (
	"log"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createFile(file, data string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)
	if err2 != nil {
		panic(err)
	}
}
