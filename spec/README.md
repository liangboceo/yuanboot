# Yuanboot 框架规格文档

本目录包含 Yuanboot 框架的完整规格说明文档。

## 文档列表

| 文档 | 描述 |
|------|------|
| [yuanboot.md](./yuanboot.md) | 框架完整使用说明 |

## 框架概述

Yuanboot 是一个简单、轻量、快速、基于依赖注入的 Go 微服务框架。

### 核心特性

- **高性能**: 基于 fasthttp 和 net.http 的双服务器实现
- **MVC架构**: 完整的 MVC 模式支持
- **依赖注入**: 灵活的 DI 容器集成
- **中间件**: 丰富的中间件支持 (CORS, JWT, Logger, Recovery等)
- **服务发现**: 支持 Nacos, Eureka, Consul, ETCD
- **多协议**: REST API 和 gRPC 双协议支持

### 快速安装

```bash
go get github.com/liangboceo/yuanboot
```

### 版本信息

- 当前版本: v1.9.2747
- 官方网站: https://yuanboot.star2cloud.com
- GitHub: https://github.com/liangboceo/yuanboot
