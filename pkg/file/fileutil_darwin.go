//go:build darwin
// +build darwin

package fileutil

// // CopyFile
// func CopyFile(src, dst string) error {
// 	source, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer source.Close()

// 	destination, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer destination.Close()

// 	buffer := bufio.NewWriterSize(destination, 65536) // The buffer size can be customized, for example, 64KB or 1MB.

// 	_, err = io.Copy(buffer, source)
// 	if err != nil {
// 		return err
// 	}
// 	if err = buffer.Flush(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ReplaceHomeDir(path string) string {
// 	if strings.HasPrefix(FixPath(path), "~") {
// 		return strings.Replace(path, "~", HomeDir(), -1)
// 	}
// 	return path
// }

// func Exist(name string) bool {
// 	_, err := os.Stat(FixPath(name))
// 	if !os.IsNotExist(err) {
// 		return true
// 	}
// 	return false
// }

// func Mkdir(path string, mode os.FileMode) bool {
// 	err := os.Mkdir(path, mode)
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// func HomeDir() string {
// 	return os.Getenv("HOME")
// }

// func FixPath(path string) string {
// 	osType := runtime.GOOS
// 	if osType == "windows" {
// 		return strings.ReplaceAll(path, "/", "\\")
// 	}
// 	return path
// }

// func ListDir(path string) ([]*model.FileInfo, error) {
// 	nodes, err := ioutil.ReadDir(FixPath(path))
// 	if err != nil {
// 		return nil, err
// 	}
// 	var files []*model.FileInfo
// 	var dirs []*model.FileInfo
// 	osType := runtime.GOOS
// 	for _, dir := range nodes {
// 		var ct time.Time
// 		var statT *syscall.Stat_t
// 		if osType == "darwin" {
// 			statT = dir.Sys().(*syscall.Stat_t)
// 			ct = time.Unix(statT.Birthtimespec.Sec, statT.Birthtimespec.Nsec)
// 			sepa := ""
// 			if !strings.HasSuffix(path, string(filepath.Separator)) {
// 				sepa = string(filepath.Separator)
// 			}
// 			f := &model.FileInfo{
// 				Name:       dir.Name(),
// 				Path:       path + sepa + dir.Name(),
// 				IsDir:      dir.IsDir(),
// 				Mode:       dir.Mode().String(),
// 				ModTime:    dir.ModTime().Format(consts.Custom),
// 				Size:       dir.Size(),
// 				Gid:        int(statT.Gid),
// 				Uid:        int(statT.Uid),
// 				CreateTime: ct.Format(consts.Custom),
// 			}
// 			if dir.IsDir() {
// 				dirs = append(dirs, f)
// 			} else {
// 				files = append(files, f)
// 			}

// 		}
// 	}
// 	return append(dirs, files...), nil

// }

// // GetFilePath 获取文件路径
// func GetFilePath(file string) string {
// 	return filepath.Dir(file)
// }
