package omicafe

import (
	"os"
	"path/filepath"
	"strings"
)

type FileManager struct {
	BaseDir string
}

// 初始化文件管理器
func NewFileManager(baseDir string) *FileManager {
	_ = os.MkdirAll(baseDir, 0755) // 确保目录存在
	return &FileManager{BaseDir: baseDir}
}

// 写文件
func (fm *FileManager) WriteFile(filename string, data []byte) error {
	path := filepath.Join(fm.BaseDir, fm.sanitizeFilename(filename))
	return os.WriteFile(path, data, 0644)
}

// 读文件
func (fm *FileManager) ReadFile(filename string) ([]byte, error) {
	path := filepath.Join(fm.BaseDir, fm.sanitizeFilename(filename))
	return os.ReadFile(path)
}

// 删除文件
func (fm *FileManager) DeleteFile(filename string) error {
	path := filepath.Join(fm.BaseDir, fm.sanitizeFilename(filename))
	return os.Remove(path)
}

// 转换 URL 路径为安全的文件名
func (fm *FileManager) sanitizeFilename(filename string) string {
	return strings.ReplaceAll(filename, "/", "@")
}
