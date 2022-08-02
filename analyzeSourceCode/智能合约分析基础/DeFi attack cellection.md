前言：本文只为了交流学习，为了快速学习，已有的写的较好的资料直接照抄，并且在参考中注明来源。介绍的都非常简略，但是读者可以根据参考链接得到详细的讲解。

来源：[geth-analyze](https://github.com/learnerLj/geth-analyze)

# 汇总表格

| 发生时间             | 名字               | 类型                           | 损失金额       | 简述原因                                                     | 代码                                                         | txhash                                                       |
| -------------------- | ------------------ | ------------------------------ | -------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 2022 年 07 月 25 日  | LPC                | BNB                            | 约45,715 美元  | 实际的账户余额在后续操作中变化了，但是变量 recipientBalance 没有再次更新，导致后续计算余额偏大。 |                                                              |                                                              |
| 2022 年 06 月 27 日  | **XCarnival NFT**  | ETH，业务逻辑漏洞              | 约 380 万美元  | 未考虑抵押物被取出不能借款的业务逻辑，主要是 isWithdraw(bool) 未在borrow过程中考虑 | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu8F7kEfcAwJijxpxQ?e=bkBbg4) | 0x51cbfd46f21afb44da4fa971f220bd28a14530e1d5da5009cfbdfee012e57e35，0x20396765bbeacb5062edf7a18d90016dfaf61a05229e17f9c432c1a2af429dde，0x2d785de160e2cd2148041926d643b9fa39f6724cfa82d9aabe4628711f096844，0xcc3fda1e5540486de15f707ccc82a6f9c8c78e0ef3ef02e4318b3bea24ace701，0xd46e881c78ad0cc18b8076063fb6dcf29346eabe49dc9d1435e6a47af953747f，0xa7affbd744f7a96715fa4904e4b9f8cab0357cefc269e3d9fe6b68524482f703，0xc129fefaa52c65bc31d5b5581f70ffe891deb518e2c79e8e6683bb3b58b3e008，0xc129fefaa52c65bc31d5b5581f70ffe891deb518e2c79e8e6683bb3b58b3e008，0xf12eafb1c48f0b0d793e9e88ba2a3ad0a90112771e580ec92b07daa6e2ddc5d0 |
| 2022 年 04 月 15 日  | Rikkei Finance     | BSC, oracle，鉴权              |                | oracle未鉴权，导致恶意设置oracle，导致代币汇率恶意调整       |                                                              | 0x93a9b022df260f1953420cd3e18789e7d1e095459e36fe2eb534918ed1687492 |
| 2022 年 03 月 29 日  | Ronin Network      | 侧链，私钥泄露                 | 6.1 亿美元     | 为了方便而设置的私钥未及时收回权限，后来私钥被黑，并且权限控制没做好，导致黑客通过一个私钥获取到了5/9私钥的权限，大额取出资金。 |                                                              |                                                              |
| 2022 年 3 月 21 日   | OneRing Finance    | Fantom，汇率                   | 约 146 万美元  | 经济学模型漏洞，根据pool的储备量计算价格，导致巨额转账操纵汇率。 | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu8HlNtgQ2HsDotpeg?e=uFLhaH) | 0xca8dd33850e29cf138c8382e17a19e77d7331b57c7a8451648788bbb26a70145 |
| 2021 年 08 月 10 日  | Poly Network       | 跨链协议，鉴权                 | 6.1 亿美元     | 验证区块头的合约处理执行交易的跨链请求，但是这个合约具备修改中继验证者的权限，导致可以通过构造哈希碰撞，修改了验证者。 |                                                              |                                                              |
| 2022 年 03 月 16 日  | Hundred Finance    | xdai稳定币链，重入             |                | 借贷合约采用ERC667，但是未加互斥锁，导致重入，借走大量钱。   |                                                              |                                                              |
| 2022 年 03 月 13 日  | Paraluni           | BSC，重入                      | 约 170 万美元  | LP token 具备 _pid标识，但是未校验一致性（小问题）。主要问题是，计算流动行未考虑重入。 |                                                              |                                                              |
| 2022 年 03 月 03 日  | TreasureDAO NFT    | Arbitrum                       |                | ERC721和ERC1155交易在同一个函数中进行，忽视转账NFT数量为0时，ERC721转账仍然进行。 | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu8ItYWq3RuVnGLUKw?e=lhTSK3) | 0x82a5ff772c186fb3f62bf9a8461aeadd8ea0904025c3330a4d247822ff34bc02 |
| 2022 年 01 月 28 日  | **QBridge**        | ETH                            | 约 8000 万美金 | ERC token 对应一个ID，用于选取handler，但是bytes->address的映射如果key不存在，handler为默认的零地址，绕过了转账流程，零地址也恰好绕过了三重检查。 |                                                              | 0x478d83f2ad909c64a9a3d807b3d8399bb67a997f9721fc5580ae2c51fab92acf， |
| 2021 年 11 ⽉ 30 ⽇  | **MonoX Finance**  | ETH, polygon，经济学模型，鉴权 | 约 3100 万美元 | 消除流动性未鉴权；反复 MONO swap MONO ，导致 MONO 价格极端高。 | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu86C9EWmwyh3zTQew?e=3JFoZf) | 0x9f14d093a2349de08f02fc0fb018dadb449351d0cdb7d0738ff69cc6fef5f299 |
| 2021 年 10 月 27 日  | **Cream Finance**  | ETH                            | 约 1.3 亿美元  | 流动性股价计算是 (总代币资产+股份资产)/股票数量，攻击者利用闪电贷和杠杆添加大量流动性股票，然后转入大量代币资产拉高股价。 |                                                              | 0x0fe2542079644e107cbf13690eb9c2c65963ccb79089ff96bfaf8dced2331c92 |
| 2021 年 9 月 12 日   | Zabu Finance       | Avalanche                      |                | 抵押存在手续费，但是合约中记录的抵押没有先扣除手续费，导致提现的代币比实际抵押的多。代币总量减少使得抵押奖励极大 |                                                              |                                                              |
| 2021 年 08 月 30 日  | **Cream Finance**  | ETH                            | 约 1800 万美元 | 攻击者从crETH借贷后，借AMP合约重入crETH合约，同时借贷大量AMP代币。 |                                                              | 0xa9a1b8ea288eb9ad315088f17f7c7386b9989c95b4d13c81b69d5ddad7ffe61e |
| 2021 年 08 月 12 日  | **DAO Maker**      | ETH、私钥泄露                  |                | 私钥泄露导致管理员权限授予攻击合约。                         |                                                              | 0x33d1e2700eb2b5626e09a89008a7b445b34dfc72112c0500e80b97057688e236, 0x10caabf9130b1264ed2094c1dc6d35cee7587b35387cdac703e32e41c8ccfea9 |
| 2021 年 08 月 04 日  | **Popsicle**       | ETH                            | 超 2100 万美元 | LP token 可以转让，合约按照LP token 和已支付款项发放奖励，导致LP token 转让后能多次更新奖励 | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu8_9yByfj7L-V9qYg?e=SsUudM) | 0xcd7dae143a4c0223349c16237ce4cd7696b1638d116a72755231ede872ab70fc |
| 2021 年 08 月 04 日  | Wault.Finance      | bsc                            | 93 万美元      | WEX 池子似乎是 WEX 数量越少，则 WEX 价格越高。而攻击者利用闪电贷大量质押，导致价格变动。 |                                                              | 0x31262f15a5b82999bf8d9d0f7e58dcb1656108e6031a2797b612216a95e1670e |
| 2021 年 07 月 17 日  | PancakeBunny       | Polygon                        |                | 闪电贷之后质押权益人设置为 VaultSushiFlipToFlip，导致合约balance异常大，与用户持有的LP相乘后得到异常大值，最后铸币出额外的 polyBUNNY 代币。 | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu9ENkaAl2THB-ymQA?e=NMDYCC) | 0x25e5d9ea359be7fc50358d18c2b6d429d27620fe665a99ba7ad0ea460e50ae55 |
| 2021 年 06 月 28 日  | SafeDollar         | Polygon                        |                | 项目参考  SUSHI 的 MasterChef 合约，它每次存入代币都会 burn一部分，但是取出是存入的数量。反复抵押提现使得PL及其小，使得铸币奖励异常多 |                                                              |                                                              |
| 2021年 06 月 21 日   | Impossible Finance | BSC, 不可信外部调用            |                | 项目为了省gas实现了必须由Router才能调用的cheapSwap，但是攻击者通过创建自己的Pair，在pair中就swap，而后续的cheapSwap未检查K | [onedrive](https://1drv.ms/u/s!At0_LwVPvookhu9GovQTE7HWAzvGwg?e=r36rQB) |                                                              |
| 2021 年 06 月 16 日  | **Alchemix**       | ETH，项目方操作错误            |                | 操作流程失误导致资金损失。                                   |                                                              | *0x3cc071f9f40294bb250fc7b9aa6b2d7e6ca5707ce4d6d222157d7a0feef618b3* |
| 2021 年 05 月 28 日  | BurgerSwap         | BSC                            |                | pair 合约完全依赖  PlatForm（类似router）的检查，本身未检查K值，重入  swapExactTokensForTokens 滑点不变 |                                                              | 0xac8a739c1f668b13d065d56a03c37a686e0aa1c9339e79fcbc5a2d0a6311e333 |
| 2021  年 05 月 20 日 | PancakeBunny       | BSC                            |                | LP token 的价值计算依赖池子中的代币比例，而这较为容易被闪电贷操控。 |                                                              |                                                              |
| 2021年5月20日        | **Rari**           | ETH                            | 近 1500 万美元 | `ibETH.work` 函数可以调用任何合约，导致中途充值ETH，使得 totalETH 增大，导致 ibETH代币价值升高，然后在同一笔交易中取出。 |                                                              | 0x171072422efb5cd461546bfe986017d9b5aa427ff1c07ebe8acc064b13a7b7be |
| 2021 年 03 月 08 日  | **DODO**           | ETH                            | 212 万美元     | 资金池合约初始化函数没有任何鉴权以及防止重复调用初始化的限制，导致攻击者利用闪电贷将真币借出，然后通过重新对合约初始化将资金池代币对替换为攻击者创建的假币，从而绕过闪电贷资金归还检查将真币收入囊中 |                                                              | 0x395675b56370a9f5fe8b32badfa80043f5291443bd6c8273900476880fb5221e |
| 2021 年 03 月 06 日  | **Paid Network**   | ETH                            |                | 私钥泄露或者项目内部原因，导致任意 mint。                    |                                                              | 0x4bb10927ea7afc2336033574b74ebd6f73ef35ac0db1bb96229627c9d77555a0 |
| 2021 年 02 月 27 日  | **Furucombo**      | ETH                            | 超 1500 万美元 | 未初始化代理合约，导致被黑客初始化，更改了 implementation 合约 |                                                              | 0x6a14869266a1dcf3f51b102f44b7af7d0a56f1766e5b1908ac80a6a23dbaf449 |
| 2021 年 1 月 27 日   | **SushiSwap**      | ETH                            |                | DIGG 没有设置兑换路径，而默认的是 WETH，于是攻击者创建交易对，控制价格，导致手续费兑换过程中产生了巨大的滑点。 |                                                              | 0x0af5a6d2d8b49f68dcfd4599a0e767450e76e08a5aeba9b3d534a604d308e60b |











# 攻击分析

## LPC

参考官方的报导即可：

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Brief analysis of <a href="https://twitter.com/search?q=%24LPC&amp;src=ctag&amp;ref_src=twsrc%5Etfw">$LPC</a> flashloan attack:<br/>The attacker first borrowed 1,353,900 <a href="https://twitter.com/search?q=%24LPCs&amp;src=ctag&amp;ref_src=twsrc%5Etfw">$LPCs</a> via flashloan from Pancake, then called the _transfer function in the LPC contract to transfer to himself. <a href="https://t.co/Fxn9O0MTn8">https://t.co/Fxn9O0MTn8</a></p>&mdash; Beosin Alert (@BeosinAlert) <a href="https://twitter.com/BeosinAlert/status/1551535854681718784?ref_src=twsrc%5Etfw">July 25, 2022</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>





##  XCarnival NFT

详细介绍可见：[熊市新考验 —— XCarnival NFT 借贷协议漏洞分析 ](https://mp.weixin.qq.com/s/F2hpBNRzhZmfCn7BQanGOA)和 [XCarnical 攻击事件分析](https://s3cunda.github.io/2022/06/28/XCarnical-%E6%94%BB%E5%87%BB%E4%BA%8B%E4%BB%B6%E5%88%86%E6%9E%90.html)，我来补充代码逻辑。

1. 用户A需要抵押 cellection 合约上的 NFT。cellection 必须是白名单。
2. 每个xtoken 合约都是一个 pool，有自己的借贷容量。
3. 每个抵押物都会生成 order，这是核心数据结构。后面借款操作都是根据order进行。
4. 拍卖是清算的唯一手段，抵押者也可以参与拍卖，这个叫做赎回。
5. 另外，拍卖卖出前，抵押者可以还清贷款，取回抵押物。但是会花费一些手续费。

## Rikkei Finance

参考[慢雾](https://mp.weixin.qq.com/s/W4fn5tVOUGmX_brGLFXgeQ)

## Ronin Network

参考[慢雾](https://mp.weixin.qq.com/s/0U58Chw970X2GWcj2fvLPg)；简单了解[侧链](https://www.binance.com/zh-CN/news/top/7115028)；

## OneRing Finance

总而言之，经济学模型的漏洞，根据瞬时储备量计算价格，利用闪电贷制造巨大差值，导致汇率变化。在这个案例中，攻击者的操作相当的复杂，需要熟悉很多流程才可以。

这里的收获是：我熟悉了 uniswap 和flash loan 的基本思路。



参考：

- [慢雾：OneRing Finance 被黑分析](https://mp.weixin.qq.com/s/MyR_O8wuZJUT1S6eIMH9TA)
- [跟踪交易](https://dashboard.tenderly.co/tx/fantom/0xca8dd33850e29cf138c8382e17a19e77d7331b57c7a8451648788bbb26a70145?trace=0.2.2)
- [攻击合约](https://ftmscan.com/address/0x6a6d593ed7458b8213fa71f1adc4a9e5fd0b5a58)
- https://www.blocktempo.com/onering-finance-is-exploited/



## Poly Network

Poly Network 是由 Neo、Ontology、Switcheo 基金会共同作为创始成员，分布科技作为技术提供方共同发起的跨链组织。

如下图，通过官方的介绍我们可以清楚的看出 Poly Network 的架构设计：用户可以在源链上发起跨链交易，交易确认后由源链 Relayer 将区块头信息同步至 Poly Chain，之后由 Poly Chain 将区块头信息同步至目标链 Relayer，目标链 Relayer 将验证信息转至目标链上，随后在目标链进行区块头验证，并执行用户预期的交易。

<img src="https://mmbiz.qpic.cn/mmbiz_png/qsQ2ibEw5pLZfxqtrw89jEJeBltcLwZ9FkqoWFD4ibLnRzrYGE05YE26R2icdByne7md9HjOE1eL1a6alWBT0wG7g/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1" alt="图片"  />

攻击核心： `EthCrossChainManager `合约用于验证 Poly Chain 同步来的区块头以确认跨链信息的真实。`EthCrossChainData` 合约用于存储跨链数据，中继链验证人 (即 Keeper) 的公钥也存储在这个合约中。`LockProxy` 则用于资产管理。`EthCrossChainManager` 执行交易时没有检验 to 和 _method，导致构造哈希碰撞，修改了中继链验证人.

参考：

- https://mp.weixin.qq.com/s/MyR_O8wuZJUT1S6eIMH9TA
- https://www.blocktempo.com/onering-finance-is-exploited/



##   Hundred Finance

Hundred Finance 是一个去中心化应用程序（DApp），它支持加密货币的借贷。它是一种多链协议，与 Chainlink 预言机集成，以确保市场健康和稳定，同时专门为长尾资产提供市场。

见参考链接

参考：

- https://mp.weixin.qq.com/s/tlXn3IDSbeoxXQfNe_dH3A

##   Paraluni

Paraluni 是 BSC 链上的一个元宇宙金融（DeFi）项目。用户可以质押代币添加流动性获取收益。详细见参考链接。



参考：

- https://mp.weixin.qq.com/s/a5fFI5sFNAyuDxGqTFmC2A
- https://dashboard.tenderly.co/tx/bsc/0x70f367b9420ac2654a5223cc311c7f9c361736a39fd4e7dff9ed1b85bab7ad54



## TreasureDAO NFT

TreasureDAO 是一个基于 Arbitrum（L2）上的 NFT 项目。



参考：

- https://mp.weixin.qq.com/s/SEbXWmugJBz0C00vyzYcCw
- https://dashboard.tenderly.co/tx/arbitrum/0x82a5ff772c186fb3f62bf9a8461aeadd8ea0904025c3330a4d247822ff34bc02
- https://arbiscan.io/address/0x812cda2181ed7c45a35a691e0c85e231d218e273#code
- https://arbiscan.io/address/0x2e3b85f85628301a0bce300dee3a6b04195a15ee



## QBridge

首先了解 blockchain bridge，可以参考附录。**Tornado** Cash Classic is A fully decentralized protocol for private transactions on Ethereum.

黑客从 tonado 获得最初资产，尽力避免被追踪。

resourceID 是资产的 ID（比如某种 token），如果名单中不存在这类资产，则返回零地址（用映射实现）。

跨链资产转移包括了 ECR token 和 native token(ETH)，但是它们共用相同的事件。在 ERC token 转账时，忽视了映射中不存在的 key 对应 default value，这里时0x0，导致绕过转账检测。

---

In the QBridgeHandler contract, tokenAddress in L127 is **address(0).** There are 3 statements to ensure the correctness of the *data* parameter. However, none of them failed

-Line 128: As **address(0)** is whitelisted, Line 128 would be bypassed

-Line 134: As the amount is 190 ETH (bigger than minAmounts), Line 134 would be bypassed

-Line 135: **As address(0)** is an externally owned address (EOA), the low level call from safeTransferFrom() would return successfully

![img](https://miro.medium.com/max/1400/0*sTvpjyy3xEQcuXH7)



参考：

- https://mp.weixin.qq.com/s/PLbuI9JFxyFRlDlj9rPvmQ
- https://learnblockchain.cn/article/3649
- https://docs.tornado.cash/general/readme
- https://etherscan.io/address/0x99309d2e7265528dc7c3067004cc4a90d37b7cc3#code
- https://dashboard.tenderly.co/tx/mainnet/0xac7292e7d0ec8ebe1c94203d190874b2aab30592327b6cc875d00f18de6f3133

## MonoX Finance

MonoX 是⼀种新的 DeFi 协议，使⽤单⼀代币设计⽤于流动性池。这是通过将存⼊的代币与 vCASH 稳定币组合成⼀个虚拟交易对来实现的。其中的单⼀代币流动性池的第⼀个应⽤是⾃动做市商系统 - Monoswap，它在 2021 年 10 ⽉时推出。

详细过程见：https://www.freebuf.com/news/306917.html

主要漏洞点有两个，第一个移出流动性未鉴权，`Monoswap.sol `中的 `removeLiquidity` 函数。这导致有人可以恶意移出其他人的流动性，让自己独享收益。

第二个，`Monoswap.sol ` 中的 `swapExactTokenForToken` 函数可以影响汇率，经济模型设计问题。只要反复 swap 相同的代币，就会让这个代币的价格指数升高。



参考：

- https://mp.weixin.qq.com/s/s0tO1aqOKGlRcXjyZFU_3Q
- https://www.freebuf.com/news/306917.html
- https://etherscan.io/address/0xa1fba7f2079131acb3e2073563ec53c8d43bc144#code
- https://dashboard.tenderly.co/tx/mainnet/0x9f14d093a2349de08f02fc0fb018dadb449351d0cdb7d0738ff69cc6fef5f299

## Cream Finance

简单地说，计算合约中抵押物的价格时，会把池子中所有的代币总数作为依据，而攻击者利用闪电贷借空了池子。

具体步骤：

1. 从 DssFlash 中闪电贷借出 5 亿个 DAI，然后抵押他们获得 4.5 亿个 yDAI，再把 yDAI 作为流动性代币添加到多个池子里，得到新的流动性代币 yUSD。
2. 再 向 Cream 的 crYUSD 池子抵押 yUSD。
3. 从 AAVE 闪电贷借出约 52.4 万个 WETH，并将其抵押到 Cream 的 crETH 池子。再结合第二步抵押的代币，将 crYUSD 池子里的 yUSD 借空。
4. crYUSD 允许杠杆借贷，于是反复循环使得池子中抵押的 yUSD 极端多。
5. 攻击者另外一部分借贷的 1,873 个 ETH 先转换成 745 万个 USDC，再转换成 338 万 个 DUSD，用于赎回第一步用 yDAI 抵押的ownership token，用 ownership token 再从池子中取出 yDAI 等代币。这些代币又被转入 yUSD 的池子中（而不是抵押），直接拉高了流动性分红时股份的价格，攻击者自身抵押 yUSD 得到的**股份价格就高了**，因此可以用他们在其他池子中借出更多的资金。
6. 反复这样的操作，导致抵押物价格巨高无比，将 15 个池子都借空。

股份价格是通过 总资产/股份数量 计算的，而 总资产=代币资产+抵押资产，所以只要拉高 代币资产，而抵押资产所占有的股份不变， 就可以提高股份单价。

参考：

- https://mp.weixin.qq.com/s/ykz63ZtfbObwRs3UTE3toQ

## Zabu Finance

Zabu Finance 是 Avalanche 上的下一代去中心化金融 (DeFi) 项目。Zabu Finance 成熟的生态系统包括收益聚合、收益耕作、抵押、筹款。

详细分析见参考，写的很好。



参考：

- https://learnblockchain.cn/article/3287

## Cream Finance 

2021 年 8 月 30日的攻击。慢雾写的不怎么好，可以参考https://blockcast.it/2021/08/31/amber-group-reentrancy-attack-explained/，在最后面写的很清楚。

这是重入攻击。攻击者从 crETH 借贷后，借 AMP 合约重入crETH合约，同时借贷大量AMP代币。接着还清第一次 crtETH 的借贷，还闪电贷，最终得到 41 WETH + 9.74M AMP。

1. Exp.trigger() 函数中第 94 行先是一个UniswapV2 的闪电贷，借出了500 WETH。实际执行的流程在 uniswapV2Call() 函数 完成。

![img](https://cdni.blockcast.it/wp-content/uploads/2021/08/31143722/image8-768x159.png)

3. 由于闪电贷借的是 WETH 而 crETH 需要使用 ETH 才能铸造，因此在第 105 行，先将 WETH 换成 ETH，接下来将换出的 ETH 全数发给 crETH 合约铸造出 crETH cTokens。在 113 调用一次 `Comptroller.enterMarkets()` 将 crETH 激活以便后续的操作。

   ![img](https://cdni.blockcast.it/wp-content/uploads/2021/08/31143736/image9-768x186.png)

4. 将 crETH Tokens 抵押，得到 **19,480,000 AMP tokens**，在 crAMP.borrow() 函数（121行）中把 AMP 代币转给攻击合约。

   ![img](https://cdni.blockcast.it/wp-content/uploads/2021/08/31143358/image18-768x149.png)

5. 攻击合约收到 AMP tokens 时， `tokensReceived()`   中的回调函数会重入 crETH 合约 ，造成抵押的 crETH Tokens 双花，再次**借贷 355 个 ETH**。

6. 随后攻击者使用另外一个合约 (0x0ec3) 对已经爆仓的合约 (0x38c4) 进行清算，在 `Liquidator.trigger() ` 函数里，攻击者**用部分 AMP 代币清算了自身创造的不良资产**，获得 crETH 抵押品（第60 行），随后将 crETH 换成 **187.58  ETH**（第61 行），并发回给owner，即攻击合约。

![img](https://cdni.blockcast.it/wp-content/uploads/2021/08/31143717/image7-768x164.png)



参考：

- https://blockcast.it/2021/08/31/amber-group-reentrancy-attack-explained/
- https://mp.weixin.qq.com/s/a9s61_u30f4X8310A952_Q

## DAO Maker

参考链接讲得不错。

参考：

- https://mp.weixin.qq.com/s/N-afjgJD3R3JhlcrFxx12A



## Popsicle

Popsicle Finance是专注于自动做市（AMM）流动性提供商（LP）的下一代跨链收益提高平台。旨在成为一个完全分散的平台，由其用户（ICE治理令牌的持有者）管理。ICE令牌将用于对协议更新，资产池包含，费用管理以及协议其他关键运营方面的提案进行投票。

**攻击信息**

通过初步追踪分析，攻击者创建了3个攻击合约，进行了1笔攻击交易，共盗取资金价值超2100万美元，攻击信息如下：

*攻击者钱包地址：*https://cn.etherscan.com/address/0xf9E3D08196F76f5078882d98941b71C0884BEa52

*攻击者合约地址：*

合约1：https://cn.etherscan.com/address/0xdFb6faB7f4bc9512d5620e679E90D1C91C4EAdE6

合约2：https://cn.etherscan.com/address/0x576cf5f8ba98e1643a2c93103881d8356c3550cf

合约3：https://cn.etherscan.com/address/0xd282f740bb0ff5d9e0a861df024fcbd3c0bd0dc8

*攻击交易：*https://cn.etherscan.com/tx/0xcd7dae143a4c0223349c16237ce4cd7696b1638d116a72755231ede872ab70fc

*SorbettoFragola合**约地址：*https://cn.etherscan.com/address/0xc4ff55a4329f84f9Bf0F5619998aB570481EBB48#contracts

**攻击过程**

攻击者在一笔交易中，利用同一攻击手法进行了多次获利，下面分步解析该笔交易，方便读者更清晰的了解攻击过程。

第一步：攻击合约1从Aave利用闪电贷借出3000万枚USDT,1.3万枚ETH,1400枚BTC,3000万枚USDC,300万枚DAI,20万枚UNI。

第二步：攻击者合约1通过Uniswap V3协议使用5492枚WETH和3000枚USDT添加流动性并获取10.52枚Popsicle LP。

第三步：攻击者合约1将获取的10.52枚Popsicle LP发送至攻击者合约2，后者又将LP发送至攻击者合约3，最后将LP归还至攻击者合约1。

第四步：攻击者合约1归还获取10.52枚Popsicle LP,得到添加流动性的5492枚WETH和3000枚USDT。

第五步：攻击者合约2和合约3分别从Uniswap V3优化器SorbettoFragola得到392枚ETH和215万枚USDT（凭空获利）。

交易信息中，攻击中循环执行8次同样的攻击（第二步至第五步）。

第六步：攻击者合约2和合约3将获取到的资金全部转至攻击者合约1。

第七步：攻击者归还闪电贷的3000万枚USDT,1.3万枚ETH,1400枚BTC,3000万枚USDC,300万枚DAI,20万枚UNI及其手续费。

第八步：攻击者将获取的资金统一转至攻击者钱包。



具体成功原因见慢雾的分析。

参考：

- http://www.xilian.link/depth/75760.html
- https://mp.weixin.qq.com/s/O6gJeXVgYqodTXyh8h9FFg

## Wault.Finance

Wault Finance 是一个去中心化的金融中心平台，它将所有主要 DeFi 用例连接到一个简单的生态系统中，位于币安智能链上。简而言之，这是一个多合一的 DeFi 平台。

WUSDMaster 是一个质押 BSC_USDT 换取 WUSD 的合约，可以通过质押 (stake) BSC_USDT 来获得 WUSD， 通过赎回 (redeem) 将 WUSD 燃烧，然后换成 BSC_USDT，在这过程中一部分资金会转给金库 (Treasury)， WUSDMaster 会用 WEX 补贴给用户。

下图显示了 WAULTx（治理代币）、WEX（农业代币）和 WaultSwap（去中心化交易所）如何整合在一起形成 Wault 金融生态系统。

<img src="https://www.aqniu.com/wp-content/uploads/2021/08/image003-5-1024x1024.png" alt="img" style="zoom:67%;" />

1. 首先攻击者在 WaultSwapPair (BSC_BUSD-WUSD) 中通过闪电贷借了 16,839,004 枚 WUSD，并在 WUSDMaster合约中的赎回 (redeem) 函数，将闪电贷借到的 WUSD 燃烧掉，从而得到 15,037,230 BUSD +106,502,606 WEX。而在这里，其WEX的获得**成本约等于 0.015 BUSD/WEX**。
2. 去 PancakePair (WBNB-BSC_USDT) 中通过闪电贷借了 40,000,000 枚BSC_USDT，共使用23, 000, 000BUSD购买了 517,938,118 个WEX。在这一步，其WEX的获得**成本为 0.044 BUSD/WEX**。同时，也是这一步使得WEX的价格变得非常的高。
3. 在完成上述的准备后，攻击者 68 次调用 WUSDMaster 合约中的质押(stake)函数，质押 BUSD获得 WUSD。stake 函数会执行 `wswapRouter.swapExactTokensForTokensSupportingFeeOnTransferTokens` 会按一定比率换出 WEX 代币。这里由于第二步中攻击者将WEX的价格拉高，池子的WEX获得**成本大约为0.131BUSD/WEX**。
4. 攻击者在此时卖出手中所有的WEX，这里卖出的均价为0.041BUSD/WEX。
5. 归还闪电贷，并将获利转换为ETH跨链离场。

**能够操纵价格的原因：**

WEX 池子似乎是 WEX 数量越少，则 WEX 价格越高。而攻击者利用闪电贷大量质押，导致价格变动。

```solidity
    function stake(uint256 amount) external nonReentrant {
        require(amount <= maxStakeAmount, 'amount too high');
        usdt.safeTransferFrom(msg.sender, address(this), amount); //抵押
        if(feePermille > 0) { //扣除手续费
            uint256 feeAmount = amount * feePermille / 1000;
            usdt.safeTransfer(treasury, feeAmount); //转给金库
            amount = amount - feeAmount;
        }
        uint256 wexAmount = amount * wexPermille / 1000;
        usdt.approve(address(wswapRouter), wexAmount); // usdt 授权给 router
        wswapRouter.swapExactTokensForTokensSupportingFeeOnTransferTokens( //router 转 wex 给
            wexAmount,
            0,
            swapPath,
            address(this),
            block.timestamp
        );
        wusd.mint(msg.sender, amount);
        
        emit Stake(msg.sender, amount);
    }
```



参考：

- https://www.aqniu.com/vendor/76328.html
- https://mp.weixin.qq.com/s/aFSnSDPk4RYlcKz6Qr_CmQ
- https://bscscan.com/address/0xa79Fe386B88FBee6e492EEb76Ec48517d1eC759a#code

## PancakeBunny

这里我们来学习分析交易。

1. 首选是闪电贷，这是通过代理合约调用的，接着进入闪电贷，在 `LendingPool/contracts/protocol/lendingpool/LendingPool.sol` 的 `flashLoan` 函数借出许多代币，

2. 然后在 `UniswapV2Router02/contract/UniswapV2Router02.sol` 添加流动性，通过 `0x4b1f1e2435a9c96f7330faea190ef6a7c8d70001` 的 `mint` 函数获取 SLP.

3. 然后通过 `VaultSushiFlipToFlip.deposit`  函数抵押小部分的 LP，权益人设置为攻击合约。

4. 通过 `MiniChefV2.deposit` 函数抵押大部分 LP，权益人设置为 `VaultSushiFlipToFlip`（实际上是它的代理合约）

5. 通过 0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff 的 `swapExactTokensForTokens` 函数，将 100,000 WETH 转换成 30,027,861.276 WMATIC.
6. 通过 `VaultSushiFlipToFlip.withdrawAll` 函数，获得之前的 LP 和 polyBUNNY 代币奖励。但是 VaultSushiFlipToFlip 合约总共的 amount 被第 4 步改变了。
   1. 之后影响到 `profit`，再影响到 `performanceFee`，`_minter.mintForV2` 铸造大量的 polyBUNNY 奖励。
   2. 具体逻辑在 `BunnyMinterV2/contracts/bunny/BunnyMinterV2.sol`，调用里面的 `mintForV2`, 它调用了 `mintFor` 函数，接着进入了 `if (marketBuy == false)` 部分的 `_zapAssets` 函数。
   3. 进入 `else if` 分支，先通过调用 SushiSwap  Router 合约的 `removeLiquidity `函数进行移除流动性，然后调用 `_tokenToAsset` 将移除流动性获得 USDC 与 USDT 代币分别在 QuickSwap 中兑换成 polyBUNNY 与 WETH 代币并在 QuickSwap 中添加流动性。
   4. 然后进入  `ZapPolygon/contracts/zap/ZapPolygon.sol` 的 `zapInToken` 函数，而第 5 步 得到大量的 WMATIC 会让价格下降，因此接下来通过 `_swapMATICToFlip` 函数将 WMATIC 代币兑换成的 WETH 与 polyBUNNY 代币就会较少，导致最后转给 BUNNY_POOL 的 LP 会较少，达到减少消耗攻击者付出的 SLP 目的，最终减少了一部分攻击成本。



参考：

- https://mp.weixin.qq.com/s/f2kD_l9Cs1mHQXBQYemwvQ
- https://www.tuoniaox.com/news/p-509322.html
- https://dashboard.tenderly.co/tx/polygon/0x25e5d9ea359be7fc50358d18c2b6d429d27620fe665a99ba7ad0ea460e50ae55





## SafeDollar

看的比较大略。



参考：

- https://mp.weixin.qq.com/s/3_qOkt6rlp1seRlu6L1Hfg





## Impossible Finance

Impossible Finance 的 DEX 架构参考了 Uniswap v2，但在 Pair 的实现上有所不同。Impossible Pair 分别实现了 cheapSwap 与 swap 两个接口。cheapSwap 函数限制了只由 Router 合约可进行调用，swap 函数则是任意用户都可调用进行代币兑换操作。

为了节省用户支付的交易费，项目方需要尽量精简智能合约中的逻辑（执行的代码越少，则消耗的Gas减少，最终需要支付的交易费也随之降低）。 为了实现低交易费（相对别的平台低），Impossible Finance在Uniswap V2的基础上做了一些优化，而在此次事件中被利用的漏洞正与这些优化密切相关。

[Impossible Finance](https://impossible.finance/): 与PancakeSwap核心业务基本相同，但是有 Low Fees （交易费低）的特色。

参考：

- https://mp.weixin.qq.com/s/CXqGxmXEJ4DeSYb8qv8vVw
- https://zhuanlan.zhihu.com/p/388231258 （这个写的最好）



## Alchemix

类型特殊，不做分析。

## BurgerSwap

BurgerSwap 是一个仿 Uniswap AMM 项目，但是和 Uniswap 架构有所区别。BurgerSwap 架构总体分成【Delegate -> lpPlatForm -> Pair】。其中 Delegate 层管理了所有的 Pair 的信息，并负责创建 lpPlatForm 层。然后 lpPlatForm 层再往下创建对应的 Pair 合约。在整个架构中，lpPlatForm 层充当了 Uniswap 中 Router 的角色，负责将计算交易数据和要兑换的代币转发到 Pair 合约中，完成兑换。

1. **次攻击开始于 Pancake 的闪电贷**，攻击者从 Pancake 中借出了大量的 WBNB
2. 将这些 WBNB 通过 BurgerSwap  兑换成 Burger 代币
3. 攻击者使用自己控制的代币(攻击合约本身) 和 Burger 代币通过 Delegate 层创建了一个交易对并添加流动性。
4. 通过 PaltForm 层的 swapExactTokensForTokens 函数发起了兑换，兑换路径为【攻击者自己控制的代币 -> Burger -> WBNB】
5. 攻击者创建的交易对的 `_innerTransferFrom` 函数会调用攻击者控制的代币合约，于是攻击者从 `_innerTransferFrom` 函数中重入 `swapExactTokensForTokens` 函数。
6. 因为 pair 合约没有在兑换后根据恒定乘积公式检查兑换后的数值，完全依赖了 PlatForm 层的数据进行兑换。在重入的兑换过程中，兑换的数量没有因为滑点的关系而导致兑换数量的减少。



参考：

- https://mp.weixin.qq.com/s/p16-rCxvqQaxj3SWvw0hXw
- 将 LP 抵押至 PancakeBunny 项目的 VaultFlipToFlip 合约中





## PancakeBunny



1. 使用 0.5 个 WBNB 与约 189 个 USDT 在 PancakeSwap 中添加流动性并获取对应的 LP。将 LP 抵押至 PancakeBunny 项目的 VaultFlipToFlip 合约中。

2. 攻击者再次发起另一笔交易，从 PancakeSwap 的多个流动性池子中闪电贷借出约 232 万枚 BNB 代币，并从 Fortube 项目的闪电贷模块借出 296万 USDT 代币。

3. 借来的 296 万 USDT 代币与 7744 WBNB 代币在 PancakeSwap 的 WBNB-USDT 池子添加流动性，并把获得的 14 万 LP 留在 WBNB-USDT 池子中。

4. 将步骤 2 中借到的剩余 231.5 万枚 WBNB在 PancakeSwap 池中兑换为382万枚USDT。（**这里已经可以控制价格，大量BNB代币进入池子，BNB价格变的极低**）

5. 因为步骤 1 抵押了 LP，攻击者直接调用 VaultFlipToFlip 合约的 getReward 函数，来获取 BUNNY 代币奖励。 `getReward` 不仅会取出质押的 LP，还会取出存放在 pair 中的 LP。以下都是 getReward 方法调用的相关过程：

   ![image-20220729224402316](http://blog-blockchain.xyz/202207292244403.png)

   ![image-20220729225221172](http://blog-blockchain.xyz/202207292252242.png)

   1. 进入 mintForV2 操作，其会先将一定量 (performanceFee) 的 LP 转至 WBNB-USDT 池子中移除流动性，在 `_zapAssetsToBunnyBNB` 函数，取出步骤 3 中加入的流动性资金 296 万枚USDT和 7744 枚BNB。
   2. 在完成移除流动性后，调用 `zapBSC.zapInToken` 函数，分别把上一步中收到的 WBNB 与 USDT 代币转入 zapBSC 合约中。
   3. 在 zapInToken 操作中，其会通过 `_swapTokenForBNB` 在 PancakeSwap 的 WBNB-USDT 池子中把转入的 USDT 兑换成 WBNB。
   4. 再通过 `_swapBNBToFlip` 将刚才得到的 WBNB 的一半，在 PancakeSwap 的 WBNB-BUNNY 池子中兑换成 BUNNY 代币，并将得到的 BUNNY 代币与剩余的 WBNB 代币在 WBNB-BUNNY 池子中添加流动性获得 LP， receiver 是 mintForV2 合约。此时 WBNB 的数量非常多。
   5. 接着，`_zapAssetsToBunnyBNB` 函数 会计算 BunnyMinterV2 合约当前收到的 WBNB-BUNNY LP 数量，并将其返回给 mintForV2。
   6. 随后将会调用 `priceCalculator.valueOfAsset` 函数来计算这些 LP 的价值，计算公式为 `amount.mul(IBEP20(WBNB).balanceOf(address(asset))).mul(2).div(IPancakePair(asset).totalSupply())`，单价是BNB数量 \* 2 / 总LP数量，由于 BNB 数量非常多，LP 对 BNB 的价格很高。
   7. 然后使用 LP 铸币，意外出现很多很多  BUNNY 代币。



参考：

- https://mp.weixin.qq.com/s/O2j5OyUh2qJZSRhnMD5KTg
- https://www.freebuf.com/articles/compliance/274621.html
- https://dashboard.tenderly.co/tx/bsc/0x897c2de73dd55d7701e1b69ffb3a17b0f4801ced88b0c75fe1551c5fcce6a979



## Rari

The Rari Capital Ethereum Pool deposits ETH into Alpha Finance’s ibETH token as one of our yield-generating strategies. This strategy tracks the value of its ibETH as `ibETH.totalETH() / ibETH.totalSupply()`.

According to Alpha Finance, `ibETH.totalETH()` is vulnerable to manipulation inside the `ibETH.work` function, and a user of `ibETH.work` can call any contract it wants to inside `ibETH.work`, including the Rari Capital Ethereum Pool deposit and withdrawal functions.

The attacker repeatedly executed the following steps inside of `ibETH.work` (simplified explanation):
1. Flashloan ETH from dYdX.
2. Deposit that ETH into the Rari Capital Ethereum Pool.
3. Manipulate the value of `ibETH.totalETH()` by pushing it artificially high.
4. Withdraw more ETH from the Rari Capital Ethereum Pool than the attacker deposited because the Rari Capital Ethereum Pool’s balances are artifically inflated (because `ibETH.totalETH()` is artificially inflated).
5. At the end of `ibETH.work`, the value of `ibETH.totalETH()` returns to its true value, leading the Rari Capital Ethereum Pool’s balances to values lower than they were before the attack as a result of the attacker withdrawing more than they deposited while their balance was artificially inflated.



参考：

- https://medium.com/rari-capital/5-8-2021-rari-ethereum-pool-post-mortem-60aab6a6f8f9
- https://mp.weixin.qq.com/s/0Lwjf14hW5ahz3Om6jXRug



## DODO

资金池合约初始化函数没有任何鉴权以及防止重复调用初始化的限制，导致攻击者利用闪电贷将真币借出，然后通过重新对合约初始化将资金池代币对替换为攻击者创建的假币，从而绕过闪电贷资金归还检查将真币收入囊中



参考：

- https://mp.weixin.qq.com/s/1OrH7Ucqyt9sl7lkBBmb_g

## Paid Network

私钥泄露或者项目内部原因，导致任意 mint。

参考：

- https://mp.weixin.qq.com/s/iw4GdF1KbPmlQOm8Z3qrFA

##   SushiSwap

主要看这个即可：https://mp.weixin.qq.com/s/CUholEeD8AWL15psz1-9tQ











# 附录

## 经验

1. 报导有时候可能不准确，描述上可能省略部分细节。另外，由于调用栈非常深，文字描述有些难以体现调用的层次和顺序。
2. tenderly 是很好用的工具，结合 debugger 和 overview 能够较为清楚的了解调用顺序。但是这需要你熟悉 uniswap 的结构，清楚流动性、质押等相关内容。
3. 某些项目多次被攻击，搜索相关资料时可以限制时间段。

4. 工具网站：
   1. https://dashboard.tenderly.co/explorer
   2. https://tools.blocksec.com/tx

5. 漏洞报导主要来源：
   1. https://github.com/SunWeb3Sec/DeFiHackLabs
   2. https://www.slowmist.com/zh/news.html
   3. https://www.zhihu.com/org/ling-shi-ke-ji-43
   4. 相关项目的负责人，可能会在 medium 或者 推特上宣告漏洞细节。
   5. https://www.zhihu.com/org/cheng-du-lian-an-ke-ji-28/posts





## DeFi 涉及到的合约标准

### ERC165

判断合约是否实现某个接口，详情见 https://eips.ethereum.org/EIPS/eip-165#how-interfaces-are-identified



### ERC20

这是相对简单的代笔标准，存在时间非常久了，最主要的功能是：

- transfer tokens from one account to another
- get the current token balance of an account
- get the total supply of the token available on the network
- approve whether an amount of token from an account can be spent by a third-party account

主要方法有：

```solidity
function name() public view returns (string)
function symbol() public view returns (string)
function decimals() public view returns (uint8)
function totalSupply() public view returns (uint256)
function balanceOf(address _owner) public view returns (uint256 balance)
function transfer(address _to, uint256 _value) public returns (bool success)
function transferFrom(address _from, address _to, uint256 _value) public returns (bool success)
function approve(address _spender, uint256 _value) public returns (bool success)
function allowance(address _owner, address _spender) public view returns (uint256 remaining)
```

这些函数的含义很容易理解，特别注意：

- `transfer` 允许 value=0 的情况。
- `transferFrom` 和 `approve` 配套，`allowance` 返回授权金额。

代码实现：https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20.sol

参考：

- https://eips.ethereum.org/EIPS/eip-20
- https://ethereum.org/en/developers/docs/standards/tokens/erc-20/

### ERC721

基本介绍：We’ll use ERC721 to track items in our game, which will each have their own unique attributes. Whenever one is to be awarded to a player, it will be minted and sent to them. Players are free to keep their token or trade it with other people as they see fit, as they would any other asset on the blockchain!

存在拓展：The [`ERC721`](https://docs.openzeppelin.com/contracts/3.x/api/token/ERC721#ERC721) contract includes all standard extensions ([`IERC721Metadata`](https://docs.openzeppelin.com/contracts/3.x/api/token/ERC721#IERC721Metadata) and [`IERC721Enumerable`](https://docs.openzeppelin.com/contracts/3.x/api/token/ERC721#IERC721Enumerable)). That’s where the [`_setTokenURI`](https://docs.openzeppelin.com/contracts/3.x/api/token/ERC721#ERC721-_setTokenURI-uint256-string-) method comes from: we use it to store an item’s metadata.主要是元素据和添加 url

基本接口：



```solidity
pragma solidity ^0.4.20;

///  Note: the ERC-165 identifier for this interface is 0x80ac58cd.
interface ERC721 is ERC165 {

    event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId);
    event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId);
    event ApprovalForAll(address indexed _owner, address indexed _operator, bool _approved);
    function balanceOf(address _owner) external view returns (uint256);
    function ownerOf(uint256 _tokenId) external view returns (address);
    function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes data) external payable;
    function safeTransferFrom(address _from, address _to, uint256 _tokenId) external payable;
    function transferFrom(address _from, address _to, uint256 _tokenId) external payable;
    function approve(address _approved, uint256 _tokenId) external payable;
    function setApprovalForAll(address _operator, bool _approved) external;
    function getApproved(uint256 _tokenId) external view returns (address);
    function isApprovedForAll(address _owner, address _operator) external view returns (bool);
}

interface ERC165 {
    function supportsInterface(bytes4 interfaceID) external view returns (bool);
}
```

拓展：提供url

```solidity
interface ERC721Metadata is ERC721 {
    function name() external view returns (string _name);

    function symbol() external view returns (string _symbol);
    function tokenURI(uint256 _tokenId) external view returns (string);
}
```

拓展：

```solidity
/// @title ERC-721 Non-Fungible Token Standard, optional enumeration extension
/// @dev See https://eips.ethereum.org/EIPS/eip-721
///  Note: the ERC-165 identifier for this interface is 0x780e9d63.
interface ERC721Enumerable is ERC721  {
    function totalSupply() external view returns (uint256);
    function tokenByIndex(uint256 _index) external view returns (uint256);
    function tokenOfOwnerByIndex(address _owner, uint256 _index) external view returns (uint256);
}
```





注意：

- `from=0` 表示铸造，`to=0` 表示销毁，但是创建合约时不会触发 `transfer` 事件。转账后撤销对应 `approval`
- 已授权后，`Aproval` 还是可以多次调用，表示再确认。 
- 如果 `safeTransferFrom` 的接收者是合约，那么合约必须实现 `onERC721Received`， `bytes data` 指定调用转账接收合约的函数。重载的函数，相当于调用接收合约的 fallback 或者 receive。`transferFrom` 就没有这项检查
- `ERC165` 检查合约是否实现某个接口





参考：

- https://docs.openzeppelin.com/contracts/3.x/erc721
- https://eips.ethereum.org/EIPS/eip-721

### ERC1155

这是一个标准接口，支持开发同质化的、半同质化的、非同质化的代币和其他配置的通用智能合约。

```solidity
pragma solidity ^0.5.9;

/**
    @title ERC-1155 Multi Token Standard
    @dev See https://eips.ethereum.org/EIPS/eip-1155
    Note: The ERC-165 identifier for this interface is 0xd9b67a26.
 */
interface ERC1155 is ERC165 {
    event TransferSingle(address indexed _operator, address indexed _from, address indexed _to, uint256 _id, uint256 _value);
    
    event TransferBatch(address indexed _operator, address indexed _from, address indexed _to, uint256[] _ids, uint256[] _values);
    
    event ApprovalForAll(address indexed _owner, address indexed _operator, bool _approved);
    
    event URI(string _value, uint256 indexed _id);
    
    function safeTransferFrom(address _from, address _to, uint256 _id, uint256 _value, bytes calldata _data) external;

    function safeBatchTransferFrom(address _from, address _to, uint256[] calldata _ids, uint256[] calldata _values, bytes calldata _data) external;

    function balanceOf(address _owner, uint256 _id) external view returns (uint256);

    function balanceOfBatch(address[] calldata _owners, uint256[] calldata _ids) external view returns (uint256[] memory);

    function setApprovalForAll(address _operator, bool _approved) external;

    function isApprovedForAll(address _owner, address _operator) external view returns (bool);
}


pragma solidity ^0.5.9;

/**
    Note: The ERC-165 identifier for this interface is 0x4e2312e0.
*/
interface ERC1155TokenReceiver {

    function onERC1155Received(address _operator, address _from, uint256 _id, uint256 _value, bytes calldata _data) external returns(bytes4);

    function onERC1155BatchReceived(address _operator, address _from, uint256[] calldata _ids, uint256[] calldata _values, bytes calldata _data) external returns(bytes4);       
}
```

注意：

- ERC 1155 不同 id 表示不同的代币，这样用户可以拥有多种类型的代币。 `TransferBatch` 只是通过数组表达出 id[], amount[]，实现批量转 token

参考：

- https://eips.ethereum.org/EIPS/eip-1155

### ERC677

ERC677 在 token 进行转账之后，会回调到目标合约的 `onTokenTransfer` 方法. 简单的说就是对ERC20的补充，是的它能够向后兼容ERC223。

ERC677 除了包含了 ERC20 的所有方法和事件之外，增加了一个 `transferAndCall` 方法：

```javascript
function transferAndCall(address receiver, uint amount, bytes data) returns (bool success)
```

完成转账和记录日志之后，代币合约还会调用接收合约的 `onTokenTransfer` 方法，用来触发接收合约的逻辑。这就要就接收 ERC677 代币的合约必须实现 `onTokenTransfer` 方法，用来给代币合约调用。`onTokenTransfer` 方法定义如下：

```javascript
function onTokenTransfer(address from, uint256 amount, bytes data) returns (bool success)
```

接收合约就可以在这个方法中定义自己的业务逻辑，可以在发生转账的时候自动触发。换句话说，智能合约中的业务逻辑，可以通过代币转账的方式来触发自动运行。这就给了智能合约的应用场景有了很大的想象空间。

---

Allow tokens to be transferred to contracts and have the contract trigger logic for how to respond to receiving the tokens within a single transaction. This adds a new function to [ERC20](https://github.com/ethereum/EIPs/issues/20) token contracts, `transferAndCall` which can be called to transfer tokens to a contract and then call the contract with the additional data provided. Once the token is transferred, the token contract calls the receiving contract's function `onTokenTransfer(address,uint256,bytes)` and triggers an event `Transfer(address,address,uint,bytes)`, following the convention set in [ERC223](https://github.com/ethereum/EIPs/issues/223).

 ERC223 changes the behavior of ERC20's `transfer(address,uint256)`, specifying that it should throw if transferring to a contract that does not implement `onTokenTransfer`. This is problematic because there are deployed contracts in use that assume they can safely call `transfer(address,uint256)` to move tokens to their recipient. If one of these deployed contracts were to transfer an ERC223 token to a contract(e.g. a multisig wallet) the tokens would effectively become stuck in the transferring contract.

`transferAndCall` behaves similarly to `transfer(address,uint256,bytes)`, but allows implementers to gain the functionality without the risk of inadvertently locking up tokens in non-ERC223 compatible contracts. **It is distinct from ERC223's `transfer(address,uint256,bytes)` only in name**, but this distinction allows for easy distinguishability between tokens that are ERC223 and tokens that are simply ERC20 + ERC667.

参考：

- https://github.com/ethereum/EIPs/issues/677
- https://www.codeleading.com/article/18155158318/
- https://segmentfault.com/a/1190000022775687

### ERC1820





### ERC 2612

因为 uniswapERC20 用到了这一部分的内容，所以补上。

However, a limiting factor in this design stems from the fact that the ERC-20 `approve` function itself is defined in terms of `msg.sender`. This means that user’s *initial action* involving ERC-20 tokens must be performed by an EOA (*but see Note below*). If the user needs to interact with a smart contract, then they need to make 2 transactions (`approve` and the smart contract call which will internally call `transferFrom`). Even in the simple use case of paying another person, they need to hold ETH to pay for transaction gas costs.

**This ERC extends the ERC-20 standard with a new function `permit`, which allows users to modify the `allowance` mapping using a signed message, instead of through `msg.sender`.**



参考：

- https://eips.ethereum.org/EIPS/eip-2612



##  glossary



参考：

- https://www.futurelearn.com/links/f/2xqu0bc09qdqou8wmdjgtn57qqmbxtc



## auction

拍卖是 defi 中常遇到的基本组件，所以有必要进行一定的学习。

简单介绍：Auctions are platforms for selling goods in a public forum through open and competitive bidding. Commonly, the auction winner is the bidder who submitted the highest price, however, there are a variety of other rules to determine the winner. **The auctioneer leads the auction according to the rules and always charge a fee to the vendor for his services, usually a percentage of the gross selling price of the good**.

总而言之，需要明确三方身份 auctioneer, seller. buyer

拍卖的分类有如下：

- English auction, also called open ascending price auction, where participants bid openly against one another, with each subsequent bid required to be higher than the previous bid.
- Dutch auction, also called open descending price auction, used for perishable commodities such as ﬂowers, ﬁsh and tobacco. In the traditional Dutch auction, the auctioneer begins with a high price and then continously lowers it until a bidder is willing to accept the auctioneer’s price, or until the seller’s reserve price is met.
- First-price sealed-bid auction (FPSB), also called blind auction, where all bidders submit bids in a sealed envelop to the auctioneer, so that no bidder knows the bid of any other participant. Later, the auctioneer opens the envelope to determine the winner who submitted the higher bid.
- Vickrey auction, also known as Second-price sealed-bid auction (SPSB), which is identical to the ﬁrst-price sealed-bid auction, with the exception that the winning bidder pays the second-highest bid rather than his own.

实际拍卖可能还有时间限制，拍卖者撤回，拍卖者底价等因素。



参考：

1. Braghin, Chiara & Cimato, Stelvio & Damiani, Ernesto & Baronchelli, Michael. (2020). Designing Smart-Contract Based Auctions. 10.1007/978-3-030-16946-6_5. 

## Liquidations

基本介绍：Loans on a blockchain typically operate as follows. Lenders with a surplus of money provide assets to a lending smart contract. Borrowers then provide a security deposit, known as collateral, to borrow cryptocurrency.

当抵押物信用不够时(e.g., below 150% of the debt value )，就会采用如下三种方式收回债务:

1. a loan can be made available for liquidation by the smart contract. Liquidators then pay back the debt in exchange for receiving the collateral at a discount (i.e., liquidation spread), or the collateral is liquidated through an auction.
2. Debt can also be rescued by “topping up” the collateral, such that the loan is sufficiently collateralized.
3. Finally, the borrower can repay parts of their debt.



参考：

1. Kaihua Qin, Liyi Zhou, Pablo Gamito, Philipp Jovanovic, and Arthur Gervais. 2021. An Empirical Study of DeFi Liquidations: Incentives, Risks, and Instabilities. In ACM Internet Measurement Conference (IMC ’21), November 2–4, 2021, Virtual Event, USA. ACM, New York, NY, USA, 15 pages. https://doi.org/10.1145/3487552.3487811

## uniswap

### 基本原理

看官方文档即可：https://docs.uniswap.org/protocol/V2/concepts/protocol-overview/how-uniswap-works

### 代码实现

Uniswap has 4 smart contracts in total. They are divided into **core** and **periphery**.

1. **Core** is for storing the funds (the tokens) and exposing functions for swapping tokens, adding funds, getting rewards, etc.
2. **Periphery** is for interacting with the **core**.

**Core** consists of the following smart contracts:

<img src="https://miro.medium.com/max/1400/1*WbCK5HMsPexKYuZxujx1Rg.png" alt="img" style="zoom:50%;" />

1. **Pair** — a smart contract that implements the functionality for swapping, minting, burning of tokens. This contract is created for every exchange pair like *Dogecoin ↔ Shiba*.
2. **Factory** — creates and keeps track of all Pair contracts
3. **ERC20** — for keeping track of ownership of pool. Think of the pool as a property. When liquidity providers provide funds to the pool, they get “pool ownership tokens” in return. These ownership tokens earn rewards (by traders paying a small percentage for each trade). When liquidity providers want their funds back, they just submit the ownership tokens back and get their funds + the rewards that were accumulated. The **ERC20** contract keeps track of the ownership tokens.

**Periphery** consists of just one smart contract:

1. **Router** is for interacting with the core. Provides functions such as `swapExactETHForTokens`, `swapETHForExactTokens`, etc.

The main functionality is the following:

1. **Managing the funds** (how tokens such as Dogecoin and Shiba are managed in the pool)
2. Functions for **liquidity providers —** deposit more funds and withdraw the funds along with the rewards
3. Functions for **traders** — swapping
4. **Managing pool ownership tokens**
5. **Protocol fee** — Uniswap v2 introduced a switchable protocol fee. This protocol fee goes to the Uniswap team for their efforts in maintaining Uniswap. At the moment, this protocol fee is turned off but it can be turned on in the future. When it’s on, the traders will still pay the same fee for trading but 1/6 of this fee（1/6 of 0.3%） will now go to the Uniswap team and the rest 5/6 will go to the liquidity providers as the reward for providing their funds.

Another functionality that is not core to Uniswap but it’s a useful helper for other contracts in the Ethereum ecosystem: 

**Price oracle** — Uniswap tracks prices of tokens relative to each other and can be used as a price oracle for other smart contracts in the Ethereum ecosystem. 

#### pair

The Pair contract implements the exchange between a pair of tokens such as Dogecoin and Shiba. The full code of the Pair smart contract can be found on Github under [v2-core/contracts/UniswapV2Pair.sol](https://github.com/Uniswap/v2-core/blob/master/contracts/UniswapV2Pair.sol)

pair 中的 token 和 ERC20 的关系：

<img src="https://miro.medium.com/max/1400/1*j4T8KVmNZhDekIC8O_zNtw.png" alt="img" style="zoom:50%;" />

更多的见参考文章，写的很不错。

最推荐看白皮书以及下面的代码详解。

https://ethereum.org/es/developers/tutorials/uniswap-v2-annotated-code/

讲解视频：https://www.youtube.com/watch?v=bxV0OKPz-G4





参考：

- https://zhuanlan.zhihu.com/p/255190320
- https://zhuanlan.zhihu.com/p/269205336
- https://hackmd.io/@HaydenAdams/HJ9jLsfTz#%F0%9F%A6%84-Uniswap-Whitepaper
- https://betterprogramming.pub/uniswap-smart-contract-breakdown-ea20edf1a0ff
- https://betterprogramming.pub/uniswap-smart-contract-breakdown-part-2-b9ea2fca65d1



## Flash Loans

闪电贷：一种不需要用户抵押资金的贷款，但是必须在发放资金的**同一交易中**偿还贷款人。

可以参考



参考：

- https://zhuanlan.zhihu.com/p/360131349

- https://www.tuoluo.cn/article/detail-9991688.html

- [实践](https://www.alchemy.com/overviews/creating-a-flash-loan-using-aave#flash-i)



## DEX Aggregators

参考：https://www.lcx.com/role-of-dex-aggregators-in-defi/

写的很好。





## 代理合约

登链社区的[帖子](https://learnblockchain.cn/article/1102)写的很好，足够简单了解了。

大部分的proxy pattern里 合约状态（存储内容什么的）都在proxy合约里，implementation合约除了代码没有任何状态。



参考：

- https://learnblockchain.cn/article/1102





## blockchain bridge

A blockchain bridge, otherwise known as a cross-chain bridge, connects two blockchains and allows users to send cryptocurrency from one chain to the other. Basically, if you have bitcoin but want to spend it like Ethereum, you can do that through the bridge. Blockchain bridges solve this problem by enabling token transfers, smart contracts and data exchange, and other feedback and instructions between two independent platforms. 

This concept is a lot similar to Layer 2 solutions even though the two systems have different purposes. Layer 2 is built on top of an existing blockchain so while it does improve speed, the lack of interoperability remains. Cross-chain bridges are also independent entities that don’t belong to any blockchain.

Blockchain bridges can do a lot of cool stuff like converting smart contracts and sending data, but the most common utility is token transfer. When you have bitcoin and want to transfer some of it to Ethereum, the blockchain bridge will hold your coin and create equivalents in ETH for you to use. None of the crypto involved actually moves anywhere. Rather, the amount of BTC you want to transfer gets locked in a smart contract while you gain access to an equal amount of ETH. When you want to convert back to BTC, the ETH you had or whatever’s left of it will get burned and an equal amount of BTC goes back to your wallet.





拓展链接：

- https://blog.liquid.com/blockchain-cross-chain-bridge
- https://finance.sina.com.cn/blockchain/roll/2022-06-14/doc-imizirau8337864.shtml

## 爆仓

所谓爆仓，是指投资者保证金账户中的客户权益为负值。在市场行情发生较大变化时，如果投资者保证金账户中资金的绝大部分都被交易保证金所占用，而且交易方向又与市场走势相反时，由于保证金交易的杠杆效应，就很容易出现爆仓。

爆仓的发生实际上是投资者资金链断裂的结果。为避免这种情况的发生，需要特别控制好头寸，合理地进行资金管理，切忌像股票交易中可能出现的满仓操作；并且与股票交易不同，投资者必须对股指期货行情进行及时跟踪。如果爆仓导致了亏空且由投资者的原因引起，投资者需要将亏空补足，否则会面临法律追索。

总而言之，杠杠高，借钱投资，需要保证金，如果因为投资亏损，保证金不足，就会爆仓。

**实例分析：**

某客户帐户原有保证金200000元，8月9日，[开仓](https://baike.baidu.com/item/开仓)买进9月[沪深300指数期货合约](https://baike.baidu.com/item/沪深300指数期货合约/12803638)15手，均价1200点（每点100元），手续费为单边每手10元，[当日结算价](https://baike.baidu.com/item/当日结算价)为1195点，[保证金比例](https://baike.baidu.com/item/保证金比例)为8%。

当日开仓持仓盈亏=（1195-1200）×15×100=-7500元

手续费=10×15=150元

当日权益=200000-7500-150=192350元

保证金占用=1195×15×100×8%=143400元

资金余额（即可交易资金）=192350-143400=48950元

8月10日，该客户没有交易，但9月沪深300指数期货合约的当日结算价降为1150点，当日账户情况为：

历史持仓盈亏=（1150-1195）×15×100=-67500元

当日权益=192350-67500=124850元

保证金占用=1150×15×100×8%=138000元

资金余额（即可开仓交易资金）=124850-138000=-13150元

8月11日开盘前，客户没有将应该追加的保证金交给期货公司，而9月[股指期货合约](https://baike.baidu.com/item/股指期货合约)又以[跳空](https://baike.baidu.com/item/跳空)下跌90点以1060点开盘并继续下跌。[期货经纪公司](https://baike.baidu.com/item/期货经纪公司)将该客户的15手[多头](https://baike.baidu.com/item/多头)持仓强制[平仓](https://baike.baidu.com/item/平仓)，[成交价](https://baike.baidu.com/item/成交价)为1055点。这样，该账户的情况为：

当日[平仓盈亏](https://baike.baidu.com/item/平仓盈亏)=（1055-1150）×15×100=-142500元

手续费=10×15=150元

实际权益=124850-142500-150=-17800元

即该客户倒欠了期货经纪公司17800元。

发生爆仓时，投资者需要将亏空补足，否则会面临法律追索。为避免这种情况的发生，需要特别控制好仓位，切忌像[股票交易](https://baike.baidu.com/item/股票交易)那样满仓操作。并且对行情进行及时跟踪，不能像股票交易那样一买了之。因此[期货](https://baike.baidu.com/item/期货)实际并不适合任何投资者来做。

参考：

- https://baike.baidu.com/item/%E7%88%86%E4%BB%93%E7%8E%B0%E8%B1%A1/936585
- https://www.hfzq.com.cn/review_cms_a907173b-54ac-42d6-b7ff-641902e37b01.shtm
