package reader

import (
	"regexp"
	"testing"
)

func TestRead(t *testing.T) {
	stream := make(chan string)
	reg, _ := regexp.Compile(`(2[0-9]{3})-(0[1-9]|1[012])-([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2},\d+`)

	type args struct {
		id       uint64
		filePath string
		prefix   regexp.Regexp
		stream   chan string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "we",
			args: args{
				id:       1,
				filePath: "/Users/joyo/Desktop/error.log",
				prefix:   *reg,
				stream:   stream,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Read(tt.args.id, tt.args.filePath, tt.args.prefix, tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
