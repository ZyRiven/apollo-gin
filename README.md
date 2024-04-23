<div align="center">

# Apollo-GateWay 🚀

**一个基于串口通信的网关设备**

<p align="center">
 <img width="20%" src="logo.png" align="center" alt="banner" />
</p>

![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
</div>

## 📝 概述

本文档旨在介绍一个基于串口通信的网关设备，其主要功能是采集来自各种传感器的数据，并将其上传至云端数据中台。该网关设备作为物联网系统中的关键节点，负责连接硬件设备与云端服务，实现数据的采集、传输和存储。

## ✨ 设备功能

- 通过串口与lora设备进行无线通信
- 采集各种传感器（例如温度传感器、湿度传感器、光照传感器等）的数据。
- 将数据上传至云端数据中台

## 🖥️ 技术架构

- **硬件**:
  - 单板计算机（例如树莓派）
  - 串口（GPIO）
  - Lora
- **软件**:
  - 操作系统：Linux
  - 编程语言：Golang
  - 通讯协议：串口通信协议
  - 云端通讯协议（MQTT）

## ⚙️ 使用说明

``` sh
git clone https://github.com/ZyRiven/apollo-gin.git

cd apollo-gin

go mod tidy
# config中 MQTT开启需要配置服务端
go build && ./apollo
```

- **下载MQTT**
    - [官网](https://mqttx.app/zh)
    - 推荐使用Doctor下载（https://hub.docker.com）
    ``` sh
    docker run -d --name emqx -p 18083:18083 -p 1883:1883 emqx/emqx:latest
    ```


## ©️ 许可证

该项目遵循MIT许可证。更多信息请查看[`LICENSE`](LICENSE)文件。