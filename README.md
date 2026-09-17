# AI Video Summarizer

## 功能特性
- YouTube视频转录提取
- AI摘要生成（DeepSeek API）
- 关键点提取和时间戳索引
- 历史记录管理
- 响应式设计

## 技术栈
- 后端：Go + Gin + SQLite
- 前端：Next.js 14 + Tailwind CSS
- AI：DeepSeek API

## 安装运行

### 后端
```bash
cd D:/ai-video-summarizer
go mod tidy
go run cmd/server/main.go
```

### 前端
```bash
cd D:/ai-video-summarizer/web
npm install
npm run dev
```

### 环境变量
创建 `.env` 文件：
```
DEEPSEEK_API_KEY=your_api_key_here
PORT=8080
```

## API端点
- POST /api/transcript - 提取视频转录
- POST /api/summarize - 生成摘要
- GET /api/history - 历史记录
- GET /api/stats - 统计信息
