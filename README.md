# Go 写文件工具

### 用途

创建内容为随机字符的指定大小的文件。该文件大小为 MB 的整数倍。可选单位为 M 或 G。

### 下载编译

- go get github.com/csxuejin/gowrite
- cd $GOPATH/src/github.com/gowrite
- 如果要编译成 mac 版本的可执行文件，就运行 `make build`，如果要运行成 linux 环境下的可执行文件则运行 `make build-linux`。运行后生成 gowrite 可执行文件。

### 使用示例

- `./gowrite 100m`  在当前目录下创建大小为 100MB 的文件，其中 m 大小写皆可。
- `./gowrite 1g` 在当前目录下创建大小为 1GB 的文件，其中 g 大小写皆可。
- `./gowrite 100m hellofile` 在当前目录下创建大小为 100MB 且名称为 hellofile 的文件。
- `./gowrite 100m hellofile /tmp` 在 /tmp 目录下创建大小为 100MB 且名称为 hellofile 的文件。