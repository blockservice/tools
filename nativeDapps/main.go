package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ChungkueiBlock/tools/nativeDapps/fomo3d"
	"github.com/ChungkueiBlock/tools/nativeDapps/gas"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	pubStr   string
	affcode  [32]byte
	fomoAddr string
	timeConf int64
	team     *big.Int
	node     string
	private  string
)

func init() {

	flag.Int64Var(&timeConf, "timeLeft", 28, "剩余多少时间就触发买key操作")
	flag.StringVar(&fomoAddr, "fomoAddr", "0xa62142888aba8370742be823c1782d17a0389da1", "指定fomo3d的合约地址")
	flag.StringVar(&node, "node", "https://mainnet.infura.io/YKZGQG2QTBx0tiWoB2IF", "指定以太节点RPC地址")
	flag.StringVar(&private, "private", "0173e7a434193e8d49c477a6d6d28b6da4c09048aad580bef43956f4d975fb87", "指定私钥")
	flag.Parse()

	if pri := os.Getenv("FOMO_PRI"); pri != "" {
		private = pri
	}

	key, _ := crypto.HexToECDSA(private)
	pubStr = crypto.PubkeyToAddress(key.PublicKey).String()
	log.Printf("当前正在使用地址：%v 进行买KEY操作\n", pubStr)

	// team 默认购买熊队
	team = big.NewInt(1)
	// 开发者劳动收益， 感谢买key用开发者的返佣链接
	copy(affcode[:], "xxp")
}

func main() {

	conn, err := ethclient.Dial(node)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fomo, err := fomo3d.NewFoMo3Dlong(common.HexToAddress(fomoAddr), conn)
	if err != nil {
		log.Fatalf("实例化fomo错误：%v\n", err)
	}

	timeLeft, _ := getRoundInfo(fomo)
	timer := time.NewTimer(time.Duration((timeLeft - timeConf)) * time.Second)
	defer timer.Stop()
	fmt.Printf("等待\"%v秒\"后，再检查倒计时\n", timeLeft-timeConf)

	for {
		select {
		case <-timer.C:
			timeLeft, addr := getRoundInfo(fomo)
			if timeLeft-timeConf <= 0 && addr != pubStr {
				// TODO 做风控，允许设置买key上限
				buyXname(fomo)
				// check Pending tx packed into block
				log.Println("等待14秒。。。打包后，继续检查倒计时")
				timer.Reset(14 * time.Second)
			} else {
				log.Printf("%v秒后，再检查倒计时\n", timeLeft-timeConf)
				timer.Reset(time.Duration((timeLeft - timeConf)) * time.Second)
			}

		}
	}
}

func getRoundInfo(fomo *fomo3d.FoMo3Dlong) (t int64, addr string) {
	_, _, _, timeEnd, _, _, _, plyAddr, _, _, _, _, _, _, _ := fomo.GetCurrentRoundInfo(nil)
	timeNow := time.Now().Unix()
	timeLeft := timeEnd.Int64() - timeNow
	fmt.Printf("游戏距离结束时间剩余 %v秒, 目前最后一次买key的地址是 %v\n", timeLeft, plyAddr.String())
	return timeLeft, plyAddr.String()
}

func buyXname(fomo *fomo3d.FoMo3Dlong) {
	vaule, _ := fomo.GetBuyPrice(nil)
	// 字符格式的私钥转化为结构体
	key, _ := crypto.HexToECDSA(private)

	// 绑定私钥，填写tx选项
	opts := bind.NewKeyedTransactor(key)

	// 比正常的Price价格高1倍，加速我们的交易被打包
	gasPrice := gas.Price()
	gasPrice.Mul(gasPrice, big.NewInt(2))
	opts.GasPrice = gasPrice
	opts.Value = vaule
	opts.GasLimit = 500000

	tx, err := fomo.BuyXname(opts, affcode, team)
	if err != nil {
		fmt.Println("buy key error", err)
	}

	fmt.Println("买key FIRE !!! Pending Tx ID is: ", tx.Hash().String())
}
