package models

import "mime/multipart"

type FileOut struct {
	Filename       string `json:"filename"`
	Size           int    `json:"size"`
	UploadDate     int64  `json:"uploadDate"`
	UploadedUserId int    `json:"userId"`
	FileURL        string `json:"fileUrl"`
}

type FileUploadInput struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	UserId   int    `json:"userId"`
	FileData multipart.File
}

type FileUploadLogInput struct {
	Size        int64
	UploadDate  int64
	Filename    string
	UserId      int
	StorageLink string
}
