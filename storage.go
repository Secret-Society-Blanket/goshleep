package main

import (
  "fmt"
  "io"
  "io/ioutil"
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
  writeToFile("stored.yaml", string(bytes))
  log.Println(fi)
  log.Println(string(bytes))

  return (prefix + suffix)

}

func Load(out *[]Verb) {
  txt := loadFromFile("stored.yaml")
  if err := yaml.Unmarshal(txt, out); err != nil {
    fmt.Println(err)
  }
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
