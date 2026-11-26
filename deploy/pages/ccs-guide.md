<div class="guide-hero">
  <h1>零基础接入教程</h1>
  <p>用 TokenPilot 的 API 密钥配置 Codex 或 DeepSeek Harness，发送第一条消息。</p>
</div>

## 先选客户端

| 你的网络情况 | 推荐客户端 | 接入方式 |
|---|---|---|
| 可以访问外网，能打开 OpenAI 官网并下载 Codex | **Codex** | 在 TokenPilot 点击「导入到 CCS」，由 CC-Switch 自动填入供应商配置 |
| 无法访问外网或打不开 OpenAI 官网 | **DeepSeek Harness（DSH）** | 在 CC-Switch 的 DeepSeek Harness 页签手动添加供应商，或直接在 DSH 页面添加自定义模型 API。确保电脑能访问 DSH 下载页和 TokenPilot 服务 |

任选一条路线即可。已有 OpenAI 分组的密钥（例如绑定 `codex0.4` 的密钥），可以用同一把密钥接入 DSH。

## 准备 TokenPilot API 密钥

1. 登录 TokenPilot，打开 [API 密钥](/keys)。已经有绑定目标分组的密钥，可以直接使用。
2. 需要新建时，点击「创建密钥」，选择 **OpenAI** 平台下你要用的分组，例如 `codex0.4`（以账号显示的分组为准）。
3. 手动配置 DSH 时，在密钥列表点击该密钥旁边的复制图标，备用。
4. 在 [模型广场](/model-plaza) 核对该分组实际支持的模型。配置时填写模型的 **ID**。`codex0.4` 是分组名称。

## Codex · 通过 CC-Switch 自动导入

### 1. 安装 Codex 和 CC-Switch

