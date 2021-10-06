package paste

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

func assertNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func assertEqual(t *testing.T, expect interface{}, got interface{}) {
	if expect != got {
		t.Fatal(fmt.Sprintf("expect %+v, got %+v", expect, got))
	}
}

func TestTemporaryGet(t *testing.T) {
	paste := Temporary{AbstractPaste: &AbstractPaste{}}
	paste.ExpireCount = 2
	paste.ExpireSecond = 10086

	assertNil(t, paste.Save())
	assertNil(t, paste.Get(""))
	assertNil(t, paste.Get(""))
	assertEqual(t, gorm.ErrRecordNotFound, paste.Get(""))
}

func TestTemporaryAutoDelete(t *testing.T) {
	var expireSecond uint64 = 1

	paste := Temporary{AbstractPaste: &AbstractPaste{}}
	paste.ExpireCount = 10086
	paste.ExpireSecond = expireSecond

	assertNil(t, paste.Save())
	assertNil(t, paste.Get(""))
	key := paste.Key
	paste.Key = "a1b2c3d4"
	time.Sleep(time.Second * time.Duration(expireSecond + 1))
	assertEqual(t, false, exist(key, &paste))
	assertEqual(t, gorm.ErrRecordNotFound, (&Temporary{AbstractPaste: &AbstractPaste{Key: key}}).Get(""))
}

func TestTemporaryConcurrentGet(t *testing.T) {
	var (
		expireCount   uint64 = 3
		concurrentCnt        = 20
		getCnt        uint64 = 0
	)

	paste := Temporary{AbstractPaste: &AbstractPaste{}}
	paste.ExpireCount = expireCount
	paste.ExpireSecond = 300
	assertNil(t, paste.Save())

	key := paste.Key

	ch := make([]chan bool, concurrentCnt)

	for i := 0; i < concurrentCnt; i++ {
		ch[i] = make(chan bool)
		go func(i int) {
			err := (&Temporary{AbstractPaste: &AbstractPaste{Key: key}}).Get("")
			ch[i] <- err == nil
		}(i)
	}

	for i := 0; i < concurrentCnt; i++ {
		get := <-ch[i]
		if get {
			getCnt += 1
			t.Logf("coroutine %d got the paste", i)
		}
	}

	assertEqual(t, expireCount, getCnt)
}

func TestMain(m *testing.M) {
	m.Run()
}
