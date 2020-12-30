package model

type FileDetails struct {
	FileId		string `json:"file_id"`
	Name 		string	`json:"name"`
	Size		int64 	`json:"size"`
}