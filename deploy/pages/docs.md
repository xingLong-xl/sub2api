## 基础信息

| 客户端 | Base URL | 协议 |
| --- | --- | --- |
| Claude Code | `https://tokenpilot.site` | Anthropic Messages |
| Codex CLI | `https://tokenpilot.site/v1` | OpenAI Responses |

注意：

- Claude Code 的 Base URL **不要带 `/v1`**，客户端会自己追加 `/v1/messages`。
- Codex 的 Base URL **要带 `/v1`**，客户端会自己追加 `/responses`。
- 同一把 Key 能否同时用于 Claude Code 和 Codex，取决于 Key 所属分组。Anthropic 分组用于 Claude 模型，OpenAI 分组用于 GPT/Codex 模型；不确定时请在模型广场确认，或分别创建对应分组的 Key。

## 开始之前

1. 登录 TokenPilot 控制台。
2. 进入「API 密钥」（`/keys`），创建一个新的 Key。
3. 在下面的环境变量或配置文件中，把 `sk-你的密钥` 替换成真实 Key。
4. 修改配置后，重新打开终端，或完全退出并重启 Claude Code / Codex。

> 不要把 API Key 提交到 Git 仓库或公开分享。如果怀疑泄露，请立即删除旧 Key 并重新创建。

## Claude Code

适合使用 Claude 系列模型进行项目开发。配置后，在项目目录执行 `claude` 即可进入交互式编码环境。

环境变量适合当前终端快速验证；`settings.json` 适合长期使用。修改后请重启 Claude Code。

#### Windows

##### 安装 Claude Code

```powershell
winget install OpenJS.NodeJS.LTS
```

安装 Node.js 后重新打开 PowerShell，再执行：

```powershell
npm install -g @anthropic-ai/claude-code
claude --version
```

##### 写入当前终端环境

```powershell
$env:ANTHROPIC_AUTH_TOKEN="sk-你的密钥"
$env:ANTHROPIC_BASE_URL="https://tokenpilot.site"
```

##### 可选：写入持久配置

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.claude" | Out-Null
@'
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "sk-你的密钥",
    "ANTHROPIC_BASE_URL": "https://tokenpilot.site"
  }
}
'@ | Set-Content "$env:USERPROFILE\.claude\settings.json"
```

##### 启动或重启 Claude Code

```powershell
claude
```

#### Linux

##### 安装 Claude Code

```bash
# Ubuntu / Debian：先安装 Node.js 20/22
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt-get install -y nodejs

npm install -g @anthropic-ai/claude-code
claude --version
```

其他发行版请先安装 Node.js 20+，再执行 `npm install -g @anthropic-ai/claude-code`。

##### 写入当前终端环境

```bash
export ANTHROPIC_AUTH_TOKEN="sk-你的密钥"
export ANTHROPIC_BASE_URL="https://tokenpilot.site"
```

##### 可选：写入持久配置

```bash
mkdir -p ~/.claude
cat > ~/.claude/settings.json <<'EOF'
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "sk-你的密钥",
    "ANTHROPIC_BASE_URL": "https://tokenpilot.site"
  }
}
EOF

chmod 600 ~/.claude/settings.json
```

##### 启动或重启 Claude Code

```bash
claude
```

#### macOS

##### 安装 Claude Code

```bash
# 已安装 Node.js 18+ 可跳过
brew install node
node -v

npm install -g @anthropic-ai/claude-code
claude --version
```

没有 Homebrew 时，可以从 [nodejs.org](https://nodejs.org/) 下载 LTS 安装包。

##### 写入当前终端环境

```bash
export ANTHROPIC_AUTH_TOKEN="sk-你的密钥"
export ANTHROPIC_BASE_URL="https://tokenpilot.site"
```

##### 可选：写入持久配置

```bash
mkdir -p ~/.claude
cat > ~/.claude/settings.json <<'EOF'
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "sk-你的密钥",
    "ANTHROPIC_BASE_URL": "https://tokenpilot.site"
  }
}
EOF

chmod 600 ~/.claude/settings.json
```

##### 启动或重启 Claude Code

```bash
claude
```

## Codex

适合使用 GPT 系列模型进行代码生成、修改和审查。Codex 使用 OpenAI Responses 协议，因此：

- `base_url` 填 `https://tokenpilot.site/v1`；
- 在 `~/.codex/config.toml` 中定义 `model_provider`；
- API Key 通过 `OPENAI_API_KEY` 环境变量读取；
- 下面的 `model = "gpt-5.5"` 只是示例，必须换成模型广场里当前分组可用的模型。

修改配置后请重启 Codex。

#### Windows

##### 安装 Codex CLI

```powershell
winget install OpenJS.NodeJS.LTS
```

安装 Node.js 后重新打开 PowerShell，再执行：

```powershell
npm install -g @openai/codex
codex --version
```

##### 写入 config.toml

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.codex" | Out-Null
@'
model_provider = "tokenpilot"
model = "gpt-5.5"
review_model = "gpt-5.5"
model_reasoning_effort = "high"
disable_response_storage = true

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
env_key = "OPENAI_API_KEY"
wire_api = "responses"
requires_openai_auth = false
supports_websockets = false
'@ | Set-Content "$env:USERPROFILE\.codex\config.toml"
```

##### 写入 API Key

```powershell
setx OPENAI_API_KEY "sk-你的密钥"
```

`setx` 不会影响已经打开的窗口，执行后请重新打开 PowerShell。

##### 启动或重启 Codex

```powershell
codex
```

#### Linux

##### 安装 Codex CLI

```bash
# Ubuntu / Debian：先安装 Node.js 20/22
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt-get install -y nodejs

