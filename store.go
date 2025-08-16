/*
 * Copyright 2024 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
Package barky provides hierarchical configuration storage and path parsing utilities.

Features:
- Structured key-value storage with nested path support and conflict detection.
- Path parsing (SplitPath) and construction (JoinPath).
- Two path types:
  - Key (e.g., "user.name") for map access
  - Index (e.g., "[0]") for list access

- Maintains a tree structure (treeNode) for type-safe hierarchy management.

Use cases:
- Accessing values in JSON/YAML/TOML-like configs
- Managing nested configuration data
- Validating structure and detecting conflicts

Notes:
- Path syntax follows common patterns (e.g., "users[0].profile.age").
- Strict distinction between key and index segments.
*/
package barky

import (
	"errors"
	"fmt"
	"sort"
)

// treeNode represents a node in the hierarchical key path tree.
// Each node tracks the type of its path segment and its child nodes.
type treeNode struct {
	Type PathType
	Data map[string]*treeNode
}

// ValueInfo stores a value and the file it came from.
type ValueInfo struct {
	File  int8
	Value string
}

// Storage maintains hierarchical key-value pairs using a tree structure.
// - `data` stores flat key-value mappings.
// - `file` maps filenames to IDs.
// - `root` tracks the hierarchical path tree.
// It supports nested paths and detects type conflicts.
type Storage struct {
	root *treeNode            // Root node of the hierarchy tree.
	data map[string]ValueInfo // Flat key-value map.
	file map[string]int8      // File-to-ID mapping.
}

// NewStorage creates a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]ValueInfo),
		file: make(map[string]int8),
	}
}

// RawData returns the internal key-value map.
// Warning: exposes internal state directly.
func (s *Storage) RawData() map[string]ValueInfo {
	return s.data
}

// AddFile adds a file to the storage and returns its ID.
func (s *Storage) AddFile(file string) int8 {
	idx, ok := s.file[file]
	if !ok {
		idx = int8(len(s.file))
		s.file[file] = idx
	}
	return idx
}

// RawFile returns the raw file-to-ID map.
func (s *Storage) RawFile() map[string]int8 {
	return s.file
}

// SubKeys returns the immediate subkeys under the given key path.
// It walks the tree structure and returns child elements if the path exists.
// Returns an error if there's a type conflict along the path or if the path is not found.
func (s *Storage) SubKeys(key string) (_ []string, err error) {
	var path []Path
	if key != "" {
		if path, err = SplitPath(key); err != nil {
			return nil, err
		}
	}

	if s.root == nil {
		return nil, nil
	}

	n := s.root
	for i, pathNode := range path {
		if n == nil || pathNode.Type != n.Type {
			return nil, fmt.Errorf("property conflict at path %s", JoinPath(path[:i+1]))
		}
		v, ok := n.Data[pathNode.Elem]
		if !ok {
			return nil, nil
		}
		n = v
	}

	if n == nil {
		return nil, fmt.Errorf("property path %s not found", key)
	}

	keys := make([]string, 0, len(n.Data))
	for k := range n.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

// Has checks if the given key exists in storage,
// either as a direct value or as a valid path in the hierarchy.
func (s *Storage) Has(key string) bool {
	if key == "" || s.root == nil {
		return false
	}

	if _, ok := s.data[key]; ok {
		return true
	}

	path, err := SplitPath(key)
	if err != nil {
		return false
	}

	n := s.root
	for _, node := range path {
		if n == nil || node.Type != n.Type {
			return false
		}
		v, ok := n.Data[node.Elem]
		if !ok {
			return false
		}
		n = v
	}
	return true
}

// Set inserts a key-value pair into storage.
// Builds/extends the tree structure as needed.
// Returns an error if the key is empty or if a type conflict occurs.
func (s *Storage) Set(key string, val string, file int8) error {
	if key == "" {
		return errors.New("key is empty")
	}

	path, err := SplitPath(key)
	if err != nil {
		return err
	}

	// Initialize root if empty
	if s.root == nil {
		s.root = &treeNode{
			Type: path[0].Type,
			Data: make(map[string]*treeNode),
		}
	}

	n := s.root
	for i, pathNode := range path {
		if n == nil || pathNode.Type != n.Type {
			return fmt.Errorf("property conflict at path %s", JoinPath(path[:i+1]))
		}
		v, ok := n.Data[pathNode.Elem]
		if !ok {
			if i < len(path)-1 {
				v = &treeNode{
					Type: path[i+1].Type,
					Data: make(map[string]*treeNode),
				}
			}
			n.Data[pathNode.Elem] = v
		}
		n = v
	}
	if n != nil {
		return fmt.Errorf("property conflict at path %s", key)
	}

	s.data[key] = ValueInfo{file, val}
	return nil
}
