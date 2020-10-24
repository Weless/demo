package meta

import (
	"net_disk_demo/db"
)

// FileMeta:文件元信息结构
type FileMeta struct {
	// 文件sha1值
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta:新增或者更新文件元信息
func UpdateFileMeta(fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
}

// UpdateFileMetaDB:新增或者更新文件元信息到mysql
func UpdateFileMetaDB(fMeta FileMeta) bool {
	return db.OnFileUploadFinished(fMeta.FileSha1, fMeta.FileName, fMeta.Location, fMeta.FileSize)
}

// GetFileMeta:根据fileSha1获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSh1 string) (FileMeta, error) {
	data, err := db.GetFileMeta(fileSh1)
	if err != nil {
		return FileMeta{}, err
	}
	return FileMeta{
		FileSha1: data.FileHash,
		FileName: data.FileName.String,
		FileSize: data.FileSize.Int64,
		Location: data.FileAddr.String,
	}, nil
}

// RemoveFileMeta:根据fileSha1删除文件
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
