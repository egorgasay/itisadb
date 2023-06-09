package storage

import (
	"github.com/dolthub/swiss"
	"os"
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
				name: "index/innner/inner2",
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
				name:  "index/innner/inner3",
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
		wantErr bool
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
				dst: "index11/inner1",
				src: "index22",
			},
		},
		{
			name: "ok",
			args: args{
				dst: "index678/inner1/inner2/inner3",
				src: "index23/inner1",
			},
		},
		{
			name: "not found",
			args: args{
				dst: "index99",
				src: "index98",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := s.CreateIndex(tt.args.dst)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}

				err = s.CreateIndex(tt.args.src)
				if err != nil {
					t.Fatalf("CreateIndex() error = %v", err)
				}
			}
			if err := s.AttachToIndex(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
				t.Fatalf("AttachToIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				original, err := s.findIndex(tt.args.src)
				if err != nil {
					t.Fatalf("findIndex() error = %v", err)
				}

				split := strings.Split(tt.args.src, "/")
				if len(split) == 0 {
					t.Fatalf("index name error, fix it!")
				}

				attached, err := s.findIndex(tt.args.dst + "/" + split[len(split)-1])
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
				name: "index/inner",
			},
		},
		{
			name: "ok",
			args: args{
				name: "index/inner/inner2/inner3/inner4",
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

func TestStorage_GetIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "index",
			},
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name: "ok#2",
			args: args{
				name: "index66",
			},
			want: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			},
		},
		{
			name: "ok#3",
			args: args{
				name: "index6/inner",
			},
			want: map[string]string{
				"key":  "value",
				"key1": "value1",
				"key2": "value2",
			},
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
			want: map[string]string{},
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

				for k, v := range tt.want {
					index.Set(k, v)
				}
			}
			got, err := s.GetIndex(tt.args.name, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIndex() got = %v, want %v", got, tt.want)
			}
		})
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
				name: "index678/qwe",
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
			split := strings.Split(tt.args.name, "/")
			if !tt.wantOk && len(split) > 1 {
				path := strings.Join(split[:len(split)-1], "/")
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

func TestStorage_GetFromDisk(t *testing.T) {
	s := Storage{
		ramStorage: ramStorage{Map: swiss.NewMap[string, string](10), RWMutex: &sync.RWMutex{}, path: "C:\\tmp2"},
	}
	err := s.Set("test_key", "test_value", false)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	err = s.Set("test_key2", "test_value2", false)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	err = s.ramStorage.save()
	if err != nil {
		t.Fatalf("indexes.save() error = %v", err)
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
				key: "test_key",
			},
			want: "test_value",
		},
		{
			name: "ok",
			args: args{
				key: "test_key2",
			},
			want: "test_value2",
		},
		{
			name: "notFound",
			args: args{
				key: "test_key3",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFromDisk(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromDisk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromDisk() got = %v, want %v", got, tt.want)
			}
		})
	}

	if err := os.RemoveAll("C:\\tmp2"); err != nil {
		t.Fatalf("RemoveAll() error = %v", err)
	}
}

func TestStorageGetFromDiskIndex(t *testing.T) {
	s := Storage{
		indexes: indexes{Map: swiss.NewMap[string, ivalue](10), RWMutex: &sync.RWMutex{}, path: "C:\\tmp3"},
	}

	err := s.CreateIndex("test_index")
	if err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	err = s.CreateIndex("test_index/inner")
	if err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	err = s.SetToIndex("test_index", "test_attr_key", "test_value", false)
	if err != nil {
		t.Fatalf("SetToIndex() error = %v", err)
	}

	err = s.SetToIndex("test_index/inner", "test_attr_key2", "test_value2", false)
	if err != nil {
		t.Fatalf("SetToIndex() error = %v", err)
	}

	err = s.indexes.save()
	if err != nil {
		return
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
				name: "test_index",
				key:  "test_attr_key",
			},
			want: "test_value",
		},
		{
			name: "inner",
			args: args{
				name: "test_index/inner",
				key:  "test_attr_key2",
			},
			want: "test_value2",
		},
		{
			name: "notFound",
			args: args{
				name: "test_index/inne3r",
				key:  "test_attr_k3ey2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFromDiskIndex(tt.args.name, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromDiskIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromDiskIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}
