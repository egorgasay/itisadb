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

func TestStorage_GetFromIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name: "index",
				key:  "key",
			},
			want: "val",
		},
		{
			name: "ok",
			args: args{
				name: "index.innner.inner2",
				key:  "key2",
			},
			want: "val2",
		},
		{
			name: "not found",
			args: args{
				name: "index",
				key:  "key3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Errorf("CreateIndex() error = %v", err)
				}

				index, err := s.findIndex(tt.args.name)
				if err != nil {
					t.Errorf("findIndex() error = %v", err)
				}

				index.Set(tt.args.key, tt.want)
			}
			got, err := s.GetFromIndex(tt.args.name, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_SetToIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name:  "index",
				key:   "key",
				value: "val",
			},
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				name:  "index.innner.inner3",
				key:   "key2",
				value: "val2",
			},
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				name:  "index44",
				key:   "key",
				value: "val",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Errorf("CreateIndex() error = %v", err)
				}
			}
			if err := s.SetToIndex(tt.args.name, tt.args.key, tt.args.value, false); (err != nil) != tt.wantErr {
				t.Errorf("SetToIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				index, err := s.findIndex(tt.args.name)
				if err != nil {
					t.Errorf("findIndex() error = %v", err)
				}
				getValue, err := index.Get(tt.args.key)
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

func TestStorage_AttachToIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				dst: "index1",
				src: "index2",
			},
		},
		{
			name: "ok",
			args: args{
				dst: "index11.inner1",
				src: "index22",
			},
		},
		{
			name: "ok",
			args: args{
				dst: "index678.inner1.inner2.inner3",
				src: "index23.inner1",
			},
		},
		{
			name: "notFound",
			args: args{
				dst: "index99",
				src: "index98",
			},
			wantErr: ErrIndexNotFound,
		},
		{
			name: "circle",
			args: args{
				dst: "index1.inner1",
				src: "index1",
			},
			wantErr: ErrCircularAttachment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.wantErr, ErrIndexNotFound) {
				err := s.CreateIndex(tt.args.dst)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}

				err = s.CreateIndex(tt.args.src)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}
			}
			if err := s.AttachToIndex(tt.args.dst, tt.args.src); !errors.Is(err, tt.wantErr) {
				t.Fatalf("AttachToIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				original, err := s.findIndex(tt.args.src)
				if err != nil {
					t.Fatalf("findIndex() error = %v", err)
				}

				split := strings.Split(tt.args.src, ".")
				if len(split) == 0 {
					t.Fatalf("index name error, fix it!")
				}

				attached, err := s.findIndex(tt.args.dst + "." + split[len(split)-1])
				if err != nil {
					t.Fatalf("findIndex() error = %v", err)
				}

				attached.Set("key", "value")
				original.Set("key", "value")
				attached.Set("key1", "value1")
				original.Set("key1", "value1")

				originalMap := make(map[string]string, 10)

				original.Iter(func(k string, v ivalue) bool {
					originalMap[k] = v.GetValue()
					return false
				})

				attachedMap := make(map[string]string, 10)

				attached.Iter(func(k string, v ivalue) bool {
					attachedMap[k] = v.GetValue()
					return false
				})

				if !reflect.DeepEqual(originalMap, attachedMap) {
					t.Errorf("AttachToIndex() originalMap = %v, attachedMap = %v", originalMap, attachedMap)
				}
			}
		})
	}
}

