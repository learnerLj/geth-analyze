## 应用二进制接口(ABI)

在地址类型的介绍中提到了底层调用的 call函数，这里将会介绍它的用法，以及和ABI函数的配合。

ABI 全名 Application Binary Interface，翻译为应用二进制接口。它定义了与合约交互的规范，因此底层函数 (如 call) 直接给合约发消息前，需要了解 ABI。

ABI 是由合约生成的，规定与合约交互方式的规则，它是一个接口，常被 web3 等库调用。熟悉 REST API 的读者应该能很快理解。

### 接口含义

例如下面是智能合约：

```solidity
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

/**
 * @title Owner
 * @dev 设置和改变所有者
 */
contract Owner {

    address private owner;
    
    // 设置事件
    event OwnerSet(address indexed oldOwner, address indexed newOwner);
    
    // 函数修饰器，限制调用者必须是所有者
    modifier isOwner() {
        // 如果require 第一个参数为false，就回滚，并且日志中包含作为错误信息的第二个参数。它常用于限制合约调用是否合法。
        require(msg.sender == owner, "Caller is not owner");
        _;
    }
    
    /**
     * @dev 构造函数默认部署者为所有者
     */
    constructor() {
        owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
        emit OwnerSet(address(0), owner);
    }

    
    function changeOwner(address newOwner) public isOwner {
        emit OwnerSet(owner, newOwner);
        owner = newOwner;
    }

    function getOwner() external view returns (address) {
        return owner;
    }
}
```

它的 ABI 如下：

```json
[
			{
				"inputs": [],
				"stateMutability": "nonpayable",
				"type": "constructor"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "oldOwner",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "newOwner",
						"type": "address"
					}
				],
				"name": "OwnerSet",
				"type": "event"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "newOwner",
						"type": "address"
					}
				],
				"name": "changeOwner",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "getOwner",
				"outputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"stateMutability": "view",
				"type": "function"
			}
		]
```

#### 对于函数：

- `type`: `"function"`, `"constructor"` (可以省略，默认是 `"function"`; 也可能是`"fallback"` );
- `name`: 函数的名字;
- `constant`: `true` 表示该函数调用不修改区块链状态，只读或者只生成调用后销毁的 memory 变量;
- `payable`: `true`表示可以接收以太币, 默认是 `false`;
- `stateMutability`:四种结果，`pure` (不读取也不修改状态), `view` (不修改状态，和上面的 `constant` 是等价的), `nonpayable` and `payable` (否，是接收以太币 );
- `inputs`: 对象的数组，包括:
  - `name`: 参数的名字
  - `type`: 参数的类型
- `outputs`: 和 `inputs` 类似, 如果没有输出可为空.



**对于下面这个函数的解读：**

input 输入变量是内置类型 (internalType) 中的地址类型 (address)，类型 (type) 是 地址 (address)。

output 返回值为空。

该函数的属性标记 (stateMutability) 是不可转账 (nonpayable) ，类型是函数 (function)

			{
				"inputs": [
					{
						"internalType": "address",
						"name": "newOwner",
						"type": "address"
					}
				],
				"name": "changeOwner",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},

#### 对于事件:

- `type`:  `"event"`
- `name`:事件名字;
- `inputs`:对象的数组，包括:
  - `name`: 参数的名字
  - `type`: 参数的类型
  - `indexed`: `true` 表示是特殊结构`topics`的一部分（见事件的 indexed 修饰）, `false` 表示日志文件.
- `anonymous`: `true` 表示事件被声明为 `anonymous`.

对下面这个事件的解读：

非匿名，输出参数有两个，花括号类型标记，其余与函数差别不大。

			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "oldOwner",
						"type": "address"
					},
					{
						"indexed": true,
						"internalType": "address",
						"name": "newOwner",
						"type": "address"
					}
				],
				"name": "OwnerSet",
				"type": "event"
			},

至于全局变量`address private owner`，由于 private 限制访问，所以不在 ABI 中。

从上面的格式中，可以看到 ABI 和编码很相关，发送的数据应当 ABI 的方式组织，同样的也需要对应的编码格式。

### 函数选择器

函数选择器也和接口强相关，因为在 call 之类的底层调用中，需要根据函数签名匹配函数。函数调用依靠调用数据的前四个字节匹配函数，这四个字节是函数签名的哈希值的前四个字节。调用函数的数据的编码格式按照顺序是函数名和带圆括号的参数类型列表，参数之间只靠逗号分隔。注意函数返回值并不参与哈希，这样可以进一步解耦，更灵活地重载。

