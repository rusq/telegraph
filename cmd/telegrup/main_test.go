// Command telegrup is the command line telegra.ph file uploader.
// Run with `-h` to get some help.
package main

import (
	"bytes"
	"errors"
	"testing"
)

type testError struct {
	Err string
}

func (e *testError) Error() string {
	return e.Err
}

func Test_printResults(t *testing.T) {
	type args struct {
		results []result
		asJson  bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"normal output",
			args{
				[]result{
					{Num: 42, Path: "/file/path"},
				},
				false,
			},
			"42: OK: https://telegra.ph/file/path\n",
			false,
		},
		{"error output",
			args{
				[]result{
					{Num: 42, Path: "", Err: errors.New("too many bits in your bytes")},
				},
				false,
			},
			"42: ERROR: too many bits in your bytes",
			false,
		},
		{"json output",
			args{
				[]result{
					{Num: 42, Path: "/file/path", Err: &testError{"boo boo"}},
				},
				true,
			},
			`[{"num":42,"path":"/file/path","error":{"Err":"boo boo"}}]
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := printResults(w, tt.args.results, tt.args.asJson); (err != nil) != tt.wantErr {
				t.Errorf("printResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("printResults() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
