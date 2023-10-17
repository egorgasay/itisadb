package commands

import (
	"reflect"
	"testing"
)

func TestParseSet(t *testing.T) {
	type args struct {
		action string
		split  []string
	}
	tests := []struct {
		name    string
		args    args
		wantSc  SetCommand
		wantErr bool
	}{
		{
			name: "simple_Set",
			args: args{
				action: Set,
				split:  []string{`key`, `"value"`},
			},
			wantSc: SetCommand{
				action:   "set",
				key:      "key",
				value:    "value",
				server:   0,
				unique:   false,
				readOnly: false,
				level:    0,
			},
			wantErr: false,
		},
		{
			name: "simple_RSet_with_spaces",
			args: args{
				action: "rset",
				split:  []string{`key`, `"value we qjdjq dfwj "`},
			},
			wantSc: SetCommand{
				action:   Set,
				key:      "key",
				value:    "value we qjdjq dfwj ",
				server:   0,
				unique:   false,
				readOnly: true,
				level:    0,
			},
			wantErr: false,
		},
		{
			name: "Set_with_quotes",
			args: args{
				action: "set",
				split:  []string{`key`, `"qdqwf we """ qjdjq dfwj "`},
			},
			wantSc: SetCommand{
				action:   Set,
				key:      "key",
				value:    `qdqwf we """ qjdjq dfwj `,
				server:   0,
				unique:   false,
				readOnly: false,
				level:    0,
			},
			wantErr: false,
		},
		{
			name: "set_on_server",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `on`, `1`},
			},
			wantSc: SetCommand{
				action:   Set,
				key:      "key",
				value:    "value",
				server:   1,
				unique:   false,
				readOnly: false,
				level:    0,
			},
			wantErr: false,
		},
		{
			name: "set_on_server_with_level",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `on`, `1`, `level`, `2`},
			},
			wantSc: SetCommand{
				action:   Set,
				key:      "key",
				value:    "value",
				server:   1,
				unique:   false,
				readOnly: false,
				level:    2,
			},
			wantErr: false,
		},
		{
			name: "set_with_level_on_server",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `level`, `1`, `on`, `2`},
			},
			wantSc: SetCommand{
				action:   Set,
				key:      "key",
				value:    "value",
				server:   2,
				unique:   false,
				readOnly: false,
				level:    1,
			},
			wantErr: false,
		},
		{
			name: "missing_last_quote",
			args: args{
				action: Set,
				split:  []string{`key`, `"value`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "missing_last_quote#2",
			args: args{
				action: Set,
				split:  []string{`key`, `"value on 6`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "missing_first_quote",
			args: args{
				action: Set,
				split:  []string{`key`, `value"`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "missing_quotes",
			args: args{
				action: Set,
				split:  []string{`key`, `value`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "missing_quotes#2",
			args: args{
				action: Set,
				split:  []string{`key`, ``},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "no_value",
			args: args{
				action: Set,
				split:  []string{`key`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "no_key",
			args: args{
				action: Set,
				split:  []string{},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "unexpected_symbol_after_value",
			args: args{
				action: Set,
				split:  []string{`key`, `""b`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "unexpected_symbol_before_value",
			args: args{
				action: Set,
				split:  []string{`key`, `bghg g""`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "set_on_server_bad_format",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `on`, `we`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "set_on_server_with_level_bad_format",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `on`, `1`, `level`, `we`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "set_on_server_with_level_no_digit",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `on`, `level`, `we`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "set_on_server_with_level_no_digit#2",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `on`, "2", `level`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
		{
			name: "set_on_server_with_level_unexpected_keyword",
			args: args{
				action: "set",
				split:  []string{`key`, `"value"`, `wewrwrgw`},
			},
			wantSc:  SetCommand{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSc, err := ParseSet(tt.args.action, tt.args.split)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSc, tt.wantSc) {
				t.Errorf("ParseSet() gotSc = %v, want %v", gotSc, tt.wantSc)
			}
		})
	}
}
