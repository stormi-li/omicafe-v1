package omicafe

func NewFileCache(dir string, size int) *FileCache {
	return newFileCache(dir, size)
}
