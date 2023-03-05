# go-chatgpt-api

Unofficial ChatGPT API (web based, not the GPT-3 API).

Java Swing GUI application: [ChatGPT-Swing](https://github.com/linweiyuan/ChatGPT-Swing)

Golang TUI application: [go-chatgpt](https://github.com/linweiyuan/go-chatgpt)

---

## APIs

If not mentioned, the `Authorization: accessToken` header is mandatory (no need to pass Bearer, API will do).

### POST /user/login

```json
{
  "username": "",
  "password": ""
}
```

This is different from the official API which uses form data in login, all in JSON in this project.

No need accessToken.

---

### GET /auth/session

Renew accessToken.

No need accessToken, but need cookies from login API.

---

### GET /conversations

Return a list of conversations, currently hard-coded 100 (max) and only return these, good enough for normal use cases.

---

### POST /conversation

```json
{
  "message_id": "",
  "parent_message_id": "",
  "conversation_id": "",
  "content": ""
}
```

`message_id`: always a new UUID.

`parent_message_id`: if you start a new conversation, pass a new UUID, otherwise, get this id from response.

`conversation_id`: if you start a new conversation, pass null, otherwise, get this id from response.

---

### POST /conversation/gen_title/:id

If you start a new conversation, the conversation id will be returned, then pass this returned id to gen title.

---

### GET /conversation/:id

Get all messages of this conversation.

---

### PATCH (POST) /conversation/:id

Rename or delete this conversation.

If rename:

```json
{
  "title": ""
}
```

If delete:

```json
{
  "is_visible": false
}
```

---

### POST /conversation/message_feedback

```json
{
  "message_id": "",
  "conversation_id": "",
  "rating": ""
}
```

If like, pass "thumbsUp" for `rating`, otherwise, "thumbsDown".

---

### PATCH (POST) /conversations

Clear all conversations.

```json
{
  "is_visible": false
}
```