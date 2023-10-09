package utils

import (
	"github.com/yueluoa/infrastructure/gerror"
	"sync"
	"time"
)

var once sync.Once
var idGenerator Generator

const (
	epoch        = int64(1683648000000) // 这是从2017年1月1日到当前时间的毫秒数，用于在ID中占用42位时间戳
	workerIDBits = uint(10)             // 节点ID的位数
	sequenceBits = uint(12)             // 序列号的位数
	maxWorkerID  = -1 ^ (-1 << workerIDBits)
	maxSequence  = -1 ^ (-1 << sequenceBits)
)

type Generator interface {
	Generate() int64
}

type IDGenerator struct {
	workerID      int64
	sequence      int64
	lastTimestamp int64
	mu            sync.Mutex
}

func init() {
	generator, err := NewIDGenerator(1)
	if err != nil {
		panic(err)
	}
	idGenerator = generator
}

func NewIDGenerator(workerID int64) (Generator, error) {
	var (
		ig  = &IDGenerator{}
		err error
	)

	once.Do(func() {
		if workerID < 0 || workerID > maxWorkerID {
			err = gerror.New("worker ID out of range")
			return
		}
		ig.workerID = workerID
		ig.sequence = 0
	})

	return ig, err
}

func (g *IDGenerator) Generate() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1000000 // 将当前时间转换为毫秒
	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			timestamp = g.waitNextMillis(timestamp)
		}
	} else {
		g.sequence = 0
	}

	if timestamp < g.lastTimestamp {
		return 0
	}

	g.lastTimestamp = timestamp

	return ((timestamp - epoch) << (workerIDBits + sequenceBits)) | (g.workerID << sequenceBits) | g.sequence
}

func (g *IDGenerator) waitNextMillis(currentTimestamp int64) int64 {
	for {
		timestamp := time.Now().UnixNano() / 1000000
		if timestamp > currentTimestamp {
			return timestamp
		}
	}
}

func GetIDGenerator() Generator {
	return idGenerator
}
