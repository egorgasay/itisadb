package servernumber

import (
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	type args struct {
		server int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				server: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Set(tt.args.server); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			f, err := os.Open("server_number")
			if err != nil {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer f.Close()

			var b = make([]byte, 100)
			n, err := f.Read(b)
			if err != nil {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			b = b[:n]
			if string(b) != "1" {
				t.Errorf("Set() want [1], got [%s]", string(b))
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				dir: ".",
			},
			want: 1,
		},
		{
			name: "error",
			args: args{
				dir: "/root",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.dir)
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
