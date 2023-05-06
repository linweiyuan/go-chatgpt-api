# go-chatgpt-api

### [英文文档](README.md)

## 一个尝试绕过 `Cloudflare 403` 和 `Access denied` 的正向代理程序

### 实验性质项目，不保证稳定性和向后兼容，使用风险自负

| 版本 | 分支       | 镜像                                 | 特点                                                       |
|----|----------|------------------------------------|----------------------------------------------------------|
| 新版 | `main`   | `linweiyuan/go-chatgpt-api:latest` | `API` 直连，仅需跑一个容器即可（不算 `warp` 和 `cookies`）                | 
| 旧版 | `legacy` | `linweiyuan/go-chatgpt-api:legacy` | 基于浏览器，需要额外跑 `linweiyuan/chatgpt-proxy-server`（不算 `warp`） | 

新版是趋势，但是如果用新版遇到问题，可以撤回旧版，旧版还能用，并且基于浏览器应该是比较终极的解决方案，理论上很长一段时间都可以用（旧版资源占用会多一点）

---

### 使用的过程中遇到问题应该如何解决

汇总贴：https://github.com/linweiyuan/go-chatgpt-api/issues/74

---

### 支持的 API（URL 和参数基本保持着和官网一致，部分接口有些许改动）

---

### ChatGPT APIs

---

- `ChatGPT` 登录（返回 `accessToken`）（目前仅支持 `ChatGPT` 账号，谷歌或微软账号没有测试）

`POST /chatgpt/login`

<details>

```json
{
  "username": "email",
  "password": "password"
}
```

</details>

---

- 获取对话列表（历史记录）

`GET /chatgpt/conversations?offset=0&limit=20`

`offset` 不传默认为 0, `limit` 不传默认为 20 (最大为 100)

---

- 获取对话内容

`GET /chatgpt/conversation/{conversationID}`

---

- 新建对话

`POST /chatgpt/conversation`

<details>

```json
{
  "action": "next",
  "messages": [
    {
      "id": "message id",
      "author": {
        "role": "user"
      },
      "content": {
        "content_type": "text",
        "parts": [
          "Hello World"
        ]
      }
    }
  ],
  "parent_message_id": "parent message id",
  "conversation_id": "conversation id",
  "model": "text-davinci-002-render-sha",
  "timezone_offset_min": -480,
  "history_and_training_disabled": false
}
```

</details>

---

- 生成对话标题

`POST /chatgpt/conversation/gen_title/{conversationID}`

<details>

```json
{
  "message_id": "role assistant response message id"
}
```

</details>

---

- 重命名对话标题

`PATCH /chatgpt/conversation/{conversationID}`

<details>

```json
{
  "title": "new title"
}
```

</details>

---

- 删除单个对话

`PATCH /chatgpt/conversation/{conversationID}`

<details>

```json
{
  "is_visible": false
}
```

</details>

---

- 删除所有对话

`PATCH /chatgpt/conversations`

<details>

```json
{
  "is_visible": false
}
```

</details>

---

- 消息反馈

`POST /chatgpt/conversation/message_feedback`

<details>

```json
{
  "message_id": "message id",
  "conversation_id": "conversation id",
  "rating": "thumbsUp/thumbsDown"
}
```

</details>

---

### Platform APIs

---

- `platform` 登录（返回 `sessionKey`）

`POST /platform/login`

<details>

```json
{
  "username": "email",
  "password": "password"
}
```

</details>

---

