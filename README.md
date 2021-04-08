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

![](https://oss.codery.cn/images/2020/07/12/20200712131203.png)

## 启动
cd /cmd

go run main.go //启动后台管理接口

go run client.go //启动RabbitMQ写入数据库客户端

go run spike.go //启动秒杀系统，支持横行扩展


## 测试

### 测试没有使用集群，只是一个服务器

### 这里我使用测试工具是jmeter

设置：

![](https://oss.codery.cn/images/2020/07/12/20200712090309.png)


测试结果：

![](https://oss.codery.cn/images/2020/07/12/20200712085943.png)


RabbitMQ：

![](https://oss.codery.cn/images/2020/07/12/20200712090446.png)

mysql:
并没有超卖，测试添加了1000个库存

![](https://oss.codery.cn/images/2020/07/12/20200712090604.png)

![](https://oss.codery.cn/images/2020/07/12/20200712090629.png)


