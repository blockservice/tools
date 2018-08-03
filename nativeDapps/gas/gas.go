package gas

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/robfig/cron"
)

var gWei int64 = 1000000001
var gasPriceCache = big.NewInt(3 * gWei)
var cacheLock sync.RWMutex

func init() {
	c := cron.New()
	c.AddFunc("@every 5m", suggestGasPrice)
	c.Start()
}

// response of  https://ethgasstation.info/json/ethgasAPI.json
type gasAPI struct {
	Fastest float32 `json:"fastest"`
	Average float32 `json:"average"`
	SafeLow float32 `json:"safeLow"`
}

// Price get gasPrice from Cache
func Price() *big.Int {
	cacheLock.RLock()
	cached := gasPriceCache
	cacheLock.RUnlock()

	return cached
}

// SetGasPrice 后台管理需要的Set接口
func SetGasPrice(gasPrice *big.Int) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	gasPriceCache = gasPrice
}

func suggestGasPrice() {
	url := "https://ethgasstation.info/json/ethgasAPI.json"

	gasClient := http.Client{
		Timeout: time.Second * 8, // Maximum of 8 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	res, getErr := gasClient.Do(req)
	if getErr != nil {
		log.Println(getErr)
		return
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(readErr)
		return
	}

	gas := gasAPI{}
	jsonErr := json.Unmarshal(body, &gas)
	if jsonErr != nil {
		log.Println(jsonErr)
		return
	}

	// Average = gWei ×　１０, if gas.Average < 1, will be 0
	gasPrice := big.NewInt(int64(gas.Average) * gWei / 10)

	// if gasPrice is broken, will set the last value.
	if gasPrice.Cmp(big.NewInt(0)) < 0 {
		log.Println("Got the gasPirce is Broken!", gasPrice)
		return
	}
	log.Println("next Ethereum gasPric is Wei ", gasPrice)
	SetGasPrice(gasPrice)
}
