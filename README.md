# wbutil-go

Go 语言工具库

---

## 📋 注意事项

### 1. 设置环境变量

> ⚠️ 如果不设置，公司内部的库将无法拉取

```shell
go env -w GOPRIVATE='git.intra.weibo.com'
go env -w GONOPROXY='git.intra.weibo.com'
go env -w GONOSUMDB='git.intra.weibo.com'
```

### 2. Git 配置

```shell
git config --global url.ssh://git@git.intra.weibo.com:2222/.insteadOf https://git.intra.weibo.com/
```

### 3. 测试环境 Mesh 服务

测试环境需要增加 mesh 服务：

```shell
docker run -d \
  -p 9981:9981 \
  --restart=on-failure:3 \
  --name motan \
  -e MESH_DETECT_BACKEND_URL=http://127.0.0.1:8080/ \
  -e MESH_DETECT_OK_STATUS=200 \
  -e MESH_DETECT_INTERVAL=1 \
  -v /data1/motan/logs:/data1/motan/logs \
  registry.api.weibo.com/openapi_rd/weibo-mesh-with-mesh-confs:0.1.65-0.1.304 \
  -application wb-search-app_finder-web
```

---

## 🚀 API Request 示例

### 单个请求

```go
params := map[string]any{
    "uid": 123456,
    "ids": "5261460160645774",
}

// 超时控制
other := &api.OtherOption{
    Timeout: 3 * time.Second,
}

result, err := api.Request(ctx, "mblog_motan", params, other)
```

### 批量请求

```go
result := api.MRequest(
    ctx,
    []*api.BatchRequestItem{
        {
            URLKey: "finder_ad_banner",
            Params: map[string]any{
                "uid": 123456,
                "ids": "5261460160645774",
            },
        },
        {
            URLKey: "mblog_motan",
            Params: map[string]any{
                "uid": 123456,
                "ids": "5261460160645774",
            },
        },
        {
            URLKey: "da_search",
            Params: map[string]any{
                "uid":   2655350737,
                "query": "快递停运",
                "sid":   "search_fe_openapi",
            },
        },
    },
)
```

---

## 🗄️ Redis 示例

### 获取连接

```go
import "git.intra.weibo.com/search_fe/wbutil-go/facades/db/redis"

// 基本连接
conn, err := redis.GetConnect(ctx,redis.ConnectArgs{
    BusKey: "fe_m_redis",
})

// 带 Hash 分片的连接
conn, err := redis.GetConnect(ctx,redis.ConnectArgs{
    BusKey:  "user_read",
    HashKey: "user_123",  // 用于 hash 分片
    Timeout: 500,         // 超时时间（毫秒）
    DB:      1,           // 数据库编号，-1 使用配置默认值
})

// 指定节点连接
conn, err := redis.GetConnect(ctx,redis.ConnectArgs{
    BusKey: "counter_redis",
    HashNo: 2,  // 指定节点编号（从1开始）
})
```

### 字符串操作

```go
// 设置值（带过期时间）
err := conn.Set(ctx, "key", "value", 10*time.Minute)

// 获取值
value, err := conn.Get(ctx, "key")

// 自增/自减
count, err := conn.Incr(ctx, "counter")
count, err := conn.Decr(ctx, "counter")
count, err := conn.IncrBy(ctx, "counter", 10)
```

### Hash 操作

```go
// 设置 hash 字段
err := conn.HSet(ctx, "user:123", "name", "张三", "age", "25")

// 获取单个字段
name, err := conn.HGet(ctx, "user:123", "name")

// 获取所有字段
data, err := conn.HGetAll(ctx, "user:123")
// data: map[string]string{"name": "张三", "age": "25"}

// 删除字段
err := conn.HDel(ctx, "user:123", "age")
```

### 有序集合操作

