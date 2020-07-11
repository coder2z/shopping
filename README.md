# go开发电商网站高并发秒杀系统

这个系统的主要目的在于秒杀，所有其他地方都做的很简单。功能不多！

## 技术栈：
> web框架：gin

> 消息队列：RabbitMQ

> 分布式方案：hash环

> orm: gorm

> 限流器：tollbooth

> 登录验证：jwt

## 架构

![](https://gitee.com/myxy99/pic/raw/master/img/blog/2020/07/11/20200711230052.png)

## 启动
cd /cmd

go run main.go //启动后台管理接口

go run client.go //启动RabbitMQ写入数据库客户端

go run spike.go //启动秒杀系统，支持横行扩展
