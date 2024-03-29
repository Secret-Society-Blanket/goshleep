package goshleep

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Store a collection of verbs
func Store(toStore []Verb) string {
	// fmt.Println("Storing...")
	bytes, err := json.Marshal(toStore)

	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println("Storing...")
	suffix := "dump.yaml"
	t := time.Now()
	prefix := viper.GetString("storage") + "backups/"
	prefix += t.Format("15:04-2-1-06|")
	writeToFile(prefix+suffix, string(bytes))
	writeToFile(viper.GetString("storage")+"stored.json", string(bytes))

	return (prefix + suffix)

}

// Load a collection of verbs
func Load(out *[]Verb) {
	txt := loadFromFile(viper.GetString("storage") + "stored.json")
	if err := json.Unmarshal(txt, out); err != nil {fmt.Println(err)}
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

func loadFromFile(filename string) []byte {
	txt, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return txt

}
