package converterr

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"itisadb/internal/constants"
	"testing"
)

func TestConvertToGRPC(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			wantErr: nil,
		},
		{
			name: "notFound",
			args: args{
				err: constants.ErrNotFound,
			},
			wantErr: status.Error(codes.NotFound, constants.ErrNotFound.Error()),
		},
		{
			name: "ErrObjectNotFoundWithErrorf",
			args: args{
				err: fmt.Errorf("aaaa %w", constants.ErrObjectNotFound),
			},
			wantErr: status.Error(codes.ResourceExhausted, fmt.Errorf("aaaa %w", constants.ErrObjectNotFound).Error()),
		},
		{
			name: "ErrObjectNotFoundWithJoin",
			args: args{
				err: errors.Join(constants.ErrObjectNotFound, errors.New("aaaa")),
			},
			wantErr: status.Error(codes.ResourceExhausted, errors.Join(constants.ErrObjectNotFound, errors.New("aaaa")).Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertToGRPC(tt.args.err); !errors.Is(err, tt.wantErr) {
				t.Errorf("ConvertToGRPC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
