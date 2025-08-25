# Scratchor Crawler

一个用于爬取 Scratchor 题库的工具。

## 使用方法

### 爬虫命令 (crawler)

```bash
# 使用默认限流器并发数 (2)
go run cmd/crawler/main.go

# 自定义限流器并发数
go run cmd/crawler/main.go --limiter 5
```

### 答案更新命令 (answer_updater)

```bash
# 使用默认限流器并发数 (2)
go run cmd/answer_updater/main.go

# 自定义限流器并发数
go run cmd/answer_updater/main.go --limiter 3
```

## 参数说明

- `--limiter`: 限流器并发数，控制同时进行的HTTP请求数量，默认为2
  - 较小的值可以减少对服务器的压力，但爬取速度较慢
  - 较大的值可以加快爬取速度，但可能增加被服务器限制的风险

## 注意事项

- 请合理设置限流器并发数，避免对目标服务器造成过大压力
- 建议从较小的值开始测试，根据实际情况调整
