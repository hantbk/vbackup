package fileutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExist(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_exist_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "existing_file",
			args: tmpFile.Name(),
			want: true,
		},
		{
			name: "non_existing_file",
			args: "/tmp/non_existent_file",
			want: false,
		},
		{
			name: "existing_directory",
			args: "/tmp",
			want: true,
		},
		{
			name: "empty_path",
			args: "",
			want: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exist(tt.args); got != tt.want {
				t.Errorf("Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixPath(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "normal_path",
			args: "/path/to/file",
			want: "/path/to/file",
		},
		{
			name: "with_double_slash",
			args: "/path//to/file",
			want: "/path//to/file",
		},
		{
			name: "empty_path",
			args: "",
			want: "",
		},
		{
			name: "relative_path",
			args: "./path/to/file",
			want: "./path/to/file",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FixPath(tt.args); got != tt.want {
				t.Errorf("FixPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "absolute_path",
			args: "/path/to/file.txt",
			want: "/path/to",
		},
		{
			name: "relative_path",
			args: "path/to/file.txt",
			want: "path/to",
		},
		{
			name: "file_in_root",
			args: "/file.txt",
			want: "/",
		},
		{
			name: "current_dir",
			args: "file.txt",
			want: ".",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilePath(tt.args); got != tt.want {
				t.Errorf("GetFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeDir(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tests := []struct {
		name     string
		homeEnv  string
		want     string
	}{
		{
			name: "valid_home",
			homeEnv: "/home/testuser",
			want: "/home/testuser",
		},
		{
			name: "empty_home",
			homeEnv: "",
			want: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.homeEnv)
			if got := HomeDir(); got != tt.want {
				t.Errorf("HomeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_listdir_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files and directories
	files := []struct {
		name string
		content string
		mode os.FileMode
	}{
		{"file1.txt", "content1", 0644},
		{"file2.txt", "content2", 0644},
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f.name)
		if err := os.WriteFile(path, []byte(f.content), f.mode); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.Mkdir(filepath.Join(tmpDir, "testdir"), 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		want    int
		wantErr bool
	}{
		{
			name: "valid_directory",
			path: tmpDir,
			want: 3, // 2 files + 1 directory
			wantErr: false,
		},
		{
			name: "non_existent_directory",
			path: "/non/existent/path",
			want: 0,
			wantErr: true,
		},
		{
			name: "empty_path",
			path: "",
			want: 0,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("ListDir() got %v files, want %v", len(got), tt.want)
			}
		})
	}
}

func TestMkdir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_mkdir_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		mode os.FileMode
		want bool
	}{
		{
			name: "new_directory",
			path: filepath.Join(tmpDir, "newdir"),
			mode: 0755,
			want: true,
		},
		{
			name: "existing_directory",
			path: tmpDir,
			mode: 0755,
			want: false,
		},
		{
			name: "invalid_path",
			path: "/non/existent/path/newdir",
			mode: 0755,
			want: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mkdir(tt.path, tt.mode); got != tt.want {
				t.Errorf("Mkdir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceHomeDir(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	testHome := "/home/testuser"
	os.Setenv("HOME", testHome)

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "path_with_tilde",
			path: "~/documents/file.txt",
			want: "/home/testuser/documents/file.txt",
		},
		{
			name: "absolute_path",
			path: "/absolute/path/file.txt",
			want: "/absolute/path/file.txt",
		},
		{
			name: "empty_path",
			path: "",
			want: "",
		},
		{
			name: "tilde_in_middle",
			path: "/path/~/file.txt",
			want: "/path/~/file.txt",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceHomeDir(tt.path); got != tt.want {
				t.Errorf("ReplaceHomeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_copyfile_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcContent := []byte("test content")
	srcFile := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(srcFile, srcContent, 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{
			name: "successful_copy",
			src: srcFile,
			dst: filepath.Join(tmpDir, "dest.txt"),
			wantErr: false,
		},
		{
			name: "source_not_exist",
			src: filepath.Join(tmpDir, "nonexistent.txt"),
			dst: filepath.Join(tmpDir, "dest2.txt"),
			wantErr: true,
		},
		{
			name: "invalid_destination",
			src: srcFile,
			dst: "/non/existent/path/dest.txt",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyFile(tt.src, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				destContent, err := os.ReadFile(tt.dst)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(destContent, srcContent) {
					t.Errorf("CopyFile() content mismatch: got %v, want %v", destContent, srcContent)
				}
			}
		})
	}
}