package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

// Store a collection of verbs
func Store(toStore []Verb) string {
	// fmt.Println("Storing...")
	bytes, err := yaml.Marshal(toStore)

	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println("Storing...")
	suffix := "dump.yaml"
	t := time.Now()
	prefix := t.Format("15:04-2-1-06|")
	fi := writeToFile(prefix+suffix, string(bytes))
	log.Println(fi)
	log.Println(string(bytes))

	return (prefix + suffix)

}

func Load(toLoad string) []Verb {

}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}

	return file.Sync()
}
