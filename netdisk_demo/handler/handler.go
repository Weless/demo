package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net_disk_demo/meta"
	"net_disk_demo/utils"
	"os"
	"strings"
	"time"
)

// UploadHandler:上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传的html页面
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		// 1. 接收文件
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "/Users/joey/joey/go_project/netdisk_demo/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 本地创建一个新文件用来接收文件流
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("fail to copy file to newfile, err:%s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		fmt.Println("filehash is " + fileMeta.FileSha1)

		//meta.UpdateFileMeta(fileMeta)

		ok := meta.UpdateFileMetaDB(fileMeta)
		if ok {
			fmt.Println("upload file successfully !!!")
		} else {
			fmt.Println("failed to upload file")
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler:上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler:获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("failed to parse form, err:%s\n", err.Error())
		return
	}
	fileHash := r.Form.Get("filehash")
	//fMeta := meta.GetFileMeta(fileHash)
	fMeta, _ := meta.GetFileMetaDB(fileHash)
	fmt.Printf("fileMeta is %+v\n", fMeta)

	data, err := json.Marshal(fMeta)
	if err != nil {
		fmt.Printf("failed to marshal fMeta, err:%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler:下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("failed to parse form, err:%s\n", err.Error())
		return
	}
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		fmt.Printf("failed to open file, err:%s\n", err.Error())
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("failed to read data from file, err:%s\n", err.Error())
		return
	}

	// 浏览器下载
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")

	w.Write(data)
}

// FileMetaUpdateHandler:更新元信息接口（重命名）
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName

	// 重命名
	pathParts := strings.Split(curFileMeta.Location, "/")
	newFileLocation := strings.Join(pathParts[:len(pathParts)-1], "/") + "/" + newFileName

	fmt.Printf("newFilePath is %s\n", newFileLocation)

	err := os.Rename(curFileMeta.Location, newFileLocation)
	if err != nil {
		fmt.Printf("failed to rename old file to new file, err:%s\n", err.Error())
		return
	}

	curFileMeta.Location = newFileLocation

	// 更新文件元信息
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		fmt.Printf("failed to marshal data, err:%s\n", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler:删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	// 先删除文件
	fMeta := meta.GetFileMeta(fileSha1)
	fmt.Println("file location is " + fMeta.Location)
	err := os.Remove(fMeta.Location)
	if err != nil {
		fmt.Printf("failed to delete file, err: %s\n", err.Error())
		return
	}

	// 删除元信息
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete successful!"))
}
