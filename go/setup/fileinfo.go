package setup

type FileType int64

const (
	FileTypeUnkonwn FileType = iota
	FileTypeFile
	FileTypeDir
	FileTypeSoftlink
)

type FileInfo struct {
	Type     FileType
	Name     string
	Linkname string
	Size     int64
	Mode     int64
}
