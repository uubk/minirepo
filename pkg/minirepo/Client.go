package minirepo

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"github.com/uubk/minirepo/pkg/minirepo/types"
)

// Minirepo client
type Minirepo struct {
	// Local cache directory
	localCache string
	// Remote origin URL
	remote     string
	// ASCII-armored pubkey for signatures
	signingKey string
	// Parsed repository metadata, if available
	meta       *types.RepoInfo
}

// Creates a new minirepo client.
//  - localCache is expected to contain a directory where the repository downloads should be cached
//  - url is expected to contain a repositories upstream
//  - key is expected to contain a full ASCII-armored GPG public key which was used for signing the metadata
// This method will also fetch, validate and decode metadata
func NewRepoClient(localCache, url, key string) (*Minirepo, error) {
	obj := Minirepo{
		localCache: localCache,
		remote:     url,
		signingKey: key,
	}
	err := obj.fetchMeta()
	if err != nil {
		return nil, err
	}
	err = obj.decodeMeta()
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// Returns the *latest* version of a file, that is, it deletes a local copy before download, should it exist
func (m *Minirepo) GetFileLatest(filePath ...string) (bool, string, error) {
	fileFullPath := []string{m.localCache}
	fileFullPath = append(fileFullPath, filePath...)
	fileRef := path.Join(fileFullPath...)
	_, err := os.Stat(fileRef)
	if err == nil {
		err = os.Remove(fileRef)
		if err != nil {
			return true, fileRef, errors.New("deletion of old file failed")
		}
		fileRef, err := m.GetFile(filePath...)
		return true, fileRef, err
	} else {
		fileRef, err := m.GetFile(filePath...)
		return false, fileRef, err
	}
}

// Returns a local path to the file requested if possible
func (m *Minirepo) GetFile(filePath ...string) (string, error) {
	if m.meta == nil {
		return "", errors.New("no metadata available")
	}

	if len(filePath) == 0 {
		return "", errors.New("no path specified")
	}

	var curEntry *types.DirEntry
	for _, item := range filePath {
		if item == "" {
			return "", errors.New("invalid path part: empty string")
		}
 		if curEntry == nil {
			for _, otherItem := range m.meta.Contents {
				if otherItem.Name == item {
					curEntry = &otherItem
					break
				}
			}
			if curEntry == nil {
				// Didn't find anything
				return "", errors.New("file not found")
			}
		} else {
			found := false
			for _, otherItem := range curEntry.Children {
				if otherItem.Name == item {
					curEntry = &otherItem
					found = true
					break
				}
			}
			if !found {
				return "", errors.New("file not found")
			}
		}
	}
	// curEntry contains full metadata and the file should exist

	// Did we already fetch this file?
	fileFullPath := []string{m.localCache}
	fileFullPath = append(fileFullPath, filePath...)
	fileRef := path.Join(fileFullPath...)
	_, err := os.Stat(fileRef)
	if err == nil {
		return fileRef, nil
	}

	// Nope, file does not exist
	fileUrl := m.remote
	for _, segment := range filePath {
		fileUrl += "/" + url.PathEscape(segment)
	}
	response, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("file download failed: %s", err)
	}
	defer response.Body.Close()
	os.MkdirAll(path.Dir(fileRef), 0700)
	fileContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("file download failed: %s", err)
	}

	hashSum := hex.EncodeToString(sha256.New().Sum(fileContent))
	if hashSum != curEntry.Hash || hashSum == "" {
		return "", fmt.Errorf("checksum mismatch")
	}
	return fileRef, ioutil.WriteFile(fileRef, fileContent, 0600)
}

// Decode a local copy of the metadata file
func (m *Minirepo) decodeMeta() error {
	metaBin, err := ioutil.ReadFile(path.Join(m.localCache, "meta.yml"))
	if err != nil {
		return fmt.Errorf("metadata read failed: %s", err)
	}
	m.meta = &types.RepoInfo{}
	return yaml.Unmarshal(metaBin, &m.meta)
}

// Fetch current metadata and write it to disk if and only if the signature is valid
func (m *Minirepo) fetchMeta() error {
	response, err := http.Get(m.remote + "/meta.yml")
	if err != nil {
		return fmt.Errorf("metadata download failed: %s", err)
	}
	defer response.Body.Close()
	metaYml, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("metadata download failed: %s", err)
	}
	response, err = http.Get(m.remote + "/meta.asc")
	if err != nil {
		return fmt.Errorf("metadata download failed: %s", err)
	}
	defer response.Body.Close()
	metaAsc, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("metadata download failed: %s", err)
	}

	// We have both metadata and signature. Verify signature before opening metadata file!!
	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(m.signingKey))
	if err != nil {
		return fmt.Errorf("keyring decode failed: %s", err)
	}
	_, err = openpgp.CheckArmoredDetachedSignature(keyring, bytes.NewReader(metaYml), bytes.NewReader(metaAsc))
	if err != nil {
		return fmt.Errorf("signature invalid or check failed: %s", err)
	}

	return ioutil.WriteFile(path.Join(m.localCache, "meta.yml"), metaYml, 0600)
}
