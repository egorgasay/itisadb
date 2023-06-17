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

func Test_WriteSet(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		time.Sleep(1 * time.Second)
		if err = os.RemoveAll(PATH); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteSet("test"+fmt.Sprint(i), "test")
	}

	time.Sleep(1 * time.Second)

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

func Test_WriteDelete(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err = os.RemoveAll(PATH); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteDelete("test" + fmt.Sprint(i))
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
		want := strutil.Base64Encode([]byte(fmt.Sprintf("%v %s %s", Delete, "test"+fmt.Sprint(i), "")))
		if !reflect.DeepEqual(v, want) {
			t.Errorf("want %s, got %s", want, v)
		}
	}

	if !ok {
		t.Error("Scan failed")
	}
}

func Test_WriteSetToIndex(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		time.Sleep(1 * time.Second)
		if err = os.RemoveAll(PATH); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteSetToIndex("test"+fmt.Sprint(i), "test"+fmt.Sprint(i), "test")
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(PATH)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("indexes is empty")
		return
	}

	for _, f := range files {
		func() {
			findex, err := os.Open(PATH + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer findex.Close()

			sc := bufio.NewScanner(findex)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", SetToIndex, "test"+fmt.Sprint(i)+".test"+fmt.Sprint(i), "test")
				v := sc.Bytes()
				want := strutil.Base64Encode([]byte(w))
				dec, err := strutil.Base64Decode(v)
				if err != nil {
					t.Error(err)
					return
				}
				if !reflect.DeepEqual(v, want) {
					t.Errorf("want %s, got %s\n decoded: %v\n want decoded: %v", want, v, string(dec), w)
				}
			}

			if !ok {
				t.Error("Scan failed")
			}
		}()

	}

}
