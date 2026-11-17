<div class="guide-hero">
  <div class="guide-eyebrow">QUICK START · 约 10 分钟</div>
  <h1>零基础接入 CC-Switch</h1>
  <p>从选择客户端、安装软件，到导入 TokenPilot 密钥、配置模型和余额查询，按顺序完成第一次接入。</p>
  <div class="guide-meta"><span>推荐 Codex</span><span>Windows / macOS / Linux</span></div>
</div>

> 这份教程面向第一次使用 CC-Switch 的同学。建议从第一步开始按顺序操作，不要跳步。

## 先选客户端：推荐 Codex

CC-Switch 支持多种 AI 编程客户端，但接入成本不同：

| 客户端 | 是否推荐 | 说明 |
|---|---|---|
| **Codex** | **推荐** | 使用 OpenAI 账号即可，接入门槛低，适合大多数人 |
| Claude Code | 谨慎 | 需要自己购买 Anthropic 账号或订阅，成本较高 |
| Gemini CLI | 可选 | 需要 Gemini 相关账号 |
| Grok | 可选 | 需要对应账号和分组 |

如果你是第一次使用，建议直接选择 **Codex**。

> **为什么不优先推荐 Claude？**
>
> Claude Code 需要你自己的 Anthropic 账号或订阅，账号注册和付费门槛更高。Codex 的接入和使用成本更低，更适合先跑通流程。

## 第一步：安装你要使用的客户端

### 选择 Codex

打开官方页面：

<a href="https://openai.com/zh-Hans-CN/codex/" target="_blank" rel="noopener noreferrer">openai.com/zh-Hans-CN/codex/</a>

按照页面提示安装 Codex 桌面客户端，或者使用你已有的 Codex CLI。

### 选择 Claude Code

打开官方文档：

<a href="https://code.claude.com/docs/en/desktop" target="_blank" rel="noopener noreferrer">code.claude.com/docs/en/desktop</a>

按照文档安装 Claude Desktop / Claude Code。

> 安装完客户端后先不用登录或配置 API，下一步会用 CC-Switch 接管配置。

## 第二步：安装 CC-Switch

打开官方下载页：

<a href="https://ccswitch.io/zh/download" target="_blank" rel="noopener noreferrer">ccswitch.io/zh/download</a>

根据系统选择安装包：

| 系统 | 安装包 |
|---|---|
| Windows | `Windows.msi` |
| macOS | `macOS.dmg` |
| Linux Debian/Ubuntu | `.deb` |
| 其他 Linux | `.AppImage` |

安装完成后先启动一次 CC-Switch，让系统注册 `ccswitch://` 协议。

> **注意**
>
> 请使用安装版，不要使用 Portable 便携版。便携版通常不会注册协议，点击「导入到 CCS」时会提示未安装。

## 第三步：在 CC-Switch 里先选客户端

打开 CC-Switch，在顶部或主界面先选择你要使用的客户端：

```text
Codex → 选择 Codex
Claude Code → 选择 Claude
Gemini CLI → 选择 Gemini
```

**先选客户端，再去 TokenPilot 导入密钥。**

这样做的原因是：CC-Switch 会按当前客户端决定把供应商写入哪一份配置。顺序反了，可能出现密钥导入到了错误客户端的情况。

## 第四步：在 TokenPilot 创建 API 密钥

登录 TokenPilot，进入：

**<a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥</a> → 创建密钥**

根据刚才选择的客户端选择分组：

| 客户端 | 分组平台 |
|---|---|
| Codex | OpenAI |
| Claude Code | Anthropic |
| Gemini CLI | Gemini |
| Grok | Grok |

填写一个容易识别的名称，例如：

```text
Codex 日常
```

然后创建。

> **提示**
>
> 密钥只显示一次，创建后请立即复制保存。丢失后只能重新创建。

## 第五步：导入到 CC-Switch

在 <a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥列表</a>里找到刚创建的密钥，点击：

**导入到 CCS**

浏览器会询问是否打开 CC-Switch，选择「允许」。

CC-Switch 会自动填入：

- 供应商名称
- API 请求地址
- API Key
- 对应的客户端类型
- 用量查询脚本

> **如果提示「CC-Switch 未安装或协议处理程序未注册」**
>
> 1. 确认 CC-Switch 已经安装并启动过
> 2. 重启浏览器后重试
> 3. 不要使用 Portable 便携版
> 4. Windows 可在「默认应用」里搜索 `ccswitch` 检查协议

## 第六步：编辑供应商并获取模型列表

导入成功后，**不要直接开始使用**。

在 CC-Switch 里找到刚导入的供应商，点击「编辑」，然后点击：

**获取模型列表**

等待模型列表拉取完成。

这一步会通过你填写的 API 请求地址和 API Key 查询可用模型。成功后，模型会出现在 Codex 的 `/model` 列表里。

> 如果获取失败，先检查 API Key、API 请求地址和分组平台是否正确，再重试。

### 选择默认模型

模型列表拉取成功后：

1. 从列表中选择你要使用的模型
2. 列表第一行通常就是默认模型
3. 保存配置

Codex 导入时默认会填：

```text
gpt-5.5
```

如果你的分组里配置了其他模型，以「获取模型列表」实际返回的模型为准。

