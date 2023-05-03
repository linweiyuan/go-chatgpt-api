# go-chatgpt-api

### [中文文档](README_zh.md)

## A forward proxy program attempting to bypass `Cloudflare 403` and `Access Denied`.

| Version | Branch   | Image                              | Features                                                                                                       |
|---------|----------|------------------------------------|----------------------------------------------------------------------------------------------------------------|
| New     | `main`   | `linweiyuan/go-chatgpt-api:latest` | Direct connection to the `API`, only requires one container to run (excluding `warp` and `cookies`)            | 
| Old     | `legacy` | `linweiyuan/go-chatgpt-api:legacy` | Based on the browser, requires an additional `linweiyuan/chatgpt-proxy-server` image to run (excluding `warp`) | 

The new version is the trend, but if there are problems with the new version, you can switch back to the old version.
The old version is still usable and based on the browser, which should be an ultimate solution that can be used for a
long time (although it may consume slightly more resources).

---

### Troubleshooting

English countries does not have the "Great Firewall", so many issues are gone.

More details: https://github.com/linweiyuan/go-chatgpt-api/issues/74

---

### Supported APIs (URL and parameters are mostly consistent with the official website, with slight modifications to some interfaces).

---

- get conversation list

`GET /conversations?offset=0&limit=20`

`offset` defaults to 0, `limit` defaults to 20 (max 100).

---

- get conversation content

`GET /conversation/{conversationID}`

---

- create conversation

`POST /conversation`

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

- generate conversation title

`POST /conversation/gen_title/{conversationID}`

<details>

```json
{
  "message_id": "role assistant response message id"
}
```

</details>

---

- rename conversation

`PATCH /conversation/{conversationID}`

<details>

```json
{
  "title": "new title"
}
```

</details>

---

- delete conversation

`PATCH /conversation/{conversationID}`

<details>

```json
{
  "is_visible": false
}
```

</details>

---

- delete all conversations

`PATCH /conversations`

<details>

```json
{
  "is_visible": false
}
```

</details>

---

- feedback message

`POST /conversation/message_feedback`

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

- login (currently only support ChatGPT accounts)

`POST /auth/login`

<details>

```json
{
  "username": "email",
  "password": "password"
}
```

</details>

---

- chat completion (apiKey)

`POST /v1/chat/completions`

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

### Configuration

To set a proxy, you can use the environment variable `GO_CHATGPT_API_PROXY`, such
as `GO_CHATGPT_API_PROXY=http://127.0.0.1:20171` or `GO_CHATGPT_API_PROXY=socks5://127.0.0.1:20170`. If it is commented
out or left blank, it will not be enabled.

To use with `warp`: `GO_CHATGPT_API_PROXY=socks5://chatgpt-proxy-server-warp:65535`. Since the scenario that requires
setting up `warp` can directly access the `ChatGPT` website by default, using the same variable will not cause
conflicts.

---

`docker-compose.yaml`:

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

I only develop and test on `Arch Linux`, which is a `rolling` release version, meaning that everything on the system is
`up-to-date`. If you encounter a `yaml` error while using it, you can add `version: '3'` in front of `services:`.

If you encounter an `Access denied` error, but your server is indeed
in [Supported countries and territories](https://platform.openai.com/docs/supported-countries), try this
configuration (it is not guaranteed to solve the problem, for example, if your server is in `Zone A`, but `Zone A`
is not on the list of supported countries, even if you use `warp` to change to a `Cloudflare IP`, the result will still
be
`403`):

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

If you want to make sure the image is always latest, try this:

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

<summary>AD</summary>

`Vultr` Referral Program: https://www.vultr.com/?ref=7372562

</details>
