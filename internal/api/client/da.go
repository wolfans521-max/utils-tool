// Package client 提供API客户端实现
// DA协议是基于TCP的二进制协议，用于与DA服务进行通信
package client

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/bytedance/sonic"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"time"
)

const (
	// DACommonHeadSize CommonHead结构体大小
	DACommonHeadSize = 32
	// DACMDHeadSize CMDHead结构体大小
	DACMDHeadSize = 8
	// DAMagicNumber DA协议魔数
	DAMagicNumber = 0xa3b9f14d
)

// DACommonHead DA协议通用头部
type DACommonHead struct {
	Version   uint8  // 版本号
	Flag      uint8  // 标志位
	Cbit      uint8  // 压缩位
	Pad       uint8  // 填充位
	SeqID     uint32 // 序列ID
	AsynSeqID uint32 // 异步序列ID
	BodyLen   uint32 // 消息体长度
	Timestamp uint64 // 时间戳（微秒）
	Magic     uint32 // 魔数
	Reserved  uint32 // 保留字段
}

// DACMDHead DA协议命令头部
type DACMDHead struct {
	Cmd uint32 // 命令ID（28位）+ 版本（4位）
	Len uint32 // 数据长度
}

// DATiming DA请求时间记录结构
type DATiming struct {
	TotalTime   time.Duration // 总耗时
	ConnectTime time.Duration // 连接时间
	SendTime    time.Duration // 请求发送时间
	WaitTime    time.Duration // 等待响应时间
	ReadTime    time.Duration // 响应读取时间
}

// DAClient DA协议客户端
type DAClient struct {
	host    string
	port    int
	timeout time.Duration
	version int
}

// NewDAClient 创建DA客户端
// timeout参数为总超时时间
func NewDAClient(host string, port int, timeout time.Duration, version int) *DAClient {
	if version == 0 {
		version = 1
	}
	return &DAClient{
		host:    host,
		port:    port,
		timeout: timeout,
		version: version,
	}
}

// RequestWithTiming 执行DA请求并返回时间记录
func (c *DAClient) RequestWithTiming(args map[string]interface{}) ([]byte, *DATiming, error) {
	timing := &DATiming{}

	// 记录开始时间
	startTime := time.Now()

	// 设置总超时的截止时间
	deadline := startTime.Add(c.timeout)

	// 建立TCP连接（使用总超时时间）
	connectStart := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.host, c.port), c.timeout)
	if err != nil {
		timing.TotalTime = time.Since(startTime)
		timing.ConnectTime = time.Since(connectStart)
		// 如果连接失败，其他时间都为0
		return nil, timing, fmt.Errorf("连接DA服务失败: %w", err)
	}
	defer conn.Close()

	// 记录连接时间
	timing.ConnectTime = time.Since(connectStart)

	// 设置总超时的截止时间
	conn.SetDeadline(deadline)

	// 编码请求
	requestData, err := c.encode(args)
	if err != nil {
		timing.TotalTime = time.Since(startTime)
		// 编码失败，send/wait/read都为连接后的等待时间
		return nil, timing, fmt.Errorf("编码DA请求失败: %w", err)
	}

	// 发送请求
	sendStart := time.Now()
	_, err = conn.Write(requestData)
	if err != nil {
		timing.TotalTime = time.Since(startTime)
		timing.SendTime = time.Since(sendStart)
		// 发送失败，wait/read都为0
		return nil, timing, fmt.Errorf("发送DA请求失败: %w", err)
	}
	timing.SendTime = time.Since(sendStart)

	// 读取响应头 - 使用 io.ReadFull 确保完整读取
	waitStart := time.Now() // 等待响应开始
	headBuf := make([]byte, DACommonHeadSize)
	_, err = io.ReadFull(conn, headBuf)
	if err != nil {
		timing.TotalTime = time.Since(startTime)
		timing.WaitTime = time.Since(waitStart)
		// 读取响应头失败，read为0
		return nil, timing, fmt.Errorf("读取DA响应头失败: %w", err)
	}
	timing.WaitTime = time.Since(waitStart)

	// 解析响应头
	var head DACommonHead
	buf := bytes.NewReader(headBuf)
	binary.Read(buf, binary.LittleEndian, &head)

	// 读取响应体
	if head.BodyLen > 0 {
		bodyBuf := make([]byte, head.BodyLen)
		// 使用 io.ReadFull 确保完整读取响应体
		readStart := time.Now()
		_, err = io.ReadFull(conn, bodyBuf)
		if err != nil {
			timing.TotalTime = time.Since(startTime)
			timing.ReadTime = time.Since(readStart)
			return nil, timing, fmt.Errorf("读取DA响应体失败: %w", err)
		}
		timing.ReadTime = time.Since(readStart)

		// 根据版本跳过不同长度的头部
		// version=1: 跳过CMDHead（8字节）
		// version=2: 跳过OSCMDHead（28字节）+ 额外的12字节 = 40字节
		// version=4: 跳过3个CMDHead（3*8=24字节）
		skipBytes := DACMDHeadSize
		if c.version == 2 {
			skipBytes = 40 // OS格式需要跳过更多字节
		} else if c.version == 4 {
			skipBytes = 24 // 新意图格式跳过3个CMDHead
		}

		if len(bodyBuf) >= skipBytes {
			timing.TotalTime = time.Since(startTime)
			return bodyBuf[skipBytes:], timing, nil
		}
	}

	timing.TotalTime = time.Since(startTime)
	return nil, timing, fmt.Errorf("DA响应为空")
}