详细选择函数过程需要深入执行过程，可见[博客](https://medium.com/@hayeah/how-to-decipher-a-smart-contract-method-call-8ee980311603)，后面我也会学习。

### ABI 的参数编码

**编码规则如下**：

设 `X` 是编码前的值，对于静态类型 `a` (内置的类型)，定义 `len(a)`是 `a` 转化成二进制数后的位数(注：所有类型底层最终是由二进制数表示)；对于动态类型 `a` (如数组、元组，`bytes`、`string`、`T[k]`)，我们常用编码后的长度 `len(enc(a))`。`enc()` 是我们定义的函数，它输入参数是类型 (包括静态类型和动态类型)，返回值是二进制序列 (一串二进制数，但是含义不在于数值而是字符排列顺序)。

我们的核心就在于如何定义编码函数 `enc()`。首先设定编码的基本格式

1. 对于元组，表示如下。不同的 `head()` 放在一块，表示二进制代码直接拼接。

`enc(X) = head(X(1)) ... head(X(k)) tail(X(1)) ... tail(X(k))` ，函数的参数列表就是元组。

定义 `head` 和 `tail` 如下：

- 若 `X` 是静态类型，`head(X(i)) = enc(X(i))` ， `tail(X(i)) = ""` （空字符串）。因为静态类型是唯一的，可以直接编码，无需额外的参数说明。
- 若 `X` 是动态类型: ``head(X) = enc(len(head(X) tail(X)))``，``tail(X(i)) = enc(X(i))``，即需要在实际编码值前面添加编码后的长度。一来方便读取，也明确了动态的类型的确切类型（如 `T` 类型的数组 `T[k]`，确切类型是长度为 `k`，类型为 `T` 的数组）。

2. 对于一般变量的编码规则

   1. `T[k]` 对于任意 `T` 和 `k`：数组当作同类型变量凑在一起的元组

      `enc(X) = enc((X[0], ..., X[k-1]))`

   2. `T[]` 当 `X` 有 `k` 个元素（`k` 默认是类型 `uint256`）：不定长数组多了元素的个数作为前缀。

      `enc(X) = enc(k) enc([X[1], ..., X[k]])`

   3. 具有 `k` （呈现为类型 `uint256`）个字节的 `bytes`：不定长数组多了元素的个数作为前缀，如果 `byte` 数组，然后直接抄下来。

      `enc(X) = enc(k) pad_right(X)`，`pad_right(x)` 的意思是在左边把原来的字节序列添加上去，填充在右边，`enc(k)` 参靠下面 `uint` 类型的编码方式。

   4. `string`：先把 `string` 类型转成 `bytes` 类型，注意 `k` 是指转化后的字节数。

      `enc(X) = enc(enc_utf8(X))`，`enc_utf8(X)` 指将`string` 类型转成 `bytes` 类型。

   5. `uint<M>`：`enc(X)` 是在 `X` 的高位补充若干 0 直到长度为 32 字节。

   6. `address`：与 `uint160` 的情况相同。

   7. `int<M>`：`enc(X)` 是在 `X` 补码的高位添加若干字节，直到长度为 32 字节；

      - 如果  `X` 是负数，在高位一直添加值为 `0xff` （16 进制转二进制，实际上就是 8 个 1。注意  `int` 和 `uint` 这两类的位数都是 8 的倍数）
      - 对于 `X` 是非负数，在高位添加 `0x00` 值（即 8 位全为 0），直到为32个字节。

   8. `bool`：与 `uint8` 的情况相同，`1` 用来表示 `true`，`0` 表示 `false`。

   9. `fixed<M>x<N>`：`enc(X)` 就是 `enc(X * 10**N)`，其中 `X * 10**N` 可以理解为 `int256`。

   10. `fixed`：与 `fixed128x18` 的情况相同。

   11. `ufixed<M>x<N>`：`enc(X)` 就是 `enc(X * 10**N)`，其中 `X * 10**N` 可以理解为 `uint256`。

   12. `ufixed`：与 `ufixed128x18` 的情况相同。

   13. `bytes<M>`：`enc(X)` 就是 `X` 的字节序列加上为使长度成为 32 字节而添加的若干 0 值字节。

### ABI编码的实际格式

在编写合约时，或者需要验证签名时，对于多个参数往往需要打包多个参数。调用函数时发送的数据从第五个字节开始就是参数的编码了。格式如下：``function_selector(f) enc((a_1, ..., a_n))``。

函数选择器即函数签名的 keccak256 哈希，函数签名即是我们写代码时见到的函数名和参数类型，如``myFunc(string, uint8, uint8, uint8)``。函数选择器即这个哈希值的前 4 个字节。具体可以用 web3 库或者合约中的密码学函数实现。

**具体例子：**

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;

contract Foo {
    function bar(bytes3[2] memory) public pure {}
    function baz(uint32 x, bool y) public pure returns (bool r) { r = x > 32 || y; }
    function sam(bytes memory, bool, uint[] memory) public pure {}
}
```

可得 `baz` 函数选择器：`0xcdcd77c0`

参数都是静态类型，所以  

```
enc((uint32 39,bool true)) 
= head(uint32) head(bool) tail(uint32) tail(bool)
= enc(uint32 39) enc(bool true) 
= enc(uint32 39) enc(uint32 1)
=0x0000000000000000000000000000000000000000000000000000000000000045+0x0000000000000000000000000000000000000000000000000000000000000001 //注释：注意实际没有加号+，这里为了方便说明两段是如何连接的

最终 ABI 编码如下：0xcdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001`
```

更多的例子可见[文档](https://docs.soliditylang.org/en/latest/abi-spec.html#examples)。

### 事件

事件是从 EVM 日志中提取出来的片段，为了方便解析它，事件类似于函数的 ABI 。事件有事件名和参数列表，编码时把参数列表分成两份，一份是带有 `indexed` 标识的参数列表（对于非匿名事件里面至多三个参数，对于匿名事件至多4个参数，在编写合约时也有这样的限制），另一部分则是无这个标识的参数列表。标有 `indexed` 的参数列表会和事件签名的 Keccak 哈希共同构成日志中的特殊数据结构 `topics`(这种数据结构便于检索)。无 `indexed` 标识的参数列表会根据普通类型的编码规则，生成序列。

详细地，事件的结构如下：

- `address`：由 EVM 自动提供的事件所在合约的地址；

- `topics[0]`：`keccak(事件名+"("+EVENT_ARGS.map(canonical_type_of).join(",")+")")`

  `EVENT_ARGS.map(canonical_type_of).join(",")` 表示事件的每个参数对应的类型，类型之间用逗号分开。例如，对 `event myevent(uint indexed foo,int b)`，那么`topics[0]=keccak(myevent(uint,int))`。

  如果事件被声明为 `anonymous`，那么 `topics[0]` 不会被生成；

- `topics[n]`：`EVENT_INDEXED_ARGS[n - 1]` 

  `EVENT_INDEXED_ARGS[n-1] ` 是带有 `indexed` 标识的参数列表中下标为 `n-1` 的参数，即第 `n` 个参数；这个式子表示每个 topics 结构里面的内容是什么。

- `data`：`abi_serialise(EVENT_NON_INDEXED_ARGS)` 

  `EVENT_NON_INDEXED_ARGS` 是不带有 `indexed` 标识的事件参数，`abi_serialise()` 把参数列表 ABI 编码，相当于前面提到的 `enc()` 编码函数。

​	

关于设计原理的说明：

对于复杂的类型（超过 32 个字节或者时动态类型），比如结构体、`bytes`、`string`，编码的前面有 `keccak` 能够高效的检索这样的类型变量，但是也增加了解析的复杂度。因此，要精心设计将需要检索的参数标上 `indexed`，不需要检索而定位后直接获取的变量就不带 `indexed`。当然也可以制造冗余，每个变量都设置一个 带 `indexed` 的参数和不带 `indexed` 的参数，但是部署合约时 gas 会更高，调用消耗也会更高。

### 错误处理的编码

当人为设置回滚函数时，有时需要提供错误的描述性信息。这样的参数也会参与 ABI 的编码。

如下所示，`error` 是新增的类型，和事件类似，但是用于 `revert` 操作，提供提示信息。这里触发的 `error` 类型，将会以函数的编码方式编码，在以后可能会改变，将 `error` 的函数选择器，改成 错误选择器 error selector，固定为 四个字节的全 0  (`0x00000000`) 或者全 1  (`0xffffffff`) 。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

contract TestToken {
    error InsufficientBalance(uint256 available, uint256 required);
    function transfer(address /*to*/, uint amount) public pure {
        revert InsufficientBalance(0, amount);
    }
}
```
