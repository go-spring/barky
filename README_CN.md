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

`Barky` 是一个用于处理分层键值数据结构的 Go 工具包，主要面向配置文件（如 JSON、YAML、TOML）等场景。
它可以将嵌套的数据结构展开成扁平化的 `map[string]string`，同时保留路径信息、冲突检测能力和多文件来源的追踪。

---

## ✨ 核心功能

### 1. 扁平化（Flattening）

- 支持将嵌套的 `map`、`slice`、`array` 转换为扁平化的 `map[string]string`。
- 使用 **点号** 表示 map 键，使用 **中括号** 表示数组/切片索引。

### 2. 路径处理（Path Handling）

* `Path` 类型将层级键解析为一组 **路径段**（map 键或数组索引）。
* 提供 `SplitPath` 将 `"foo.bar[0]"` 转换为结构化路径，`JoinPath` 将路径段重新拼接成字符串。
* 支持语法校验，避免非法键（如连续的点号、未闭合的中括号）。

### 3. 存储（Storage）

* `Storage` 类型管理一组扁平化的键值对，同时在内部构建一棵树，用于检测结构冲突。
* 每个值都带有来源文件索引（`file index`），方便在多文件合并时追踪来源。
* 提供的方法包括：

    * `Set`：设置键值，带冲突检测。
    * `Get`：获取值，支持默认值。
    * `Has`：判断键是否存在。
    * `SubKeys`：列出某个路径下的子键。
    * `Keys`：获取所有已存储的键（排序后）。

---

## 📦 安装方式

```bash
go get github.com/go-spring/barky
```

---

## 🛠 使用示例

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

输出：

```
Keys: [server.hosts[0].ip server.hosts[1].ip]
SubKeys of server.hosts: [0 1]
Get server.hosts[0].ip: 192.168.0.1
```

---

## 📖 使用场景

* **配置管理**
  将不同格式的配置文件转换成统一的扁平化键值对，便于比较和合并。

* **查询与检索**
  使用简单的路径字符串（如 `"server.hosts[0].ip"`）直接访问数据，而无需手动遍历嵌套结构。

* **多文件合并**
  同时处理多个配置文件，追踪每个键值来自哪个文件，并在冲突时提供检测。

* **数据转换**
  作为桥梁，将层级化数据结构转换为更易处理的键值形式，用于差异分析、测试或下游系统处理。

---

## ⚠️ 注意事项

* 路径中不允许包含 **空格**、**连续的点号**、**未闭合的方括号** 等非法格式
* 如果路径在树结构中存在类型冲突（如 `"user.name"` 与 `"user[0]"`），会返回错误
* `RawData` 与 `RawFile` 方法直接暴露内部存储，使用时需谨慎

---

## 📜 许可证

Apache 2.0 License