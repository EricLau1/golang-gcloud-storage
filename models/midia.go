package models

import "mime/multipart"

type Midia struct {
	Name string
	Type string
	Link string
	Size int64
	File multipart.File
}
