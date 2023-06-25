# go-chatgpt-api

## 一个尝试绕过 `Cloudflare` 来使用 `ChatGPT` 接口的程序

### 实验性质练手项目，不保证稳定性和向后兼容，使用风险自负

---

### 相关博客（程序更新很多次，文章的内容可能和现在的不一样，供参考）：[ChatGPT](https://linweiyuan.github.io/categories/ChatGPT/)

- [如何生成 GPT-4 arkose_token](https://linweiyuan.github.io/2023/06/24/%E5%A6%82%E4%BD%95%E7%94%9F%E6%88%90-GPT-4-arkose-token.html)
- [利用 HTTP Client 来调试 go-chatgpt-api](https://linweiyuan.github.io/2023/06/18/%E5%88%A9%E7%94%A8-HTTP-Client-%E6%9D%A5%E8%B0%83%E8%AF%95-go-chatgpt-api.html)
- [一种解决 ChatGPT Access denied 的方法](https://linweiyuan.github.io/2023/04/15/%E4%B8%80%E7%A7%8D%E8%A7%A3%E5%86%B3-ChatGPT-Access-denied-%E7%9A%84%E6%96%B9%E6%B3%95.html)
- [ChatGPT 如何自建代理](https://linweiyuan.github.io/2023/04/08/ChatGPT-%E5%A6%82%E4%BD%95%E8%87%AA%E5%BB%BA%E4%BB%A3%E7%90%86.html)
- [一种取巧的方式绕过 Cloudflare v2 验证](https://linweiyuan.github.io/2023/03/14/%E4%B8%80%E7%A7%8D%E5%8F%96%E5%B7%A7%E7%9A%84%E6%96%B9%E5%BC%8F%E7%BB%95%E8%BF%87-Cloudflare-v2-%E9%AA%8C%E8%AF%81.html)

---

### 如果有疑问而不是什么程序出错其实可以在 [Discussions](https://github.com/linweiyuan/go-chatgpt-api/discussions) 里发而不是新增 Issue

---

### 使用的过程中遇到问题应该如何解决

汇总贴：https://github.com/linweiyuan/go-chatgpt-api/issues/74

---

### 范例（URL 和参数基本保持着和官网一致，部分接口有些许改动）

部分例子，不是全部，**理论上**全部基于文本传输的接口都支持

https://github.com/linweiyuan/go-chatgpt-api/tree/main/example （需安装 `HTTP Client` 插件）

---

### 配置

如需设置代理，可以设置环境变量 `GO_CHATGPT_API_PROXY`，比如 `GO_CHATGPT_API_PROXY=http://127.0.0.1:20171`
或者 `GO_CHATGPT_API_PROXY=socks5://127.0.0.1:20170`，注释掉或者留空则不启用

如需配合 `warp` 使用：`GO_CHATGPT_API_PROXY=socks5://chatgpt-proxy-server-warp:65535`，因为需要设置 `warp`
的场景已经默认可以直接访问 `ChatGPT` 官网，因此共用一个变量不冲突（国内 `VPS` 不在讨论范围内）

GPT-4 相关模型需要验证 `arkose_token`，配置环境变量 `GO_CHATGPT_API_ARKOSE_TOKEN_URL`
为 `https://arkose-token.linweiyuan.com` 即可获取不是很合法的 `arkose_token` （测试过很多账号都可以 403 -> 200，部分依旧
403）

---

`docker-compose` 配置文件：

```yaml
services:
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - TZ=Asia/Shanghai
      - GO_CHATGPT_API_PROXY=
      - GO_CHATGPT_API_ARKOSE_TOKEN_URL=
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
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - TZ=Asia/Shanghai
      - GO_CHATGPT_API_PROXY=socks5://chatgpt-proxy-server-warp:65535
      - GO_CHATGPT_API_ARKOSE_TOKEN_URL=
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

如果你知道什么是 `teams-enroll-token`，可以通过环境变量 `TEAMS_ENROLL_TOKEN` 设置它的值

然后利用这条命令来检查是否生效:

`docker-compose exec chatgpt-proxy-server-warp warp-cli --accept-tos account | awk 'NR==1'`

```
Account type: Free （没有生效）

Account type: Team （设置正常）
```

---

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

这个只会更新新镜像，旧的镜像如果没手动删除还会在本地，如果新镜像不适用，将 `<none>` 镜像重新打 `tag`
即可，比如：`docker tag <IMAGE_ID> linweiyuan/go-chatgpt-api`，这样就完成了回滚

---

### 如何集成其他第三方客户端

- [moeakwak/chatgpt-web-share](https://github.com/moeakwak/chatgpt-web-share)

环境变量

```
CHATGPT_BASE_URL=https://go-chatgpt-api.linweiyuan.com/chatgpt/backend-api/
```

- [lss233/chatgpt-mirai-qq-bot](https://github.com/lss233/chatgpt-mirai-qq-bot)

`config.cfg`

```
[openai]
browserless_endpoint = "https://go-chatgpt-api.linweiyuan.com/chatgpt/backend-api/"
```

- [Kerwin1202/chatgpt-web](https://github.com/Kerwin1202/chatgpt-web) | [Chanzhaoyu/chatgpt-web](https://github.com/Chanzhaoyu/chatgpt-web)

环境变量

```
API_REVERSE_PROXY=https://go-chatgpt-api.linweiyuan.com/chatgpt/backend-api/conversation
```

- [pengzhile/pandora](https://github.com/pengzhile/pandora)（不完全兼容）

环境变量

```
go-chatgpt-api: GO_CHATGPT_API_PANDORA=1

pandora: CHATGPT_API_PREFIX=https://go-chatgpt-api.linweiyuan.com
```

---

- [1130600015/feishu-chatgpt](https://github.com/1130600015/feishu-chatgpt)

`application.yaml`

```yaml
proxy:
  url: https://go-chatgpt-api.linweiyuan.com
```

### 如何控制打包行为
Fork 此项目后，可以在 `Settings-Secrets and variables-Actions` 下控制如下行为：
`Secrets` 页添加 `DOCKER_HUB_TOKEN` 即可自行打包推送到个人的 Dockerhub 账户下（[如何申请 token](https://docs.docker.com/docker-hub/access-tokens/)）

`Variables` 页添加 `USE_GHCR=1` 即可推送到个人的 GHCR 仓库（[需要开启仓库的写入权限](https://stackoverflow.com/questions/75926611/github-workflow-to-push-docker-image-to-ghcr-io)）
`Variables` 页添加 `PLATFORMS=linux/amd64,linux/arm64` 即可同时打包 amd64 和 arm64 的架构的镜像（缺省情况下只会打包 amd64）

---

### 最后感谢各位同学

<!--suppress HtmlRequiredAltAttribute -->
<a href="https://github.com/linweiyuan/go-chatgpt-api/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=linweiyuan/go-chatgpt-api" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

---

<details>

<summary>广告位</summary>

---

个人微信（没有验证，谁都能加，添加即通过，不用打招呼，直接把问题发出来，日常和私人问题不聊，不进群；可以解答程序使用问题，但最好自己要有一定的基础；可以远程调试，仅限 `SSH`
或`ToDesk`，但不保证能解决）：

![](https://linweiyuan.github.io/about/mmqrcode.png)

---

微信赞赏码（经济条件允许的可以考虑支持下）：

![](https://linweiyuan.github.io/about/mm_reward_qrcode.png)

</details>