按官方页面提示安装 Codex 桌面客户端，也可以使用已有的 Codex CLI。[Codex 官网](https://openai.com/zh-Hans-CN/codex/)

根据自己的系统下载并安装 CC-Switch。安装完成后先启动一次，让系统注册 `ccswitch://` 协议。[CC-Switch 下载页](https://ccswitch.io/zh/download)

| 系统 | 安装包 |
|---|---|
| Windows | `.msi` |
| macOS | `.dmg` |
| Linux | 按发行版选 `.deb`、`.AppImage` 等可用版本 |

### 2. 从 TokenPilot 导入供应商

1. 在 TokenPilot 的 [API 密钥](/keys) 列表中，找到绑定目标 OpenAI 分组的密钥。
2. 点击 **「导入到 CCS」**，在浏览器弹窗中允许打开 CC-Switch，并在 CC-Switch 中确认导入。
3. 在 CC-Switch 的 **Codex** 页签里找到刚导入的供应商。导入链接根据密钥所属分组的平台指定客户端，OpenAI 分组固定导入到 Codex。

供应商名称、API 地址、API Key 和用量查询脚本会随链接传入。

### 3. 选择模型并发送第一条消息

1. 在 CC-Switch 里编辑刚导入的供应商，点击「获取模型列表」。
2. 从返回的列表中选一个模型。预填的模型可能与分组不一致，以你实际可用的模型为准。
3. 保存并启用该供应商。若 Codex 已打开，请完全退出后重新启动。
4. 打开 Codex 新会话，选中刚配置的模型，发送「你好，请回复接入成功」。使用 CLI 时，可以输入 `/model` 查看并切换模型。

### 余额和用量查询（可选）

导入时会带上用量查询脚本。若余额查询失败，在 CC-Switch 中打开该供应商的用量查询配置，选择 **Sub2API**，核对 API 地址和密钥后保存。

### 在 CC-Switch 中手动添加 Codex 供应商（可选）

在 Codex 页签新建供应商，填写以下字段。如果已经自动导入，可直接编辑已有供应商。

| 字段 | 填写内容 |
|---|---|
| 名称 | `TokenPilot-Codex`，也可以自定义 |
| 备注 / 官网链接 | 可选。官网链接可填 `https://tokenpilot.site` |
| API Key | TokenPilot 创建的密钥，CC-Switch 会据此生成认证配置 |
| API 请求地址 | `https://tokenpilot.site/v1`（手动填写 Codex 的 OpenAI Responses Base URL 时使用） |
| 管理与测速 | 使用同一个 API 地址 |
| 模型列表 / 默认模型 | 点击「获取模型列表」，从该分组可用的模型中选默认模型 |
| 子代理模型 / 默认强度 / 并发上限 | 初次接入可以留空，使用客户端默认设置 |
| 自动审批模型 | 仅在使用 Codex 的自动审批功能时按需设置，选择当前模型列表里的模型 |

## DeepSeek Harness（DSH）· 手动接入

### 下载并打开 DSH

打开官方页面，点击「Download Desktop」，选择适合你电脑的安装包，安装后启动 DeepSeek Harness。[DeepSeek Harness 官网](https://www.deepseek.com/en/harness/)

DSH 有两种手动接入方式，**任选一种即可**。推荐方式一，通过 CC-Switch 手动导入。以下用 `codex0.4` 这类 OpenAI 分组举例。

### 方式一 · 在 CC-Switch 中添加 DSH 供应商

如果你已经使用 CC-Switch，可以在里面统一管理 DSH 的供应商。需要安装时，可前往 [CC-Switch 下载页](https://ccswitch.io/zh/download)。

1. 打开 CC-Switch，选择 **DeepSeek Harness** 选项。
2. 点击 **添加供应商**，选择 **自定义配置 · 手动填写所有字段**。
3. 按下表填写配置。

| CC-Switch 字段 | 填写内容 |
|---|---|
| 供应商标识 | 可自定义，例如 `tokenpilot-codex04`。必须唯一，只能使用小写字母、数字和连字符 |
| 名称 | 可自定义，例如 `TokenPilot · codex0.4`。选一个方便在 DSH 中识别的名称 |
| 备注 | 可选，例如「日常编程」 |
| 官网链接 | 可选，可填 `https://tokenpilot.site` |
| 接口格式 | 选择 **OpenAI Responses**，与当前 OpenAI 分组匹配 |
| API Key | 粘贴 TokenPilot 中绑定目标分组的密钥，下方配置会自动填充 |
| Base URL | `https://tokenpilot.site/v1` |
| 请求头 | 按默认设置留空即可 |
| 额外选项 | 按默认设置留空即可 |
| 模型配置 | 点击 **获取模型列表**，选择该分组支持的模型。获取失败时可手动填写模型 ID，在 [模型广场](/model-plaza) 核对 |

至少配置一个模型，然后保存供应商，完成 DSH 配置写入。接着打开 DSH 客户端页面，在模型选择器中选择刚添加的供应商和模型。已经打开 DSH 的用户，可以刷新页面后再选择。

### 方式二 · 直接在 DSH 页面添加供应商

打开 DSH 的 **设置 → 模型 → 添加模型提供商**，切换到 **自定义模型 API**。

| DSH 字段 | 填写内容 |
|---|---|
| Provider ID | `tokenpilot`（小写，保存后固定。添加多个分组时，每个配置使用不同 ID，例如 `tokenpilot-codex04`） |
| 显示名称 | `TokenPilot · codex0.4`，也可以自定义 |
| API 地址（Base URL） | `https://tokenpilot.site/v1` |
| API 协议 | **OpenAI Responses** |
| API 密钥 | 粘贴 TokenPilot 中绑定目标分组的密钥 |

1. 在 **模型目录**里点「获取可用模型」，勾选要用的模型，再点「更新列表」。如果获取失败或没有返回，手动添加该分组支持的模型 ID，可在 [模型广场](/model-plaza) 核对。
2. 至少添加一个模型，然后保存提供商。
3. 回到会话页面，在模型选择器中选择刚配置的供应商和模型。

两种方式的 Base URL 都填到 `/v1`，DSH 会追加 `/responses`。

### 发送第一条消息

新建会话，在你刚添加的供应商（例如 `TokenPilot · codex0.4`）下选择一个可用模型，发送「你好，请回复接入成功」。

## 常见问题

### 点击 TokenPilot「导入到 CCS」后，如何在 DSH 使用同一把密钥？

TokenPilot 当前会把 OpenAI 分组的自动导入配置发送到 CC-Switch 的 Codex 页签。要接入 DSH，可以在 CC-Switch 中选择 **DeepSeek Harness → 添加供应商 → 自定义配置**，使用同一把密钥填写上面的字段并保存。也可以直接在 DSH 页面添加自定义模型 API。

### 在 CC-Switch 的 DeepSeek Harness 中保存后，DSH 里找不到模型？

先确认该供应商的模型配置里已经添加模型并保存，再刷新 DSH 页面，打开模型选择器查看。

### 点击「导入到 CCS」没有反应？

确认 CC-Switch 已安装并启动过，检查浏览器是否拦截了打开应用的提示。用便携版时，如果系统没有注册 `ccswitch://`，请换用安装版再试。

### 获取模型失败，或者测试消息发不出去？

1. 确认客户端能访问 TokenPilot，API 地址和协议与教程一致。若报 `404`，检查是否重复填写了 `/v1`，或在 Base URL 末尾填了 `/responses`。
2. 确认密钥已启用、绑定了正确的分组，并且账号有可用额度。
3. 核对模型 ID 是否属于该分组。模型列表获取失败时，可以在客户端手动填写模型 ID，再发送测试消息。
4. 使用 Codex 时，在 CC-Switch 中确认供应商已启用，完全退出 Codex 后再重新打开。

### 如何接入其他客户端？

从 API 密钥列表触发导入时，TokenPilot 会按密钥所属的平台选择目标。Anthropic 分组通常导入 Claude Code，Gemini 分组导入 Gemini，Grok 分组导入 Grok。需要手动配置 Claude Code 或 Codex CLI 时，可查阅 [配置文档](/docs)。

| 客户端 | 常用手动配置地址 |
|---|---|
| Codex / DSH（OpenAI Responses） | `https://tokenpilot.site/v1` |
| Claude Code（Anthropic Messages） | `https://tokenpilot.site` |
| Gemini CLI | `https://tokenpilot.site` |
| Grok | `https://tokenpilot.site/v1` |

请妥善保管 API Key，避免发到聊天、截图或公开仓库。如果密钥泄露，请在 [API 密钥](/keys) 页面停用旧密钥并重新创建。
