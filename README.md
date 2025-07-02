# Retell 📚🎧

> 基于 Telegram Bot 的智能英语学习助手，让英语学习更高效、更有趣！

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Telegram](https://img.shields.io/badge/Telegram-Bot-blue.svg)](https://telegram.org/apps)

## ✨ 项目简介

传统的英语学习往往局限于机械式的背单词，这种方法存在诸多弊端：

❌ **传统背单词的问题**
- 单词零散，缺乏语境联系，遗忘率高
- 只锻炼机械记忆，听说读写能力无法提升
- 学习效率低下，无法形成语言思维

✅ **我们的解决方案：短文背诵**
- 📖 **深度理解**：掌握单词读音、含义和用法
- 🏗️ **语法强化**：分析句子结构，理解语法规则
- ✍️ **综合提升**：通过默写练习，同时提高语法理解和写作能力

Retell 正是为了解决听力训练缺失这一痛点而诞生的项目，通过 AI 语音合成技术，为您的英语学习文章配上标准发音，让学习更加立体化。

## 🚀 核心功能

### 📝 智能文章管理
- **文章导入**：支持添加英语学习文章
- **AI 语音合成**：使用 Azure 语音服务，将文章转换为高质量音频
- **智能摘要**：集成智谱 AI，自动生成文章缩略图和摘要

### 📚 学习历史追踪
- **历史记录**：完整的学习文章历史管理
- **多媒体体验**：支持文字阅读和音频播放
- **便捷管理**：随时删除不需要的文章，保持学习库整洁

### 🤖 Telegram 集成
- **即时互动**：通过 Telegram Bot 随时随地学习
- **用户友好**：简洁的界面设计，操作简单直观

## 📦 快速开始

### 环境要求

- Go 1.19+
- Telegram Bot Token
- Azure 语音服务 API Key
- 智谱 AI API Key

### 安装步骤

1. **克隆项目**

   ```bash
   git clone https://github.com/usual2970/retell.git
   cd retell
   ```

2. **编译项目**

   ```bash
   go build -o retell
   ```

3. **配置环境变量**

   ```bash
   # Telegram 机器人 Token
   export TG_TOKEN="your_telegram_bot_token"
   
   # 智谱 AI API Key（用于生成文章缩略图）
   export ZHIPU_API_KEY="your_zhipu_api_key"
   
   # Azure 语音服务 API Key（用于文本转语音）
   export AZURE_SPEECH_KEY="your_azure_speech_key"
   
   # Azure 语音服务区域（可选，默认为 eastus）
   export AZURE_SPEECH_REGION="eastus"
   ```

4. **启动服务**

   ```bash
   ./retell serve
   ```

### 配置说明

| 环境变量 | 说明 | 是否必需 |
|---------|------|---------|
| `TG_TOKEN` | Telegram Bot Token | ✅ 必需 |
| `ZHIPU_API_KEY` | 智谱 AI API Key | ✅ 必需 |
| `AZURE_SPEECH_KEY` | Azure 语音服务密钥 | ✅ 必需 |
| `AZURE_SPEECH_REGION` | Azure 服务区域 | ❌ 可选（默认：eastus） |

## 📱 使用演示

### 主菜单界面

![主菜单](assets/image-1.png)

*简洁直观的操作界面*

### 文章列表管理

![文章列表](assets/image-2.png)

*历史文章一目了然，支持快速访问*

### 文章详情页面

![文章详情](assets/image-3.png)

*支持文字阅读和音频播放的完整学习体验*

## 🏗️ 技术架构

- **后端框架**：Go + Clean Architecture
- **数据库**：PocketBase（嵌入式数据库）
- **AI 服务**：
  - Azure 语音服务（文本转语音）
  - 智谱 AI（内容摘要生成）
- **即时通讯**：Telegram Bot API
- **部署**：支持 Railway、Docker 等平台

## 🤝 贡献指南

我们欢迎任何形式的贡献！

1. Fork 本项目
2. 创建您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目遵循 [MIT License](LICENSE) 开源协议。您可以自由地使用、复制、修改和分发本软件。

## 🔗 相关链接

- [Telegram Bot API 文档](https://core.telegram.org/bots/api)
- [Azure 语音服务](https://azure.microsoft.com/zh-cn/services/cognitive-services/speech-services/)
- [智谱 AI](https://open.bigmodel.cn/)

---

<div align="center">
  <h3>🌟 开始您的高效英语学习之旅！</h3>
  <p>如有任何问题或建议，欢迎提交 Issue 或 Pull Request</p>
  <p>让我们一起让 Retell 变得更加出色！🚀</p>
</div>
