package models

type SourceType int

const (
	SourceTypeURL SourceType = iota
	SourceTypeBase64
	SourceTypeFile
	SourceTypeFolder
)
