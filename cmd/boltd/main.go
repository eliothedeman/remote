package main

import (
	"flag"
	"log"

	"github.com/eliothedeman/remote"
)

var (
	listen = flag.String("listen", ":9988", "Address to listen on.")
	path   = flag.String("path", "bolt.db", "Path of the boltdb file.")
)

func main() {
	flag.Parse()

	s, err := remote.OpenServer(*path)
	if err != nil {
		panic(err)
	}

	log.Println("Opened bolt database at path", *path)

	log.Println("Serving bolt database at", *listen)
	err = s.ServeTCP(*listen)
	if err != nil {
		panic(err)
	}
}
