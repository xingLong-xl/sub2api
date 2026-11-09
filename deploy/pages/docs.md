## 一分钟接入

| 你需要准备 | 在哪里获取 |
| --- | --- |
| 一个 TokenPilot 账号 | [tokenpilot.site](https://tokenpilot.site) |
| 一个 API Key | 控制台 →「API 密钥」→「创建密钥」 |
| 一个模型名 | [模型广场](https://tokenpilot.site/model-plaza) |
| 一个 Base URL | `https://tokenpilot.site/v1` |

> 第一次使用，建议按下面的「快速开始」完整走一遍，大约 3 分钟。

## 快速开始

### 1. 注册并登录

打开 [tokenpilot.site](https://tokenpilot.site)，点击右上角「登录」。没有账号请先注册，按页面提示完成验证后进入控制台。

### 2. 创建 API Key

1. 在左侧菜单进入「API 密钥」（`/keys`）。
2. 点击「创建密钥」，填写一个容易辨认的名称，例如 `my-app`。
3. 复制生成的密钥（以 `sk-` 开头）。**密钥只在创建时完整展示一次**，请立即保存。

> 不要把 API Key 提交到 Git 仓库、公开论坛或截图里。如果怀疑泄露，请立即删除旧密钥并重新创建。

### 3. 充值或兑换额度

- 「兑换中心」（`/redeem`）：输入兑换码增加余额。
- 「购买订阅」（`/purchase`）：站点开启在线支付时，可以直接购买套餐。
- 余额、订阅进度和消费明细可以在控制台首页和「用量记录」（`/usage`）中查看。

### 4. 选择模型

打开「模型广场」（`/model-plaza`），确认三件事：

- **模型名**：请求里的 `model` 必须和广场中展示的完全一致；
- **分组**：你的 API Key 属于哪个分组，决定能用哪些模型；
- **价格**：实付价格 = 渠道价 × 分组倍率，或你的专属倍率。

### 5. 发出第一个请求

```bash
export TOKENPILOT_API_KEY="sk-你的密钥"

curl https://tokenpilot.site/v1/chat/completions \
  -H "Authorization: Bearer $TOKENPILOT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "用三句话介绍你自己"}]
  }'
```

如果响应里的 `choices[0].message.content` 有内容，说明接入成功。把 `model` 换成模型广场里的其他模型，就可以调用对应模型。

## 接口地址

TokenPilot 同时兼容三种主流协议，按你使用的 SDK 选择 Base URL：

| 协议 | Base URL | 适用客户端 |
| --- | --- | --- |
| OpenAI Chat Completions / Responses | `https://tokenpilot.site/v1` | OpenAI SDK、Codex CLI、NextChat、Cherry Studio |
| Anthropic Messages | `https://tokenpilot.site` | Anthropic SDK、Claude Code |
| Gemini | `https://tokenpilot.site` | Gemini SDK、Gemini CLI（SDK 会自动追加 `v1beta`） |

### 常用端点

| 端点 | 用途 |
| --- | --- |
| `POST /v1/chat/completions` | OpenAI 兼容对话 |
| `POST /v1/responses` | OpenAI Responses |
| `POST /v1/messages` | Anthropic 消息 |
| `POST /v1/messages/count_tokens` | 预估 token 数 |
| `POST /v1beta/models/{model}:generateContent` | Gemini 生成内容 |
| `GET /v1/models` | 查询当前 Key 可用的模型 |
| `GET /v1/usage` | 查询当前 Key 的用量 |
| `POST /v1/images/generations` | 图片生成（需分组支持） |
| `POST /v1/videos` | 视频生成（需分组支持） |

### 鉴权方式

三种请求头任选其一，推荐第一种：

| 请求头 | 示例 | 说明 |
| --- | --- | --- |
| `Authorization` | `Bearer sk-...` | 所有协议通用，推荐 |
| `x-api-key` | `sk-...` | Anthropic 客户端常用 |
| `x-goog-api-key` | `sk-...` | Gemini 客户端常用 |

Gemini 的 `?key=sk-...` 查询参数仍然可用，但只建议在无法设置请求头的客户端中使用；旧的 `?api_key=` 写法已废弃。

> 每个响应都会带上 `X-Client-Request-ID`。联系管理员排查问题时，请提供这个 ID，可以快速定位到具体请求。

## 调用示例

> 下面的示例统一用 `gpt-4o-mini`、`claude-sonnet-4-5`、`gemini-2.0-flash` 演示。实际调用时，请替换成[模型广场](https://tokenpilot.site/model-plaza)中可用的模型名。

### OpenAI 兼容（Python）

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-你的密钥",
    base_url="https://tokenpilot.site/v1",
)

resp = client.chat.completions.create(
    model="gpt-4o-mini",  # 换成模型广场里的模型名
    messages=[{"role": "user", "content": "你好"}],
)
print(resp.choices[0].message.content)
```

### OpenAI 兼容（Node.js）

```js
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.TOKENPILOT_API_KEY,
  baseURL: "https://tokenpilot.site/v1",
});

const resp = await client.chat.completions.create({
  model: "gpt-4o-mini",
  messages: [{ role: "user", content: "你好" }],
});
console.log(resp.choices[0].message.content);
```

### 流式输出

把 `stream` 设为 `true`，服务端会以 SSE 返回 `data:` 事件，最后以 `data: [DONE]` 结束：

```bash
curl -N https://tokenpilot.site/v1/chat/completions \
  -H "Authorization: Bearer $TOKENPILOT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "写一首短诗"}],
    "stream": true
  }'
