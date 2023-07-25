# go-chatgpt-api

## 一个尝试绕过 `Cloudflare` 来使用 `ChatGPT` 接口的程序

---

### 支持接口

- https://chat.openai.com/auth/login 登录返回 `accessToken`（谷歌和微软账号暂不支持登录，但可正常使用其他接口）
- 模型和插件查询
- `GPT-3.5` 和 `GPT-4` 对话增删改查及分享
- https://platform.openai.com/playground 登录返回 `apiKey`
- `apiKey` 余额查询
- 等等 ...
- 支持 `ChatGPT` 转 `API`，接口 `/imitate/v1/chat/completions`，利用 `accessToken` 模拟 `apiKey`，实现伪免费使用 `API`
  ，从而支持集成仅支持 `apiKey` 调用的第三方客户端项目

范例（URL 和参数基本保持着和官网一致，部分接口有些许改动），部分例子，不是全部，**理论上**全部基于文本传输的接口都支持

https://github.com/linweiyuan/go-chatgpt-api/tree/main/example （需安装 `HTTP Client` 插件）

---

### 使用的过程中遇到问题应该如何解决

汇总贴：https://github.com/linweiyuan/go-chatgpt-api/issues/74

如果有疑问而不是什么程序出错其实可以在 [Discussions](https://github.com/linweiyuan/go-chatgpt-api/discussions) 里发而不是新增
Issue

群聊：https://github.com/linweiyuan/go-chatgpt-api/discussions/197

---

### 配置

如需设置代理，可以设置环境变量 `PROXY`，比如 `PROXY=http://127.0.0.1:20171`
或者 `PROXY=socks5://127.0.0.1:20170`，注释掉或者留空则不启用

如果代理需账号密码验证，则 `http://username:password@ip:port` 或者 `socks5://username:password@ip:port`

如需配合 `warp` 使用：`PROXY=socks5://chatgpt-proxy-server-warp:65535`，因为需要设置 `warp`
的场景已经默认可以直接访问 `ChatGPT` 官网，因此共用一个变量不冲突（国内 `VPS` 不在讨论范围内，请自行配置网络环境，`warp`
服务在魔法环境下才能正常工作）

家庭网络无需跑 `warp` 服务，跑了也没用，会报错，仅在服务器需要

`CONTINUE_SIGNAL=1`，开启/imitate接口自动继续会话功能，留空关闭，默认关闭

---

`GPT-4` 相关模型目前需要验证 `arkose_token`
，可以配合 [linweiyuan/chatgpt-arkose-token-api](https://github.com/linweiyuan/chatgpt-arkose-token-api)
一起使用（需要自己有一定的基础和动手能力，不提供技术支持）

---

根据你的网络环境不同，可以展开查看对应配置，下面例子是基本参数，更多参数查看 [compose.yaml](https://github.com/linweiyuan/go-chatgpt-api/blob/main/compose.yaml)

<details>

<summary>直接利用现成的服务</summary>

服务器不定时维护，不保证高可用，利用这些服务导致的账号安全问题，与本项目无关

- https://go-chatgpt-api.linweiyuan.com
- https://api.tms.im

</details>

<details>

<summary>网络在直连或者通过代理的情况下可以正常访问 ChatGPT</summary>

```yaml
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped
```

</details>

<details>

<summary>服务器访问 ChatGPT 提示 "Unable to load site"</summary>

如何验证：`curl https://chat.openai.com | grep '<p>' | awk '{$1=$1;print}'`

```yaml
  go-chatgpt-api:
    container_name: go-chatgpt-api
    image: linweiyuan/go-chatgpt-api
    ports:
      - 8080:8080
    environment:
      - TZ=Asia/Shanghai
      - PROXY=socks5://chatgpt-proxy-server-warp:65535
    depends_on:
      - chatgpt-proxy-server-warp
    restart: unless-stopped

  chatgpt-proxy-server-warp:
    container_name: chatgpt-proxy-server-warp
    image: linweiyuan/chatgpt-proxy-server-warp
    restart: unless-stopped
```

</details>

---

目前 `warp` 容器检测到流量超过 1G 会自动重启，如果你知道什么是 `teams-enroll-token`
（不知道就跳过），可以通过环境变量 `TEAMS_ENROLL_TOKEN`
设置它的值，然后利用这条命令来检查是否生效

`docker-compose exec chatgpt-proxy-server-warp warp-cli --accept-tos account | awk 'NR==1'`

```
Account type: Free （没有生效）

Account type: Team （设置正常）
```

### Render部署

点击下面的按钮一键部署，缺点是免费版本冷启动比较慢

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/linweiyuan/go-chatgpt-api)

---

如果要让运行的镜像总是保持最新，可以配合这个一起使用（[官方文档](https://containrrr.dev/watchtower/arguments/#without_updating_containers)
可以设置哪些容器不更新，请自行查看）：

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

（本项目没有计划写更新日志和做历史版本管理）

---

### 如何集成其他第三方客户端

- [moeakwak/chatgpt-web-share](https://github.com/moeakwak/chatgpt-web-share)

环境变量

```
CHATGPT_BASE_URL=http://go-chatgpt-api:8080/chatgpt/backend-api/
```

- [lss233/chatgpt-mirai-qq-bot](https://github.com/lss233/chatgpt-mirai-qq-bot)

`config.cfg`

```
[openai]
browserless_endpoint = "http://go-chatgpt-api:8080/chatgpt/backend-api/"
```

- [Kerwin1202/chatgpt-web](https://github.com/Kerwin1202/chatgpt-web) | [Chanzhaoyu/chatgpt-web](https://github.com/Chanzhaoyu/chatgpt-web)

环境变量

```
API_REVERSE_PROXY=http://go-chatgpt-api:8080/chatgpt/backend-api/conversation
```

- [pengzhile/pandora](https://github.com/pengzhile/pandora)（不完全兼容）

环境变量

```
CHATGPT_API_PREFIX=http://go-chatgpt-api:8080
```

---

- [1130600015/feishu-chatgpt](https://github.com/1130600015/feishu-chatgpt)

`application.yaml`

```yaml
proxy:
  url: http://go-chatgpt-api:8080
```

---

- [Yidadaa/ChatGPT-Next-Web](https://github.com/Yidadaa/ChatGPT-Next-Web)

环境变量

```
BASE_URL=http://go-chatgpt-api:8080/imitate
```

### 如何控制打包行为

Fork 此项目后，可以在 `Settings-Secrets and variables-Actions` 下控制如下行为：
`Secrets` 页添加 `DOCKER_HUB_TOKEN` 即可自行打包推送到个人的 Dockerhub
账户下（[如何申请 token](https://docs.docker.com/docker-hub/access-tokens/)）

`Variables` 页添加 `USE_GHCR=1` 即可推送到个人的 GHCR
仓库（[需要开启仓库的写入权限](https://stackoverflow.com/questions/75926611/github-workflow-to-push-docker-image-to-ghcr-io)）
`Variables` 页添加 `PLATFORMS=linux/amd64,linux/arm64` 即可同时打包 amd64 和 arm64 的架构的镜像

---

### 相关博客（程序更新很多次，文章的内容可能和现在的不一样，仅供参考）：[ChatGPT](https://linweiyuan.github.io/categories/ChatGPT/)

- [如何生成 GPT-4 arkose_token](https://linweiyuan.github.io/2023/06/24/%E5%A6%82%E4%BD%95%E7%94%9F%E6%88%90-GPT-4-arkose-token.html)
- [利用 HTTP Client 来调试 go-chatgpt-api](https://linweiyuan.github.io/2023/06/18/%E5%88%A9%E7%94%A8-HTTP-Client-%E6%9D%A5%E8%B0%83%E8%AF%95-go-chatgpt-api.html)
- [一种解决 ChatGPT Access denied 的方法](https://linweiyuan.github.io/2023/04/15/%E4%B8%80%E7%A7%8D%E8%A7%A3%E5%86%B3-ChatGPT-Access-denied-%E7%9A%84%E6%96%B9%E6%B3%95.html)
- [ChatGPT 如何自建代理](https://linweiyuan.github.io/2023/04/08/ChatGPT-%E5%A6%82%E4%BD%95%E8%87%AA%E5%BB%BA%E4%BB%A3%E7%90%86.html)
- [一种取巧的方式绕过 Cloudflare v2 验证](https://linweiyuan.github.io/2023/03/14/%E4%B8%80%E7%A7%8D%E5%8F%96%E5%B7%A7%E7%9A%84%E6%96%B9%E5%BC%8F%E7%BB%95%E8%BF%87-Cloudflare-v2-%E9%AA%8C%E8%AF%81.html)

---

### 最后感谢各位同学

<!--suppress HtmlRequiredAltAttribute -->
<a href="https://github.com/linweiyuan/go-chatgpt-api/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=linweiyuan/go-chatgpt-api&max=-1"  alt=""/>
</a>

Made with [contrib.rocks](https://contrib.rocks).

---

<details>

<summary>广告位</summary>

---

[Vultr 推荐注册](https://www.vultr.com/?ref=7372562)

通过这条链接进去注册登录，并充值 `$10`，然后创建服务器用至少一个月，我可以得到佣金

---

微信赞赏码（经济条件允许的可以考虑支持下）：

![](https://linweiyuan.github.io/about/mm_reward_qrcode.png)

</details>
