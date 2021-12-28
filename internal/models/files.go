package models

type FileOut struct {
	Filename string `json:"filename",bson:"filename"`
	Size     int    `json:"size",bson:"size"`
	Date     int64  `json:"uploadDate",bson:"date"`
	UserId   string `json:"userId",bson:"userId"`
	Url      string `json:"url",bson:"url"`
}

type FileUploadInput struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	UserId   string `json:"userId"`
	FileData []byte
}

type FileUploadLogInput struct {
	Size       int64  `bson:"size"`
	UploadDate int64  `bson:"date"`
	Filename   string `bson:"filename"`
	UserId     string `bson:"userId"`
	Url        string `bson:"url"`
}