```

> 使用 `curl` 时加 `-N` 关闭缓冲，否则可能看起来“没有输出”。

### Anthropic / Claude（Python）

```python
from anthropic import Anthropic

client = Anthropic(
    api_key="sk-你的密钥",
    base_url="https://tokenpilot.site",
)

msg = client.messages.create(
    model="claude-sonnet-4-5",  # 以模型广场为准
    max_tokens=1024,
    messages=[{"role": "user", "content": "你好"}],
)
print(msg.content[0].text)
```

### Claude Code

```bash
export ANTHROPIC_BASE_URL="https://tokenpilot.site"
export ANTHROPIC_AUTH_TOKEN="sk-你的密钥"
export CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1
claude
```

也可以写入 `~/.claude/settings.json`：

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://tokenpilot.site",
    "ANTHROPIC_AUTH_TOKEN": "sk-你的密钥",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1"
  }
}
```

> `ANTHROPIC_BASE_URL` 不要带 `/v1`，Claude Code 会自己追加 `/v1/messages`。如果模型名和默认值不一致，再设置 `ANTHROPIC_MODEL`、`ANTHROPIC_DEFAULT_SONNET_MODEL` 等变量。

### Codex CLI

Codex CLI 使用 OpenAI Responses 协议。推荐通过 `~/.codex/config.toml` 指定 TokenPilot；只设置 `OPENAI_BASE_URL` 时，新版本 Codex 仍可能走 ChatGPT 登录。

#### 1. 安装

```bash
npm install -g @openai/codex
codex --version
```

#### 2. 创建 `~/.codex/config.toml`

macOS / Linux 路径为 `~/.codex/config.toml`，Windows 为 `%USERPROFILE%\.codex\config.toml`：

```toml
model_provider = "tokenpilot"
model = "gpt-5.1-codex"          # 换成模型广场里你分组可用的 Codex 模型
review_model = "gpt-5.1-codex"
disable_response_storage = true

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
env_key = "TOKENPILOT_API_KEY"
wire_api = "responses"
requires_openai_auth = false
supports_websockets = false
```

说明：

