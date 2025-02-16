package main

import (
	"log"
	"os"
)

func main() {
	dirents, err := os.ReadDir("../../")
	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Create("modules.generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = fh.Close()
	}()

	_, err = fh.WriteString("package builtins\n\nimport (\n")
	if err != nil {
		log.Fatal(err)
	}
	for _, dirent := range dirents {
		if !dirent.IsDir() {
			continue
		}

		_, err = fh.WriteString("\t_ \"github.com/Doridian/fox/modules/" + dirent.Name() + "\"\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = fh.WriteString(")\n")
	if err != nil {
		log.Fatal(err)
	}
}
