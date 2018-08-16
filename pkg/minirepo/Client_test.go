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

package minirepo

import (
	"testing"
	"os"
	"path"
	"github.com/uubk/minirepo/internal/minirepo"
	"io/ioutil"
	"crypto/rand"
	"net/http"
	"net"
	"golang.org/x/sys/unix"
	"strings"
)

func generateTestAssets(dir string) {
	err := os.MkdirAll(path.Join(dir, "repo"), 0700)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(path.Join(dir, "repo", "a_dir"), 0700)
	if err != nil {
		panic(err)
	}

	repoRoot := path.Join(dir, "repo")
	svc := minirepo.NewServer(dir, repoRoot, "Unittest Server")
	svc.GenerateKeypair()
	svc.LoadKeypair()

	randomData := make([]byte, 256)
	rand.Read(randomData)
	ioutil.WriteFile(path.Join(repoRoot, "a_dir", "testfile"), randomData, 0700)
	svc.UpdateMetadata()
}

func provideTestServer(root string) (string, *http.Server) {
	http.Handle("/", http.FileServer(http.Dir(root)))
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	srv := &http.Server {
	}
	go func () {
		srv.Serve(listener)
	} ()

	return listener.Addr().String(), srv
}

func InitTestClient() (*Minirepo, *http.Server, error) {
	testPath, err := ioutil.TempDir("", "minirepo-unittest")
	if err != nil {
		return nil, nil, err
	}
	generateTestAssets(testPath)
	pubkeyBin, err := ioutil.ReadFile(path.Join(testPath, "pub.asc"))
	if err != nil {
		return nil, nil, err
	}
	bindAddr, server := provideTestServer(path.Join(testPath, "repo"))

	clientPath := path.Join(testPath, "client")
	os.Mkdir(clientPath, 0700)
	client, err := NewRepoClient(clientPath, "http://" +bindAddr, string(pubkeyBin))
	if err != nil {
		return nil, nil, err
	}


	return client, server, nil
}

func TestMetaDownload(t *testing.T) {
	_, server, err := InitTestClient()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	server.Shutdown(nil)
}

func TestFileDownload(t *testing.T) {
	client, server, err := InitTestClient()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer server.Shutdown(nil)

	filePath, err := client.GetFile("a_dir", "testfile")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	if filePath == "" {
		t.Fatal("Unexpectedly empty path")
	}
}

func TestFileDownloadFresh(t *testing.T) {
	client, server, err := InitTestClient()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer server.Shutdown(nil)

	refreshed, filePath, err := client.GetFileLatest("a_dir", "testfile")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	if filePath == "" {
		t.Fatal("Unexpectedly empty path")
	}
	if refreshed {
		t.Fatal("First download shouldn't have 'refreshed' flag set")
	}

	refreshed, filePath, err = client.GetFileLatest("a_dir", "testfile")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	if filePath == "" {
		t.Fatal("Unexpectedly empty path")
	}
	if !refreshed {
		t.Fatal("Second download should have 'refreshed' flag set")
	}
}

func TestFileDownloadErrors(t *testing.T) {
	err, client, server := InitTestClient()
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	defer server.Shutdown(nil)

	filePath, err := client.GetFile("a_dir", "no_file")
	if err == nil {
		t.Fatal("Expected error missing")
	} else if err.Error() != "file not found" {
		t.Fatal("Unxpected error: ", err)
	}
	if filePath != "" {
		t.Fatal("Unexpectedly got a path?")
	}

	filePath, err = client.GetFile("no_dir", "no_file")
	if err == nil {
		t.Fatal("Expected error missing")
	} else if err.Error() != "file not found" {
		t.Fatal("Unxpected error: ", err)
	}
	if filePath != "" {
		t.Fatal("Unexpectedly got a path?")
	}

	filePath, err = client.GetFile()
	if err == nil {
		t.Fatal("Expected error missing")
	} else if err.Error() != "no path specified" {
		t.Fatal("Unxpected error: ", err)
	}
	if filePath != "" {
		t.Fatal("Unexpectedly got a path?")
	}

	filePath, err = client.GetFile("", "")
	if err == nil {
		t.Fatal("Expected error missing")
	} else if err.Error() != "invalid path part: empty string" {
		t.Fatal("Unxpected error: ", err)
	}
	if filePath != "" {
		t.Fatal("Unexpectedly got a path?")
	}

	// Make directory unwritable
	// First try should work
	filePath, err = client.GetFile("a_dir", "testfile")
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}
	if filePath == "" {
		t.Fatal("Unexpectedly empty path")
	}
	os.Chmod(path.Dir(filePath), unix.O_RDONLY)
    // Second try shouldn't
	_, filePath, err = client.GetFileLatest("a_dir", "testfile")
	if err == nil {
		t.Fatal("Expected error missing")
	} else if !strings.HasSuffix(err.Error(), "a_dir/testfile: permission denied") {
		t.Fatal("Unxpected error: ", err)
	}
	if filePath == "" {
		t.Fatal("Unexpectedly empty path")
	}

	// Directory should be missing and also not be creatable
	client.localCache = "/proc"
	filePath, err = client.GetFile("a_dir", "testfile")
	if err == nil {
		t.Fatal("Expected error missing")
	} else if err.Error() != "open /proc/a_dir/testfile: no such file or directory" {
		t.Fatal("Unxpected error: ", err)
	}
	if filePath == "" {
		t.Fatal("Unexpectedly empty path")
	}
}