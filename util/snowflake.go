package util

import (
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	// 雪花算法位分配(短ID版, 总计46bit, 生成ID不超过14位数字)
	// 1 bit  符号位(始终为0)
	// 40 bit 时间戳(毫秒级, 自纪元起可用约34.8年, 至2058年)
	// 3 bit  机器ID(0~7, 支持最多8个实例)
	// 3 bit  序列号(0~7, 同一毫秒内自增, 单机8000个/秒)
	// ID上限 2^46-1 = 70368744177663(14位), 在 JS Number.MAX_SAFE_INTEGER(2^53)之内
	timestampBits = 40
	workerBits    = 3
	sequenceBits  = 3

	maxWorkerId    = -1 ^ (-1 << workerBits)   // 7
	maxSequence    = -1 ^ (-1 << sequenceBits) // 7
	workerShift    = sequenceBits              // 3
	timestampShift = sequenceBits + workerBits // 6

	// 自定义纪元: 2024-01-01 00:00:00 UTC
	epoch = 1704067200000
)

type snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	sequence  int64
}

var sf *snowflake
var once sync.Once

// initSnowflake 初始化雪花生成器
func initSnowflake() {
	workerId := int64(1)
	if v := os.Getenv("SNOWFLAKE_WORKER_ID"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil && id >= 0 && id <= maxWorkerId {
			workerId = id
		}
	}
	sf = &snowflake{
		workerId: workerId,
	}
}

// NextId 生成雪花ID(线程安全)
// 时钟回退时沿用最后时间戳, 序列用尽时伪推进至下一毫秒, 保证ID严格唯一且单调递增
func NextId() int64 {
	once.Do(initSnowflake)

	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixMilli() - epoch

	// 时钟回退: 沿用最后时间戳, 避免生成重复ID
	if now < sf.timestamp {
		now = sf.timestamp
	}

	if now == sf.timestamp {
		sf.sequence = (sf.sequence + 1) & maxSequence
		if sf.sequence == 0 {
			// 当前毫秒序列号用尽: 推进至下一毫秒(不依赖真实时钟, 保证持续吞吐)
			now = sf.timestamp + 1
		}
	} else {
		sf.sequence = 0
	}

	sf.timestamp = now
	return (now << timestampShift) | (sf.workerId << workerShift) | sf.sequence
}
