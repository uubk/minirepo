package main

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	minirepo2 "github.com/uubk/minirepo/internal/minirepo"
)

func main() {
	verbose := flag.Bool("verbose", true, "Enable verbose output")
	root := flag.String("root", "~/.minirepo", "Minirepo root directory")
	repo := flag.String("repo", "~/.minirepo/repo", "Minirepo repository directory")
	name := flag.String("Name", "minirepo", "Minirepo repository Name")

	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	rootDir, err := homedir.Expand(*root)
	if err != nil {
		log.WithError(err).WithField("root", *root).Fatal("Couldn't expand root directory")
	}
	repoDir, err := homedir.Expand(*repo)
	if err != nil {
		log.WithError(err).WithField("repo", *repo).Fatal("Couldn't expand repo directory")
	}

	os.Mkdir(rootDir, 0700)

	svc := minirepo2.NewServer(rootDir, repoDir, *name)

	// Ensure public/private keys exists
	pubkeyFile := path.Join(rootDir, "pub.asc")
	_, err = os.Stat(pubkeyFile)
	if err != nil {
		// File does not exist -> generate new keys
		svc.GenerateKeypair()
	}
	log.Info("Loading keys")
	svc.LoadKeypair()
	log.Info("Updating metadata")
	svc.UpdateMetadata()
}
