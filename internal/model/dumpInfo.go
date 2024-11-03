package model

type DumpInfo struct {
	Filename string `json:"filename"`
	Type     string `json:"type"` // directory, file
	Mode     int    `json:"mode"` // file mode
}
