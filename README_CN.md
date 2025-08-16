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

Barky 是一个 **层级化配置存储与路径解析工具库**，支持类似 JSON / YAML / TOML 的嵌套访问方式。

通过 `Path` 结构体与 `Storage` 存储引擎，Barky 提供了 **类型安全** 的路径管理、子键查询、冲突检测等功能，方便在配置解析与管理中使用。

---

## ✨ 功能特性

* **路径解析与构建**

    * `SplitPath`: 将字符串路径解析为结构化的 `Path` 数组
    * `JoinPath`: 将 `Path` 数组还原为字符串路径

* **两类路径类型**

    * **Key**: 用于字典/对象访问（如 `"user.name"`)
    * **Index**: 用于数组访问（如 `"[0]"`)

* **层级化存储结构**

    * 使用树结构（`treeNode`）维护路径层次，保证类型一致性
    * 支持子键查询 (`SubKeys`)
    * 支持路径存在性检测 (`Has`)
    * 写入时自动构建路径，检测并阻止冲突 (`Set`)

* **冲突检测**

    * 当路径层级类型不一致时（例如某处既被当作对象又被当作数组），会抛出错误

* **文件来源追踪**

    * 每个键值对会记录其所属文件及文件 ID，方便多文件合并场景下定位来源

---

## 📦 安装

```bash
go get github.com/go-spring/barky
```

---

## 🛠 使用示例

### 路径解析与构建

```go
path, _ := barky.SplitPath("users[0].profile.name")
// path => [ {Key:"users"}, {Index:"0"}, {Key:"profile"}, {Key:"name"} ]

joined := barky.JoinPath(path)
// joined => "users[0].profile.name"
```

### 存储与查询

```go
s := barky.NewStorage()
fileID := s.AddFile("config.yaml")

// 写入配置
s.Set("users[0].profile.name", "Alice", fileID)
s.Set("users[0].profile.age", "30", fileID)

// 检查路径是否存在
exists := s.Has("users[0].profile.name")
// true

// 查询子键
subs, _ := s.SubKeys("users[0].profile")
// subs => ["age", "name"]

// 查看底层存储
raw := s.RawData()
// map["users[0].profile.name"] => ValueInfo{File:0, Value:"Alice"}
```

---

## 📖 使用场景

* **配置管理**

    * 加载和合并多文件配置（YAML/JSON/TOML）
    * 路径冲突检测，避免结构不一致

* **结构化数据访问**

    * 以类型安全的方式访问嵌套结构
    * 统一处理键/数组下标

* **验证与调试**

    * 检查配置文件路径是否存在
    * 获取某路径下所有子键

---

## ⚠️ 注意事项

* 路径中不允许包含 **空格**、**连续的点号**、**未闭合的方括号** 等非法格式
* 如果路径在树结构中存在类型冲突（如 `"user.name"` 与 `"user[0]"`），会返回错误
* `RawData` 与 `RawFile` 方法直接暴露内部存储，使用时需谨慎

---

## 📜 许可证

Apache 2.0 License