- `base_url` 必须带 `/v1`，Codex 会自己追加 `/responses`；
- `env_key` 填环境变量名，不要把 `sk-...` 直接写进 `config.toml`；
- `requires_openai_auth = false` 表示使用 API Key，不要求 ChatGPT OAuth 登录；
- `supports_websockets = false` 使用 HTTP/SSE，避免部分网关不支持 WebSocket；
- `model` 必须是当前 API Key 所属分组可用的模型，可在[模型广场](https://tokenpilot.site/model-plaza)确认。

#### 3. 设置环境变量并启动

macOS / Linux：

```bash
export TOKENPILOT_API_KEY="sk-你的密钥"
codex
```

PowerShell：

```powershell
$env:TOKENPILOT_API_KEY="sk-你的密钥"
codex
```

Windows CMD：

```bat
set TOKENPILOT_API_KEY=sk-你的密钥
codex
```

#### 4. 验证

```bash
codex exec "用一句话介绍当前目录"
```

常见问题：

| 现象 | 检查项 |
| --- | --- |
| 提示 401 / Unauthorized | `TOKENPILOT_API_KEY` 是否已设置，Key 是否被删除或停用 |
| 提示 model not found | `model` 是否在当前 Key 的分组中，名称是否与模型广场完全一致 |
| 仍然弹出 ChatGPT 登录 | `config.toml` 路径是否正确，`requires_openai_auth` 是否为 `false` |
| WebSocket / 连接失败 | 保留 `supports_websockets = false`，并确认 Base URL 带 `/v1` |

> 也可以直接使用控制台「API 密钥」页面的「使用密钥」功能，选择 Codex CLI 后复制生成的 `config.toml` 和环境变量，避免手动拼写。

### Gemini

```bash
curl "https://tokenpilot.site/v1beta/models/gemini-2.0-flash:generateContent" \
  -H "x-goog-api-key: $TOKENPILOT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"你好"}]}]}'
```

Gemini CLI：

```bash
export GOOGLE_GEMINI_BASE_URL="https://tokenpilot.site"
export GEMINI_API_KEY="sk-你的密钥"
export GEMINI_MODEL="gemini-2.0-flash"
gemini
```

### 图形界面客户端

Cherry Studio、NextChat、Chatbox、Open WebUI 等客户端，选择「OpenAI 兼容」或「自定义 API」：

| 配置项 | 填写内容 |
| --- | --- |
| API 地址 / Base URL | `https://tokenpilot.site/v1` |
| API Key | `sk-...` |
| 模型 | 模型广场中展示的完整模型名 |

> 部分客户端会在 Base URL 后自动追加 `/v1`，此时填 `https://tokenpilot.site` 即可。以客户端实际请求路径为准。

## 模型与计费

### 查询可用模型

```bash
curl https://tokenpilot.site/v1/models \
  -H "Authorization: Bearer $TOKENPILOT_API_KEY"
```

返回的是当前 Key 所在分组可用的模型列表（OpenAI 格式）。如果列表为空，通常是 Key 还没有绑定分组，或分组下暂时没有可用渠道。

### 计费规则

- **按 Token**：输入、输出、缓存写入、缓存读取分别计价，单位为「元 / 1M token」。
- **按次 / 按图片 / 按视频**：部分模型按请求次数或媒体张数计费，以模型广场为准。
- **分组倍率**：实付 = 渠道单价 × 分组倍率；如果管理员给你配置了专属倍率，会覆盖分组倍率。
- **长上下文阶梯**：超过阈值后按更高档位计价，具体看模型广场的档位说明。
- **分时 / 高峰倍率**：部分分组在指定时段或高峰时段有额外倍率，模型广场会标注。

### 查询用量

```bash
curl "https://tokenpilot.site/v1/usage?days=7" \
  -H "Authorization: Bearer $TOKENPILOT_API_KEY"
```

返回当前 Key 的今日/累计用量、费用、RPM/TPM，以及按模型统计的用量。`days` 可选 1–90。

### 不登录查询密钥用量

打开「密钥用量」（`/key-usage`），输入 API Key 即可查看余额、配额和近期用量。该页面不会保存你的密钥。

### 余额与订阅

- 余额不足时请求会返回 `403 billing_error`，请先充值或兑换。
- 订阅型分组会优先消耗订阅额度；超出后按站点规则处理。
- 具体余额、订阅进度和账单以控制台为准。

## 错误与排查

### 常见状态码

| 状态码 | 错误类型 | 常见原因 | 处理方式 |
| --- | --- | --- | --- |
| 400 | `invalid_request_error` | JSON 格式错误、缺少 `model`、参数不合法 | 检查请求体，确认 `model` 已填写 |
| 401 | `authentication_error` | 没有带 Key、Key 错误或已删除 | 检查 `Authorization: Bearer sk-...` |
| 403 | `permission_error` | Key 没有该分组/模型的权限 | 在模型广场确认分组，或联系管理员开通 |
| 403 | `billing_error` | 余额不足、订阅额度用完 | 充值/兑换，或购买订阅 |
| 404 | `not_found_error` | 模型名或路径不存在 | 对照模型广场的完整模型名 |
| 429 | `rate_limit_exceeded` | 并发、RPM/TPM 或窗口额度超限 | 按 `Retry-After` 退避后重试 |
| 5xx | `api_error` | 上游模型服务异常 | 稍后重试；持续出现请联系管理员 |

> 不同协议的响应结构略有差异：OpenAI 是 `{"error": {...}}`，Anthropic 是 `{"type": "error", "error": {...}}`，Gemini 是 `{"error": {"code", "message", "status"}}`。排查时优先看 `message` 和 `X-Client-Request-ID`。

### 429 怎么处理

1. 看响应头 `Retry-After`，按它建议的秒数退避，不要立刻重试。
2. 降低并发或请求频率；批量任务加一点随机抖动。
3. 确认不是 Key 的 5 小时 / 1 天 / 7 天窗口额度用完（控制台「API 密钥」可以看到额度）。
4. 如果业务确实需要更高并发，联系管理员调整分组或账号容量。

### 401 / 403

- 401：优先检查请求头。可以直接用下面的命令验证：

```bash
curl https://tokenpilot.site/v1/models \
  -H "Authorization: Bearer sk-你的密钥"
```

- 403：Key 有效，但当前分组没有这个模型，或者余额/订阅额度不足。先看错误 `code` 是 `permission_error` 还是 `billing_error`。

### 模型不存在

模型名必须和模型广场完全一致，包括大小写、日期后缀和分隔符。例如 `gpt-4o-mini` 和 `gpt-4o` 是两个不同的模型。可以先用 `GET /v1/models` 看当前 Key 能调用哪些模型。

### 流式响应很慢或中断

- 首 Token 慢通常是上游排队或网络问题；可以换一个模型或分组再试。
- 流式过程中断时，客户端通常会收到不完整的 SSE。建议加超时和重试，并记录 `X-Client-Request-ID`。
- 可以在「渠道状态」页面查看当前渠道是否可用。

### Claude Code 连不上

按顺序检查：

1. `ANTHROPIC_BASE_URL` 是否写成 `https://tokenpilot.site`（不要带 `/v1`，也不要带 `/v1/messages`）。
2. `ANTHROPIC_AUTH_TOKEN` 是否是完整的 `sk-...`，有没有多余空格或引号。
3. 当前 Key 的分组是否包含 Claude 模型；可以在模型广场确认。
4. 如果公司网络需要代理，先确认终端能访问 `https://tokenpilot.site`。
5. 仍然失败时，用 `curl` 直接请求 `/v1/messages`，把 `X-Client-Request-ID` 和错误信息发给管理员。

## 安全建议

- 一个应用一个 Key，不要把同一个 Key 用在多个项目里。
- 不要把 Key 写进前端代码、Git 仓库或公开截图；用环境变量或密钥管理服务。
- 定期在控制台轮换 Key；怀疑泄露时立即删除旧 Key。
- 如果 Key 支持额度或过期时间，按最小权限设置。

## 联系我们

遇到文档没有覆盖的问题：

1. 先看上面的「错误与排查」和「常见问题」。
2. 登录控制台，在「个人资料」（`/profile`）页面查看管理员配置的联系方式。
3. 反馈时请附上：请求时间、模型名、完整错误信息、`X-Client-Request-ID`（记得隐去 API Key）。

---

本文档适用于 [tokenpilot.site](https://tokenpilot.site)。模型、价格和可用性以控制台与模型广场的实时展示为准。
