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