- [List models](https://platform.openai.com/docs/api-reference/models/list)

`GET /platform/v1/models`

---

- [Retrieve model](https://platform.openai.com/docs/api-reference/models/retrieve)

`GET /platform/v1/models/{model}`

---

- [Create completion](https://platform.openai.com/docs/api-reference/completions/create)

`POST /platform/v1/completions`

<details>

```json
{
  "model": "text-davinci-003",
  "prompt": "Say this is a test",
  "max_tokens": 7,
  "temperature": 0,
  "stream": true
}
```

</details>

---

- [Create chat completion](https://platform.openai.com/docs/api-reference/chat/create)

`POST /platform/v1/chat/completions`

<details>

```json
{
  "messages": [
    {
      "role": "user",
      "content": "Hello World"
    }
  ],
  "model": "gpt-3.5-turbo",
  "stream": true
}
```

</details>

---

- [Create edit](https://platform.openai.com/docs/api-reference/edits/create)

`POST /platform/v1/edits`

<details>

```json
{
  "model": "text-davinci-edit-001",
  "input": "What day of the wek is it?",
  "instruction": "Fix the spelling mistakes"
}
```

</details>

---

- [Create image](https://platform.openai.com/docs/api-reference/images/create)

`POST /platform/v1/images/generations`

<details>

```json
{
  "prompt": "A cute dog",
  "n": 2,
  "size": "1024x1024"
}
```

</details>

---

- [Create embeddings](https://platform.openai.com/docs/api-reference/embeddings/create)

`POST /platform/v1/embeddings`

<details>

```json
{
  "model": "text-embedding-ada-002",
  "input": "The food was delicious and the waiter..."
}
```

</details>

---
 
- [List files](https://platform.openai.com/docs/api-reference/files/list)

`GET /platform/v1/files`

---

- 获取 `credit grants` （只能传 `sessionKey`）

`GET /platform/dashboard/billing/credit_grants`

---

- 获取 `subscription` （只能传 `sessionKey`）

`GET /platform/dashboard/billing/subscription`

---

- 获取 `api keys` （只能传 `sessionKey`）

`GET /platform/dashboard/user/api_keys`

---

如需设置代理，可以设置环境变量 `GO_CHATGPT_API_PROXY`，比如 `GO_CHATGPT_API_PROXY=http://127.0.0.1:20171`
或者 `GO_CHATGPT_API_PROXY=socks5://127.0.0.1:20170`，注释掉或者留空则不启用

如需配合 `warp` 使用：`GO_CHATGPT_API_PROXY=socks5://chatgpt-proxy-server-warp:65535`，因为需要设置 `warp`
的场景已经默认可以直接访问 `ChatGPT` 官网，因此共用一个变量不冲突

---

`docker-compose` 配置文件：

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - GIN_MODE=release
      - GO_CHATGPT_API_PROXY=
    restart: unless-stopped
```

我仅仅在 `Arch Linux` 上进行开发和测试，这是一个滚动更新的版本，意味着系统上所有东西都是最新的，如果你在使用的过程中 `yaml`
报错了，则可以加上 `version: '3'` 在 `services:` 前面

如果遇到 `Access denied`，但是你的服务器确实在[被支持的国家或地区](https://platform.openai.com/docs/supported-countries)
，尝试一下这个配置（不保证能解决问题，比如你的服务器在 A 地区，但 A 地不在支持列表内，即使用上了 `warp` 后是 `Cloudflare IP`
，结果也会是 `403`）：

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - GIN_MODE=release
      - GO_CHATGPT_API_PROXY=socks5://chatgpt-proxy-server-warp:65535
    depends_on:
      - chatgpt-proxy-server-warp
    restart: unless-stopped

  chatgpt-proxy-server-warp:
    container_name: chatgpt-proxy-server-warp
    image: linweiyuan/chatgpt-proxy-server-warp
    environment:
      - LOG_LEVEL=OFF
    restart: unless-stopped
```

如果要让运行的镜像总是保持最新，可以配合这个一起使用：

```yaml
services:
  watchtower:
    container_name: watchtower
    image: containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 3600
    restart: unless-stopped
```

<details>

<summary>广告位</summary>

`Vultr` 推荐链接：https://www.vultr.com/?ref=7372562

---

个人微信（没有验证，谁都能加，但是不聊日常和私人问题，不进群；可以解答程序使用问题，但最好自己要有一定的基础；可以远程调试，仅限 `SSH`
或`ToDesk`，但不保证能解决）：

![](https://linweiyuan.github.io/about/mmqrcode.png)

---

微信赞赏码（经济条件允许的可以考虑支持下）：

![](https://linweiyuan.github.io/about/mm_reward_qrcode.png)

</details>
