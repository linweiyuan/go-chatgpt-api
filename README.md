# go-chatgpt-api

Unofficial ChatGPT API.

- get conversation list
- get conversation content
- start conversation
- gen conversation title
- rename conversation
- delete conversation
- delete all conversations
- feedback message
- chat completion (apiKey)

---

**No need to run `chatgpt-proxy-server` anymore.**

`compose.yaml`:
```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - GIN_MODE=release
    restart: unless-stopped
```
