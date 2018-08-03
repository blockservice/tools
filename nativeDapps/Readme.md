# 赢fomo3d大奖的工具

## 使用说明
``` bash
➜  tools go  run main.go -h
Usage of /tmp/go-build286277460/b001/exe/main:
  -fomoAddr string
    	指定fomo3d的合约地址 (default "0xa62142888aba8370742be823c1782d17a0389da1")
  -node string
    	指定以太节点地址 (default "http://127.0.0.1:18545")
  -private string
    	指定私钥 (default "0173e7a434193e8d49c477a6d6d28b6da4c09048aad580bef43956f4d975fb87")
  -timeLeft int
    	剩余多少时间就触发买key操作 (default 28)
```

## 使用例子

指定地址到倒计时小于20秒时，进行买1个key操作。

```bash
go run main.go -private 0173e7a434193e8d49c477a6d6d28b6da4c09048aad580bef43956f4d975fb87 -timeLeft 20
```