// Request 执行DA请求
func (c *DAClient) Request(args map[string]interface{}) ([]byte, error) {
	result, _, err := c.RequestWithTiming(args)
	return result, err
}

// encode 编码DA请求
func (c *DAClient) encode(args map[string]interface{}) ([]byte, error) {
	// 将参数转换为JSON
	jsonData, err := sonic.Marshal(args)
	if err != nil {
		return nil, err
	}

	// 根据版本选择编码方式
	switch c.version {
	case 1:
		return c.encodeCommon(jsonData)
	case 2:
		return c.encodeOS(jsonData)
	case 3:
		return c.encodeIntention(args)
	case 4:
		return c.encodeIntentionNew(args)
	default:
		return c.encodeCommon(jsonData)
	}
}

// encodeCommon 编码通用格式（version=1）
func (c *DAClient) encodeCommon(jsonData []byte) ([]byte, error) {
	// 构造CMDHead
	cmdHead := DACMDHead{
		Cmd: 1, // cmd=1, ver=0
		Len: uint32(len(jsonData)),
	}
	cmdHeadBuf := new(bytes.Buffer)
	binary.Write(cmdHeadBuf, binary.LittleEndian, cmdHead)
	cmdData := append(cmdHeadBuf.Bytes(), jsonData...)

	// 构造CommonHead
	seqID := c.generateSeqID()
	commonHead := DACommonHead{
		Version:   2,
		Flag:      0,
		Cbit:      0,
		Pad:       0,
		SeqID:     seqID,
		AsynSeqID: 0,
		BodyLen:   uint32(len(cmdData)),
		Timestamp: uint64(time.Now().UnixNano() / 1000), // 微秒
		Magic:     DAMagicNumber,
		Reserved:  0,
	}

	// 序列化CommonHead
	commonHeadBuf := new(bytes.Buffer)
	binary.Write(commonHeadBuf, binary.LittleEndian, commonHead)

	// 组合完整请求
	return append(commonHeadBuf.Bytes(), cmdData...), nil
}

// encodeOS 编码OS格式（version=2）
func (c *DAClient) encodeOS(jsonData []byte) ([]byte, error) {
	// OS格式的CMDHead结构不同，包含多个cmd字段
	// 这里简化实现，实际使用时需要根据具体需求调整
	type OSCMDHead struct {
		Cmds uint32
		Cmd0 uint32
		Cmd1 uint32
		Cmd2 uint32
		Cmd3 uint32
		Cmd4 uint32
		Len  uint32
	}

	cmdHead := OSCMDHead{
		Cmds: 5,
		Cmd0: 1,
		Cmd1: 2,
		Cmd2: 3,
		Cmd3: 4,
		Cmd4: 5,
		Len:  uint32(len(jsonData)),
	}
	cmdHeadBuf := new(bytes.Buffer)
	binary.Write(cmdHeadBuf, binary.LittleEndian, cmdHead)
	cmdData := append(cmdHeadBuf.Bytes(), jsonData...)

	// 构造CommonHead
	seqID := c.generateSeqID()
	tm := time.Now().UnixNano() / 1000 // 微秒

	type OSCommonHead struct {
		Version   uint32
		SeqID     uint32
		AsynSeqID uint32
		BodyLen   uint32
		TimeHigh  uint32
		TimeLow   uint32
		Magic     uint32
		Reserved  uint32
	}

	commonHead := OSCommonHead{
		Version:   2,
		SeqID:     seqID,
		AsynSeqID: 0,
		BodyLen:   uint32(28 + len(jsonData)),
		TimeHigh:  uint32(tm >> 32),
		TimeLow:   uint32(tm & 0xFFFFFFFF),
		Magic:     DAMagicNumber,
		Reserved:  0,
	}

	commonHeadBuf := new(bytes.Buffer)
	binary.Write(commonHeadBuf, binary.LittleEndian, commonHead)

	return append(commonHeadBuf.Bytes(), cmdData...), nil
}

