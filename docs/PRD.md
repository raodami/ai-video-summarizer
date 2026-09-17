# AI Video Summarizer

## 项目概述
YouTube视频自动摘要工具 - 输入视频链接，AI自动生成摘要、关键点、时间戳索引

## 核心功能
- YouTube视频转录文本提取
- AI摘要生成（DeepSeek API）
- 关键点提取和时间戳索引
- 多种格式导出（Markdown/JSON/PDF）
- 历史记录管理
- 批量处理支持

## 技术栈
- 后端：Go + Gin + SQLite
- 前端：Next.js 14 + Tailwind CSS
- AI：DeepSeek API
- 数据源：YouTube Transcript API

## API端点
- POST /api/transcript - 提取视频转录
- POST /api/summarize - 生成摘要
- GET /api/history - 历史记录
- GET /api/export/:id - 导出摘要
- GET /api/stats - 统计信息

## UI设计
- Stripe风格深色主题 (#533afd)
- 响应式布局
- 实时进度显示
- 代码高亮显示转录文本

## 环境变量
- DEEPSEEK_API_KEY
- PORT
