package fileutil

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestListDir(t *testing.T) {
	infos, err := ListDir("/Users/hant/Github/vbackup")
	if err != nil {
		return
	}
	marshal, err := json.Marshal(infos)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
func TestFixPath(t *testing.T) {
	fmt.Println(FixPath("/Users/hant/Github/vbackup"))
}
func TestGetFilePath(t *testing.T) {
	res := GetFilePath("../hant/go")
	fmt.Println(res)
}

func TestCopyFile(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{src: "/Users/hant/Github/vbackup/examples/conf/app.yml", dst: "/Users/hant/Github/vbackup/examples/conf/app1.yml"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFile(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