// encodeIntention 编码意图格式（version=3）
func (c *DAClient) encodeIntention(args map[string]interface{}) ([]byte, error) {
	query, _ := args["query"].(string)
	uid, _ := args["uid"].(int64)
	if uid == 0 {
		if uidFloat, ok := args["uid"].(float64); ok {
			uid = int64(uidFloat)
		}
	}

	// IntentionCMDHead结构
	type IntentionCMDHead struct {
		UID     uint64
		Flag    uint16
		BodyLen uint16
	}

	intentionHead := IntentionCMDHead{
		UID:     uint64(uid),
		Flag:    0,
		BodyLen: uint16(len(query)),
	}

	// 构造CMDHead
	cmdHead := DACMDHead{
		Cmd: 15,                      // cmd=15 for intention
		Len: uint32(12 + len(query)), // IntentionCMDHead(12字节) + query
	}

	cmdHeadBuf := new(bytes.Buffer)
	binary.Write(cmdHeadBuf, binary.LittleEndian, cmdHead)
	binary.Write(cmdHeadBuf, binary.LittleEndian, intentionHead)
	cmdData := append(cmdHeadBuf.Bytes(), []byte(query)...)

	// 构造CommonHead
	seqID := c.generateSeqID()
	commonHead := DACommonHead{
		Version:   2,
		Flag:      0,
		Cbit:      0,
		Pad:       0,
		SeqID:     seqID,
		AsynSeqID: 0,
		BodyLen:   uint32(len(cmdData)),
		Timestamp: uint64(time.Now().UnixNano() / 1000),
		Magic:     DAMagicNumber,
		Reserved:  0,
	}

	commonHeadBuf := new(bytes.Buffer)
	binary.Write(commonHeadBuf, binary.LittleEndian, commonHead)

	return append(commonHeadBuf.Bytes(), cmdData...), nil
}

// encodeIntentionNew 编码新意图格式（version=4）
func (c *DAClient) encodeIntentionNew(args map[string]interface{}) ([]byte, error) {
	query, _ := args["query"].(string)
	uid, _ := args["uid"].(int64)
	if uid == 0 {
		if uidFloat, ok := args["uid"].(float64); ok {
			uid = int64(uidFloat)
		}
	}
	sceneType := 128
	if t, ok := args["type"].(int); ok {
		sceneType = t
	}

	// 构造JSON结构
	jsonStruct := map[string]interface{}{
		"querys": []map[string]interface{}{
			{
				"type":         sceneType,
				"intent":       query,
				"origin_query": query,
				"uid":          uid,
			},
		},
	}
	jsonData, err := sonic.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}

	// DAQueryParseReqNew结构
	type DAQueryParseReqNew struct {
		UID     uint32
		Flag    uint16
		BodyLen uint16
	}

	parseReq := DAQueryParseReqNew{
		UID:     uint32(uid),
		Flag:    0,
		BodyLen: uint16(len(jsonData)),
	}

	// 构造CMDHead
	cmdHead := DACMDHead{
		Cmd: 1,
		Len: uint32(8 + len(jsonData)), // DAQueryParseReqNew(8字节) + json
	}

	cmdHeadBuf := new(bytes.Buffer)
	binary.Write(cmdHeadBuf, binary.LittleEndian, cmdHead)
	binary.Write(cmdHeadBuf, binary.LittleEndian, parseReq)
	cmdData := append(cmdHeadBuf.Bytes(), jsonData...)

	// 构造CommonHead
	seqID := c.generateSeqID()
	commonHead := DACommonHead{
		Version:   2,
		Flag:      0,
		Cbit:      0,
		Pad:       0,
		SeqID:     seqID,
		AsynSeqID: 0,
		BodyLen:   uint32(len(cmdData)),
		Timestamp: uint64(time.Now().UnixNano() / 1000),
		Magic:     DAMagicNumber,
		Reserved:  0,
	}

	commonHeadBuf := new(bytes.Buffer)
	binary.Write(commonHeadBuf, binary.LittleEndian, commonHead)

	return append(commonHeadBuf.Bytes(), cmdData...), nil
}

// generateSeqID 生成序列ID
func (c *DAClient) generateSeqID() uint32 {
	// 生成一个随机的序列ID
	// 格式：随机数(1-3) + 进程ID后4位 + 时间戳后5位
	pid := os.Getpid()
	now := time.Now().UnixNano() / 1000000 // 毫秒

	randNum := rand.IntN(3) + 1
	pidPart := pid % 10000
	timePart := now % 100000

	seqID := uint32(randNum*100000000 + pidPart*10000 + int(timePart))
	return seqID
}
