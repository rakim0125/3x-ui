[English](/README.md) | [فارسی](/README.fa_IR.md) | [العربية](/README.ar_EG.md) |  [中文](/README.zh_CN.md) | [Español](/README.es_ES.md) | [Русский](/README.ru_RU.md)

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./media/3x-ui-dark.png">
    <img alt="3x-ui" src="./media/3x-ui-light.png">
  </picture>
</p>

[![Release](https://img.shields.io/github/v/release/mhsanaei/3x-ui.svg)](https://github.com/MHSanaei/3x-ui/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/mhsanaei/3x-ui/release.yml.svg)](https://github.com/MHSanaei/3x-ui/actions)
[![GO Version](https://img.shields.io/github/go-mod/go-version/mhsanaei/3x-ui.svg)](#)
[![Downloads](https://img.shields.io/github/downloads/mhsanaei/3x-ui/total.svg)](https://github.com/MHSanaei/3x-ui/releases/latest)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/mhsanaei/3x-ui/v2.svg)](https://pkg.go.dev/github.com/mhsanaei/3x-ui/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhsanaei/3x-ui/v2)](https://goreportcard.com/report/github.com/mhsanaei/3x-ui/v2)

**3X-UI** — 一个基于网页的高级开源控制面板，专为管理 Xray-core 服务器而设计。它提供了用户友好的界面，用于配置和监控各种 VPN 和代理协议。

> [!IMPORTANT]
> 本项目仅用于个人使用和通信，请勿将其用于非法目的，请勿在生产环境中使用。

作为原始 X-UI 项目的增强版本，3X-UI 提供了更好的稳定性、更广泛的协议支持和额外的功能。

## 功能亮点

- **OpenAPI** — 提供 REST API，使用 API Key 鉴权进行程序化运维
- **注册中心节点** — 支持节点自动注册与心跳（分布式部署）
- **API Key 管理** — 在 Web 面板中创建、列出、删除 API Key
- **首次初始化默认入站** — 若数据库中尚无任何入站，系统会自动创建一条 **VLESS** 入站：端口 **443**、传输 **TCP**、安全 **REALITY**、流控 **xtls-rprx-vision**，目标站点 **www.microsoft.com:443**（SNI 一致）、uTLS 指纹 **chrome**；REALITY 公私钥与 Short ID 在初始化时自动生成，可在面板中修改或删除

完整接口说明见 [apidoc.md](/apidoc.md)。

## 配置说明

### 注册中心地址

注册根 URL 按以下优先级读取：

1. 环境变量 `XUI_REGISTRY_NODES` 或 `REGISTRY_NODES`（逗号分隔多个地址）
2. 环境变量 `XUI_REGISTRY_NODES_FILE` 指定的文件（每行一个 URL，以 `#` 开头的行为注释）
3. 数据目录下 `{XUI_DB_FOLDER}/registry_nodes` 文件
4. 若以上均不存在，使用程序内置默认列表

详见 [`.env.example`](/.env.example) 与 [`config/registry_nodes.example`](/config/registry_nodes.example)。

## Linux 直装包（从源码构建）

生成 amd64 / arm64 的 musl 静态二进制及 `tar.gz` 安装包：

```bash
bash packaging/build-linux-direct.sh
```

产物位于 `dist/`，文件名形如 `x-ui-linux-<架构>-direct-<git 描述>.tar.gz`。

- 在 **macOS** 上交叉编译需要已安装 [Zig](https://ziglang.org/)（例如 `brew install zig`）。
- 在 **Linux** 上若只打本机架构，可设置 `ONLY_NATIVE=1`。

解压后执行 `sudo bash install-binary-linux.sh` 安装（安装过程可能会联网下载 Xray 与 geo 数据）。

## 快速开始

```
bash <(curl -Ls https://raw.githubusercontent.com/mhsanaei/3x-ui/master/install.sh)
```

完整文档请参阅 [项目Wiki](https://github.com/MHSanaei/3x-ui/wiki)。

## 特别感谢

- [alireza0](https://github.com/alireza0/)

## 致谢

- [Iran v2ray rules](https://github.com/chocolate4u/Iran-v2ray-rules) (许可证: **GPL-3.0**): _增强的 v2ray/xray 和 v2ray/xray-clients 路由规则，内置伊朗域名，专注于安全性和广告拦截。_
- [Russia v2ray rules](https://github.com/runetfreedom/russia-v2ray-rules-dat) (许可证: **GPL-3.0**): _此仓库包含基于俄罗斯被阻止域名和地址数据自动更新的 V2Ray 路由规则。_

## 支持项目

**如果这个项目对您有帮助，您可以给它一个**:star2:

<a href="https://www.buymeacoffee.com/MHSanaei" target="_blank">
<img src="./media/default-yellow.png" alt="Buy Me A Coffee" style="height: 70px !important;width: 277px !important;" >
</a>

</br>
<a href="https://nowpayments.io/donation/hsanaei" target="_blank" rel="noreferrer noopener">
   <img src="./media/donation-button-black.svg" alt="Crypto donation button by NOWPayments">
</a>

## 随时间变化的星标数

[![Stargazers over time](https://starchart.cc/MHSanaei/3x-ui.svg?variant=adaptive)](https://starchart.cc/MHSanaei/3x-ui) 
