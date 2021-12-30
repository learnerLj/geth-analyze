任务清单

- [ ] 理解合约部署的详细信息。
- [ ] 阅读完 oracle 部分源码
- [ ] 开始阅读 account 的 abi 部分，理解账号和合约的数据结构以及源码中的封装。



# contract部分

### 源码中如何部署合约

传入的参数有部署合约的交易，合约的 ABI，合约的字节码，与合约交互的方法

```GO
func DeployContract(opts *TransactOpts, abi abi.ABI, bytecode []byte, backend ContractBackend, params ...interface{}) (common.Address, *types.Transaction, *BoundContract, error) {
	// Otherwise try to deploy the contract
	c := NewBoundContract(common.Address{}, abi, backend, backend, backend)

	input, err := c.abi.Pack("", params...)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	tx, err := c.transact(opts, nil, append(bytecode, input...))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	c.address = crypto.CreateAddress(opts.From, tx.Nonce())
	return c.address, tx, c, nil
}
```



## Oracle

​	oracle 翻译是预言机，英文中的意思是预卜先知，知晓消息的意思。在区块链里用于合约获取链外的数据。例如你想把比特币转换成美元，如果在链上进行，那么就需要从链外获取比特币和美元的汇率，例如[price feed oracles](https://developer.makerdao.com/feeds/)。但是以太坊是封闭的系统，直接与外界交互很容易破坏 EVM 安全性，因此才用了预言机作为中间层，沟通链上和链外。详细可见[chainlink的文档](https://chain.link/education/blockchain-oracles)和[官方文档](https://ethereum.org/en/developers/docs/oracles/)。

​	在以太坊上，oracle 是已经部署的智能合约和链外组件，它可以查询 API 提供的信息，然后给其他合约发消息，更新合约的数据。但是只相信唯一的数据源也是很不可靠的方式，通常是多个数据源。我们可以自己创建，也可以直接使用服务商提供的服务。

​	一般 oracle 机制如下：

1. 到了需要链外数据的时候，合约触发事件。
2. 链外的接口监听事件的日志。
3. 链外接口处理事务，然后交易的方式返回数据给合约。

​	实际上，多个可信的数据来源在链上处理是比较耗费 gas 的，因此提出了通过密码学手段，在链外汇总数据，然后发给合约。

### oracle 的例子

​	下面是一个例子，从网络导入合约库，获取接口信息，然后创建合约类型 `AggregatorV3Interface` 的变量 `priceFeed`，然后结合获取的接口信息，在构造函数里创建在特定地址已经部署好的合约实例，调用函数`priceFeed.latestRoundData()`，返回的是元组，因此用多个数据接收。这样就获得了最新的 ETH 和 USD 的汇率。而我们导入的合约`priceFeed` 以及它在链外的配套接口，被称作预言机 oracle。类似的，我们也可以通过 oracle 解决链上难以产生可靠的随机数的问题。  

​	更多的例子可以看 chainlink 这些提供商，提供的文档，详细地说明了流程。

```js
// This example code is designed to quickly deploy an example contract using Remix.

pragma solidity ^0.6.7;

import "https://github.com/smartcontractkit/chainlink/blob/master/evm-contracts/src/v0.6/interfaces/AggregatorV3Interface.sol";

contract PriceConsumerV3 {

    AggregatorV3Interface internal priceFeed;

    /**
     * Network: Kovan
     * Aggregator: ETH/USD
     * Address: 0x9326BFA02ADD2366b30bacB125260Af641031331
     */
    constructor() public {
        priceFeed = AggregatorV3Interface(0x9326BFA02ADD2366b30bacB125260Af641031331);
    }

    /**
     * Returns the latest price
     */
    function getLatestPrice() public view returns (int) {
        (
            uint80 roundID, 
            int price,
            uint startedAt,
            uint timeStamp,
            uint80 answeredInRound
        ) = priceFeed.latestRoundData();
        return price;
    }
}
```



## checkpoint

​	检查点实际上是一个标记，用于确认这个状态和之前的状态是可信的。在区块链上，检查点往往是有足够的可信实体共同签名后，正式生成。它意味着检查点的状态是不可逆的，无条件可信的。这也是区块链防止造假的手段之一。

