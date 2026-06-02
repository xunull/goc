# traverse_v2 使用文档

`traverse_v2` 是 `traverse` 的重写版本，专为大规模目录遍历设计。相比 v1，它修正了完成判定的并发 race、避免每条目的 lstat 系统调用，并支持**层级化签单**（per-dir 完成事件）。

> v1 没有被移除，仍可继续使用。v2 是独立包：`github.com/xunull/goc/traverse_v2`。

---

## 目录

- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [API 参考](#api-参考)
- [Options 配置](#options-配置)
- [使用模式](#使用模式)
  - [一行拿全部文件路径（GetAllPaths）](#一行拿全部文件路径getallpaths)
  - [按语言统计文件（GetFileCount）](#按语言统计文件getfilecount)
  - [收集 + 集合查询（GetFileList）](#收集--集合查询getfilelist)
  - [统计文件数量（cloc 风格）](#统计文件数量cloc-风格)
  - [收集所有匹配后缀的文件路径](#收集所有匹配后缀的文件路径)
  - [per-dir 完成事件（签单）](#per-dir-完成事件签单)
  - [全局完成事件](#全局完成事件)
  - [带取消的遍历](#带取消的遍历)
- [Item 字段说明](#item-字段说明)
- [错误处理](#错误处理)
- [性能调优](#性能调优)
- [与 v1 的差异](#与-v1-的差异)
- [常见陷阱](#常见陷阱)

---

## 快速开始

```go
package main

import (
    "fmt"
    "sync/atomic"

    "github.com/xunull/goc/traverse_v2"
)

func main() {
    var fileCount atomic.Int64

    trv := traverse_v2.New("/path/to/scan", func(item *traverse_v2.Item) {
        if !item.IsDir {
            fileCount.Add(1)
        }
    })

    if err := trv.Run(); err != nil {
        fmt.Println("traverse err:", err)
    }

    fmt.Println("files:", fileCount.Load())
}
```

三行核心：
1. `New(path, callback, opts...)` 构造遍历器
2. `Run()` 阻塞执行，遍历完成才返回
3. callback 由 worker pool 并发调用，**必须是 goroutine-safe**

---

## 核心概念

### 签单模型（Hierarchical Signing Sheet）

每个目录持有一个内部 `dirNode`，记录两个原子计数：

| 字段 | 含义 |
|---|---|
| `expected` | 该目录直接子项（子目录 + 文件）数量 |
| `done` | 已完成的子项数量 |

一个目录的子目录"完成"是指它**整棵子树**都完成；一个文件"完成"是指它的 callback 返回。完成事件会**沿父指针向上冒泡**：叶子完成 → 父目录 done++ → 父目录可能完成 → 祖父目录 done++ → ... → root 完成 → `Run()` 返回。

### 两个独立 worker pool

| Pool | 默认并发 | 用途 |
|---|---|---|
| dir pool | 64 | 跑 `os.ReadDir`（IO 密集） |
| file pool | `NumCPU × 2` | 跑用户 callback |

两边独立调参，互不挤兑。pool 队列满时会以"溢出 goroutine"形式立刻执行任务——既不丢任务也不会因递归提交而死锁。

---

## API 参考

### 构造

```go
func New(path string, onItem func(*Item), opts ...Option) *Traverse
```

- `path`：根目录绝对路径
- `onItem`：每发现一个 file 或非根 dir 时被调用；可为 `nil`
- `opts`：见下面 [Options 配置](#options-配置)

### 执行

```go
func (t *Traverse) Run() error
```

阻塞直到整个子树遍历完毕、所有 file callback 返回。返回值非 nil 表示**期间出现过 `ReadDir` 错误**，包装的是第一个错误，全部错误可通过 `Errors()` 取回。

### 错误与取消

```go
func (t *Traverse) Errors() []error      // 所有累积的错误副本
func (t *Traverse) HasErrors() bool      // 快速判断
func (t *Traverse) Cancel()              // 触发取消（next traverseDir 进入会立刻关闭签单返回）
```

### 全局完成信号

```go
func (t *Traverse) Done() <-chan struct{}  // 全部完成时关闭的 channel,适合 select
```

`Done()` 在构造时就创建，调用 `Run()` 之前也能安全引用。配合 `WithOnComplete(fn)` 是消费者侧感知"全部结束"的两条独立路径——见 [全局完成事件](#全局完成事件)。

---

## Options 配置

| Option | 默认 | 说明 |
|---|---|---|
| `WithDirWorkers(n)` | 64 | dir pool worker 数 |
| `WithFileWorkers(n)` | `NumCPU*2` | file pool worker 数 |
| `WithWorkerCount(n)` | — | v1 兼容,同时设 dir/file 两个 pool 为 n |
| `WithQueueScale(n)` | 16 | 队列大小 = workers × scale |
| `WithMaxDepth(d)` | 0（无限） | 最大递归深度，root 是 0 |
| `WithDepth(d)` | — | v1 兼容别名,等价 `WithMaxDepth` |
| `WithSyncMode()` | off | v1 兼容,等价 `WithDirWorkers(1)` |
| `WithSyncFileOpMode()` | off | v1 兼容,等价 `WithFileWorkers(1)` |
| `WithTargetExt(".go")` | — | 只对该扩展名文件触发 callback |
| `WithExcludePrefix("_", ...)` | — | 忽略指定前缀的文件名 |
| `WithExcludeSuffix(".bak", ...)` | — | 忽略指定后缀的文件名 |
| `WithExcludeDir("vendor", "src/gen")` | — | 忽略指定目录（按 basename 或相对路径匹配） |
| `WithSkipDotEntries()` | off | 跳过所有以 `.` 开头的文件和目录 |
| `WithSkipKnownIgnoreDirs()` | off | 跳过 `lang_ext.CommonExcludeDir`（node_modules、vendor、dist 等） |
| `WithSkipKnownBinaryFiles()` | off | 跳过 `lang_ext.CommonExcludeFileExt`（.exe、.so、.pyc 等） |
| `WithSensibleDefaults()` | off | 一次性打开上面三个 skip |
| `WithDefaultExclude()` | off | v1 兼容别名,等价 `WithSensibleDefaults` |
| `WithOnlyDir()` | off | 只对目录触发 callback，文件不触发 |
| `WithOnDirComplete(fn)` | — | 每个目录子树完成时调用一次（见下文） |
| `WithOnComplete(fn)` | — | 整个遍历结束时调用一次（在 Done() 关闭之前） |

---

## 使用模式

### 一行拿全部文件路径（GetAllPaths）

`GetAllPaths` 是 v1 `DirTraverse.GetAllPath` 的对应物，封装了 `New + Run + 收集` 三步。返回相对于 `dir` 的文件路径（不含目录），分隔符固定 `/`。

```go
paths, err := traverse_v2.GetAllPaths("/some/dir",
    traverse_v2.WithSensibleDefaults(),
    traverse_v2.WithTargetExt(".go"),
)
if err != nil {
    log.Println("partial:", err) // 仍可能拿到部分结果
}
fmt.Println(len(paths), "files")
```

返回的顺序**不保证**，需要稳定顺序就自己 `sort.Strings(paths)`。

### 按语言统计文件（GetFileCount）

`GetFileCount` 是 v1 `GetFileCount` 的对应物，返回总数 + 按语言分组：

```go
stats, err := traverse_v2.GetFileCount(root,
    traverse_v2.WithSensibleDefaults(),
)
fmt.Println("total files:", stats.Total)
for lang, n := range stats.ByLanguage {
    fmt.Printf("  %s: %d\n", lang, n)
}
```

`ByLanguage` 的 key 来自 `lang_ext.CommonLanguageExt`（如 "Golang"、"Python"、"Markdown"）。未识别扩展名的文件**会计入 `Total` 但不进入 `ByLanguage`**。

> 与 v1 差异：v1 的 `Count` 因为 callback 同时收文件和目录会**多计目录**，v2 的 `Total` 只数文件。v1 的 `TargetCount` 字段被去掉了（与 `WithTargetExt` 配合时 `Total` 就是它）。

### 收集 + 集合查询（GetFileList）

`GetFileList` 等价于 `GetAllPaths` 加一个 O(1) 查询集合，方便后续判断"某路径是否在结果里"：

```go
res, err := traverse_v2.GetFileList(root, traverse_v2.WithTargetExt(".go"))
fmt.Println("count:", len(res.List))
if _, exists := res.Set["pkg/foo/main.go"]; exists {
    // ...
}
```

> 与 v1 差异：v1 的 `GetFileList` 在**没设 `WithTargetExt`** 时会返回空 List（callback 里的 bug）；v2 在没设过滤时返回**所有文件**。

### 统计文件数量（cloc 风格）

```go
type Stats struct {
    Total    atomic.Int64
    ByExt    sync.Map // map[string]*atomic.Int64
}

var stats Stats

trv := traverse_v2.New(root, func(item *traverse_v2.Item) {
    if item.IsDir {
        return
    }
    stats.Total.Add(1)
    counter, _ := stats.ByExt.LoadOrStore(item.Ext, &atomic.Int64{})
    counter.(*atomic.Int64).Add(1)
},
    traverse_v2.WithSensibleDefaults(),
)

if err := trv.Run(); err != nil {
    log.Println("partial errors:", trv.Errors())
}
```

### 收集所有匹配后缀的文件路径

```go
var (
    mu    sync.Mutex
    files []string
)

trv := traverse_v2.New(root, func(item *traverse_v2.Item) {
    if item.IsDir {
        return
    }
    mu.Lock()
    files = append(files, item.FullPath)
    mu.Unlock()
},
    traverse_v2.WithTargetExt(".go"),
    traverse_v2.WithSkipKnownIgnoreDirs(),
)

_ = trv.Run()
fmt.Printf("found %d .go files\n", len(files))
```

> 高并发下 slice 写入要加锁。如果性能敏感可换 channel 推流 + 单 reader 聚合。

### per-dir 完成事件（签单）

每个目录的整棵子树完成时回调一次，可用于：
- 把每个目录的文件统计聚合成"目录视图"
- 增量输出（边遍历边把已完成的目录写到磁盘）
- 进度展示（每完成一个顶层目录就 +1）

```go
// 每个目录的行数累加器
type DirStat struct {
    Lines int64
}
dirStats := sync.Map{} // map[string]*DirStat

trv := traverse_v2.New(root, func(item *traverse_v2.Item) {
    if item.IsDir {
        // 初始化 dir 的统计
        dirStats.Store(item.Path, &DirStat{})
        return
    }
    // 找到所属目录,把行数累加上去
    dirPath := filepath.Dir(item.Path)
    if v, ok := dirStats.Load(dirPath); ok {
        atomic.AddInt64(&v.(*DirStat).Lines, countLines(item.FullPath))
    }
},
    traverse_v2.WithOnDirComplete(func(item *traverse_v2.Item) {
        // 此目录的整棵子树都处理完了 → 安全地读取累加结果
        if v, ok := dirStats.Load(item.Path); ok {
            fmt.Printf("%s: %d lines\n", item.Path, v.(*DirStat).Lines)
        }
    }),
)

_ = trv.Run()
```

**关键保证**：`OnDirComplete` fire 时，该目录的所有子目录已经各自 fire 过 `OnDirComplete`，且该目录下所有文件的 `onItem` 都已经返回。root 一定**最后**触发。

### 全局完成事件

当**调用 `Run()` 的不是消费者**（消费者在另一个 goroutine 里）时，用以下两种方式之一拿到"全部完成"信号。

#### 方式 A：回调式 `WithOnComplete`

```go
results := make(chan Stat, 1)

trv := traverse_v2.New(root,
    func(item *traverse_v2.Item) {
        // 处理 item,累加到某个共享 stats
    },
    traverse_v2.WithOnComplete(func() {
        // 整个遍历结束时调用一次,可以安全读 stats
        results <- collectStats()
    }),
)

go trv.Run()

stat := <-results
fmt.Println(stat)
```

#### 方式 B：channel 式 `Done()`

```go
trv := traverse_v2.New(root, onItem)

go trv.Run()

select {
case <-trv.Done():
    fmt.Println("traverse finished")
case <-time.After(30 * time.Second):
    trv.Cancel()
    <-trv.Done()  // 取消后仍要等收尾
    fmt.Println("timeout, cancelled")
}
```

#### 信号触发顺序（重要）

完成时按这个**严格顺序**触发，所有观察者看到的状态都是自洽的：

1. root 目录的 `WithOnDirComplete(...)` 触发（如果注册了）
2. `WithOnComplete(...)` 触发（如果注册了）
3. `Done()` 返回的 channel 被关闭
4. `Run()` 返回

也就是说：**`Done()` 关闭时，所有 callback 都已经返回**——你不会观察到"Done 关了但回调还在跑"的中间态。这是签单冒泡到 root 才关闭 channel 保证的（`node.go:48-72`）。



```go
trv := traverse_v2.New(root, func(item *traverse_v2.Item) {
    if shouldStop(item) {
        trv.Cancel() // 安全地让 Traverse 内部自己感知到
    }
})

// 单独的 watchdog
go func() {
    time.Sleep(5 * time.Second)
    trv.Cancel()
}()

_ = trv.Run() // 取消后仍会等已经提交的任务收尾,但不再发起新的目录读取
```

`Cancel()` 后 `Run()` 仍会等待已提交的任务执行完才返回——不会强杀 callback。

---

## Item 字段说明

```go
type Item struct {
    Path     string      // 相对根的路径,"/"分隔; root 自己的 Path 是 ""
    FullPath string      // 绝对路径,使用 OS 分隔符
    Name     string      // basename
    Ext      string      // 扩展名(含点),空字符串表示无扩展
    Mode     fs.FileMode // 来自 DirEntry.Type(); 不是 os.Stat 那种完整 mode
    IsDir    bool
    Depth    int         // root 为 0; 一级子项为 1
}
```

注意 `Mode`：只包含 `entry.Type()` 暴露的 bit（type 信息 + 平台允许的 perm bit）。如果需要 size/mtime/owner，需要再调一次 `os.Stat(item.FullPath)`。

---

## 错误处理

- **`Run()` 返回非 nil** ⇒ 遍历期间至少有一次 `os.ReadDir` 失败（权限不足、目录消失等）。**遍历仍然完成**，没失败的部分都处理过了
- 完整错误列表用 `trv.Errors()` 取
- 用户 callback **不应 panic**——内部 worker 没有 recover，panic 会让 worker 死掉、签单永远等不到，整个 `Run()` 阻塞

```go
trv := traverse_v2.New(root, func(item *traverse_v2.Item) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("callback panic on %s: %v", item.Path, r)
        }
    }()
    // ... user logic ...
})
```

---

## 性能调优

### 默认配置适用场景

- 默认 `DirWorkers=64`, `FileWorkers=NumCPU*2`
- 大多数 SSD + 中等规模仓库（< 100 万文件）跑得很好

### 不同硬件/工作负载的建议

| 场景 | DirWorkers | FileWorkers | 其他 |
|---|---|---|---|
| SSD + 浅而宽的树（多顶层目录） | 128+ | NumCPU\*4 | — |
| HDD（机械硬盘） | 4-8 | NumCPU | 避免 seek 风暴 |
| Callback 是纯 CPU（如解析 AST） | 64 | NumCPU | 别开太多 file worker |
| Callback 是慢 IO（如网络上传） | 64 | 200+ | file worker 多多益善 |
| 内存吃紧 | 32 | NumCPU | 把 `QueueScale` 调低（如 4） |

### 测过哪些数据

在 M1 Max + macOS + APFS 上，4 层深、每层 10 branch、每个目录 10 文件（共 ~11K dirs + ~100K files）：
- v1: ~229 ms
- v2: ~208 ms（**约快 10%**，主要来自去掉 `entry.Info()` 的 lstat）

Linux + ext4 + 旋转盘上预期差距会更大（lstat 开销更高）。

### 别做的事

- 别在 callback 里做长阻塞（持锁、网络同步请求等）——会卡死 file worker，签单冒泡停滞
- 别把 `FileWorkers` 设到几千——goroutine 多了 Go 调度器自身开销会反过来咬性能
- 别在内层 callback 里嵌套调用 `Run()`——会重复 close `t.done` panic

---

## 与 v1 的差异

| 维度 | v1 (`traverse`) | v2 (`traverse_v2`) |
|---|---|---|
| 完成判定 | `Count == OverCount` 全局相等 | 签单原子三元组 + CAS + 父指针冒泡 |
| 并发 race | 有（跨目录交错时可能误判完成） | 无（race detector 干净） |
| 系统调用 | 每条目额外一次 `entry.Info()`(lstat) | 只用 `entry.Type()`,零额外调用 |
| Pool | 一个 pool 同时跑目录和回调 | 两个独立 pool,分别调参 |
| Submit 死锁风险 | 多通道间接调度 | submit 不阻塞(满则 overflow) |
| 资源清理 | `Close()` → `pool.Release()` 是空函数(泄漏) | `Run()` 自动 close pool |
| per-dir 完成事件 | 无 | `WithOnDirComplete` |
| Cancellation | `cancel context` 但 pool 不感知 | 同左 + 取消后签单仍正确冒泡 |

### 迁移：v1 → v2

API 差不多。最常见的几处变化：

| v1 | v2 |
|---|---|
| `NewDirTraverse(p, fn)` + `Handle()` + `WorkSheet.Wait()` + `Close()` | `New(p, fn).Run()` |
| `DirTraverse.GetAllPath(opts...)` | `GetAllPaths(dir, opts...)`（顶层函数） |
| `traverse.GetFileCount(dir, opts...)` | `GetFileCount(dir, opts...)`（结构字段更新，见用例小节） |
| `traverse.GetFileList(dir, opts...)` | `GetFileList(dir, opts...)`（修复 v1 没 TargetExt 时空 List 的 bug） |
| `WithWorkerCount(n)` | `WithWorkerCount(n)` 已保留 / 或用 `WithDirWorkers + WithFileWorkers` 精细调 |
| `WithDefaultExclude()` | `WithDefaultExclude()` 已保留为别名 / 推荐用 `WithSensibleDefaults()` |
| `WithDepth(d)` | `WithDepth(d)` 已保留为别名 / 推荐用 `WithMaxDepth(d)` |
| `WithSyncMode()` | `WithSyncMode()` 已保留，等价 `WithDirWorkers(1)` |
| `WithSyncFileOpMode()` | `WithSyncFileOpMode()` 已保留，等价 `WithFileWorkers(1)` |
| `WithExcludeUnknown(bool)` | 未实现（v1 设了从来没读过，是死字段） |
| `WithProgressBarOut()` | 未实现（v1 也没真正接上 pb 库；自行用 `WithOnDirComplete` + 计数实现） |
| `DirTraverse.WaitOver()` | 用 `Run()` 阻塞,或另起 goroutine + `<-Done()` |
| `DirTraverse.Close()` | 不需要,`Run()` 自动清理 pool |
| `DirTraverse.SetOption(opts...)` | 不支持运行时改 option,在 `New(...)` 时传入 |
| `item.FileInfo.Size()` 等 | v2 没有 FileInfo,需要时 `os.Stat(item.FullPath)` |

---

## 常见陷阱

### 1. callback 内的 `item` 别带出 callback 作用域以外

```go
var saved []*traverse_v2.Item
trv := traverse_v2.New(root, func(item *traverse_v2.Item) {
    saved = append(saved, item) // ❌ 数据竞争
})
```

`saved` 没保护 + 多个 worker 并发写 = data race。要么加锁，要么用 channel 推流。

### 2. 期待 callback 调用顺序

callback 是**并发**的，文件顺序、目录顺序都不保证。需要顺序就把结果收集起来在 `Run()` 之后排序。

### 3. 期待 `OnDirComplete` 严格按 DFS post-order

通常它接近 post-order（叶子先，父后），但当多个兄弟目录并发完成时，事件顺序取决于哪个先完成。**唯一保证**是：父目录的 `OnDirComplete` 一定在其所有子目录的 `OnDirComplete` 之后。

### 4. callback panic 卡死 Run

如前所述，callback 里要么自己 recover，要么确保不会 panic。

### 5. 把 `Run()` 当幂等接口

`Run()` 只能调一次。Traverse 实例不可复用（内部 channel 关闭后就废了）。要再跑用新的 `New(...)`。

---

## 完整最小示例

```go
package main

import (
    "fmt"
    "log"
    "sync/atomic"

    "github.com/xunull/goc/traverse_v2"
)

func main() {
    var (
        files atomic.Int64
        dirs  atomic.Int64
        lines atomic.Int64
    )

    trv := traverse_v2.New("/Users/quincy/code/myrepo",
        func(item *traverse_v2.Item) {
            if item.IsDir {
                dirs.Add(1)
                return
            }
            files.Add(1)
            // ... 真实业务: 例如 cloc 行数累加
            lines.Add(countLinesIn(item.FullPath))
        },
        traverse_v2.WithSensibleDefaults(),
        traverse_v2.WithTargetExt(".go"),
        traverse_v2.WithOnDirComplete(func(item *traverse_v2.Item) {
            if item.Depth == 1 {
                // 每个顶层目录完成时打一行
                fmt.Printf("done: %s\n", item.Path)
            }
        }),
    )

    if err := trv.Run(); err != nil {
        log.Printf("partial errors: %v", trv.Errors())
    }

    fmt.Printf("dirs=%d files=%d lines=%d\n",
        dirs.Load(), files.Load(), lines.Load())
}

func countLinesIn(path string) int64 {
    // ... 简化省略
    return 0
}
```
