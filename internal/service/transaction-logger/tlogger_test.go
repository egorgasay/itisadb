package transactionlogger

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"itisadb/internal/grpc-storage/storage"
	"modernc.org/strutil"
)

func TestMain(m *testing.M) {
	DefaultPath = "testTL"
	os.Exit(m.Run())
}

func Test_WriteSet(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteSet("test"+fmt.Sprint(i), "test")
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", Set, "test"+fmt.Sprint(i), "test")
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

func Test_WriteDelete(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteDelete("test" + fmt.Sprint(i))
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", Delete, "test"+fmt.Sprint(i), "")
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

func Test_WriteSetToObject(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteSetToObject("test"+fmt.Sprint(i), "test"+fmt.Sprint(i), "test")
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", SetToObject, "test"+fmt.Sprint(i)+".test"+fmt.Sprint(i), "test")
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

func Test_WriteDeleteAttr(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteDeleteAttr("test"+fmt.Sprint(i), "test"+fmt.Sprint(i))
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", DeleteAttr, "test"+fmt.Sprint(i)+".test"+fmt.Sprint(i), "")
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

func Test_WriteAttach(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteAttach("dst"+fmt.Sprint(i), "src"+fmt.Sprint(i))
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", Attach, "dst"+fmt.Sprint(i), "src"+fmt.Sprint(i))
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

func Test_WriteDeleteObject(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 200; i++ {
		tl.WriteDeleteObject("object" + fmt.Sprint(i))
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", DeleteObject, "object"+fmt.Sprint(i), "")
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

func Test_WriteCreateObject(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	tl.Run()
	for i := 0; i < 100; i++ {
		tl.WriteCreateObject("object" + fmt.Sprint(i))
	}

	time.Sleep(3 * time.Second)

	files, err := os.ReadDir(DefaultPath)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) == 0 {
		t.Error("objects is empty")
		return
	}

	for _, f := range files {
		func() {
			fobject, err := os.Open(DefaultPath + "/" + f.Name())
			if err != nil {
				t.Error(err)
				return
			}
			defer fobject.Close()

			sc := bufio.NewScanner(fobject)
			i := 0
			var ok bool

			for ; sc.Scan(); i++ {
				ok = true
				w := fmt.Sprintf("%v %s %s", CreateObject, "object"+fmt.Sprint(i), "")
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

func Test_Restore(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		tl.Stop()
		time.Sleep(2 * time.Second)
		if err = os.RemoveAll(DefaultPath); err != nil {
			t.Errorf("remove failed: %s", err)
		}
	}()

	count := 119
	tl.Run()
	for i := 0; i < count; i++ {
		tl.WriteSet("test"+fmt.Sprint(i), "test"+fmt.Sprint(i))
	}

	time.Sleep(5 * time.Second)

	st, err := storage.New()
	if err != nil {
		t.Error(err)
		return
	}

	if err = tl.Restore(st); err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < count; i++ {
		v, err := st.Get("test" + fmt.Sprint(i))
		if err != nil {
			t.Error(err)
		}

		if v != "test"+fmt.Sprint(i) {
			t.Errorf("want %s, got %s", "test"+fmt.Sprint(i), v)
		}
	}
}
