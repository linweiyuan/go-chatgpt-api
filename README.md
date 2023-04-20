# go-chatgpt-api

Bypass Cloudflare using [undetected_chromedriver](https://github.com/ultrafunkamsterdam/undetected-chromedriver) to use ChatGPT API.

- [一种取巧的方式绕过 Cloudflare v2 验证](https://linweiyuan.github.io/2023/03/14/%E4%B8%80%E7%A7%8D%E5%8F%96%E5%B7%A7%E7%9A%84%E6%96%B9%E5%BC%8F%E7%BB%95%E8%BF%87-Cloudflare-v2-%E9%AA%8C%E8%AF%81.html)
- [ChatGPT 如何自建代理](https://linweiyuan.github.io/2023/04/08/ChatGPT-%E5%A6%82%E4%BD%95%E8%87%AA%E5%BB%BA%E4%BB%A3%E7%90%86.html)
- [一种解决 ChatGPT Access denied 的方法](https://linweiyuan.github.io/2023/04/15/%E4%B8%80%E7%A7%8D%E8%A7%A3%E5%86%B3-ChatGPT-Access-denied-%E7%9A%84%E6%96%B9%E6%B3%95.html)

---

Also support official API (the way which using API key):

- Chat completion

---

Use both ChatGPT mode and API mode:

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - GIN_MODE=release
      - CHATGPT_PROXY_SERVER=http://chatgpt-proxy-server:9515
#      - NETWORK_PROXY_SERVER=http://host:port
#      - NETWORK_PROXY_SERVER=socks5://host:port
    depends_on:
      - chatgpt-proxy-server
    restart: unless-stopped

  chatgpt-proxy-server:
    container_name: chatgpt-proxy-server
    image: linweiyuan/chatgpt-proxy-server
    restart: unless-stopped
```
---

Use API mode only:

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - GIN_MODE=release
#      - NETWORK_PROXY_SERVER=http://host:port
#      - NETWORK_PROXY_SERVER=socks5://host:port
    restart: unless-stopped
```
---

If your IP is blocked, like "Access denied", try this (with [Cloudflare WARP](https://developers.cloudflare.com/warp-client/get-started/linux)):

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - GIN_MODE=release
      - CHATGPT_PROXY_SERVER=http://chatgpt-proxy-server:9515
      - NETWORK_PROXY_SERVER=socks5://chatgpt-proxy-server-warp:65535
    depends_on:
      - chatgpt-proxy-server
      - chatgpt-proxy-server-warp
    restart: unless-stopped

  chatgpt-proxy-server:
    container_name: chatgpt-proxy-server
    image: linweiyuan/chatgpt-proxy-server
    environment:
      - LOG_LEVEL=INFO
    restart: unless-stopped

  chatgpt-proxy-server-warp:
    container_name: chatgpt-proxy-server-warp
    image: linweiyuan/chatgpt-proxy-server-warp
    environment:
      - LOG_LEVEL=INFO
    restart: unless-stopped
```

---

### Client

Java Swing GUI application: [ChatGPT-Swing](https://github.com/linweiyuan/ChatGPT-Swing)

Golang TUI application: [go-chatgpt](https://github.com/linweiyuan/go-chatgpt)