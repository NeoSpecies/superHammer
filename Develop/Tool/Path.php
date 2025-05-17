<?php
namespace Develop\Tool;

/**
 * 路径工具类
 * 提供各种系统路径的获取方法
 * 用于统一管理应用程序中的文件路径
 */
class Path
{
    /**
     * 获取应用程序根目录路径
     * 
     * @return string 应用程序根目录的绝对路径
     */
    public static function getRoot()
    {
        // 获取当前脚本所在目录的父级的父级
        return dirname(dirname(__DIR__));
    }

    /**
     * 获取路由配置文件路径
     * 
     * @return string 路由配置文件的绝对路径
     */
    public static function getRouter()
    {
        $rootPath = Path::getRoot();
        $routes = $rootPath . '/config/router.json';
        return $routes;
    }

    /**
     * 获取应用程序配置文件路径
     * 
     * @return string 配置文件的绝对路径
     */
    public static function getConfig()
    {
        $rootPath = Path::getRoot();
        $config = $rootPath . '/config/config.json';
        return $config;
    }

    /**
     * 获取PHP Socket文件路径
     * 
     * @return string PHP Socket文件的绝对路径
     */
    public static function getSocketFile()
    {
        $rootPath = Path::getRoot();
        $socketFile = $rootPath . '/runtime/phpSocket.sock';
        return $socketFile;
    }

    /**
     * 获取主Socket文件路径
     * 
     * @return string 主Socket文件的绝对路径
     */
    public static function getSocketMainFile()
    {
        $rootPath = Path::getRoot();
        $mainSocket = $rootPath . '/runtime/mainSocket.sock';
        return $mainSocket;
    }

    /**
     * 获取业务进程PID文件路径
     * 
     * @return string PID文件的绝对路径
     */
    public static function getBussinessId()
    {
        $rootPath = Path::getRoot();
        $duangBussiness = $rootPath . '/runtime/duangBussiness.pid';
        return $duangBussiness;
    }
}