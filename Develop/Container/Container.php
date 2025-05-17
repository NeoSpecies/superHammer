<?php
namespace Develop\Container;

use Exception;
use ReflectionClass;

/**
 * 依赖注入容器类
 * 负责管理服务的注册、解析和实例化
 * 支持：
 * 1. 服务注册
 * 2. 服务解析
 * 3. 自动依赖注入
 * 4. 单例管理
 */
class Container {
    /**
     * 服务定义映射
     * 存储所有注册的服务定义
     * @var array
     */
    private $definitions = [];

    /**
     * 服务实例映射
     * 存储已实例化的服务
     * @var array
     */
    private $instances = [];

    /**
     * 注册服务到容器
     * 
     * @param string $name 服务名称
     * @param mixed $definition 服务定义（可以是闭包、类名或简单值）
     * @return void
     */
    public function set($name, $definition) {
        $this->definitions[$name] = $definition;
    }

    /**
     * 从容器获取服务实例
     * 
     * @param string $name 服务名称
     * @return mixed 服务实例
     * @throws Exception 当服务未找到或解析失败时抛出异常
     */
    public function get($name) {
        if (!isset($this->definitions[$name])) {
            throw new \Exception("Service $name not found.");
        }
    
        // 如果该服务已经被实例化，则直接返回实例
        if (isset($this->instances[$name])) {
            return $this->instances[$name];
        }
    
        // 检查服务定义是否为非复杂类型（如字符串、数组等），如果是，则直接返回
        $definition = $this->definitions[$name];
        if (!is_callable($definition) && !(is_string($definition) && class_exists($definition))) {
            // 对于简单类型，直接返回定义
            return $definition;
        }
    
        // 对于复杂类型，通过resolve方法处理
        $this->instances[$name] = $this->resolve($definition);
        return $this->instances[$name];
    }

    /**
     * 解析服务定义
     * 
     * @param mixed $definition 服务定义
     * @return mixed 解析后的服务实例
     * @throws Exception 当解析失败时抛出异常
     */
    private function resolve($definition) {
        if (is_callable($definition)) {
            try {
                return $definition($this);
            } catch (\Exception $e) {
                throw new \Exception("Error resolving callable service with message: " . $e->getMessage());
            }
        } elseif (is_string($definition) && class_exists($definition)) {
            try {
                return $this->instantiateClass($definition);
            } catch (\Exception $e) {
                throw new \Exception("Error instantiating class for service with message: " . $e->getMessage());
            }
        }
    
        throw new \Exception("Unable to resolve service definition.");
    }

    /**
     * 实例化类
     * 
     * @param string $className 类名
     * @return object 类的实例
     * @throws Exception 当实例化失败时抛出异常
     */
    private function instantiateClass($className) {
        $reflectionClass = new ReflectionClass($className);
        $constructor = $reflectionClass->getConstructor();

        if ($constructor === null) {
            return new $className();
        }

        $parameters = $constructor->getParameters();
        $dependencies = $this->resolveDependencies($parameters);

        return $reflectionClass->newInstanceArgs($dependencies);
    }

    /**
     * 解析构造函数依赖
     * 
     * @param array $parameters 构造函数参数列表
     * @return array 解析后的依赖列表
     * @throws Exception 当依赖无法解析时抛出异常
     */
    private function resolveDependencies(array $parameters) {
        $dependencies = [];
    
        foreach ($parameters as $parameter) {
            // PHP 8.0以上版本需要检查参数类型
            $type = $parameter->getType();
            if ($type && !$type->isBuiltin()) {
                $typeName = $type instanceof \ReflectionNamedType ? $type->getName() : (string)$type;
                $dependencies[] = $this->get($typeName);
            } else {
                // 处理非类依赖，如常量值或变量
                if ($parameter->isDefaultValueAvailable()) {
                    $dependencies[] = $parameter->getDefaultValue();
                } else {
                    throw new \Exception("Unresolvable dependency encountered.");
                }
            }
        }
    
        return $dependencies;
    }
}
