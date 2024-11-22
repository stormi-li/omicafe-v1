package omicafe

import (
	"sync"
)

type FileCache struct {
	lock           sync.RWMutex // 改为读写锁
	MaxSize        int
	FileMgr        *FileManager
	LRUManager     *LRUManager
	CacheHitCount  int // 缓存命中次数
	CacheMissCount int // 缓存未命中次数
	CacheClearNum  int // 缓存清除次数
}

// 初始化文件缓存
func newFileCache(cacheDir string, maxSize int) *FileCache {
	return &FileCache{
		MaxSize:    maxSize,
		FileMgr:    NewFileManager(cacheDir),
		LRUManager: NewLRUManager(),
	}
}

// 设置缓存文件
func (fc *FileCache) Set(key string, data []byte) {
	size := len(data)
	if size > fc.MaxSize {
		return // 文件太大，直接丢弃
	}

	fc.lock.Lock()
	defer fc.lock.Unlock()

	// 确保容量充足
	for fc.LRUManager.Size+size > fc.MaxSize {
		oldest := fc.LRUManager.RemoveOldest()
		if oldest != nil {
			_ = fc.FileMgr.DeleteFile(oldest.Key)
			fc.CacheClearNum++ // 清除计数
		}
	}

	// 写入文件并更新 LRU
	if err := fc.FileMgr.WriteFile(key, data); err == nil {
		fc.LRUManager.Add(key, size)
	}
}

// 获取缓存文件
func (fc *FileCache) Get(key string) ([]byte, bool) {
	fc.lock.RLock() // 使用读锁
	_, found := fc.LRUManager.Get(key)
	if !found {
		fc.lock.RUnlock() // 在获取写锁前释放读锁
		fc.lock.Lock()    // 获取写锁清理逻辑
		defer fc.lock.Unlock()

		fc.CacheMissCount++ // 未命中计数
		return nil, false
	}

	data, err := fc.FileMgr.ReadFile(key)
	fc.lock.RUnlock() // 提前释放读锁
	if err != nil {
		fc.lock.Lock()         // 升级为写锁，清除不存在的缓存项
		defer fc.lock.Unlock() // 确保锁的正确释放
		fc.LRUManager.Remove(key)
		fc.CacheMissCount++ // 未命中计数
		return nil, false
	}

	fc.lock.Lock()     // 写锁保护计数器更新
	fc.CacheHitCount++ // 命中计数
	fc.lock.Unlock()
	return data, true
}

// 删除缓存文件
func (fc *FileCache) Del(key string) {
	fc.lock.Lock()
	defer fc.lock.Unlock()

	fc.LRUManager.Remove(key)
	_ = fc.FileMgr.DeleteFile(key)
}

// 当前缓存使用大小
func (fc *FileCache) CurrentSize() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.LRUManager.Size
}

// 获取命中次数
func (fc *FileCache) GetCacheHitCount() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.CacheHitCount
}

// 获取未命中次数
func (fc *FileCache) GetCacheMissCount() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.CacheMissCount
}

// 获取清除次数
func (fc *FileCache) GetCacheClearCount() int {
	fc.lock.RLock()
	defer fc.lock.RUnlock()
	return fc.CacheClearNum
}
