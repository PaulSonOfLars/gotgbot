package gotgbot

import (
	"mime/multipart"
	"strconv"
	"testing"
)

func Test_getFieldContents(t *testing.T) {
	var testString = "test"
	var testInt = 42

	type args struct {
		v any
		k string
		w *multipart.Writer
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "regular string",
			args: args{
				v: testString,
			},
			want:    testString,
			wantErr: false,
		}, {
			name: "string pointer",
			args: args{
				v: &testString,
			},
			want:    testString,
			wantErr: false,
		}, {
			name: "integer",
			args: args{
				v: testInt,
			},
			want:    strconv.Itoa(testInt),
			wantErr: false,
		}, {
			name: "integer pointer",
			args: args{
				v: &testInt,
			},
			want:    strconv.Itoa(testInt),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldContents(tt.args.v, tt.args.k, tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFieldContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFieldContents() got = %v, want %v", got, tt.want)
			}
		})
	}
}
