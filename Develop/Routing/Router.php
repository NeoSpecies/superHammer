<?php

namespace Develop\Routing;

use ReflectionMethod;
use ReflectionParameter;

class Router
{
    private $routes;

    public function __construct($container)
    {
        $filePath = $container->get('routesFile');
        $this->routes = $this->loadRoutes($filePath);
    }

    private function loadRoutes($filePath)
    {
        $json = file_get_contents($filePath);
        return json_decode($json, true)['routes'];
    }

    public function match($path)
    {
        foreach ($this->routes as $route) {
            if ($route['path'] === $path && $route['language'] === 'php') {
                return $route;
            }
        }
        return null;
    }

    public function dispatch($route, $params)
    {
        list($classPath, $methodName) = explode('::', $route['command']);
        $className = "\\Develop\\Controller\\" . str_replace('/', '\\', $classPath);
        if (class_exists($className) && method_exists($className, $methodName)) {
            $controller = new $className();

            // 使用反射获取方法所需参数
            $method = new ReflectionMethod($className, $methodName);
            $methodParams = $method->getParameters();
            $args = [];
            // var_dump($methodParams);
            foreach ($methodParams as $param) {
                // 获取参数名称
                $paramName = $param->getName();

                // 检查params数组中是否有该参数，如果有则使用，如果没有则检查是否有默认值
                if (array_key_exists($paramName, $params)) {
                    $args[] = $params[$paramName];
                } elseif ($param->isDefaultValueAvailable()) {
                    $args[] = $param->getDefaultValue();
                } else {
                    // 如果没有默认值，可以选择抛出异常或传递null
                    $args[] = null;
                }
            }

            // 动态调用方法并传递参数
            return call_user_func_array([$controller, $methodName], $args);
        }
        return "Command not found.";
    }
}
