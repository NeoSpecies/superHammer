package agilitymemdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Item 数据项结构
// 用于存储键值对中的值
type Item struct {
	// Value 存储的值
	Value string
}

// Transaction 事务结构
// 用于存储事务中的操作
type Transaction struct {
	// Operations 事务操作映射
	// key: 操作的键
	// value: 操作的数据项
	Operations map[string]*Item
}

// AgilityMemDB 内存数据库结构
// 提供内存数据存储和事务支持
type AgilityMemDB struct {
	// Data 数据存储映射
	// key: 数据的键
	// value: 数据项
	Data        map[string]*Item
	// lock 读写锁
	// 用于保护并发访问
	lock        sync.RWMutex
	// transaction 当前事务
	// 存储事务中的操作
	transaction *Transaction
	// FilePath 数据文件路径
	// 用于持久化存储
	FilePath    string
}

// NewAgilityMemDB 创建新的内存数据库实例
// 功能：
// 1. 初始化数据存储映射
// 2. 设置文件路径
// 参数：
//   - filePath string: 数据文件路径
// 返回值：
//   - *AgilityMemDB: 新创建的内存数据库实例
func NewAgilityMemDB(filePath string) *AgilityMemDB {
	return &AgilityMemDB{
		Data:     make(map[string]*Item),
		FilePath: filePath,
	}
}

// LoadData 从文件加载数据
// 功能：
// 1. 读取数据文件
// 2. 解析JSON数据
// 3. 更新内存数据
// 参数：无
// 返回值：
//   - error: 加载过程中的错误信息，如果成功则返回nil
func (db *AgilityMemDB) LoadData() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	bytes, err := os.ReadFile(db.FilePath) // Updated from ioutil.ReadFile
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		return nil // File is empty
	}

	return json.Unmarshal(bytes, &db.Data)
}

// Get 获取指定键的值
// 功能：
// 1. 从内存数据库中查找键
// 2. 返回对应的值
// 参数：
//   - key string: 要查找的键
// 返回值：
//   - string: 找到的值
//   - bool: 是否找到
func (db *AgilityMemDB) Get(key string) (string, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	// 尝试从数据库中获取键对应的值
	item, ok := db.Data[key]
	if !ok {
		fmt.Printf("Key not found: %s\n", key)
		return "", false // 键不存在时返回空字符串和false
	}

	// // 如果找到了键，则打印其对应的值
	return item.Value, true // 返回找到的值和true
}

// Put 存储键值对
// 功能：
// 1. 如果有活动事务，将操作添加到事务中
// 2. 否则直接更新内存数据
// 参数：
//   - key string: 要存储的键
//   - value string: 要存储的值
// 返回值：
//   - error: 存储过程中的错误信息，如果成功则返回nil
func (db *AgilityMemDB) Put(key string, value string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.transaction != nil {
		db.transaction.Operations[key] = &Item{Value: value}
	} else {
		db.Data[key] = &Item{Value: value}
	}
	return nil
}

// Delete 删除指定键
// 功能：
// 1. 如果有活动事务，从事务中删除操作
// 2. 否则直接从内存数据中删除
// 参数：
//   - key string: 要删除的键
// 返回值：无
func (db *AgilityMemDB) Delete(key string) {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.transaction != nil {
		delete(db.transaction.Operations, key)
	} else {
		delete(db.Data, key)
	}
}

// Persist 将数据持久化到文件
// 功能：
// 1. 将内存数据转换为JSON
// 2. 写入文件
// 参数：无
// 返回值：
//   - error: 持久化过程中的错误信息，如果成功则返回nil
func (db *AgilityMemDB) Persist() error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	bytes, err := json.Marshal(db.Data)
	if err != nil {
		return err
	}

	return os.WriteFile(db.FilePath, bytes, 0644) // Updated from ioutil.WriteFile
}

// BeginTransaction 开始新的事务
// 功能：
// 1. 检查是否有活动事务
// 2. 创建新的事务
// 参数：无
// 返回值：
//   - error: 开始事务过程中的错误信息，如果成功则返回nil
func (db *AgilityMemDB) BeginTransaction() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.transaction != nil {
		return errors.New("transaction already in progress")
	}

	db.transaction = &Transaction{
		Operations: make(map[string]*Item),
	}
	return nil
}

// CommitTransaction 提交当前事务
// 功能：
// 1. 检查是否有活动事务
// 2. 将事务中的操作应用到内存数据
// 3. 清除事务
// 参数：无
// 返回值：
//   - error: 提交事务过程中的错误信息，如果成功则返回nil
func (db *AgilityMemDB) CommitTransaction() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.transaction == nil {
		return errors.New("no active transaction")
	}

	for key, item := range db.transaction.Operations {
		db.Data[key] = item
	}

	db.transaction = nil
	return nil
}

// RollbackTransaction 回滚当前事务
// 功能：
// 1. 清除当前事务
// 2. 放弃所有未提交的更改
// 参数：无
// 返回值：无
func (db *AgilityMemDB) RollbackTransaction() {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.transaction = nil
}
