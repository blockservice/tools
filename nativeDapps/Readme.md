# 赢fomo3d大奖的工具

## 使用说明
``` bash
➜  nativeDapps git:(master) ✗ go run main.go -h
Usage of /tmp/go-build777560664/b001/exe/main:
  -fomoAddr string
    	指定fomo3d的合约地址 (default "0xa62142888aba8370742be823c1782d17a0389da1")
  -node string
    	指定以太节点RPC地址 (default "https://mainnet.infura.io/YKZGQG2QTBx0tiWoB2IF")
  -private string
    	指定私钥 (default "0173e7a434193e8d49c477a6d6d28b6da4c09048aad580bef43956f4d975fb87")
  -timeLeft int
    	剩余多少时间就触发买key操作 (default 30)
exit status 2
```

## 使用例子

指定买key地址的私钥地址，当倒计时小于20秒时，进行买1个key操作。

```bash
go run main.go -private 0173e7a434193e8d49c477a6d6d28b6da4c09048aad580bef43956f4d975fb87 -timeLeft 20
```

## TODO

 感谢贡献者参与完成下面的功能，实现后，可做到整个过程机器人只须1个key就可以拿走大奖的（如果所有机器人里，只有我们一个这样的话。。。）

- 设置默认买key消耗上限，默认当买key的ETH（包含交易手续费）小于大奖金额时，就买起来。。。
- 默认开启下一轮游戏，直到满足上面的条件
- 判断是否有pending交易，避免不必要的消耗
- 添加node池，多节点广播，加速打包
- 判断买key地址的余额，避免空伤悲