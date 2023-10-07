package transactionlogger

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"itisadb/internal/storage"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestStorage_RestoreObjects(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int32
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			want: map[string]int32{
				"test":  1,
				"test2": 2,
				"test3": 3,
				"test4": 4,
			},
		},
		{
			name: "emptyMap",
			args: args{
				ctx: context.Background(),
			},
			want: map[string]int32{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage.Storage{
				mu: &sync.RWMutex{},
			}
			if !tt.wantErr {
				err := os.Mkdir(".objects", 0755)
				if err != nil && !os.IsExist(err) {
					t.Errorf("Mkdir() error = %v", err)
					return
				}
				for indx, serv := range tt.want {
					if err = s.SaveObjectLoc(tt.args.ctx, indx, serv); err != nil {
						t.Errorf("SaveObjects() error = %v", err)
					}
				}
				defer func() {
					err = os.RemoveAll(".objects")
					if err != nil {
						t.Errorf("RemoveAll() error = %v", err)
					}
				}()
			}
			got, err := s.RestoreObjects(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestoreObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestoreObjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_SaveObjects(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int32
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			want: map[string]int32{
				"test":  1,
				"test2": 2,
				"test3": 3,
				"test4": 4,
				"test5": 1,
				"test6": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage.Storage{
				mu: &sync.RWMutex{},
			}

			err := os.Mkdir(".objects", 0755)
			if err != nil && !os.IsExist(err) {
				t.Errorf("Mkdir() error = %v", err)
				return
			}

			var got = make(map[string]int32)
			defer func() {
				err = os.RemoveAll(".objects")
				if err != nil {
					t.Errorf("RemoveAll() error = %v", err)
				}
			}()

			for indx, serv := range tt.want {
				if err = s.SaveObjectLoc(tt.args.ctx, indx, serv); (err != nil) != tt.wantErr {
					t.Errorf("SaveObjects() error = %v, wantErr %v", err, tt.wantErr)
				}

				if !tt.wantErr {
					f, err := os.OpenFile(fmt.Sprintf(".objects/%d", serv), os.O_RDWR|os.O_CREATE, 0755)
					if err != nil {
						t.Errorf("OpenFile() error = %v", err)
					}
					defer f.Close()

					scanner := bufio.NewScanner(f)

					for scanner.Scan() {
						var object = scanner.Text()
						if object != "" {
							got[scanner.Text()] = serv
						}

						if errors.Is(err, io.EOF) {
							break
						} else if err != nil {
							t.Errorf("Fscanln() error = %v", err)
						}
					}
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveObjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}
