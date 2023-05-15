package storage

import (
	"github.com/dolthub/swiss"
	"sync"
	"testing"
)

func Test_value_Get(t *testing.T) {
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
			if gotVal := v.Get(); gotVal != tt.wantVal {
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
