package storage

import (
	"errors"
	"fmt"
	"github.com/dolthub/swiss"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestStorage_Set(t *testing.T) {
	s := Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		key    string
		val    string
		unique bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				key:    "key",
				val:    "val",
				unique: false,
			},
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				key:    "key2",
				val:    "val2",
				unique: false,
			},
			wantErr: false,
		},
		{
			name: "unique error",
			args: args{
				key:    "key",
				val:    "val",
				unique: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Set(tt.args.key, tt.args.val, tt.args.unique); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, ok := s.ramStorage.Map.Get(tt.args.key); !ok {
					t.Errorf("Set() pair not found")
				}
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	s := Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				key: "key",
			},
			want: "val",
		},
		{
			name: "ok",
			args: args{
				key: "key2",
			},
			want: "val2",
		},
		{
			name: "not found",
			args: args{
				key: "key3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				s.ramStorage.Put(tt.args.key, tt.want)
			}

			got, err := s.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetFromObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
				key:  "key",
			},
			want: "val",
		},
		{
			name: "ok",
			args: args{
				name: "object.innner.inner2",
				key:  "key2",
			},
			want: "val2",
		},
		{
			name: "not found",
			args: args{
				name: "object",
				key:  "key3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Errorf("CreateObject() error = %v", err)
				}

				object, err := s.findObject(tt.args.name)
				if err != nil {
					t.Errorf("findObject() error = %v", err)
				}

				object.Set(tt.args.key, tt.want)
			}
			got, err := s.GetFromObject(tt.args.name, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_SetToObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name  string
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name:  "object",
				key:   "key",
				value: "val",
			},
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				name:  "object.innner.inner3",
				key:   "key2",
				value: "val2",
			},
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				name:  "object44",
				key:   "key",
				value: "val",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Errorf("CreateObject() error = %v", err)
				}
			}
			if err := s.SetToObject(tt.args.name, tt.args.key, tt.args.value, false); (err != nil) != tt.wantErr {
				t.Errorf("SetToObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				object, err := s.findObject(tt.args.name)
				if err != nil {
					t.Errorf("findObject() error = %v", err)
				}
				getValue, err := object.Get(tt.args.key)
				if err != nil {
					t.Errorf("GetValue() error = %v", err)
				}
				if getValue != tt.args.value {
					t.Errorf("GetValue() got = %v, want %v", getValue, tt.args.value)
				}
			}
		})
	}
}

func TestStorage_AttachToObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		dst string
		src string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				dst: "object1",
				src: "object2",
			},
		},
		{
			name: "ok",
			args: args{
				dst: "object11.inner1",
				src: "object22",
			},
		},
		{
			name: "ok",
			args: args{
				dst: "object678.inner1.inner2.inner3",
				src: "object23.inner1",
			},
		},
		{
			name: "notFound",
			args: args{
				dst: "object99",
				src: "object98",
			},
			wantErr: ErrObjectNotFound,
		},
		{
			name: "circle",
			args: args{
				dst: "object1.inner1",
				src: "object1",
			},
			wantErr: ErrCircularAttachment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.wantErr, ErrObjectNotFound) {
				err := s.CreateObject(tt.args.dst)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}

				err = s.CreateObject(tt.args.src)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}
			}
			if err := s.AttachToObject(tt.args.dst, tt.args.src); !errors.Is(err, tt.wantErr) {
				t.Fatalf("AttachToObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				original, err := s.findObject(tt.args.src)
				if err != nil {
					t.Fatalf("findObject() error = %v", err)
				}

				split := strings.Split(tt.args.src, ".")
				if len(split) == 0 {
					t.Fatalf("object name error, fix it!")
				}

				attached, err := s.findObject(tt.args.dst + "." + split[len(split)-1])
				if err != nil {
					t.Fatalf("findObject() error = %v", err)
				}

				attached.Set("key", "value")
				original.Set("key", "value")
				attached.Set("key1", "value1")
				original.Set("key1", "value1")

				originalMap := make(map[string]string, 10)

				original.Iter(func(k string, v iObject) bool {
					originalMap[k] = v.GetValue()
					return false
				})

				attachedMap := make(map[string]string, 10)

				attached.Iter(func(k string, v iObject) bool {
					attachedMap[k] = v.GetValue()
					return false
				})

				if !reflect.DeepEqual(originalMap, attachedMap) {
					t.Errorf("AttachToObject() originalMap = %v, attachedMap = %v", originalMap, attachedMap)
				}
			}
		})
	}
}

