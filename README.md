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

**Barky** is a lightweight library for **hierarchical configuration storage and path parsing**, inspired by
JSON/YAML/TOML access patterns.

It provides **type-safe path handling**, structured storage, subkey lookup, and conflict detection, making it ideal for
managing nested configuration data.

---

## ✨ Features

* **Path Parsing & Construction**

    * `SplitPath`: Parse a string path into structured `Path` segments
    * `JoinPath`: Build a string path from structured segments

* **Two Path Types**

    * **Key**: For map/object access (e.g., `"user.name"`)
    * **Index**: For array/list access (e.g., `"[0]"`)

* **Hierarchical Storage**

    * Maintains a tree structure (`treeNode`) for type-safe paths
    * Query subkeys with `SubKeys`
    * Check if a path exists with `Has`
    * Insert key-value pairs with automatic tree building and conflict detection via `Set`

* **Conflict Detection**

    * Detects type mismatches (e.g., when a path segment is used both as a key and as an index)

* **File Source Tracking**

    * Each value is associated with its source file and a file ID, useful when merging multiple configs

---

## 📦 Installation

```bash
go get github.com/go-spring/barky
```

---

## 🛠 Usage Examples

### Path Parsing & Building

```go
path, _ := barky.SplitPath("users[0].profile.name")
// path => [ {Key:"users"}, {Index:"0"}, {Key:"profile"}, {Key:"name"} ]

joined := barky.JoinPath(path)
// joined => "users[0].profile.name"
```

### Storage Operations

```go
s := barky.NewStorage()

// Insert values
s.Set("users[0].profile.name", "Alice", "config.yaml")
s.Set("users[0].profile.age", "30", "config.yaml")

// Check existence
exists := s.Has("users[0].profile.name")
// true

// Get subkeys
subs, _ := s.SubKeys("users[0].profile")
// subs => ["age", "name"]

// Inspect raw storage
raw := s.RawData()
// map["users[0].profile.name"] => ValueInfo{File:0, Value:"Alice"}
```

---

## 📖 Use Cases

* **Configuration Management**

    * Load and merge multiple config files (YAML/JSON/TOML)
    * Detect and prevent structural conflicts

* **Structured Data Access**

    * Safely access nested values with clear key/index separation
    * Unified path handling across maps and arrays

* **Validation & Debugging**

    * Check if a configuration path exists
    * Retrieve all subkeys under a path

---

## ⚠️ Notes

* Invalid path formats (spaces, consecutive dots, unclosed brackets, etc.) will return an error
* Structural conflicts (e.g., `"user.name"` vs `"user[0]"`) are detected and rejected
* `RawData` and `RawFile` expose internal state directly — use with caution

---

## 📜 License

Apache 2.0 License