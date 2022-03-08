# Go Ethereum Analyze

本仓库主要面向需要阅读以太坊源码的读者，暂时只关注研究智能合约安全必需的合约执行机制。网络和共识部分也欢迎朋友补充。

本仓库复制官方的 `go-ethereum`，只进行解读，不做实际源码修改，供大家学习交流，请勿用于生产环境。

主要目标：

- [ ]  深度理解以太坊底层运行机制
- [ ]  深度理解以太坊字节码，能够对字节码进行安全分析。
- [ ]  提供修改源码的参考，能够修改源码以完成测试、插桩，甚至新建自己的链。
- [ ]  提供交易动态数据分析方法，最终实现区块链平台威胁感知和攻击检测的模型。（由于需要发表论文，需要暂时保密，时机成熟后公开）

注意：

- 只解释部分细节，源码注释中不作过于详细的介绍，只提及所需的基础和学习的方向。
- 对应的理论基础存放在 `analyzeSourceCode` 文件夹内，针对对应的内容，梳理脉络和实现流程。
- 保留了英文注释，对于容易理解的英文注释并未翻译成中文。

Done list

- [X]  `core/types/transaction.go`
- [X]  `core/types/transaction_signing.go`
- [x]  `core/genesis.go`
- [x]  `core/types/legacy_tx.go`
- [x]  `core/types/receipt.go`
- [x]  `core/types/access_list_tx.go`
- [x]  `core/types/dynamic_fee_tx.go`

Todo list:

- [ ]  `core/bloombits`
- [ ]  `core/forkchoice`
- [x]  `core/types/log.go`
- [ ]  区块链学习路线文章
- [ ]  以太坊核心数据结构文章
- [ ]  `core/tx_pool.go`
- [ ]  `core/tx_list.go`

---

## 环境准备

为了方便修改源码后进行调试，建议在 Linux 系统运行，阅读源码时用主机系统可能相对方便。

**安装虚拟机**：建议 Ubuntu20.04，具体安装教程可见其他教程。希望读者有一定的 Linux 基础，熟悉常用命令，理解 Linux 配置文件的思想，会阅读命令行提示信息。安装好虚拟机后可设置代理（学习区块链必须要学会设置代理），自行寻找教程。

**准备环境**：安装 nodejs、npm、goland、typora、中文输入法（可以用百度输入法）、vscode、goland、git，新手自行寻找教程，很容易找到。

建议使用 godoc，然后输入命令 `godoc --http localhost:6060`，可以方便地看到自动生成的文档。

**配置go编译器**：ubuntu apt包管理工具的go编译器版本太低，需要手动升级。

