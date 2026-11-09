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

1. 登录 TokenPilot [控制台](/dashboard)。
2. 进入 [API 密钥](/keys)，创建一个新的 Key。
3. 在下面的环境变量或配置文件中，把 `sk-你的密钥` 替换成真实 Key。
4. 修改配置后，重新打开终端，或完全退出并重启 Claude Code / Codex。

> 不要把 API Key 提交到 Git 仓库或公开分享。如果怀疑泄露，请立即删除旧 Key 并重新创建。

## Claude Code

适合使用 Claude 系列模型进行项目开发。配置后，在项目目录执行 `claude` 即可进入交互式编码环境。

环境变量适合当前终端快速验证；`settings.json` 适合长期使用。修改后请重启 Claude Code。

#### Windows

##### 打开 Windows PowerShell

1. 点击开始菜单
2. 搜索 PowerShell
3. 点击打开 Windows PowerShell 或 Terminal

##### 安装 Node.js LTS

```powershell
winget install OpenJS.NodeJS.LTS
```

##### 安装 Claude Code

```powershell
npm install -g @anthropic-ai/claude-code
```

##### 写入用户环境变量

```powershell
setx ANTHROPIC_AUTH_TOKEN "sk-你的密钥"
setx ANTHROPIC_BASE_URL "https://tokenpilot.site"
```

##### 可选：写入当前 PowerShell 会话

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
npm install -g @anthropic-ai/claude-code
```

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
```

##### 启动或重启 Claude Code

```bash
claude
```

#### macOS

##### 安装 Claude Code

```bash
npm install -g @anthropic-ai/claude-code
```

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
```

##### 启动或重启 Claude Code

```bash
claude
```

## Codex

适合使用 GPT 系列模型进行代码生成、修改和审查。配置后，在项目目录执行 `codex` 即可开始使用。

Codex 以 `config.toml` 定义 provider，API Key 继续通过 `OPENAI_API_KEY` 环境变量读取。修改后请重启 Codex。

#### Windows

##### 打开 Windows PowerShell

1. 点击开始菜单
2. 搜索 PowerShell
3. 点击打开 Windows PowerShell 或 Terminal

##### 安装 Node.js LTS

```powershell
winget install OpenJS.NodeJS.LTS
```

##### 安装 Codex CLI

```powershell
npm install -g @openai/codex
```

##### 推荐：写入持久配置

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.codex" | Out-Null
@'
model = "gpt-5.5"
model_provider = "tokenpilot"
model_reasoning_effort = "xhigh"

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
wire_api = "responses"
env_key = "OPENAI_API_KEY"
requires_openai_auth = false
supports_websockets = false
'@ | Set-Content "$env:USERPROFILE\.codex\config.toml"
```

##### 推荐：持久写入 API Key

```powershell
setx OPENAI_API_KEY "sk-你的密钥"
```

##### 可选：写入当前 PowerShell 会话

```powershell
$env:OPENAI_API_KEY="sk-你的密钥"
```

##### 启动或重启 Codex

```powershell
codex
```

#### Linux

##### 安装 Codex CLI

```bash
npm install -g @openai/codex
```

##### 推荐：写入持久配置

```bash
mkdir -p ~/.codex
cat > ~/.codex/config.toml <<'EOF'
model = "gpt-5.5"
model_provider = "tokenpilot"
model_reasoning_effort = "xhigh"

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
wire_api = "responses"
env_key = "OPENAI_API_KEY"
requires_openai_auth = false
supports_websockets = false
EOF
```

##### 推荐：持久写入 API Key

```bash
echo 'export OPENAI_API_KEY="sk-你的密钥"' >> ~/.bashrc
source ~/.bashrc
```

##### 可选：写入当前终端环境

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
npm install -g @openai/codex
```

##### 推荐：写入持久配置

```bash
mkdir -p ~/.codex
cat > ~/.codex/config.toml <<'EOF'
model = "gpt-5.5"
model_provider = "tokenpilot"
model_reasoning_effort = "xhigh"

[model_providers.tokenpilot]
name = "TokenPilot"
base_url = "https://tokenpilot.site/v1"
wire_api = "responses"
env_key = "OPENAI_API_KEY"
requires_openai_auth = false
supports_websockets = false
EOF
```

##### 推荐：持久写入 API Key

```bash
echo 'export OPENAI_API_KEY="sk-你的密钥"' >> ~/.zshrc
source ~/.zshrc
```

##### 可选：写入当前终端环境

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

