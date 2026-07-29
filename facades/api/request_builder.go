package api

// RequestBuilder 批量请求构造器
type RequestBuilder struct {
	requests []*BatchRequestItem
	keys     []string
	idx      map[string]int // key -> requests 下标，O(1) 查询与同名去重
}

// NewRequestBuilder 创建一个新的 RequestBuilder。
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{idx: make(map[string]int)}
}

// Add 添加一个带命名 key 的请求项，返回自身以便链式调用。
// req 为 nil 时跳过（同时不登记 key），便于直接传入可能返回 nil 的 Option 构造函数；
// 同名重复注册时覆盖原有请求项。
func (rb *RequestBuilder) Add(key string, req *BatchRequestItem) *RequestBuilder {
	if req == nil {
		return rb
	}
	if pos, ok := rb.idx[key]; ok {
		rb.requests[pos] = req
		return rb
	}
	rb.idx[key] = len(rb.requests)
	rb.keys = append(rb.keys, key)
	rb.requests = append(rb.requests, req)
	return rb
}

// Build 返回构建好的请求切片，传给 MRequest。顺序为 Add 的注册顺序。
func (rb *RequestBuilder) Build() []*BatchRequestItem {
	return rb.requests
}

// IndexOf 返回指定 key 在请求切片中的位置索引，不存在则返回 -1。
func (rb *RequestBuilder) IndexOf(key string) int {
	if i, ok := rb.idx[key]; ok {
		return i
	}
	return -1
}

// Has 判断指定 key 是否存在（语义对应旧代码的 xxPos >= 0）。
func (rb *RequestBuilder) Has(key string) bool {
	_, ok := rb.idx[key]
	return ok
}

// Len 返回请求数量。
func (rb *RequestBuilder) Len() int {
	return len(rb.requests)
}
