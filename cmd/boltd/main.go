package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
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

	logrus.WithFields(logrus.Fields{
		"file_path": *path,
		"address":   *listen,
	}).Info("Starting server")

	err = s.ServeTCP(*listen)
	if err != nil {
		panic(err)
	}
}
