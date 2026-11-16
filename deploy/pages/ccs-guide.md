<div class="guide-hero">
  <div class="guide-eyebrow">QUICK START · 约 5 分钟</div>
  <h1>零基础接入 CC-Switch</h1>
  <p>从安装客户端到在 TokenPilot 导入密钥，一步一步完成你的第一次接入。跟着下面的顺序操作即可。</p>
  <div class="guide-meta"><span>4 个步骤</span><span>Windows / macOS / Linux</span></div>
</div>

> 这份教程面向第一次使用 CC-Switch 的同学。按顺序完成，大约 5 分钟。

## 你将完成什么

CC-Switch 是一个桌面端的 AI 编程客户端供应商切换工具。接入完成后，你可以在 CC-Switch 里一键切换 TokenPilot 的 Claude Code、Codex、Gemini CLI 等配置。

整个流程只有四步：

1. 安装 CC-Switch
2. 在 TokenPilot <a href="/keys" target="_blank" rel="noopener noreferrer">创建 API 密钥</a>
3. 点击「导入到 CCS」
4. 在 CC-Switch 里切换并开始使用

## 第一步：安装 CC-Switch

打开官方下载页：

<a href="https://ccswitch.io/zh/download" target="_blank" rel="noopener noreferrer">ccswitch.io/zh/download</a>

根据系统选择安装包：

| 系统 | 安装包 |
|---|---|
| Windows | `Windows.msi` |
| macOS | `macOS.dmg` |
| Linux Debian/Ubuntu | `.deb` |
| 其他 Linux | `.AppImage` |

安装完成后先启动一次 CC-Switch，让系统注册 `ccswitch://` 协议，然后可以关闭窗口。

> **注意**
>
> 请使用安装版，不要用 Portable 便携版。便携版通常不会注册协议，点击「导入到 CCS」会提示未安装。

## 第二步：创建 API 密钥

登录 TokenPilot，进入：

**<a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥</a> → 创建密钥**

建议按你要使用的客户端选择分组：

| 客户端 | 分组平台 |
|---|---|
| Claude Code | Anthropic |
| Codex | OpenAI |
| Gemini CLI | Gemini |
| Grok | Grok |

填写一个容易识别的名称，例如「Claude Code 日常」，然后创建。

> **提示**
>
> 密钥只显示一次，创建后请立即复制保存。丢失后只能重新创建。

## 第三步：导入到 CC-Switch

在 <a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥列表</a>里找到刚才创建的密钥，点击：

**导入到 CCS**

浏览器会询问是否打开 CC-Switch，选择「允许」。

CC-Switch 会自动填入：

- 供应商名称
- API 端点地址
- API Key
- 对应的客户端类型
- 用量查询脚本

> **如果提示「CC-Switch 未安装或协议处理程序未注册」**
>
> 1. 确认 CC-Switch 已经安装并启动过
> 2. 重启浏览器后重试
> 3. 不要使用 Portable 便携版
> 4. Windows 可在「默认应用」里搜索 `ccswitch` 检查协议

## 第四步：在 CC-Switch 里使用

打开 CC-Switch，确认导入的供应商信息：

- Endpoint 是否正确
- API Key 是否完整
- 客户端类型是否匹配

然后选中这个供应商，点击切换。

切换成功后，打开对应的命令行客户端即可：

```text
Claude Code → claude
Codex → codex
Gemini CLI → gemini
```

## 地址说明

不同客户端的 Base URL 规则不同：

| 客户端 | 地址 |
|---|---|
| Claude Code | `https://tokenpilot.site`，不要带 `/v1` |
| Codex / OpenAI | `https://tokenpilot.site/v1` |
| Gemini CLI | `https://tokenpilot.site` |
| Grok | `https://tokenpilot.site/v1` |

导入功能会自动按分组平台选择正确格式，一般不需要手动改。

## 常见问题

### 点击导入没有反应？

浏览器拦截了自定义协议。检查地址栏右侧是否有关闭的弹窗提示，允许打开 CC-Switch。

### 导入后请求失败？

按顺序检查：

1. 密钥状态是否为「启用」
2. 分组是否绑定了渠道
3. 渠道是否配置了对应模型
4. 账号是否有可用额度
5. CC-Switch 里的 Endpoint 是否和本文地址说明一致

### 想换一个客户端？

不要重复创建密钥。回到 <a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥页</a>，重新点击「导入到 CCS」，选择另一个客户端即可。

### 可以手动配置吗？

可以。在 CC-Switch 里新建供应商，填入：

```text
Base URL: 按上面的地址说明填写
API Key: 你的 TokenPilot 密钥
Model: 按分组内配置的模型填写
```

## 完成

到这里就接好了。后续 TokenPilot 的密钥、分组、模型有变化时，只需要在 <a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥页</a>重新导入一次。

祝你使用顺利。
