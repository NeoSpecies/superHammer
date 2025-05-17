package di

import (
	"fmt"
	"reflect"
	"sync"
)

// Lifecycle 定义服务的生命周期类型
// 用于控制服务实例的创建和复用策略
type Lifecycle int

const (
	// Singleton 单例模式
	// 服务只会被创建一次，后续请求都返回同一个实例
	Singleton Lifecycle = iota
	// Prototype 原型模式
	// 每次请求都会创建新的服务实例
	Prototype
)

// Container 依赖注入容器
// 负责管理服务的注册、解析和生命周期
type Container struct {
	// services 存储所有注册的服务
	// key: 服务名称
	// value: 服务实例信息
	services map[string]serviceInstance
	// mutex 用于保护services的并发访问
	mutex    sync.RWMutex
}

// serviceInstance 服务实例信息
// 存储服务的构造函数、实例和生命周期信息
type serviceInstance struct {
	// Constructor 服务构造函数
	// 用于创建服务实例的函数
	Constructor interface{}
	// Instance 服务实例
	// 存储已创建的服务实例（仅用于Singleton模式）
	Instance    interface{}
	// Lifecycle 服务生命周期
	// 决定服务实例的创建和复用策略
	Lifecycle   Lifecycle
}

// NewContainer 创建新的依赖注入容器
// 功能：
// 1. 初始化服务映射表
// 2. 初始化互斥锁
// 参数：无
// 返回值：
//   - *Container: 新创建的容器实例
func NewContainer() *Container {
	return &Container{
		services: make(map[string]serviceInstance),
	}
}

// Register 注册服务到容器
// 功能：
// 1. 验证构造函数类型
// 2. 将服务信息存储到容器中
// 参数：
//   - serviceName string: 服务名称
//   - constructor interface{}: 服务构造函数
//   - lifecycle Lifecycle: 服务生命周期
// 返回值：
//   - error: 注册过程中的错误信息，如果成功则返回nil
func (c *Container) Register(serviceName string, constructor interface{}, lifecycle Lifecycle) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if reflect.TypeOf(constructor).Kind() != reflect.Func {
		return fmt.Errorf("constructor must be a function")
	}

	c.services[serviceName] = serviceInstance{
		Constructor: constructor,
		Lifecycle:   lifecycle,
	}
	return nil
}

// Resolve 从容器中解析服务实例
// 功能：
// 1. 查找服务信息
// 2. 根据生命周期策略创建或返回实例
// 3. 处理并发访问
// 参数：
//   - serviceName string: 服务名称
// 返回值：
//   - interface{}: 服务实例
//   - error: 解析过程中的错误信息，如果成功则返回nil
func (c *Container) Resolve(serviceName string) (interface{}, error) {
	c.mutex.RLock()
	instance, ok := c.services[serviceName]
	c.mutex.RUnlock() // 立即释放读锁，避免死锁

	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	if instance.Lifecycle == Singleton {
		if instance.Instance != nil {
			return instance.Instance, nil
		}

		// 获取写锁前的双重检查
		c.mutex.Lock()
		// 再次检查以确保在等待锁的过程中没有其他goroutine完成了实例化
		if instance.Instance == nil {
			constructorValue := reflect.ValueOf(instance.Constructor)
			resultValues := constructorValue.Call(nil) // 假定这个构造函数不接受参数

			if len(resultValues) != 1 {
				c.mutex.Unlock()
				return nil, fmt.Errorf("constructor for %s did not return one value", serviceName)
			}

			instance.Instance = resultValues[0].Interface()
			c.services[serviceName] = instance
		}
		result := instance.Instance
		c.mutex.Unlock()
		return result, nil
	}

	// 对于Prototype生命周期，每次都创建新实例
	constructorValue := reflect.ValueOf(instance.Constructor)
	resultValues := constructorValue.Call(nil) // 假定构造函数不接受参数且返回一个值

	if len(resultValues) != 1 {
		return nil, fmt.Errorf("constructor for %s did not return one value", serviceName)
	}

	return resultValues[0].Interface(), nil
}

// RegisterAuto 自动注册服务到容器
// 功能：
// 1. 获取服务类型信息
// 2. 使用类型名称作为服务名称
// 3. 注册服务到容器
// 参数：
//   - impl interface{}: 服务实现
//   - lifecycle Lifecycle: 服务生命周期
// 返回值：
//   - error: 注册过程中的错误信息，如果成功则返回nil
func (c *Container) RegisterAuto(impl interface{}, lifecycle Lifecycle) error {
	serviceType := reflect.TypeOf(impl).Elem()
	serviceName := serviceType.Name()

	return c.Register(serviceName, func() interface{} { return impl }, lifecycle)
}

// ResolveType 根据类型解析服务实例
// 功能：
// 1. 获取类型名称
// 2. 使用类型名称解析服务实例
// 参数：
//   - serviceType reflect.Type: 服务类型
// 返回值：
//   - interface{}: 服务实例
//   - error: 解析过程中的错误信息，如果成功则返回nil
func (c *Container) ResolveType(serviceType reflect.Type) (interface{}, error) {
	serviceName := serviceType.Name()
	return c.Resolve(serviceName)
}