> **模型列表说明**
>
> 名字是官方 GPT 模型的行，会按官方条目生成。行里的显示名、上下文窗口和思考等级不生效，以官方定义为准。

## 第七步：配置余额和用量查询

CC-Switch 可以查询 TokenPilot 的余额和用量。

在 CC-Switch 里找到当前供应商的用量查询/余额查询配置，点击配置：

1. 选择 **Sub2API**
2. 确认 API 请求地址和 API Key
3. 保存

保存后，CC-Switch 就会显示当前密钥的余额和用量。

> 导入时通常已经自动带上了用量查询脚本。如果余额显示失败，重新按上面的步骤选择 Sub2API 并保存一次。

## 第八步：重启客户端，开始使用

**每次在 CC-Switch 中修改配置并保存后，都要重启对应的客户端。**

包括修改：

- API 请求地址
- API Key
- 模型列表
- 默认模型
- 子代理模型
- 自动审批模型
- 用量查询配置

重启方式：

```text
Codex → 完全退出 Codex，再重新打开
Claude Code → 完全退出 Claude Desktop / Claude Code，再重新打开
```

只关闭窗口不算重启，要确保进程完全退出。

然后打开客户端，选择刚配置的供应商即可。

## 手动添加供应商

除了从 TokenPilot 导入，也可以在 CC-Switch 里手动新建供应商。

下面是 Codex 的配置字段说明。

### 名称

必填，用于区分不同供应商。

```text
示例：TokenPilot-Codex
```

### 备注

可选，只用于自己识别。

```text
示例：公司专用账号
```

### 官网链接

可选，填写你的站点地址：

```text
https://tokenpilot.site
```

### API Key

填写 TokenPilot 创建的密钥：

```text
sk-...
```

只需要填这里，下方 `auth.json` 会自动填充。

### API 请求地址

填写兼容 OpenAI Response 格式的服务端点：

```text
https://tokenpilot.site/v1
```

注意：

- Codex / OpenAI 使用 `/v1`
- Claude Code 使用：

```text
https://tokenpilot.site
```

不要给 Claude 的地址加 `/v1`。

### 管理与测速

填写同一个 API 地址，用于管理和测速：

```text
https://tokenpilot.site/v1
```

### 模型列表

点击「获取模型列表」，从 TokenPilot 拉取可用模型。

也可以手动添加。模型会出现在 Codex 的 `/model` 列表里，第一行是默认模型。

### 默认模型

Codex 默认请求的模型，必须来自上面的模型列表。

```text
示例：gpt-5.5
```

具体以你所在分组实际返回的模型为准。

### 子代理模型

Codex 派出的子代理默认使用的模型和推理强度。

建议：

```text
留空，跟随主模型
```

只填模型、不填强度时，使用该模型自己的默认强度。

### 默认强度

子代理的默认推理强度。

建议：

```text
留空，使用模型默认强度
```

### 子代理并发上限

一个 Codex 会话里最多同时开启多少个子代理，不含主代理。

建议：

```text
留空，使用 Codex 默认值
```

最小为 1，不支持不限。新对话才会生效。

### 自动审批模型

Codex 开启「Approve for me」自动审批时，用来判断操作能否放行的模型。

注意：

- 只能选择当前模型列表里的模型
- 审批请求和主对话发给同一家供应商
- 修改后需要重启 Codex

如果不使用自动审批，可以留空。

## 地址速查

| 客户端 | API 请求地址 |
|---|---|
| Codex / OpenAI | `https://tokenpilot.site/v1` |
| Claude Code | `https://tokenpilot.site`，不要带 `/v1` |
| Gemini CLI | `https://tokenpilot.site` |
| Grok | `https://tokenpilot.site/v1` |

从 TokenPilot 点击「导入到 CCS」时，系统会自动按分组平台填好地址，一般不需要手动改。

## 常见问题

### 点击导入没有反应？

浏览器拦截了自定义协议。检查地址栏右侧是否有弹窗提示，允许打开 CC-Switch。

### 导入后请求失败？

按顺序检查：

1. 密钥状态是否为「启用」
2. 分组平台是否和客户端一致
3. 分组是否绑定了渠道
4. 渠道是否配置了对应模型
5. 账号是否有可用额度
6. API 请求地址是否符合上面的地址速查

### 获取模型列表失败？

1. 确认 API 请求地址带不带 `/v1` 正确
2. 确认 API Key 有权限
3. 确认分组绑定了渠道
4. 确认渠道支持 `/v1/models`
5. 保存后重启客户端再试

### 余额不显示？

1. 打开 CC-Switch 的用量查询配置
2. 选择 Sub2API
3. 确认 API Key 和地址
4. 保存
5. 重启客户端

### 修改配置后不生效？

这是最常见的问题。

```text
在 CC-Switch 保存配置 → 完全退出客户端 → 重新打开客户端
```

不要只关闭窗口，要确保后台进程也退出。

### 想换一个客户端？

不要重复创建密钥。回到 <a href="/keys" target="_blank" rel="noopener noreferrer">API 密钥页</a>，重新点击「导入到 CCS」，选择另一个客户端即可。

## 完成

到这里就接好了。

以后再修改 CC-Switch 配置时，记住这个顺序：

```text
选择客户端 → 修改配置 → 保存 → 重启客户端
```

祝你使用顺利。
