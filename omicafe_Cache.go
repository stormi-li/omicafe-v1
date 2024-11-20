package omicafe

import (
	"container/list"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CacheItem struct {
	filename string // URL 路径
	filepath string // 文件在磁盘上的路径
	size     int    // 文件大小
}

type FileCache struct {
	lock         sync.RWMutex
	curSize      int
	maxCacheSize int
	cacheDir     string
	itemMap      map[string]*list.Element
	itemList     *list.List // 用于 LRU 策略
}

// 初始化文件缓存系统，读取现有缓存目录内容
func newFileCache(cacheDir string, maxSize int) *FileCache {
	fileCache := &FileCache{
		itemMap:  make(map[string]*list.Element),
		itemList: list.New(),
	}

	// 设置缓存目录和最大容量
	fileCache.cacheDir = cacheDir
	fileCache.maxCacheSize = maxSize

	// 创建缓存目录（如果不存在）
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(err)
	}

	// 遍历缓存目录并加载文件信息
	err := filepath.WalkDir(cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录本身
		if d.IsDir() {
			return nil
		}

		// 获取文件信息并构建缓存项
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		// 文件大小检查
		fileSize := int(fileInfo.Size())
		if fileSize > maxSize {
			return nil // 跳过超出容量的文件
		}

		filename := filepath.Base(path)
		cacheItem := &CacheItem{filename: filename, filepath: path, size: fileSize}
		elem := fileCache.itemList.PushFront(cacheItem)
		fileCache.itemMap[filename] = elem
		fileCache.curSize += fileSize
		// 超出容量时，清理最旧的缓存文件
		for fileCache.curSize > maxSize {
			fileCache.removeOldest()
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
	return fileCache
}

// 将文件添加到磁盘缓存，超过容量时使用 LRU 策略清理
func (fileCache *FileCache) Set(key string, data []byte) {
	if len(data) == 0 {
		return
	}
	filename := strings.ReplaceAll(key, "/", "@")

	cachePath := filepath.Join(fileCache.cacheDir, filename)

	// 文件大小超过缓存最大容量时直接返回
	fileSize := len(data)
	if fileSize > fileCache.maxCacheSize {
		return
	}

	// 清理超出容量的旧文件，确保有足够空间
	fileCache.lock.Lock()
	for fileCache.curSize+fileSize > fileCache.maxCacheSize {
		fileCache.removeOldest()
	}
	fileCache.lock.Unlock()

	// 写入文件到磁盘
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return
	}

	// 新建缓存项并添加到 LRU 列表
	fileCache.lock.Lock()
	if fileCache.itemMap[filename] == nil {
		cacheItem := &CacheItem{filename: filename, filepath: cachePath, size: fileSize}
		elem := fileCache.itemList.PushFront(cacheItem)
		fileCache.itemMap[filename] = elem
		fileCache.curSize += fileSize
	}
	fileCache.lock.Unlock()
}

// 读取缓存，如果命中则返回 true
func (fileCache *FileCache) Get(key string) []byte {
	fileCache.lock.RLock()
	defer fileCache.lock.RUnlock()
	filename := strings.ReplaceAll(key, "/", "@")
	// 转换路径并检查是否存在于缓存中
	elem, found := fileCache.itemMap[filename]
	if !found {
		return []byte{} // 缓存未命中
	}
	// 移动缓存项到列表前端
	fileCache.itemList.MoveToFront(elem)

	cacheItem := elem.Value.(*CacheItem)
	data, err := os.ReadFile(cacheItem.filepath)
	if err != nil {
		return []byte{}
	}
	return data
}

// 删除指定缓存文件
func (fileCache *FileCache) Del(key string) {
	filename := strings.ReplaceAll(key, "/", "@")
	// 检查缓存项是否存在
	fileCache.lock.Lock()
	elem, found := fileCache.itemMap[filename]
	if !found {
		return // 缓存项不存在
	}

	// 删除文件及缓存项
	cacheItem := elem.Value.(*CacheItem)
	os.Remove(cacheItem.filepath) // 删除磁盘上的文件
	fileCache.curSize -= cacheItem.size
	delete(fileCache.itemMap, filename)
	fileCache.itemList.Remove(elem)
	fileCache.lock.Unlock()
}

func (fileCache *FileCache) IsExist(key string) bool {
	return fileCache.itemMap[key] != nil
}

func (fileCache *FileCache) CurrentSize() int {
	return fileCache.curSize
}

func (fileCache *FileCache) MaxSize() int {
	return fileCache.maxCacheSize
}

// 移除最旧的缓存文件
func (fileCache *FileCache) removeOldest() {
	oldest := fileCache.itemList.Back()
	if oldest == nil {
		return
	}
	cacheItem := oldest.Value.(*CacheItem)
	fileCache.curSize -= cacheItem.size
	os.Remove(cacheItem.filepath) // 删除磁盘上的缓存文件
	delete(fileCache.itemMap, cacheItem.filename)
	fileCache.itemList.Remove(oldest)
}
