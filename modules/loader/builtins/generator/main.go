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
	defer fh.Close()

	fh.WriteString("package builtins\n\nimport (\n")
	for _, dirent := range dirents {
		if !dirent.IsDir() {
			continue
		}

		fh.WriteString("\t_ \"github.com/Doridian/fox/modules/" + dirent.Name() + "\"\n")
	}
	fh.WriteString(")\n")
}