npm install -g @openai/codex
codex --version
```

其他发行版请先安装 Node.js 20+，再执行 `npm install -g @openai/codex`。

##### 写入 config.toml

```bash
mkdir -p ~/.codex
cat > ~/.codex/config.toml <<'TOML'
model_provider = "tokenpilot"
model = "gpt-5.5"
review_model = "gpt-5.5"
model_reasoning_effort = "high"
disable_response_storage = true

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
env_key = "OPENAI_API_KEY"
wire_api = "responses"
requires_openai_auth = false
supports_websockets = false
TOML
```

##### 写入 API Key

```bash
export OPENAI_API_KEY="sk-你的密钥"
```

##### 启动或重启 Codex

```bash
codex
```

#### macOS

##### 安装 Codex CLI

```bash
# 已安装 Node.js 18+ 可跳过
brew install node
node -v

npm install -g @openai/codex
codex --version
```

##### 写入 config.toml

```bash
mkdir -p ~/.codex
cat > ~/.codex/config.toml <<'TOML'
model_provider = "tokenpilot"
model = "gpt-5.5"
review_model = "gpt-5.5"
model_reasoning_effort = "high"
disable_response_storage = true

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
env_key = "OPENAI_API_KEY"
wire_api = "responses"
requires_openai_auth = false
supports_websockets = false
TOML
```

##### 写入 API Key

```bash
export OPENAI_API_KEY="sk-你的密钥"
```

##### 启动或重启 Codex

```bash
codex
```

## 配置检查

- Claude Code 的 `ANTHROPIC_BASE_URL` 是 `https://tokenpilot.site`，不要带 `/v1`。
- Codex 的 `base_url` 是 `https://tokenpilot.site/v1`，要带 `/v1`。
- Claude Code 使用 `ANTHROPIC_AUTH_TOKEN`，Codex 使用 `OPENAI_API_KEY`。
- 修改环境变量或配置文件后，需要重新打开终端或重启对应 CLI。
- 如果提示模型不存在，检查 Key 分组和模型名是否与模型广场一致。
- 如果 Codex 仍然要求 ChatGPT 登录，检查 `~/.codex/config.toml` 路径和 `requires_openai_auth = false`。

## 模型与计费

### 查询模型和用量

```bash
# 当前 Key 可用的模型
curl https://tokenpilot.site/v1/models \
  -H "Authorization: Bearer $TOKENPILOT_API_KEY"

# 最近 7 天用量
curl "https://tokenpilot.site/v1/usage?days=7" \
  -H "Authorization: Bearer $TOKENPILOT_API_KEY"
```

- `/v1/models` 返回当前 Key 分组可用的模型；列表为空通常是 Key 未绑定分组。
- `/v1/usage` 返回用量、费用和 RPM/TPM，`days` 可选 1–90。
- 也可以打开「密钥用量」（`/key-usage`）查询余额和配额，页面不会保存 Key。

### 计费规则

- 按 Token 计费，输入、输出、缓存分别计价。
- 图片、视频等按次或按张计费，以模型广场为准。
- 实付 = 渠道单价 × 分组倍率；专属倍率优先。
- 长上下文、分时/高峰可能加价；余额不足返回 `403 billing_error`。
- 订阅分组优先消耗订阅额度。

## 错误与排查

### 常见错误

| 状态码 | 原因 | 处理 |
| --- | --- | --- |
| 401 | Key 缺失、错误或已删除 | 检查 `Authorization: Bearer sk-...` |
| 403 | 没有模型权限，或余额/订阅不足 | 检查 Key 分组，或充值 |
| 404 | 模型名或路径不存在 | 对照模型广场的完整模型名 |
| 429 | 并发、RPM/TPM 或额度超限 | 按 `Retry-After` 退避，降低并发 |
| 5xx | 上游模型服务异常 | 稍后重试，持续出现联系管理员 |

> 排查时优先看错误 `message` 和响应头 `X-Client-Request-ID`。

### Claude Code 连不上

1. `ANTHROPIC_BASE_URL` 应为 `https://tokenpilot.site`，不要带 `/v1` 或 `/v1/messages`。
2. `ANTHROPIC_AUTH_TOKEN` 应是完整的 `sk-...`，没有多余空格或引号。
3. 当前 Key 分组要包含 Claude 模型。
4. 修改配置后重启终端和 Claude Code。
5. 仍然失败时，用 `curl` 请求 `/v1/messages`，把 `X-Client-Request-ID` 发给管理员。

### Codex 仍然要求登录

1. 确认配置在 `~/.codex/config.toml`（Windows 为 `%USERPROFILE%\.codex\config.toml`）。
2. 确认 `requires_openai_auth = false`。
3. 设置 `OPENAI_API_KEY` 后重新打开终端或重启 Codex；不要继续执行 `codex login`。

## 安全建议

- 一个应用一个 Key，不要多个项目共用。
- 不要把 Key 写进前端代码、Git 仓库或公开截图，优先用环境变量。
- 定期轮换 Key；怀疑泄露时立即删除旧 Key。

## 联系我们

遇到文档没有覆盖的问题：

1. 先看上面的「错误与排查」。
2. 登录控制台，在「个人资料」（`/profile`）查看管理员配置的联系方式。
3. 反馈时附上请求时间、模型名、错误信息和 `X-Client-Request-ID`，不要发送完整 API Key。

---

本文档适用于 [tokenpilot.site](https://tokenpilot.site)。模型、价格和可用性以控制台与模型广场的实时展示为准。