```go
import "github.com/go-redis/redis/v8"

// 添加成员
err := conn.ZAdd(ctx, "ranking",
    &redis.Z{Score: 100, Member: "user1"},
    &redis.Z{Score: 200, Member: "user2"},
)

// 获取范围（不带分数）
members, err := conn.ZRange(ctx, "ranking", 0, -1)

// 获取范围（带分数）
scores, err := conn.ZRangeWithScores(ctx, "ranking", 0, -1)

// 格式化有序集合数据
formatted := redis.FormatZSet(scores, "member", "score")
// 返回: [{"member": "user1", "score": 100}, {"member": "user2", "score": 200}]

// 删除成员
err := conn.ZRem(ctx, "ranking", "user1")
```

### 列表操作

```go
// 左侧推入
err := conn.LPush(ctx, "queue", "item1", "item2")

// 右侧推入
err := conn.RPush(ctx, "queue", "item3")

// 获取范围
items, err := conn.LRange(ctx, "queue", 0, -1)
```

### 集合操作

```go
// 添加成员
err := conn.SAdd(ctx, "tags", "go", "redis", "weibo")

// 获取所有成员
members, err := conn.SMembers(ctx, "tags")

// 删除成员
err := conn.SRem(ctx, "tags", "redis")
```

### 通用操作

```go
// 检查键是否存在
count, err := conn.Exists(ctx, "key1", "key2")

// 设置过期时间
err := conn.Expire(ctx, "key", 30*time.Minute)

// 删除键
err := conn.Del(ctx, "key1", "key2")

// 执行任意命令
result, err := conn.Exec(ctx, "PING")
```

### 获取原生客户端

```go
// 获取 go-redis 原生客户端，用于执行更复杂的操作
client := conn.GetClient()
```

---

## 💾 Cache 缓存示例

> 跨平台共享内存缓存，支持进程间共享

### 获取缓存实例

```go
import "git.intra.weibo.com/search_fe/wbutil-go/facades/cache"

// 方式一：通过配置文件获取单例实例
// 需要在 config/app.yaml 中配置 cache 节点
c, err := cache.GetInstance("test_memoryname")

// 方式二：手动创建缓存实例 优先用这种方式
c, err := cache.NewSharedCache(cache.SharedCacheConfig{
    Name: "my_cache",           // 缓存名称（用于生成文件名）
    Size: 10 * 1024 * 1024,     // 缓存大小（10MB）
    Dir:  "/tmp/cache",         // 缓存目录（可选，默认临时目录）
})
```

### 基本操作

```go
// 设置缓存（永不过期）
err := c.Set("key", "value")

// 设置缓存（带过期时间，单位：秒）
err := c.Set("key", "value", 60)  // 60秒后过期

// 获取缓存
value, ok := c.Get("key")
if ok {
    fmt.Println(value)
}

// 检查键是否存在
if c.Has("key") {
    // 键存在
}

// 删除缓存
err := c.Delete("key")

// 清空所有缓存
err := c.Clear()
```

### 支持的数据类型

```go
// 字符串
c.Set("string_key", "hello world")

// 整数
c.Set("int_key", 12345)

// 浮点数
c.Set("float_key", 3.14159)

// 布尔值
c.Set("bool_key", true)

// 切片
c.Set("slice_key", []string{"a", "b", "c"})

// Map
c.Set("map_key", map[string]any{
    "name": "test",
    "age":  18,
})

// 复杂嵌套结构
c.Set("complex_key", map[string]any{
    "name": "test",
    "nested": map[string]any{
        "level1": map[string]any{
            "level2": "deep_value",
        },
    },
    "array": []any{1, 2, 3, "four"},
})
```

### 辅助方法

```go
// 获取缓存项数量
count := c.Count()

// 获取所有缓存键
keys := c.Keys()

// 获取缓存文件路径
filePath := c.GetFilePath()
```

### 生命周期管理

```go
// 关闭缓存（释放资源，保留文件）
err := c.Close()

// 删除缓存文件（彻底清理）
err := c.Unlink()
```

### 配置文件示例

```yaml
# config/app.yaml
cache:
  test_memoryname:
    max_memory: 10485760  # 10MB
```
