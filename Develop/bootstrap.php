<?php
/**
 * 应用程序引导文件
 * 负责初始化应用程序的核心组件和服务
 * 包括：
 * 1. 加载自动加载器
 * 2. 初始化依赖注入容器
 * 3. 配置服务
 * 4. 启动服务器
 */

require __DIR__ . '/vendor/autoload.php';

// 引入必要的类
use Develop\Server\AsyncUnixSocketServer;    // 异步Unix Socket服务器
use Develop\Lifecycle\SocketLifecycleHandler; // Socket生命周期处理器
use Develop\Reader\UnixSocketReader;         // Unix Socket读取器
use Develop\Routing\Router;                  // 路由处理器
use Develop\Tool\Path;                       // 路径工具类
use Develop\Tool\Input;                      // 输入处理工具类
use Develop\Container\Container;             // 依赖注入容器

// 获取必要的文件路径
$socketFile = Path::getSocketFile();         // 获取Socket文件路径
$socketMainFile = Path::getSocketMainFile(); // 获取主Socket文件路径
$pidFile = Path::getBussinessId();           // 获取业务进程PID文件路径
$routesFile = Path::getRouter();             // 获取路由配置文件路径

// 创建容器实例
$container = new Container();

// 设置服务定义
// 注册基础配置
$container->set('socketFile', $socketFile);           // 注册Socket文件路径
$container->set('pidFile', $pidFile);                 // 注册PID文件路径
$container->set('socketMainFile', $socketMainFile);   // 注册主Socket文件路径
$container->set('routesFile', $routesFile);           // 注册路由配置文件路径

// 注册Unix Socket读取器服务
$container->set('unixSocketReader', function (Container $container) {
    return new UnixSocketReader($container);
});

// 注册输入处理服务
$container->set('input', function (Container $container) {
    return new Input($container);
});
$input = $container->get('input');
// $input->initialize();  // 输入服务初始化（当前被注释）

// 注册路由服务
$container->set('router', function (Container $container) {
    return new Router($container);
});

// 注册生命周期处理器服务
$container->set('lifecycleHandler', function (Container $container) {
    return new SocketLifecycleHandler($container);
});

// 注册异步Unix Socket服务器服务
$container->set('asyncUnixSocketServer', function ($container) {
    $socketFile = $container->get('socketFile');
    $pidFile = $container->get('pidFile');
    $lifecycleHandler = $container->get('lifecycleHandler');
    $server = new AsyncUnixSocketServer($socketFile, $pidFile, $container);
    $server->registerLifecycleHandler($lifecycleHandler);
    return $server;
});

// 获取服务实例并启动服务器
$server = $container->get('asyncUnixSocketServer');
$server->start();
