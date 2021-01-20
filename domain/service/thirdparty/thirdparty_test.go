package thirdparty

import (
	"bytes"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"reflect"
	"testing"
)

func Test_getTokenAttributes(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want []*types.TokenAttribute
	}{
		{
			name: "Normal case",
			args: args{
				query: `"halo"  "tes"  "123"`,
			},
			want: []*types.TokenAttribute{
				{
					Value: "halo",
				},
				{
					Value: "tes",
				},
				{
					Value: "123",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				query: `"halo"   `,
			},
			want: []*types.TokenAttribute{
				{
					Value: "halo",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				query: `"halo"  "tes"                 "123" "1234""234"`,
			},
			want: []*types.TokenAttribute{
				{
					Value: "halo",
				},
				{
					Value: "tes",
				},
				{
					Value: "123",
				},
				{
					Value: "1234",
				},
				{
					Value: "234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokenAttributes(tt.args.query); !reflect.DeepEqual(got, tt.want) {
				strbGot := bytes.Buffer{}
				for _, g := range got {
					strbGot.WriteString("\"" + g.Value + "\" ")
				}
				strbWant := bytes.Buffer{}
				for _, g := range tt.want {
					strbWant.WriteString("\"" + g.Value + "\" ")
				}
				t.Errorf("getTokenAttributes() = %v, want %v", strbGot.String(), strbWant.String())
			}
		})
	}
}