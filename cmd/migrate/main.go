package main

import (
	"log"

	"github.com/cristiansrc/hv-go-ms-resume/internal/migration"
)

func main() {
	log.Println("hv-go-ms-resume data migration tool")
	migration.Run()
}
