# Go Ethereum Analyze

参考 [登链学院以太坊技术与实现](https://learnblockchain.cn/books/geth/part1/genesis.html)、[geth 官方文档](https://geth.ethereum.org/docs/)、[官网资源](https://ethereum.org/en/)、[姚飞亮博客](https://www.yaofeiliang.com/tags/#%E5%8C%BA%E5%9D%97%E9%93%BE)、[四年前大佬的源码分析](https://github.com/ZtesoftCS/go-ethereum-code-analysis)、[最近底层大佬的源码分析](https://github.com/blockchainGuide/blockchainguide/tree/main/source_code_analysis/ethereum/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)、[18年的博客](https://knarfeh.com/2018/03/10/go-ethereum%20%E6%BA%90%E7%A0%81%E7%AC%94%E8%AE%B0%EF%BC%88%E6%A6%82%E8%A7%88%EF%BC%89/)

## 环境准备

**安装虚拟机**：



**配置环境：**





## 源码分析准备

### `cmd` 目录里自带的工具.

|    Command    | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| :-----------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  **`geth`**   | 命令行主程序，使用参考[博客1](https://knarfeh.com/2018/03/10/go-ethereum%20%E6%BA%90%E7%A0%81%E7%AC%94%E8%AE%B0%EF%BC%88cmd%20%E6%A8%A1%E5%9D%97-geth%20%E5%91%BD%E4%BB%A4%EF%BC%89/) [博客2](https://github.com/blockchainGuide/blockchainguide/blob/main/source_code_analysis/ethereum/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E5%9F%BA%E7%A1%80%E7%90%86%E8%AE%BA%E9%83%A8%E5%88%86/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E5%90%AF%E5%8A%A8%E5%8F%82%E6%95%B0%E8%AF%A6%E8%A7%A3.md)和[官方文档](https://geth.ethereum.org/docs/interface/command-line-options)。 |
|   `clef`    | 签名工具，可以在后端为`geth`签名. |
|   `devp2p`    | P2P 开发工具，不用运行全节点就可以和其他节点通信。 |
|   `abigen`    | 代码生成器，把合约封装成易用 Golang 的包. 它可以结合合约字节码把 [应用二进制接口 ABIs](https://docs.soliditylang.org/en/develop/abi-spec.html) . 它会对 [Native DApps](https://geth.ethereum.org/docs/dapp/native-bindings) 有一点改变。 |
|  `bootnode`   | 客户端的精简版，只实现了网络节点协议, 可以在私有网络中辅助寻找节点。                                                                                                                                                                                                                            |
|     `evm`     | 以太坊虚拟机 EVM 的开发程序 能够在可配置的环境中运行底层的字节码，方便细致的调试以太坊操作码，深入执行过程。                                                                                                                                                                                                           |
|   `rlpdump`   | 以以太坊协议的编码 RLP ([Recursive Length Prefix](https://eth.wiki/en/fundamentals/rlp)) 格式输出。                                                                                                                                                                                                      |
|   `puppeth`   | 创建新的以太坊网络时的引导。                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |

### 项目结构.

```json
├── accounts	//账户管理
│   ├── abi
│   │   └── bind
│   │       └── backends
│   ├── external
│   ├── keystore
│   │   └── testdata
│   │       ├── dupes
│   │       ├── keystore
│   │       │   └── foo
│   │       └── v1
│   │           └── cb61d5a9c4896fb9658090b597ef0e7be6f7b67e
│   ├── scwallet
│   └── usbwallet
│       └── trezor
├── build	//用于编译和构建
│   ├── bin
│   └── deb
│       └── ethereum
├── cmd	//命令行工具
│   ├── abidump
│   ├── abigen
│   ├── bootnode
│   ├── checkpoint-admin
│   ├── clef
│   │   ├── docs
│   │   │   └── qubes
│   │   ├── testdata
│   │   └── tests
│   ├── devp2p
│   │   └── internal
│   │       ├── ethtest
│   │       │   └── testdata
│   │       ├── v4test
│   │       └── v5test
│   ├── ethkey
│   ├── evm
│   │   ├── internal
│   │   │   ├── compiler
│   │   │   └── t8ntool
│   │   └── testdata
│   │       ├── 1
│   │       ├── 10
│   │       ├── 11
│   │       ├── 12
│   │       ├── 13
│   │       ├── 14
│   │       ├── 15
│   │       ├── 16
│   │       ├── 17
│   │       ├── 18
│   │       ├── 19
│   │       ├── 2
│   │       ├── 3
│   │       ├── 4
│   │       ├── 5
│   │       ├── 7
│   │       ├── 8
│   │       └── 9
│   ├── faucet
│   ├── geth
│   │   └── testdata
│   │       └── vcheck
│   │           ├── minisig-sigs
│   │           ├── signify-sigs
│   │           └── sigs
│   ├── p2psim
│   ├── puppeth
│   │   └── testdata
│   ├── rlpdump
│   └── utils
├── common	//工具类
│   ├── bitutil
│   ├── compiler
│   ├── fdlimit
│   ├── hexutil
│   ├── math
│   ├── mclock
│   └── prque
├── consensus //共识算法部分
│   ├── clique
│   ├── ethash
│   └── misc
├── console	//控制台
│   ├── prompt
│   └── testdata
├── contracts	//合约部分
│   └── checkpointoracle
│       └── contract
├── core	//核心数据结构，包括状态机、链式结构、虚拟机等
│   ├── asm
│   ├── bloombits
│   ├── forkid
│   ├── rawdb
│   │   └── testdata
│   ├── state
│   │   ├── pruner
│   │   └── snapshot
│   ├── types
│   └── vm
│       ├── runtime
│       └── testdata
│           └── precompiles
├── crypto	//哈希算法和密码学
│   ├── blake2b
│   ├── bls12381
│   ├── bn256
│   │   ├── cloudflare
│   │   └── google
│   ├── ecies
│   ├── secp256k1
│   │   └── libsecp256k1
│   │       ├── build-aux
│   │       │   └── m4
│   │       ├── contrib
│   │       ├── include
│   │       ├── obj
│   │       ├── sage
│   │       └── src
│   │           ├── asm
│   │           ├── java
│   │           │   └── org
│   │           │       └── bitcoin
│   │           └── modules
│   │               ├── ecdh
│   │               └── recovery
│   └── signify
├── docs	//部分说明文档
│   ├── audits
│   └── postmortems
├── eth	//以太坊协议
│   ├── catalyst
│   ├── downloader
│   ├── ethconfig
│   ├── fetcher
│   ├── filters
│   ├── gasprice
│   ├── protocols
│   │   ├── eth
│   │   └── snap
│   └── tracers
│       ├── internal
│       │   └── tracetest
│       │       └── testdata
│       │           ├── call_tracer
│       │           └── call_tracer_legacy
│       ├── js
│       │   └── internal
│       │       └── tracers
│       └── native
├── ethclient	//RPC 调用的客户端
│   └── gethclient
├── ethdb	//数据库
│   ├── dbtest
│   ├── leveldb
│   └── memorydb
├── ethstats	//网络状态显示
├── event	//处理事件部分
├── graphql	//针对 GraphQL 的部分
├── internal	//内部的一些组件
│   ├── build
│   ├── cmdtest
│   ├── debug
│   ├── ethapi
│   ├── flags
│   ├── guide
│   ├── jsre
│   │   └── deps
│   ├── syncx
│   ├── testlog
│   ├── utesting
│   └── web3ext
├── les	//轻量级子协议
│   ├── checkpointoracle
│   ├── downloader
│   ├── fetcher
│   ├── flowcontrol
│   ├── utils
│   └── vflux
│       ├── client
│       └── server
├── light //向轻量级客户端提供按需检索的功能
├── log	//日志
├── metrics	//磁盘读写相关
│   ├── exp
│   ├── influxdb
│   ├── librato
│   └── prometheus
├── miner	//挖矿部分
│   └── stress
│       ├── 1559
│       ├── clique
│       └── ethash
├── mobile //为移动端设置的封装
├── node	//节点类型
├── p2p	//P2P 网络协议
│   ├── discover
│   │   ├── v4wire
│   │   └── v5wire
│   │       └── testdata
│   ├── dnsdisc
│   ├── enode
│   ├── enr
│   ├── msgrate
│   ├── nat
│   ├── netutil
│   ├── nodestate
│   ├── rlpx
│   ├── simulations
│   │   ├── adapters
│   │   ├── examples
│   │   └── pipes
│   └── tracker
├── params	//参数规定
├── rlp	//编码部分
├── rpc	//远程方法调用部分
│   └── testdata
├── signer	//数字签名部分
│   ├── core
│   │   ├── apitypes
│   │   └── testdata
│   │       └── fuzzing
│   ├── fourbyte
│   ├── rules
│   └── storage
├── swarm	//swarm 群节点
├── tests	//测试数据
│   ├── fuzzers
│   │   ├── abi
│   │   ├── bitutil
│   │   ├── bls12381
│   │   │   └── testdata
│   │   ├── bn256
│   │   ├── difficulty
│   │   │   └── debug
│   │   ├── keystore
│   │   │   └── corpus
│   │   ├── les
│   │   │   └── debug
│   │   ├── rangeproof
│   │   │   ├── corpus
│   │   │   └── debug
│   │   ├── rlp
│   │   │   └── corpus
│   │   ├── runtime
│   │   ├── secp256k1
│   │   ├── stacktrie
│   │   │   └── debug
│   │   ├── trie
│   │   │   └── corpus
│   │   ├── txfetcher
│   │   │   └── corpus
│   │   └── vflux
│   │       └── debug
│   ├── solidity
│   │   ├── contracts
│   │   ├── migrations
│   │   └── test
│   └── testdata
└── trie	//区块的重要数据结构字典树
```





