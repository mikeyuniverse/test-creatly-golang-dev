package models

import "mime/multipart"

type FileOut struct {
	Filename       string `json:"filename"`
	Volume         int    `json:"volume"`
	UploadDate     int64  `json:"uploadDate"`
	UploadedUserId int    `json:"userId"`
	FileURL        string `json:"fileUrl"`
}

type FileUploadInput struct {
	Filename string `json:"filename"`
	Volume   int64  `json:"volume"`
	UserId   int    `json:"userId"`
	FileData multipart.File
}

type FileUploadLogInput struct {
	Volume      int64
	UploadDate  int64
	Filename    string
	UserId      int
	StorageLink string
}