func TestStorage_DeleteObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
			},
		},
		{
			name: "ok",
			args: args{
				name: "object22",
			},
		},
		{
			name: "not found",
			args: args{
				name: "object78",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}
			}
			if err := s.DeleteObject(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				_, err := s.findObject(tt.args.name)
				if err == nil {
					t.Errorf("DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestStorage_CreateObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
			},
		},
		{
			name: "ok",
			args: args{
				name: "object.inner",
			},
		},
		{
			name: "ok",
			args: args{
				name: "object.inner.inner2.inner3.inner4",
			},
		},
		{
			name:    "wrong name",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.CreateObject(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("CreateObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				_, err := s.findObject(tt.args.name)
				if err != nil {
					t.Errorf("CreateObject() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestStorage_ToJSON(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name           string
		args           args
		structOfObject map[string]string
		want           string
		wantErr        bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
			},
			structOfObject: map[string]string{
				"key": "value",
			},
			want: "{\n\t\"isObject\": true,\n\t\"name\": \"object\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t}\n\t]\n}",
		},
		{
			name: "ok#2",
			args: args{
				name: "object66",
			},
			structOfObject: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
			want: "{\n\t\"isObject\": true,\n\t\"name\": \"object66\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key3\",\n\t\t\t\"value\": \"value3\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key4\",\n\t\t\t\"value\": \"value4\"\n\t\t}\n\t]\n}",
		},
		{
			name: "ok#3",
			args: args{
				name: "object6.inner",
			},
			structOfObject: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
			want: "{\n\t\"isObject\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
		},
		{
			name: "not found",
			args: args{
				name: "object67",
			},
			wantErr: true,
		},
		{
			name: "empty object",
			args: args{
				name: "object60",
			},
			structOfObject: map[string]string{},
			want:           "{\n\t\"isObject\": true,\n\t\"name\": \"object60\",\n\t\"values\": []\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}

				object, err := s.findObject(tt.args.name)
				if err != nil {
					t.Fatalf("findObject() error = %v", err)
				}

				for k, v := range tt.structOfObject {
					object.Set(k, v)
				}
			}
			got, err := s.ObjectToJSON(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ObjectToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			fmt.Printf("got:\n%s\nwant:\n%s\n", got, tt.want)
			if !cmpWordsInJSON(got, tt.want) {
				t.Errorf("Not equals")
			}
		})
	}
}

func cmpWordsInJSON(target1, target2 string) (equals bool) {
	m1 := make(map[rune]int)
	m2 := make(map[rune]int)

	target1 = strings.Replace(target1, "\t", "", -1)
	target2 = strings.Replace(target2, "\t", "", -1)

	for _, c := range target1 {
		m1[c]++
	}

	for _, c := range target2 {
		m2[c]++
	}

	return reflect.DeepEqual(m1, m2)
}

func TestStorage_GetObject2(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}

	var want = "{\n\t\"isObject\": true,\n\t\"name\": \"qwe\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isObject\": true,\n\t\t\t\"name\": \"edc\",\n\t\t\t\"values\": [\n\t\t\t\t{\n\t\t\t\t\t\"isObject\": true,\n\t\t\t\t\t\"name\": \"rty\",\n\t\t\t\t\t\"values\": [\n\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\"isObject\": false,\n\t\t\t\t\t\t\t\"name\": \"r3g\",\n\t\t\t\t\t\t\t\"value\": \"g3f\"\n\t\t\t\t\t\t}\n\t\t\t\t\t]\n\t\t\t\t},\n\t\t\t\t{\n\t\t\t\t\t\"isObject\": false,\n\t\t\t\t\t\"name\": \"3g\",\n\t\t\t\t\t\"value\": \"3f\"\n\t\t\t\t}\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t\"isObject\": false,\n\t\t\t\"name\": \"rfg\",\n\t\t\t\"value\": \"gwf\"\n\t\t}\n\t]\n}"
	err := s.CreateObject("qwe")
	if err != nil {
		t.Fatalf("CreateObject() error = %v", err)
	}

	err = s.CreateObject("qwe.edc")
	if err != nil {
		t.Fatalf("CreateObject() error = %v", err)
	}

	err = s.SetToObject("qwe", "rfg", "gwf", false)
	if err != nil {
		t.Fatalf("SetToObject() error = %v", err)
	}

	err = s.CreateObject("qwe.edc.rty")
	if err != nil {
		t.Fatalf("CreateObject() error = %v", err)
	}

	err = s.SetToObject("qwe.edc.rty", "r3g", "g3f", false)
	if err != nil {
		t.Fatalf("SetToObject() error = %v", err)
	}

	err = s.SetToObject("qwe.edc", "3g", "3f", false)
	if err != nil {
		t.Fatalf("SetToObject() error = %v", err)
	}

	got, err := s.ObjectToJSON("qwe")
	if err != nil {
		t.Errorf("ObjectToJSON() error = %v, wantErr false", err)
		return
	}

	if got != want {
		t.Errorf("want:\n%s,\ngot:\n%s", want, got)
	}
}

