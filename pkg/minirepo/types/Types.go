package types

import "time"

type DirEntry struct {
	Name     string
	Hash     string     `yaml:"hash,omitempty"`
	Children []DirEntry `yaml:"children,omitempty"`
}

type RepoInfo struct {
	Contents  []DirEntry
	Name      string
	Timestamp time.Time
}
