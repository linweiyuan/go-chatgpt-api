# go-chatgpt-api

### [中文文档](README_zh.md)

## A forward proxy program attempting to bypass `Cloudflare 403` and `Access Denied`.

### Experimental project, with no guarantee of stability and backward compatibility, use at your own risk.

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

## ChatGPT APIs

---

- `ChatGPT` user login (`accessToken` will be returned) (currently `Google` or `Microsoft` accounts are not supported).

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

- get conversation list

`GET /chatgpt/conversations?offset=0&limit=20`

`offset` defaults to 0, `limit` defaults to 20 (max 100).

---

- get conversation content

`GET /chatgpt/conversation/{conversationID}`

---

- create conversation

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

- generate conversation title

`POST /chatgpt/conversation/gen_title/{conversationID}`

<details>

```json
{
  "message_id": "role assistant response message id"
}
```

</details>

---

- rename conversation

`PATCH /chatgpt/conversation/{conversationID}`

<details>

```json
{
  "title": "new title"
}
```

</details>

---

- delete conversation

`PATCH /chatgpt/conversation/{conversationID}`

<details>

```json
{
  "is_visible": false
}
```

</details>

---

- delete all conversations

`PATCH /chatgpt/conversations`

<details>

```json
{
  "is_visible": false
}
```

</details>

---

- feedback message

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

## Platform APIs

---

- `platform` user login (`sessionKey` will be returned)

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

- get `credit grants` (only support `sessionkey`)

`GET /platform/dashboard/billing/credit_grants`

---

- get `subscription` (only support `sessionkey`)

`GET /platform/dashboard/billing/subscription`

---

- get `api keys` (only support `sessionkey`)

`GET /platform/dashboard/user/api_keys`

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

If you know what `teams-enroll-token` is and want to set its value, you can do so through the environment variable `TEAMS_ENROLL_TOKEN`.

Run this command to check the result:

`docker-compose exec chatgpt-proxy-server-warp warp-cli --accept-tos account | awk 'NR==1'`

```
Account type: Free (wrong)

Account type: Team (correct)
```

---

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
