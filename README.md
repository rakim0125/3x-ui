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

**3X-UI** — advanced, open-source web-based control panel designed for managing Xray-core server. It offers a user-friendly interface for configuring and monitoring various VPN and proxy protocols.

> [!IMPORTANT]
> This project is only for personal usage, please do not use it for illegal purposes, and please do not use it in a production environment.

As an enhanced fork of the original X-UI project, 3X-UI provides improved stability, broader protocol support, and additional features.

## Features

- **OpenAPI Support** - Full REST API for programmatic management with API Key authentication
- **Registry Node** - Automatic node registration and heartbeat for distributed deployments
- **API Key Management** - Create, list, and delete API keys from the web panel
- **Default inbound on first install** - If the database has no inbounds yet, a **VLESS** inbound is created on port **443** with **REALITY** (TCP, `xtls-rprx-vision`, dest/SNI `www.microsoft.com`, uTLS `chrome`). X25519 keys and a short ID are generated at init; adjust or remove it in the panel as needed.

See [apidoc.md](/apidoc.md) for complete API documentation.

## Configuration

### Registry nodes

Registry root URLs can be set in this order:

1. Environment variables `XUI_REGISTRY_NODES` or `REGISTRY_NODES` (comma-separated URLs).
2. `XUI_REGISTRY_NODES_FILE` pointing to a file (one URL per line, `#` comments allowed).
3. File `{XUI_DB_FOLDER}/registry_nodes` under your data directory.
4. Built-in defaults from the binary if nothing else is present.

See [`.env.example`](/.env.example) and [`config/registry_nodes.example`](/config/registry_nodes.example).

## Linux binary package (from source)

To build static Linux binaries and tarballs (amd64 + arm64):

```bash
bash packaging/build-linux-direct.sh
```

Artifacts: `dist/x-ui-linux-<arch>-direct-<git-describe>.tar.gz`.

- On **macOS**, [Zig](https://ziglang.org/) is required for musl cross-compilation (e.g. `brew install zig`).
- On **Linux**, set `ONLY_NATIVE=1` to build only the current architecture.

After extracting a tarball, run `sudo bash install-binary-linux.sh` (the script may download Xray and geo assets during install).

## Quick Start

```bash
bash <(curl -Ls https://raw.githubusercontent.com/rakim0125/3x-ui/main/install.sh)
```

For full documentation, please visit the [project Wiki](https://github.com/MHSanaei/3x-ui/wiki).

## A Special Thanks to

- [alireza0](https://github.com/alireza0/)

## Acknowledgment

- [Iran v2ray rules](https://github.com/chocolate4u/Iran-v2ray-rules) (License: **GPL-3.0**): _Enhanced v2ray/xray and v2ray/xray-clients routing rules with built-in Iranian domains and a focus on security and adblocking._
- [Russia v2ray rules](https://github.com/runetfreedom/russia-v2ray-rules-dat) (License: **GPL-3.0**): _This repository contains automatically updated V2Ray routing rules based on data on blocked domains and addresses in Russia._

## Support project

**If this project is helpful to you, you may wish to give it a**:star2:

<a href="https://www.buymeacoffee.com/MHSanaei" target="_blank">
<img src="./media/default-yellow.png" alt="Buy Me A Coffee" style="height: 70px !important;width: 277px !important;" >
</a>

</br>
<a href="https://nowpayments.io/donation/hsanaei" target="_blank" rel="noreferrer noopener">
   <img src="./media/donation-button-black.svg" alt="Crypto donation button by NOWPayments">
</a>

## Stargazers over Time

[![Stargazers over time](https://starchart.cc/MHSanaei/3x-ui.svg?variant=adaptive)](https://starchart.cc/MHSanaei/3x-ui)
