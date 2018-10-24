package main

type CreateFileInput struct {
	FilePathPre string
	FileName    string
	FileSize    int
	FileUnit    string
}

type CreateMultiFilesInput struct {
	FilePathPre string
	FileNum     int
	FileSize    int
	FileUnit    string
}

type CreateMultiFilesOutput struct {
	FilePathPre string
	FileNames   []string
}
