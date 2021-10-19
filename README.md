# 使用go-redis实现一个令牌管理器

## 需求描述
假设我们当前的所有服务需要一个第三方的认证，认证形式为：在发送请求的时候带上第三方颁发的令牌，该令牌具有一个时效性
第三方的令牌可以通过某个接口获取，但是该接口做了单位时间内的同一ip的请求频率的限制，因此在并发的场景下，我们需要控制令牌获取接口的频率

## 项目结构
```
├── go.mod
├── go.sum
├── main.go
├── README.go
└── src
    ├── config.go
    ├── connector.go
    ├── token-manager.go
    └── token-server.go
```
* src/config.go
    管理redis配置
* src/connector.go
	连接并返回redis客户端实例的方法
* src/token-server.go
	模拟向一个第三方的服务器请求令牌
* src/token-manager.go
    令牌管理的对象，负责获取或者更新令牌

## 实现思路
* 储存在redis中的令牌的key假定为token
* 我们设定一个更新者，在redis中的key为updater
* 当多个请求同时向redis索取令牌，但是令牌过期了或者没有令牌时，尝试去提交一个updater，然后获取redis中updater的值
* 若redis中updater与自身提交的updater保持一致，说明updater设置成功，获得了更新token的权利，此时向第三方的token服务发起请求，获取一个新的令牌，并设置到redis中
* 若redis中的updater与自身提交的updater不一致，则直接进入持续获取新令牌的阶段，若获取的令牌为空，则继续尝试获取

## 效果预览
测试的时候我使用了并发100次，每产生25次并发时中间休息一秒，并且token在redis中的有效时间为5秒
```
routine c663db9b-7666-4757-885b-bf3b7d19a524 is updating token
NewToken:d153f070-a27b-422e-97ad-8c87222ac1cf  ||  Generated at 2021-10-19 21:19:33.382049711 +0800 CST m=+9.007199601, waiting for 8 second
d153f070-a27b-422e-97ad-8c87222ac1cf
d153f070-a27b-422e-97ad-8c87222ac1cf
.
.
.
d153f070-a27b-422e-97ad-8c87222ac1cf
routine 1b80d9bc-6e6f-4a55-86fa-734df39b605f is updating token
NewToken:d1b30ccf-10d8-4ed6-9565-67b6460a84db  ||  Generated at 2021-10-19 21:19:42.388564039 +0800 CST m=+18.013714046, waiting for 4 second
d1b30ccf-10d8-4ed6-9565-67b6460a84db
d1b30ccf-10d8-4ed6-9565-67b6460a84db
.
.
.
d1b30ccf-10d8-4ed6-9565-67b6460a84db
routine 4179fff5-254f-4325-82af-55f0a33336f2 is updating token
d1b30ccf-10d8-4ed6-9565-67b6460a84db
d1b30ccf-10d8-4ed6-9565-67b6460a84db
.
.
.
d1b30ccf-10d8-4ed6-9565-67b6460a84db
NewToken:a6c8baf0-25fb-40cc-84e6-39af9fd02e29  ||  Generated at 2021-10-19 21:19:50.397202332 +0800 CST m=+26.022352261, waiting for 8 second
a6c8baf0-25fb-40cc-84e6-39af9fd02e29
.
.
.
a6c8baf0-25fb-40cc-84e6-39af9fd02e29
a6c8baf0-25fb-40cc-84e6-39af9fd02e29
```
可以看到，在redis没有token的情况下，也是仅有一个线程可以获取与更新token，其它线程仅能等待token获取并设置到redis之后才可以去获取