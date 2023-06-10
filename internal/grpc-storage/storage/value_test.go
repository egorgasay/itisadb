package storage

import (
	"fmt"
	"github.com/dolthub/swiss"
	"reflect"
	"sync"
	"testing"
)

func Test_value_GetValue(t *testing.T) {
	v := &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100), isIndex: true}
	tests := []struct {
		name    string
		wantVal string
	}{
		{
			name:    "Get",
			wantVal: "Set",
		},
		{
			name:    "foo",
			wantVal: "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.value = tt.wantVal
			if gotVal := v.GetValue(); gotVal != tt.wantVal {
				t.Errorf("Get() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_value_IsEmpty(t *testing.T) {
	v := &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100), isIndex: true}
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "IsEmpty",
			want: true,
		},
		{
			name: "NotEmpty",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				v.next = nil
			} else {
				v.next = swiss.NewMap[string, ivalue](100)
			}
			if got := v.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_value_Get(t *testing.T) {
	v := &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100), isIndex: true}
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr bool
	}{
		{
			name: "simple",
			arg:  "Get",
			want: "Set",
		},
		{
			name: "ok",
			arg:  "fqqwdfqwdfqfkmk",
			want: "bafqwdfqwedfqfr",
		},
		{
			name:    "err",
			arg:     "fqqwdfqwdfqfkmk",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				v.Set(tt.name, tt.want)
			}
			got, err := v.Get(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_value_Size(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Size - 6",
			want: 6,
		},
		{
			name: "Size - 3",
			want: 3,
		},
		{
			name: "Size - 200",
			want: 200,
		},
	}
	for _, tt := range tests {
		v := &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100), isIndex: true}
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.want; i++ {
				v.Set(fmt.Sprint(i), tt.name)
			}
			if got := v.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_value_IsIndex(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "IsIndex",
			want: true,
		},
		{
			name: "false",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100), isIndex: tt.want}
			if got := v.IsIndex(); got != tt.want {
				t.Errorf("IsIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_value_AttachIndex(t *testing.T) {
	v := &value{mutex: &sync.RWMutex{}, next: swiss.NewMap[string, ivalue](100)}
	type args struct {
		name string
		val  ivalue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				val: NewIndex("foo", nil),
			},
		},
		{
			name: "ok2",
			args: args{
				val: NewIndex("qwdqdq", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.AttachIndex(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("AttachIndex() error = %v, wantErr %v", err, tt.wantErr)
			}

			vv, ok := v.next.Get(tt.args.val.Name())
			if !ok {
				t.Errorf("AttachIndex() ok = %v, wantOk %v", ok, true)
				return
			}
			if !reflect.DeepEqual(vv, tt.args.val) {
				t.Errorf("AttachIndex() = %v, want %v", vv, tt.args.val)
				return
			}
		})
	}
}

func Test_value_save(t *testing.T) {
	is := indexes{
		Map:     swiss.NewMap[string, ivalue](100),
		RWMutex: &sync.RWMutex{},
		path:    "C:\\tmp4",
	}
	s := &Storage{
		indexes: is,
	}

	for i := 0; i < 100; i++ {
		err := s.CreateIndex(fmt.Sprint(i))
		if err != nil {
			t.Errorf("CreateIndex() error = %v", err)
		}

		err = s.SetToIndex(fmt.Sprint(i), "key", fmt.Sprint(i), false)
		if err != nil {
			t.Errorf("Set() error = %v", err)
		}
	}

	if err := is.save(); err != nil {
		t.Errorf("save() error = %v", err)
		return
	}

	for i := 0; i < 100; i++ {
		if v, err := s.GetFromDiskIndex(fmt.Sprint(i), "key"); err != nil {
			t.Errorf("GetFromDiskIndex() error = %v", err)
		} else if v != fmt.Sprint(i) {
			t.Errorf("GetFromDiskIndex() = %v, want %v", v, fmt.Sprint(i))
		}
	}
}
