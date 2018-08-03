#!/bin/bash
# Author: xxp
# start ethereum privatenet, geth with v1.8.12
#

# 安装geth和ethereumwallet的方法，自行google
# genesis.json通过puppeth自动生成的

# 本目录只是init后的产物
# geth --datadir node0 init genesis.json

set -e
set -u

# 启动挖矿节点，默认启用POA(clique)共识算法
# coinbase账号默认有3780万个Ether
geth  \
    --datadir ./node0\
    --ws\
    --wsaddr 0.0.0.0\
    --wsapi "eth,net,web3,admin,personal,txpool,miner,clique,debug"\
    --wsport 18546\
    --wsorigins "*"\
    --rpc\
    --rpcapi "eth,net,web3,admin,personal,txpool,miner,clique,debug"\
    --rpccorsdomain "*"\
    --rpcaddr 0.0.0.0\
    --rpcport 18545\
    --rpcvhosts "*"\
    --mine\
    --etherbase 0xdbeb69c655b666b3e17b8061df7ea4cc2399df11\
    --unlock 0xdbeb69c655b666b3e17b8061df7ea4cc2399df11\
    --password ./password\
    --nodiscover\
    --maxpeers '50'\
    --networkid 378\
    --targetgaslimit 471238800\
    &

# 再次进入console调试
# geth attach ipc:\./node0/geth.ipc

ethereumwallet\
    --rpc ./node0/geth.ipc\
    --node geth\
    --gethpath `which geth`\
    --swarmurl="http://swarm-gateways.net"\
    &