func TestStorage_DeleteIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name: "index",
			},
		},
		{
			name: "ok",
			args: args{
				name: "index22",
			},
		},
		{
			name: "not found",
			args: args{
				name: "index78",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}
			}
			if err := s.DeleteIndex(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				_, err := s.findIndex(tt.args.name)
				if err == nil {
					t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestStorage_CreateIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name: "index",
			},
		},
		{
			name: "ok",
			args: args{
				name: "index.inner",
			},
		},
		{
			name: "ok",
			args: args{
				name: "index.inner.inner2.inner3.inner4",
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
			if err := s.CreateIndex(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				_, err := s.findIndex(tt.args.name)
				if err != nil {
					t.Errorf("CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestStorage_ToJSON(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name          string
		args          args
		structOfIndex map[string]string
		want          string
		wantErr       bool
	}{
		{
			name: "ok",
			args: args{
				name: "index",
			},
			structOfIndex: map[string]string{
				"key": "value",
			},
			want: "{\n\t\"isIndex\": true,\n\t\"name\": \"index\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t}\n\t]\n}",
		},
		{
			name: "ok#2",
			args: args{
				name: "index66",
			},
			structOfIndex: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
			want: "{\n\t\"isIndex\": true,\n\t\"name\": \"index66\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key3\",\n\t\t\t\"value\": \"value3\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key4\",\n\t\t\t\"value\": \"value4\"\n\t\t}\n\t]\n}",
		},
		{
			name: "ok#3",
			args: args{
				name: "index6.inner",
			},
			structOfIndex: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
			want: "{\n\t\"isIndex\": true,\n\t\"name\": \"inner\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key\",\n\t\t\t\"value\": \"value\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key1\",\n\t\t\t\"value\": \"value1\"\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"key2\",\n\t\t\t\"value\": \"value2\"\n\t\t}\n\t]\n}",
		},
		{
			name: "not found",
			args: args{
				name: "index67",
			},
			wantErr: true,
		},
		{
			name: "empty index",
			args: args{
				name: "index60",
			},
			structOfIndex: map[string]string{},
			want:          "{\n\t\"isIndex\": true,\n\t\"name\": \"index60\",\n\t\"values\": []\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}

				index, err := s.findIndex(tt.args.name)
				if err != nil {
					t.Fatalf("findIndex() error = %v", err)
				}

				for k, v := range tt.structOfIndex {
					index.Set(k, v)
				}
			}
			got, err := s.IndexToJSON(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexToJSON() error = %v, wantErr %v", err, tt.wantErr)
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

func TestStorage_GetIndex2(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
	}

	var want = "{\n\t\"isIndex\": true,\n\t\"name\": \"qwe\",\n\t\"values\": [\n\t\t{\n\t\t\t\"isIndex\": true,\n\t\t\t\"name\": \"edc\",\n\t\t\t\"values\": [\n\t\t\t\t{\n\t\t\t\t\t\"isIndex\": true,\n\t\t\t\t\t\"name\": \"rty\",\n\t\t\t\t\t\"values\": [\n\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\"isIndex\": false,\n\t\t\t\t\t\t\t\"name\": \"r3g\",\n\t\t\t\t\t\t\t\"value\": \"g3f\"\n\t\t\t\t\t\t}\n\t\t\t\t\t]\n\t\t\t\t},\n\t\t\t\t{\n\t\t\t\t\t\"isIndex\": false,\n\t\t\t\t\t\"name\": \"3g\",\n\t\t\t\t\t\"value\": \"3f\"\n\t\t\t\t}\n\t\t\t]\n\t\t},\n\t\t{\n\t\t\t\"isIndex\": false,\n\t\t\t\"name\": \"rfg\",\n\t\t\t\"value\": \"gwf\"\n\t\t}\n\t]\n}"
	err := s.CreateIndex("qwe")
	if err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	err = s.CreateIndex("qwe.edc")
	if err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	err = s.SetToIndex("qwe", "rfg", "gwf", false)
	if err != nil {
		t.Fatalf("SetToIndex() error = %v", err)
	}

	err = s.CreateIndex("qwe.edc.rty")
	if err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	err = s.SetToIndex("qwe.edc.rty", "r3g", "g3f", false)
	if err != nil {
		t.Fatalf("SetToIndex() error = %v", err)
	}

	err = s.SetToIndex("qwe.edc", "3g", "3f", false)
	if err != nil {
		t.Fatalf("SetToIndex() error = %v", err)
	}

	got, err := s.IndexToJSON("qwe")
	if err != nil {
		t.Errorf("IndexToJSON() error = %v, wantErr false", err)
		return
	}

	if got != want {
		t.Errorf("want:\n%s,\ngot:\n%s", want, got)
	}
}

func TestStorage_findIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name: "index",
			},
		},
		{
			name: "ok",
			args: args{
				name: "index2",
			},
		},
		{
			name: "not found",
			args: args{
				name: "index3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}
			}
			_, err := s.findIndex(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("findIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_Size(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name: "index",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				name: "index34",
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "ok",
			args: args{
				name: "index38",
			},
			want: 11,
		},
		{
			name: "not found",
			args: args{
				name: "index389",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}

				index, err := s.findIndex(tt.args.name)
				if err != nil {
					t.Fatalf("findIndex() error = %v", err)
				}

				for i := 0; uint64(i) < tt.want; i++ {
					index.Set(strconv.Itoa(i), "value")
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

func TestStorage_IsIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
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
				name: "index",
			},
			wantOk: true,
		},
		{
			name: "ok",
			args: args{
				name: "index678",
			},
			wantOk: true,
		},
		{
			name: "not ok",
			args: args{
				name: "index678.qwe",
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
				err := s.CreateIndex(path)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}

				index, err := s.findIndex(path)
				if err != nil {
					t.Fatalf("findIndex() error = %v", err)
				}

				index.Set(split[len(split)-1], "")
			} else if tt.wantOk {
				err := s.CreateIndex(tt.args.name)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}
			}

			gotOk := s.IsIndex(tt.args.name)

			if gotOk != tt.wantOk {
				t.Errorf("IsIndex() gotOk = %v, want %v", gotOk, tt.wantOk)
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