1. 下载二进制包，可在这里[下载](https://studygolang.com/dl)，然后解压、复制到 `/usr/local`，这是我们一般放软件的地方，然后设置环境变量，建议学习 Linux 的环境变量如何设置，有什么作用。这个目录就是 GOROOT 目录，GOPATH 目录可以通过 `go env` 找到。
2. 下载 go-ethereum，建议使用 apt 下载，因为不用额外的配置，并且附带了一些好用的工具。

   添加源：`sudo add-apt-repository -y ppa:ethereum/ethereum`

   更新软件表：`sudo apt-get update`

   安装稳定版：`sudo apt-get install ethereum`

   遇见问题参考 [官方教程](https://geth.ethereum.org/docs/install-and-build/installing-geth) 和博客。注意如果未设置代理，下载速度会很慢，耐心等待。关于如何设置代理，请自行查找教程。
3. 克隆仓库，开始工作。配置好自己的 Git，不会使用请自行寻找教程。然后 VSCode、goland 配置好相关环境（自行寻找教程）
4. 安装翻译插件，复制阅读英文注释，推荐 JetBrain 系的 [Transaction](https://yiiguxing.github.io/TranslationPlugin/index.html) 和 VSCode 的 Comment Translate。

---

**请注意：仓库基于较新的 geth 源码，本教程和源码都只用于学习，并且只负责解读以及如何修改，不会做实质性的改变**

## 以太坊术语

在开始解读以太坊前，先了解以太坊中常见术语和名词。以便更好的学习后续内容。大部分内容来自 [博客](https://learnblockchain.cn/books/geth/part0/term.html)，本人补充和完善。

### 专有名词

- 外部账户：EOAs（External Owned Accounts），关联个人掌握的私钥。可以用于发送交易（转移以太币或发送消息），形同一张带数字ID的储蓄卡。
- 合约账户：Contracts Accounts，可以在以太坊上存储合约代码与合约数据的账户，外部不能直接操作此账户。只能由外部账户直接或间接调用。
- 账户状态： account state，表示一个账户在以太坊中的状态。账户状态在账户数据变化时变化。账户状态包含四项信息：nonce、余额、账户存储内容根哈希值、账户代码哈希值。状态数据不直接存储在区块上。
- 账户 Nonce: 账户随机数，是账户的交易计数。以防止重放攻击。
- 智能合约：Smart Contract，是以太坊成为区块链2.0的立足点。以太坊支持通过图灵完备的高级编程语言编写智能合约代码。部署在链上后，可以接受来自外部的交易请求和事件，以触发执行特定的合约代码逻辑，进一步生成新的交易和事件。甚至调用其他的智能合约。
- 世界状态：state，管理账户地址到账户状态的映射关系。所有账户的状态构成整个区块链状态。
- 交易：Transaction，是外部与以太坊交互的唯一途径，必须由外部账户签名，矿工执行交易，最终打包到区块中。
- 交易收据：Receipt，是方便对交易进行零知识证明、索引和搜索，将交易执行过程中的一些特定信息编码为交易收据。
- 区块：block，是由一组交易和一些辅助信息（简称区块头）、其他区块头哈希构成的数据块。其他区块头哈希表示父区块或者叔区块。
- 叔块：Uncle Block，不能成为主链一部分的孤儿区块，如果有幸被后来的区块收留进区块链就变成了叔块。收留了孤块的区块有额外的奖励。孤块一旦成为叔块，该区块统一可获得奖励。通过叔块奖励机制，来降低以太坊软分叉和平衡网速慢的矿工利益。
- 随机数：nonce，记录在区块头中，努力工作的证明。
- Gas：燃料是交易打包到区块时，在 EVM 运行所消耗的资源量的一种形象化概念，比喻需要燃料才能运行 EVM。在以太坊中，将 CPU 资源、存储资源按内置的规则，统一使用 Gas 作为资源单位表达。每执行一次虚拟机指令，均消耗一定的 Gas。
- GasPrice: 燃料价格，任何交易都需要包含一个愿意支付的燃料单价，最终根据交易消耗的燃料量，计算手续费 (usedGas*gasPrice) 支付给矿工。
- 价格预测：GPO(Gas Price Oracle)，Gas 价格预测，根据历史交易的 GasPrice 预测未来 GasPrice 走势。

### 技术术语

- ZKP: Zero Knowledge Proof，零知识证明。
- EVM：Ethereum Virtual Machine，以太坊虚拟机是执行交易的一个轻量级沙盒虚拟机。
- Message：消息，是一个不能序列化的，并且只存在于以太坊运行环境中的虚拟对象，一条消息主要包括：消息的发送方、接收方、gasLimit 等等；
- 序列化：将数据使用RLP编码为一组字节数据，便于数据交换与存储。
- RLP: 递归长度前缀编码，一种能够压缩数据的数据编码协议，在以太坊中常用于序列化数据。
- MPT：默克尔压缩前缀树， Merkle Patricia Tree，是一种经过改良的、融合了默克尔树和前缀树两种树结构优点的数据结构，是以太坊中用来组织管理账户数据、生成交易集合哈希的重要数据结构。
- Patricia Trie: 一种压缩前缀树，是一种更节省空间的树，对于 trie 的每个节点，如果该节点是其父节点唯一的儿子的话，就和父节点结合；
- Merkle Tree: 默克尔树，也称为 Hash Tree，默克尔树叶子节点的value是数据项的内容，或者是数据项的哈希值；非叶子节点的value根据其孩子节点的信息，然后按照Hash算法计算而得出的。
- Whisper：密语，是一种依托于 P2P 的通信协议，通过 Whisper 协议，节点可以将信息发送给某个特定节点，实现双节点私聊和按主题在多个节点上通信。主要用于大规模的点对点数据发现、信号协商、最小传输通信、完全隐私保护的 DApp 而设计的。
- LES： Light Ethereum Subprotocol，以太坊客户端的轻量级的子协议，只需要下载区块头，其他详细信息可以按需获取；LES Wiki
- Swarm： 蜂巢，是一个分布式存储平台和内容分发服务，是以太坊 web 3 技术栈的本地基础层服务;
- LLL，Sperpent、Mutan和Solidity：用于编写智能合约代码的的编程语言，能被编译成 EVM 代码。
- ERC20: 可以理解成 Ethereum 的一个 Token 协议规范，所有基于 Ethereum 开发的 Token 合约都遵守这个规范。遵守 ERC20 协议规范的 Token 可以被各种 Ethereum 钱包支持。
- ERC721: 是在ERC20标准上建立的Token协议规范，是针对不可互换 Token(non-fungible tokens 简称 NFT)做的智能合约标准。

## 源码分析准备

### `cmd` 目录里自带的工具.


|  Command  | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| :----------: | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`geth`** | 命令行主程序，使用参考[博客1](https://knarfeh.com/2018/03/10/go-ethereum%20%E6%BA%90%E7%A0%81%E7%AC%94%E8%AE%B0%EF%BC%88cmd%20%E6%A8%A1%E5%9D%97-geth%20%E5%91%BD%E4%BB%A4%EF%BC%89/) [博客2](https://github.com/blockchainGuide/blockchainguide/blob/main/source_code_analysis/ethereum/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E5%9F%BA%E7%A1%80%E7%90%86%E8%AE%BA%E9%83%A8%E5%88%86/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E5%90%AF%E5%8A%A8%E5%8F%82%E6%95%B0%E8%AF%A6%E8%A7%A3.md)和[官方文档](https://geth.ethereum.org/docs/interface/command-line-options)。 |
|   `clef`   | 签名工具，可以在后端为`geth`签名.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
|  `devp2p`  | P2P 开发工具，不用运行全节点就可以和其他节点通信。                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
|  `abigen`  | 代码生成器，把合约封装成易用 Golang 的包.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `bootnode` | 客户端的精简版，只实现了网络节点协议, 可以在私有网络中辅助寻找节点。                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
|   `evm`   | 以太坊虚拟机 EVM 的开发程序 能够在可配置的环境中运行底层的字节码片段，方便细致的调试以太坊操作码，深入执行过程。                                                                                                                                                                                                                                                                                                                                                                                                                              |
| `rlpdump` | 以以太坊协议的编码 RLP ([Recursive Length Prefix](https://eth.wiki/en/fundamentals/rlp)) 格式输出。                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| `puppeth` | 创建新的以太坊网络时的引导。                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |

### 项目结构

```json
├── accounts	//账户管理
│   ├── abi		//实现 abi
│   │   └── bind //生成 合约的 go 语言封装
│   │       └── backends
│   ├── external
│   ├── keystore //私钥管理，采用 secp256k 加密
│   │   └── testdata
│   │       ├── dupes
│   │       ├── keystore
│   │       │   └── foo
│   │       └── v1
│   │           └── cb61d5a9c4896fb9658090b597ef0e7be6f7b67e
│   ├── scwallet
│   └── usbwallet //硬件钱包，通过 USB 插入
│       └── trezor	//硬件钱包的协议
├── build	//用于编译和构建的脚本
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
│   ├── faucet //轻量级的水龙头（用于取测试币）
│   ├── geth
│   │   └── testdata
│   │       └── vcheck
│   │           ├── minisig-sigs
│   │           ├── signify-sigs
│   │           └── sigs
│   ├── p2psim //模拟 HTTP API 的调用
│   ├── puppeth //构建私链相关
│   │   └── testdata
│   ├── rlpdump //格式化 rlp 编码，更加漂亮的输出
│   └── utils //一些转化、辅助工具
├── common	//工具类
│   ├── bitutil //快速的按位操作
│   ├── compiler//封装 Solidity 和 Vyper 的字节码
│   ├── fdlimit
│   ├── hexutil //十六进制编码
│   ├── math //整数的相关数学工具
│   ├── mclock //用于计时的固定点
│   └── prque //优先队列数据结构
├── consensus //共识算法部分
│   ├── clique //POA
│   ├── ethash //POW
│   └── misc
├── console	//控制台
│   ├── prompt
│   └── testdata
├── contracts	//合约部分
│   └── checkpointoracle
│       └── contract
├── core	//核心数据结构，包括状态机、链式结构、虚拟机等
│   ├── asm //解析汇编指令
│   ├── bloombits //布隆过滤器批量处理
│   ├── forkid //EIP-2124 (https://eips.ethereum.org/EIPS/eip-2124) 的实现
│   ├── rawdb //底层数据库访问
│   │   └── testdata
│   ├── state //以太坊状态树的缓存
│   │   ├── pruner
│   │   └── snapshot //动态状态存储
│   ├── types //共识机制中的数据类型
│   └── vm //EVM 的实现
│       ├── runtime //负责字节码的执行
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
│   ├── catalyst //RPC 相关
│   ├── downloader //全节点同步
│   ├── ethconfig //以太坊配置文件和轻节点配置文件
│   ├── fetcher //获取同步时的区块头、交易等
│   ├── filters //区块、日志、事件、交易的过滤
│   ├── gasprice
│   ├── protocols
│   │   ├── eth
│   │   └── snap
│   └── tracers //跟踪交易
│       ├── internal
│       │   └── tracetest
│       │       └── testdata
│       │           ├── call_tracer
│       │           └── call_tracer_legacy
│       ├── js //js 写的交易跟踪器
│       │   └── internal
│       │       └── tracers
│       └── native //go 写的交易跟踪器
├── ethclient	//RPC 调用的客户端
│   └── gethclient
├── ethdb	//数据库
│   ├── dbtest
│   ├── leveldb //leveldb 数据库实现
│   └── memorydb //内存映射的数据库实现
├── ethstats	//网络状态显示
├── event	//处理实时事件
├── graphql	//提供 GraphQL 的借口
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
├── les	//轻量级以太坊子协议(LES)
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
├── miner	//区块生成和挖矿
│   └── stress
│       ├── 1559 //EIP1559 的压力测试
│       ├── clique //Clique 的压测
│       └── ethash //ethash 的压测
├── mobile //为移动端设置的简化版 API
├── node	//节点协议
├── p2p	//P2P 网络协议
│   ├── discover //节点发现协议
│   │   ├── v4wire //v4 版本
│   │   └── v5wire //v5 版本
│   │       └── testdata
│   ├── dnsdisc //EIP1459 提出的发现协议
│   ├── enode
│   ├── enr //EIP778 提出的节点记录
│   ├── msgrate //估计节点吞吐量实现更平衡的传输
│   ├── nat //端口映射
│   ├── netutil
│   ├── nodestate
│   ├── rlpx //RLPx 传输协议
│   ├── simulations
│   │   ├── adapters
│   │   ├── examples
│   │   └── pipes
│   └── tracker
├── params	//参数规定
├── rlp	//RLP 序列化格式
├── rpc	//双向 JSON-RPC 2.0
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
└── trie	//区块的重要数据结构 MPT
```

## 参考

[登链学院以太坊技术与实现](https://learnblockchain.cn/books/geth/part1/genesis.html)

[geth 官方文档](https://geth.ethereum.org/docs/)

[官网资源](https://ethereum.org/en/)

[姚飞亮博客](https://www.yaofeiliang.com/tags/#%E5%8C%BA%E5%9D%97%E9%93%BE)

[四年前大佬的源码分析](https://github.com/ZtesoftCS/go-ethereum-code-analysis)

[最近底层大佬的源码分析](https://github.com/blockchainGuide/blockchainguide/tree/main/source_code_analysis/ethereum/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)

[18年的博客](https://knarfeh.com/2018/03/10/go-ethereum%20%E6%BA%90%E7%A0%81%E7%AC%94%E8%AE%B0%EF%BC%88%E6%A6%82%E8%A7%88%EF%BC%89/)

[简书博客](https://www.jianshu.com/u/572268941378)、
