package main

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	repo := flag.String("repo", "~/.minirepo/repo", "Minirepo repository directory")
	bind := flag.String("bind", "127.0.0.1:8080", "address to listen on")
	flag.Parse()

	repoDir, err := homedir.Expand(*repo)
	if err != nil {
		log.WithError(err).WithField("root", *repo).Fatal("Couldn't expand repo directory")
	}

	http.Handle("/", http.FileServer(http.Dir(repoDir)))
	if err := http.ListenAndServe(*bind, nil); err != nil {
		panic(err)
	}
}
