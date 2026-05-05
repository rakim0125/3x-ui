# API 协议文档

本文档描述当前项目所有已开放的对外 API 协议，包括：

- 面板登录后使用的 API Key 管理接口
- 通过 API Key 访问的 `OpenAPI`

## 1. 基础约定

### 1.1 基础路径

- 面板基础路径由 `webBasePath` 决定，默认是 `/`
- 文档以下路径按默认根路径编写
- 如果实际基础路径为 `/xui/`，则：
  - `/panel/api/apikeys/list` 实际为 `/xui/panel/api/apikeys/list`
  - `/openapi/inbounds/list` 实际为 `/xui/openapi/inbounds/list`

### 1.2 通用响应结构

除 API Key 鉴权失败返回 `401` 外，大部分业务接口返回结构统一为：

```json
{
  "success": true,
  "msg": "",
  "obj": {}
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `success` | `bool` | 是否成功 |
| `msg` | `string` | 提示信息；失败时通常包含错误描述 |
| `obj` | `any` | 实际业务数据 |

### 1.3 API Key 鉴权失败响应

`/openapi/*` 未提供或提供了无效 API Key 时返回：

```json
{
  "success": false,
  "msg": "missing api key: provide X-API-Key header or api_key query parameter"
}
```

或：

```json
{
  "success": false,
  "msg": "invalid api key"
}
```

### 1.4 鉴权方式

#### 面板登录态

以下接口要求浏览器或客户端已持有面板登录后的 Session Cookie：

- `POST /panel/api/apikeys/create`
- `GET /panel/api/apikeys/list`
- `POST /panel/api/apikeys/delete/:id`

#### API Key

`/openapi/*` 所有接口统一要求 API Key，支持两种传递方式：

请求头方式：

```http
X-API-Key: sk-xxxxxxxxxxxxxxxx
```

Query 参数方式：

```text
?api_key=sk-xxxxxxxxxxxxxxxx
```

## 2. 数据结构

### 2.1 API Key 对象

```json
{
  "id": 1,
  "name": "automation",
  "prefix": "sk-AbCdEf12",
  "userId": 1,
  "createdAt": 1710000000,
  "expiresAt": 0
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | API Key 主键 |
| `name` | `string` | API Key 名称 |
| `prefix` | `string` | 明文前缀，用于展示 |
| `userId` | `int` | 归属用户 ID |
| `createdAt` | `int64` | 创建时间戳 |
| `expiresAt` | `int64` | 过期时间戳，`0` 为永不过期 |

### 2.2 Inbound 对象

开放接口使用的入站对象字段与服务端模型一致，常见字段如下：

```json
{
  "id": 1,
  "up": 0,
  "down": 0,
  "total": 0,
  "allTime": 0,
  "remark": "test",
  "enable": true,
  "expiryTime": 0,
  "trafficReset": "never",
  "lastTrafficResetTime": 0,
  "clientStats": [],
  "listen": "0.0.0.0",
  "port": 443,
  "protocol": "vless",
  "settings": "{...}",
  "streamSettings": "{...}",
  "tag": "inbound-443",
  "sniffing": "{...}"
}
```

主要字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | 入站 ID |
| `remark` | `string` | 备注 |
| `enable` | `bool` | 是否启用 |
| `port` | `int` | 入站端口 |
| `protocol` | `string` | 协议，例如 `vmess`、`vless`、`trojan` |
| `listen` | `string` | 监听地址 |
| `settings` | `string` | Xray inbound `settings` JSON 字符串 |
| `streamSettings` | `string` | Xray inbound `streamSettings` JSON 字符串 |
| `sniffing` | `string` | Xray inbound `sniffing` JSON 字符串 |
| `tag` | `string` | 入站标签 |
| `clientStats` | `array` | 客户端流量与配置列表 |

### 2.3 ClientTraffic 简要结构

客户端流量结构来自 `clientStats`，常见字段包括：

```json
{
  "id": 1,
  "email": "user@example.com",
  "up": 0,
  "down": 0,
  "total": 0,
  "expiryTime": 0,
  "enable": true
}
```

### 2.4 OutboundTraffics 结构

```json
{
  "id": 1,
  "tag": "direct",
  "up": 123,
  "down": 456,
  "total": 579
}
```

## 3. API Key 管理接口

## 3.1 创建 API Key

- 方法：`POST`
- 路径：`/panel/api/apikeys/create`
- 鉴权：面板登录态
- `Content-Type`：
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`
  - `application/json`

请求参数：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `name` | `string` | 是 | API Key 名称 |
| `expiresAt` | `int64` | 否 | 过期时间戳，`0` 表示不过期 |

请求示例：

```bash
curl -X POST "http://127.0.0.1:2053/panel/api/apikeys/create" \
  -b "3x-ui=your-session-cookie" \
  -d "name=automation" \
  -d "expiresAt=0"
```

成功响应示例：

```json
{
  "success": true,
  "msg": "",
  "obj": {
    "id": 1,
    "name": "automation",
    "prefix": "sk-AbCdEf12",
    "key": "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "createdAt": 1710000000,
    "expiresAt": 0
  }
}
```

说明：

- `obj.key` 只会在创建时返回一次
- 后续列表接口不会再返回完整明文 `key`

## 3.2 查询 API Key 列表

- 方法：`GET`
- 路径：`/panel/api/apikeys/list`
- 鉴权：面板登录态

成功响应示例：

```json
{
  "success": true,
  "msg": "",
  "obj": [
    {
      "id": 1,
      "name": "automation",
      "prefix": "sk-AbCdEf12",
      "userId": 1,
      "createdAt": 1710000000,
      "expiresAt": 0
    }
  ]
}
```

## 3.3 删除 API Key

- 方法：`POST`
- 路径：`/panel/api/apikeys/delete/:id`
- 鉴权：面板登录态

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | API Key ID |

成功响应示例：

```json
{
  "success": true,
  "msg": "delete api key",
  "obj": null
}
```

## 4. OpenAPI

以下全部接口均要求 API Key。

通用请求头示例：

```http
X-API-Key: sk-your-api-key
```

## 4.1 Inbounds 配置接口

## 4.1.1 查询全部入站

- 方法：`GET`
- 路径：`/openapi/inbounds/list`
- 说明：返回当前 API Key 所属用户的全部入站

成功响应：

```json
{
  "success": true,
  "msg": "",
  "obj": [
    {
      "id": 1,
      "remark": "test",
      "enable": true,
      "port": 443,
      "protocol": "vless"
    }
  ]
}
```

## 4.1.2 查询单个入站

- 方法：`GET`
- 路径：`/openapi/inbounds/get/:id`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | 入站 ID |

## 4.1.3 新增入站

- 方法：`POST`
- 路径：`/openapi/inbounds/add`
- `Content-Type`：
  - `application/json`
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

请求字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `up` | `int64` | 否 | 已用上传流量 |
| `down` | `int64` | 否 | 已用下载流量 |
| `total` | `int64` | 否 | 总流量限制 |
| `allTime` | `int64` | 否 | 累计历史流量 |
| `remark` | `string` | 否 | 备注 |
| `enable` | `bool` | 否 | 是否启用 |
| `expiryTime` | `int64` | 否 | 到期时间戳 |
| `trafficReset` | `string` | 否 | 重置策略 |
| `lastTrafficResetTime` | `int64` | 否 | 最后重置时间 |
| `clientStats` | `array` | 否 | 客户端配置列表 |
| `listen` | `string` | 否 | 监听地址 |
| `port` | `int` | 是 | 入站端口 |
| `protocol` | `string` | 是 | 协议类型 |
| `settings` | `string` | 是 | JSON 字符串 |
| `streamSettings` | `string` | 否 | JSON 字符串 |
| `tag` | `string` | 否 | 标签，不传时服务端自动生成 |
| `sniffing` | `string` | 否 | JSON 字符串 |

请求示例：

```json
{
  "remark": "demo",
  "enable": true,
  "port": 443,
  "protocol": "vless",
  "settings": "{\"clients\":[]}",
  "streamSettings": "{}",
  "sniffing": "{}"
}
```

标准 JSON Body 示例：

```json
{
  "up": 0,
  "down": 0,
  "total": 0,
  "allTime": 0,
  "remark": "demo",
  "enable": true,
  "expiryTime": 0,
  "trafficReset": "never",
  "lastTrafficResetTime": 0,
  "clientStats": [
    {
      "id": 0,
      "inboundId": 0,
      "enable": true,
      "email": "user@example.com",
      "up": 0,
      "down": 0,
      "allTime": 0,
      "expiryTime": 0,
      "total": 0,
      "reset": 0,
      "lastOnline": 0
    }
  ],
  "listen": "0.0.0.0",
  "port": 443,
  "protocol": "vless",
  "settings": "{\"clients\":[{\"id\":\"uuid\",\"email\":\"user@example.com\",\"enable\":true}]}",
  "streamSettings": "{\"network\":\"tcp\",\"security\":\"none\"}",
  "tag": "inbound-443",
  "sniffing": "{\"enabled\":true,\"destOverride\":[\"http\",\"tls\"]}"
}
```

Body 字段详细说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `up` | `int64` | 否 | 当前累计上传流量，通常新增时传 `0` |
| `down` | `int64` | 否 | 当前累计下载流量，通常新增时传 `0` |
| `total` | `int64` | 否 | 入站总流量限制，单位字节 |
| `allTime` | `int64` | 否 | 历史累计流量，通常为 `0` |
| `remark` | `string` | 否 | 入站备注名称 |
| `enable` | `bool` | 否 | 是否启用 |
| `expiryTime` | `int64` | 否 | 入站到期时间戳，`0` 表示不过期 |
| `trafficReset` | `string` | 否 | 流量重置策略，如 `never`、`daily`、`weekly`、`monthly` |
| `lastTrafficResetTime` | `int64` | 否 | 上次流量重置时间戳 |
| `clientStats` | `array<object>` | 否 | 客户端流量信息数组 |
| `listen` | `string` | 否 | 监听地址，空值时服务端会按默认监听处理 |
| `port` | `int` | 是 | 入站端口 |
| `protocol` | `string` | 是 | 协议类型，如 `vmess`、`vless`、`trojan`、`shadowsocks` |
| `settings` | `string` | 是 | Xray inbound `settings` 的 JSON 字符串 |
| `streamSettings` | `string` | 否 | Xray inbound `streamSettings` 的 JSON 字符串 |
| `tag` | `string` | 否 | 入站标签，不传时服务端自动生成 |
| `sniffing` | `string` | 否 | Xray inbound `sniffing` 的 JSON 字符串 |

`clientStats` 子对象字段说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `id` | `int` | 否 | 数据库记录 ID，新增通常传 `0` |
| `inboundId` | `int` | 否 | 所属入站 ID，新增时通常由服务端维护 |
| `enable` | `bool` | 否 | 客户端是否启用 |
| `email` | `string` | 是 | 客户端唯一邮箱标识 |
| `up` | `int64` | 否 | 已用上传流量 |
| `down` | `int64` | 否 | 已用下载流量 |
| `allTime` | `int64` | 否 | 历史总流量 |
| `expiryTime` | `int64` | 否 | 客户端到期时间戳 |
| `total` | `int64` | 否 | 客户端总流量限制 |
| `reset` | `int` | 否 | 重置周期天数 |
| `lastOnline` | `int64` | 否 | 最后在线时间戳 |

`settings` 字段内容说明：

- 该字段本身是 `string`
- 字符串内容必须是合法 JSON
- 其具体结构取决于 `protocol`
- 常见内容是 `clients`、`decryption`、`fallbacks` 等 Xray inbound 配置

`streamSettings` 字段内容说明：

- 该字段本身是 `string`
- 字符串内容必须是合法 JSON
- 常见字段包括 `network`、`security`、`tlsSettings`、`realitySettings`、`wsSettings`、`grpcSettings`

`sniffing` 字段内容说明：

- 该字段本身是 `string`
- 字符串内容必须是合法 JSON
- 常见示例：

```json
{
  "enabled": true,
  "destOverride": ["http", "tls"]
}
```

## 4.1.4 删除入站

- 方法：`POST`
- 路径：`/openapi/inbounds/del/:id`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | 入站 ID |

## 4.1.5 更新入站

- 方法：`POST`
- 路径：`/openapi/inbounds/update/:id`
- `Content-Type`：
  - `application/json`
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | 入站 ID |

请求体字段与“新增入站”一致。

标准 JSON Body 与“新增入站”完全一致，只是必须配合路径参数 `:id` 指向已存在的入站记录。

## 4.1.6 添加客户端

- 方法：`POST`
- 路径：`/openapi/inbounds/addClient`
- `Content-Type`：
  - `application/json`
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

说明：

- 该接口直接绑定 `Inbound` 对象
- 一般需要在请求体内携带包含客户端变更信息的 `settings` / `clientStats`

标准 JSON Body 示例：

```json
{
  "id": 1,
  "settings": "{\"clients\":[{\"id\":\"uuid-1\",\"email\":\"user@example.com\",\"enable\":true}]}",
  "clientStats": [
    {
      "email": "user@example.com",
      "enable": true,
      "up": 0,
      "down": 0,
      "allTime": 0,
      "expiryTime": 0,
      "total": 0,
      "reset": 0,
      "lastOnline": 0
    }
  ]
}
```

Body 说明：

- 该接口没有单独的 `client` 对象
- 服务端直接从 `Inbound` 请求体中解析客户端相关配置
- 实际使用时通常至少传：
  - `id`
  - `settings`
  - `clientStats`
- `settings` 中的客户端内容与 `clientStats` 中的邮箱、启用状态、流量限制应保持一致

## 4.1.7 删除客户端

- 方法：`POST`
- 路径：`/openapi/inbounds/:id/delClient/:clientId`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `id` | `int` | 入站 ID |
| `clientId` | `string` | 客户端 ID |

## 4.1.8 更新客户端

- 方法：`POST`
- 路径：`/openapi/inbounds/updateClient/:clientId`
- `Content-Type`：
  - `application/json`
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `clientId` | `string` | 客户端 ID |

说明：

- 请求体仍为 `Inbound` 结构
- 由服务端结合 `clientId` 更新对应客户端

标准 JSON Body 示例：

```json
{
  "id": 1,
  "settings": "{\"clients\":[{\"id\":\"uuid-1\",\"email\":\"user@example.com\",\"enable\":true}]}",
  "clientStats": [
    {
      "email": "user@example.com",
      "enable": true,
      "expiryTime": 1710000000,
      "total": 10737418240,
      "reset": 30
    }
  ]
}
```

Body 说明：

- 请求体结构与新增客户端接口一致
- `clientId` 由路径给出
- 更新内容主要通过 `settings` 和 `clientStats` 体现

## 4.1.9 导入入站

- 方法：`POST`
- 路径：`/openapi/inbounds/import`
- `Content-Type`：
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

表单参数：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `data` | `string` | 是 | 完整入站 JSON 字符串 |

请求示例：

```bash
curl -X POST "http://127.0.0.1:2053/openapi/inbounds/import" \
  -H "X-API-Key: sk-your-api-key" \
  -d 'data={"remark":"demo","port":443,"protocol":"vless","settings":"{\"clients\":[]}","streamSettings":"{}","sniffing":"{}"}'
```

`data` 字段内部 JSON 示例：

```json
{
  "remark": "demo",
  "enable": true,
  "port": 443,
  "protocol": "vless",
  "settings": "{\"clients\":[{\"id\":\"uuid-1\",\"email\":\"user@example.com\",\"enable\":true}]}",
  "streamSettings": "{\"network\":\"tcp\",\"security\":\"none\"}",
  "sniffing": "{\"enabled\":true,\"destOverride\":[\"http\",\"tls\"]}",
  "clientStats": [
    {
      "email": "user@example.com",
      "enable": true,
      "total": 0,
      "expiryTime": 0
    }
  ]
}
```

## 4.2 Inbounds 统计与流量接口

## 4.2.1 按邮箱查询客户端流量

- 方法：`GET`
- 路径：`/openapi/inbounds/getClientTraffics/:email`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `email` | `string` | 客户端邮箱 |

## 4.2.2 按入站 ID 查询客户端流量列表

- 方法：`GET`
- 路径：`/openapi/inbounds/getClientTrafficsById/:id`

## 4.2.3 查询客户端 IP 记录

- 方法：`POST`
- 路径：`/openapi/inbounds/clientIps/:email`

说明：

- `obj` 可能返回字符串 `"No IP Record"`
- 也可能返回原始 IP 记录内容

## 4.2.4 清空客户端 IP 记录

- 方法：`POST`
- 路径：`/openapi/inbounds/clearClientIps/:email`

## 4.2.5 重置单个客户端流量

- 方法：`POST`
- 路径：`/openapi/inbounds/:id/resetClientTraffic/:email`

## 4.2.6 重置所有入站流量

- 方法：`POST`
- 路径：`/openapi/inbounds/resetAllTraffics`

## 4.2.7 重置某入站下全部客户端流量

- 方法：`POST`
- 路径：`/openapi/inbounds/resetAllClientTraffics/:id`

## 4.2.8 删除流量耗尽客户端

- 方法：`POST`
- 路径：`/openapi/inbounds/delDepletedClients/:id`

## 4.2.9 查询在线客户端

- 方法：`POST`
- 路径：`/openapi/inbounds/onlines`

成功响应示例：

```json
{
  "success": true,
  "msg": "",
  "obj": [
    "user1@example.com",
    "user2@example.com"
  ]
}
```

## 4.2.10 查询最后在线时间

- 方法：`POST`
- 路径：`/openapi/inbounds/lastOnline`

成功响应示例：

```json
{
  "success": true,
  "msg": "",
  "obj": {
    "user1@example.com": 1710000000
  }
}
```

## 4.2.11 手动更新客户端流量

- 方法：`POST`
- 路径：`/openapi/inbounds/updateClientTraffic/:email`
- `Content-Type`：`application/json`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `email` | `string` | 客户端邮箱 |

请求体：

```json
{
  "upload": 123,
  "download": 456
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `upload` | `int64` | 是 | 要写入的上传流量值 |
| `download` | `int64` | 是 | 要写入的下载流量值 |

说明：

- 该接口使用 `application/json`
- body 中只允许这两个字段
- 数值单位为字节

## 4.2.12 按邮箱删除客户端

- 方法：`POST`
- 路径：`/openapi/inbounds/:id/delClientByEmail/:email`

## 4.3 Xray 接口

## 4.3.1 获取默认 Xray 配置模板

- 方法：`GET`
- 路径：`/openapi/xray/getDefaultJsonConfig`

## 4.3.2 获取出站流量统计

- 方法：`GET`
- 路径：`/openapi/xray/getOutboundsTraffic`

返回 `obj` 为 `OutboundTraffics[]`。

## 4.3.3 获取当前 Xray 运行结果

- 方法：`GET`
- 路径：`/openapi/xray/getXrayResult`

返回示例：

```json
{
  "success": true,
  "msg": "",
  "obj": "running"
}
```

## 4.3.4 获取当前 Xray 设置

- 方法：`POST`
- 路径：`/openapi/xray/`

返回 `obj` 为字符串化 JSON，结构如下：

```json
{
  "xraySetting": {},
  "inboundTags": [],
  "outboundTestUrl": "https://www.google.com/generate_204"
}
```

## 4.3.5 更新 Xray 设置

- 方法：`POST`
- 路径：`/openapi/xray/update`
- `Content-Type`：
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

表单参数：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `xraySetting` | `string` | 是 | Xray 配置内容 |
| `outboundTestUrl` | `string` | 否 | 出站连通性测试 URL |

`xraySetting` 字段内部 JSON 示例：

```json
{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [
    {
      "tag": "inbound-443",
      "listen": "0.0.0.0",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": []
      },
      "streamSettings": {
        "network": "tcp",
        "security": "none"
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "tag": "direct"
    }
  ]
}
```

Body 说明：

- 接口本身接收的是表单，不是 JSON body
- 其中 `xraySetting` 这个表单字段的值必须是完整合法 JSON 字符串
- `outboundTestUrl` 为普通字符串 URL

## 4.3.6 重置出站流量

- 方法：`POST`
- 路径：`/openapi/xray/resetOutboundsTraffic`
- `Content-Type`：
  - `application/x-www-form-urlencoded`
  - `multipart/form-data`

表单参数：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `tag` | `string` | 是 | 出站标签 |

## 4.4 Server 接口

## 4.4.1 获取服务器状态

- 方法：`GET`
- 路径：`/openapi/server/status`

返回 `obj` 为服务器状态对象，字段由当前运行时状态决定，通常包含 CPU、内存、磁盘、网络等统计信息。

## 4.4.2 获取 CPU 历史

- 方法：`GET`
- 路径：`/openapi/server/cpuHistory/:bucket`

路径参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| `bucket` | `int` | 聚合粒度，单位秒 |

允许值：

- `2`
- `30`
- `60`
- `120`
- `180`
- `300`

错误响应示例：

```json
{
  "success": false,
  "msg": "invalid bucket (unsupported bucket)",
  "obj": null
}
```

## 4.4.3 获取当前配置 JSON

- 方法：`GET`
- 路径：`/openapi/server/getConfigJson`

返回 `obj` 为当前服务端生成并生效的配置 JSON。

## 5. 接口调用示例

## 5.1 使用 API Key 查询入站列表

```bash
curl "http://127.0.0.1:2053/openapi/inbounds/list" \
  -H "X-API-Key: sk-your-api-key"
```

## 5.2 使用 Query 方式传递 API Key

```bash
curl "http://127.0.0.1:2053/openapi/server/status?api_key=sk-your-api-key"
```

## 5.3 更新客户端流量

```bash
curl -X POST "http://127.0.0.1:2053/openapi/inbounds/updateClientTraffic/user@example.com" \
  -H "X-API-Key: sk-your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"upload":123,"download":456}'
```

## 5.4 更新 Xray 配置

```bash
curl -X POST "http://127.0.0.1:2053/openapi/xray/update" \
  -H "X-API-Key: sk-your-api-key" \
  -d "xraySetting={}" \
  -d "outboundTestUrl=https://www.google.com/generate_204"
```

## 6. 备注

- 当前 `OpenAPI` 仅开放 Inbounds、流量统计、Xray、Server 状态相关接口
- 面板设置接口 `/panel/setting/*` 未开放为对外 API
- 文档中的字段说明以当前控制器和模型为准
- 某些接口底层直接复用面板内部结构，因此请求体中会出现 JSON 字符串字段，例如 `settings`、`streamSettings`、`sniffing`
