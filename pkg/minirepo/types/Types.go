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

// Package types contains the types visible inside a repository (to avoid cyclic includes)
package types

import "time"

// DirEntry contains a directory entry. This struct represents either
//  - a single file (Name and Hash set) or
//  - a child directory (Name and Children set)
type DirEntry struct {
	// Name of this entry
	Name string
	// If this is a file, contains the SHA-256 hash of it in Hex encoding
	Hash string `yaml:"hash,omitempty"`
	// If this is a directory, contains a list of all children
	Children []DirEntry `yaml:"children,omitempty"`
}

// RepoInfo contains the base repostitory info, that is some metadata and a list of the repositories' contents
type RepoInfo struct {
	// List of content (only directories at this level)
	Contents []DirEntry
	// Name of repository
	Name string
	// Timestamp of last update
	Timestamp time.Time
}
