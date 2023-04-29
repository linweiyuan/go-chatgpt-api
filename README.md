# go-chatgpt-api

## Unofficial ChatGPT API.

### Available APIs:

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
  "message_id": "response message id"
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

**No need to run `chatgpt-proxy-server` anymore.**

---

If you need to setup a proxy, set `GO_CHATGPT_API_PROXY`, for example: `GO_CHATGPT_API_PROXY=http://127.0.0.1:20171`
or `GO_CHATGPT_API_PROXY=socks5://127.0.0.1:20170`.

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

---

If you get `Access denied`, but the server is in support countries, have a try with this:

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
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