package message

import (
	"reflect"
	"testing"

	"github.com/ftomza/go-sspvo"
)

func TestNewCLSMessage(t *testing.T) {
	type args struct {
		cls string
	}
	tests := []struct {
		name string
		args args
		want *CLSMessage
	}{
		{
			name: "ok",
			args: args{
				cls: "test",
			},
			want: &CLSMessage{
				Message: Message{
					Fields: sspvo.JWTFields{sspvo.FieldCLS: "test"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCLSMessage(tt.args.cls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCLSMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
