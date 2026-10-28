# tokenpolit 配置文档

tokenpolit 是一个统一的 AI 模型 API 网关：你只需要一个 API Key，就可以调用 Claude、GPT、Gemini 等主流模型，无需单独购买或维护多家订阅账号。

本文档将带你完成从注册到接入的全过程。

## 快速开始

### 第一步：注册并登录

1. 打开 tokenpolit 首页，点击右上角「登录」。
2. 首次使用请先注册账号，完成邮箱验证后登录控制台。

### 第二步：创建 API Key

1. 登录后进入控制台。
2. 在左侧菜单打开「API 密钥」（`/keys`）。
3. 点击「创建密钥」，填写名称后确认。
4. 复制生成的密钥（形如 `sk-xxxxxxxx`）并妥善保存。密钥只在创建时完整展示一次。

> 请勿把 API Key 提交到 Git 仓库、公开论坛或截图里。如果发现泄露，请立即在控制台删除并重新创建。

### 第三步：充值或兑换额度

- 「兑换中心」（`/redeem`）：输入兑换码即可增加余额。
- 「购买订阅」（`/purchase`）：站点开启在线支付时，可直接购买套餐。

余额与用量可以在控制台首页和「用量记录」（`/usage`）中查看。

### 第四步：选择模型

打开「模型广场」（`/model-plaza`），可以查看当前可用模型、所属分组以及计费倍率。每个分组对应一组可调用的模型，请留意分组说明。

## 接口地址

所有接口都遵循 OpenAI / Anthropic / Gemini 的兼容协议，基础地址为：

```
https://你的站点域名/v1
```

请在控制台或站点公告中确认实际域名，将下面的示例域名替换为你自己的站点地址。

支持的接口路径：

| 接口路径 | 协议 | 典型用途 |
| --- | --- | --- |
| `/v1/chat/completions` | OpenAI Chat Completions | GPT 系列、通用对话 |
| `/v1/responses` | OpenAI Responses | 新版 OpenAI 调用方式 |
| `/v1/messages` | Anthropic Messages | Claude 系列 |
| `/v1/messages/count_tokens` | Anthropic | 预估 token 数 |
| `/v1beta` | Gemini | Gemini 系列 |

## 接入示例

### OpenAI SDK（Python）

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-你的密钥",
    base_url="https://你的站点域名/v1",
)

resp = client.chat.completions.create(
    model="gpt-4o-mini",
    messages=[{"role": "user", "content": "你好，介绍一下你自己"}],
)
print(resp.choices[0].message.content)
```

### curl 调用

```bash
curl https://你的站点域名/v1/chat/completions \
  -H "Authorization: Bearer sk-你的密钥" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "你好"}]
  }'
```

### Claude Code

在终端中配置环境变量后启动 Claude Code：

```bash
export ANTHROPIC_BASE_URL="https://你的站点域名"
export ANTHROPIC_AUTH_TOKEN="sk-你的密钥"
claude
```

### Codex CLI

```bash
export OPENAI_BASE_URL="https://你的站点域名/v1"
export OPENAI_API_KEY="sk-你的密钥"
codex
```

### Cherry Studio / NextChat 等客户端

在客户端中选择「OpenAI 兼容」或「自定义 API」：

- API 地址（Base URL）：`https://你的站点域名/v1`
- API Key：`sk-你的密钥`
- 模型名称：请填写模型广场中显示的完整模型名

## 模型与计费

- 模型名称必须与模型广场中展示的名称完全一致，大小写敏感。
- 计费按实际用量统计，输入与输出 token 分开计算，具体单价以模型广场为准。
- 不同分组可能有不同的倍率与可用模型，请按需选择。
- 每次请求的消耗可在「用量记录」中查看，按时间与模型筛选即可。

## 密钥使用情况查询

未登录时也可以打开「密钥用量」（`/key-usage`）页面，输入 API Key 即可查询该密钥的余额、配额与近期用量。该页面不会保存你的密钥。

## 常见问题

### 401 Unauthorized

密钥无效或已被删除。请确认请求头格式为：

```
Authorization: Bearer sk-你的密钥
```

然后回到「API 密钥」页面核对密钥是否仍然存在。

### 403 Forbidden

当前密钥没有该模型或分组的访问权限。请在模型广场确认模型所属分组，或联系管理员开通对应权限。

### 429 Too Many Requests

触发了并发或频率限制。请降低请求频率后重试；如果你有较高的并发需求，请联系管理员。

### 余额不足

请求会返回余额不足相关错误。请前往「兑换中心」兑换额度，或在「购买订阅」页面充值。

### 模型不存在

请检查模型名称是否与模型广场完全一致，包括大小写和分隔符。部分模型只对特定分组开放。

### 响应很慢或超时

上游模型服务可能出现波动，可以稍后重试。若长时间不可用，请在「渠道状态」页面查看当前渠道状态。

## 联系客服

遇到文档未覆盖的问题时：

1. 登录控制台，点击右上角头像打开用户菜单。
2. 在菜单中查看「联系客服」一栏，使用站点管理员配置的联系方式。

反馈问题时请附上：请求时间、使用的模型名、完整的错误信息（请先隐去 API Key）。