func TestStorage_findObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
			},
		},
		{
			name: "ok",
			args: args{
				name: "object2",
			},
		},
		{
			name: "not found",
			args: args{
				name: "object3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}
			}
			_, err := s.findObject(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("findObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_Size(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				name: "object34",
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				name: "object38",
			},
			want: 11,
		},
		{
			name: "not found",
			args: args{
				name: "object389",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}

				object, err := s.findObject(tt.args.name)
				if err != nil {
					t.Fatalf("findObject() error = %v", err)
				}

				for i := 0; uint64(i) < tt.want; i++ {
					object.Set(strconv.Itoa(i), "value")
				}
			}
			got, err := s.Size(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Size() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Size() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_IsObject(t *testing.T) {
	s := Storage{
		objects: objects{Map: swiss.NewMap[string, iObject](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{
			name: "ok",
			args: args{
				name: "object",
			},
			wantOk: true,
		},
		{
			name: "ok",
			args: args{
				name: "object678",
			},
			wantOk: true,
		},
		{
			name: "not ok",
			args: args{
				name: "object678.qwe",
			},
			wantOk: false,
		},
		{
			name: "not found",
			args: args{
				name: "ind4x678",
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			split := strings.Split(tt.args.name, ".")
			if !tt.wantOk && len(split) > 1 {
				path := strings.Join(split[:len(split)-1], ".")
				err := s.CreateObject(path)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}

				object, err := s.findObject(path)
				if err != nil {
					t.Fatalf("findObject() error = %v", err)
				}

				object.Set(split[len(split)-1], "")
			} else if tt.wantOk {
				err := s.CreateObject(tt.args.name)
				if err != nil {
					t.Fatalf("CreateObject() error = %v", err)
				}
			}

			gotOk := s.IsObject(tt.args.name)

			if gotOk != tt.wantOk {
				t.Errorf("IsObject() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	s := Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				key: "key",
			},
		},
		{
			name: "ok2",
			args: args{
				key: "key2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Set(tt.args.key, "value", false)
			if err != nil {
				t.Fatalf("Set() error = %v", err)
			}
			s.Delete(tt.args.key)
			_, err = s.Get(tt.args.key)
			if err == nil {
				t.Fatalf("Get() error = %v", err)
			}
		})
	}
}
