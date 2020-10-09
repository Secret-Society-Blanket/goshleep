package main
import(
  "log"
  "fmt"
  "github.com/goccy/go-yaml"
  "time"
  "os"
  "io"
)

// Store a collection of verbs
func Store(toStore []Verb) string{
  fmt.Println("Storing...")
  bytes, err := yaml.Marshal(toStore)

  if err != nil {
    log.Println(err)
    return ""
  }
  log.Println("Storing...")
  suffix := "dump.yaml"
  prefix := time.Now().Local()
  // fi := WriteToFile(prefix + suffix, string(bytes))
  log.Println(prefix.Format())

  return (prefix + suffix)


}

func WriteToFile(filename string, data string) error {
  file, err := os.Create(filename)
  if err != nil {
    return err
  }
  defer file.Close()

  _, err = io.WriteString(file, data)
  if err!=nil {
    return err
  }

  return file.Sync()
}
