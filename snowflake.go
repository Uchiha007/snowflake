package snowflake

import (
	"errors"
	"sync"
	"time"
)

/**
 * SnowFlake
 * 64bit: 1(null)+42(time)+9(machine)+12(sequenceId)=64
 */

const (
	twepoch        = int64(0)                         //开始时间截（1970）
	workerIdBits   = uint(9)                          //机器id所占的位数
	sequenceBits   = uint(12)                         //序列所占的位数
	workerIdMax    = int64(-1 ^ (-1 << workerIdBits)) //支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //序列号掩码
	workerIdShift  = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workerIdBits      //时间戳左移位数
)

// A Snowflake struct holds the basic information needed for a snowflake generator worker
type Worker struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerId      int64
	sequence      int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerIdMax {
		return nil, errors.New("Worker ID excess of quantity")
	}

	return &Worker{
		workerId: workerId,
	}, nil
}

// Generate creates and returns a unique snowflake ID
func (w *Worker) NextId() (id int64, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6
	// the system clock goes back
	if timestamp < w.lastTimestamp {
		return id, errors.New("Clock moved backwards")
	}

	if w.lastTimestamp == timestamp {
		w.sequence = (w.sequence + 1) & sequenceMask
		if w.sequence == 0 {
			timestamp = w.tilNextMillis()
		}
	} else {
		w.sequence = 0
	}
	w.lastTimestamp = timestamp

	d := int64((timestamp-twepoch)<<timestampShift | (w.workerId << workerIdShift) | (w.sequence))
	return d, nil
}

func (w *Worker) FixedId(timestamp, sequence int64) (id int64) {
	w.sequence = sequence & sequenceMask
	id = int64((timestamp-twepoch)<<timestampShift | (w.workerId << workerIdShift) | w.sequence)
	return id
}

func (w *Worker) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano() / 1e6
	for timestamp <= w.lastTimestamp {
		timestamp = time.Now().UnixNano() / 1e6
	}
	return timestamp
}

func (w *Worker) BreakDown(id int64) (int64, int64, int64) {
	timestamp := (id >> timestampShift) + twepoch
	sequence := id ^ (timestamp-twepoch)<<timestampShift ^ w.workerId<<workerIdShift
	return timestamp, w.workerId, sequence
}
