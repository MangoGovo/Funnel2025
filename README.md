# Funnel2025
Funnel 于 2025 年的重构版本 

### Lint(代码格式检查)
#### 手动格式化+lint检测
```shell
gofmt -w .                                # gofmt的代码格式化
gci write . -s standard -s default        # 根据lint文件格式化
golangci-lint run --config .golangci.yml  # 运行lint
```
#### 集成到IDE中
[配置方法]("https://golangci-lint.run/welcome/integrations/")