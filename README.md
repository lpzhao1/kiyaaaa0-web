# kiyaaaa0-web

## Version 0

将Server抽象为Route与Start方法；

引入Context实现Read/Write，但需要用户自行创建Context；

## Version 1

实现由框架在路由时创建Context；

框架内置部分HTTP相应；

## Version 2

尝试支持RESTful API；

尝试实现一个基于Map的handler进行路由；

更改了部分方法的实现与作用位置；

## Version 3

运用责任链模式，可集成中间件；

handler based on Map改用线程安全的sync.Map；

更改了部分方法的实现与作用位置；

## Version 4

为Server抽象添加Shutdown方法，引入关闭机制；

添加ShutdownFilter，实时维护当前未完成请求数量；

运用Hook技术，实现关闭时尽可能保证请求都已完成；

## Version 5

尝试用路由树取代Map结构存储的路由表；

仅实现了最简单的路由树，未考虑路径参数与HTTP method；
