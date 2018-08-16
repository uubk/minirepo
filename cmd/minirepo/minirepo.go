/*
 * Copyright 2018 The minirepo authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	minirepo2 "github.com/uubk/minirepo/internal/minirepo"
	"os"
	"path"
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
