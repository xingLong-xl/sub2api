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
- API Key 通过 `OPENAI_API_KEY` 环境变量读取。

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
