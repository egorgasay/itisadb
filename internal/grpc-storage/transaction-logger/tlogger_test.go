package transactionlogger

import (
	"bufio"
	"fmt"
	"modernc.org/strutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	PATH = "testTL"
	os.Exit(m.Run())
}

func TestTLogger(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteSet("test"+fmt.Sprint(i), "test")
	}

	time.Sleep(1)

	kvf, err := os.Open(PATH + "/kv/1")
	if err != nil {
		t.Error(err)
		return
	}
	defer kvf.Close()

	sc := bufio.NewScanner(kvf)
	i := 0
	var ok bool

	for ; sc.Scan(); i++ {
		ok = true
		v := sc.Bytes()
		want := strutil.Base64Encode([]byte(fmt.Sprintf("%v %s %s", Set, "test"+fmt.Sprint(i), "test")))
		if !reflect.DeepEqual(v, want) {
			t.Errorf("want %s, got %s", want, v)
		}
	}

	if !ok {
		t.Error("Scan failed")
	}
}
