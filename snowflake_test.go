package snowflake

import (
	"testing"
)

var (
	testWorkerId  = int64(8)             // less than 512
	testTimestamp = int64(1577808000000) // 2020 01-01
	testWorker    *Worker
)

func TestNewWorker(t *testing.T) {
	worker, err := NewWorker(int64(1 << 9))
	if err != nil {
		t.Logf("%v,%+v", err, worker)
	}

	testWorker, err = NewWorker(testWorkerId)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkWorker_NextId(b *testing.B) {
	worker, _ := NewWorker(1)
	for i := 0; i < b.N; i++ {
		_, _ = worker.NextId()
	}
}

func TestWorker_NextId(t *testing.T) {
	_, err := testWorker.NextId()
	_, err = testWorker.NextId()
	if err != nil {
		t.Error(err)
	}
}

func TestWorker_FixedId(t *testing.T) {
	id := testWorker.FixedId(1577808000000, 1)
	if id != 3308903202816032769 {
		t.Error("not match")
	}
	timestamp, workerId, sequence := testWorker.BreakDown(id)
	if timestamp != testTimestamp || workerId != testWorkerId || sequence != 1 {
		t.Error("failed")
	}
}

//func TestWorker_BreakDown(t *testing.T) {
//	id, _ := testWorker.NextId()
//	timestamp, workerId, sequence := testWorker.BreakDown(id)
//	t.Log(timestamp, workerId, sequence)
//	if workerId != testWorkerId || sequence != 2 || timestamp > time.Now().UnixNano()/1e6 {
//		t.Error("break down error")
//	}
//}
