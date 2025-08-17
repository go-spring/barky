# Barky

<div>
   <img src="https://img.shields.io/github/license/go-spring/barky" alt="license"/>
   <img src="https://img.shields.io/github/go-mod/go-version/go-spring/barky" alt="go-version"/>
   <img src="https://img.shields.io/github/v/release/go-spring/barky?include_prereleases" alt="release"/>
   <a href="https://codecov.io/gh/go-spring/barky" > 
      <img src="https://codecov.io/gh/go-spring/barky/graph/badge.svg?token=SX7CV1T0O8" alt="test-coverage"/>
   </a>
   <a href="https://deepwiki.com/go-spring/barky"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki"></a>
</div>

[English](README.md) | [中文](README_CN.md)

`Barky` is a Go package for managing hierarchical key-value data structures, mainly designed for configuration files
such as JSON, YAML, and TOML.
It can flatten nested data structures into a `map[string]string` while preserving path information, detecting conflicts,
and tracking multi-file sources.

## ✨ Key Features

### 1. Flattening

* Supports converting nested `map`, `slice`, and `array` into a flat `map[string]string`.
* Uses **dots** to denote map keys and **brackets** to denote array/slice indices.

### 2. Path Handling

* The `Path` type parses hierarchical keys into a sequence of **path segments** (map keys or array indices).
* Provides `SplitPath` to convert `"foo.bar[0]"` into structured path segments, and `JoinPath` to join path segments
  back into a string.
* Performs syntax validation to prevent invalid keys (e.g., consecutive dots, unclosed brackets).

### 3. Storage

* The `Storage` type manages a set of flattened key-value pairs and builds an internal tree to detect structural
  conflicts.
* Each value is associated with a source file index, making it easy to track origins when merging multiple files.
* Key methods include:

    * `Set`: Set a key-value pair with conflict detection.
    * `Get`: Retrieve a value, optionally providing a default.
    * `Has`: Check whether a key exists.
    * `SubKeys`: List subkeys under a given path.
    * `Keys`: Get all stored keys in sorted order.

## 📦 Installation

```bash
go get github.com/go-spring/barky
```

## 🛠 Usage Example

```go
package main

import (
	"fmt"
	"github.com/go-spring/barky"
)

func main() {
	s := barky.NewStorage()
	fileIdx := s.AddFile("config.yaml")

	_ = s.Set("server.hosts[0].ip", "192.168.0.1", fileIdx)
	_ = s.Set("server.hosts[1].ip", "192.168.0.2", fileIdx)

	fmt.Println("Keys:", s.Keys())
	fmt.Println("SubKeys of server.hosts:", s.SubKeys("server.hosts"))
	fmt.Println("Get server.hosts[0].ip:", s.Get("server.hosts[0].ip"))
}
```

Output:

```
Keys: [server.hosts[0].ip server.hosts[1].ip]
SubKeys of server.hosts: [0 1]
Get server.hosts[0].ip: 192.168.0.1
```

## 📖 Use Cases

* **Configuration Management**
  Convert configuration files of various formats into a unified flat key-value map for easier comparison and merging.

* **Querying & Retrieval**
  Access nested data directly using simple path strings like `"server.hosts[0].ip"`, without manually traversing the
  structure.

* **Multi-File Merging**
  Handle multiple configuration files at once, track the source of each key, and detect conflicts.

* **Data Transformation**
  Flatten hierarchical structures for easier processing in testing, diffing, or downstream systems.

## ⚠️ Notes

* Paths must not contain **spaces**, **consecutive dots**, **unclosed brackets**, or other invalid formats.
* Type conflicts in the tree structure (e.g., `"user.name"` vs `"user[0]"`) will return an error.
* `RawData` and `RawFile` expose internal storage directly—use with caution.

## 📜 License

Apache 2.0 License