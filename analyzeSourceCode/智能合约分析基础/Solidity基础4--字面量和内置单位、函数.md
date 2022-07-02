## 字面常量

### 地址字面常量

地址的字面常量有EIP-55的标准（区分大写字母和小写字母来验证），只有经过校验后才能作为address变量。

### 有理数和整数的字面常量

十进制的小数字面常量都会带一个`.`比如`1.` `.1` `1.2`。。`.5*8`的结果是整型的`4`。

也支持科学计数法，但是指数部分需要是整数（防止无理数出现）。`2e-10`

为了提高可读性，数字之间可以添加下划线，编译器会自动忽略。但是下划线不允许连续出现，而且只能出现在数字之间。`1_2e345_678`

数值字面常量支持任意精度和长度，也支持对应类型的所有运算，字面常量的运算结果还是字面常量。但是如果出现在变量表达式就会发生隐式转化，并且不同类型的字面常量和变量运算不能通过编译。也就是说`2**800+1-2**800`在字面常量中是允许的

（在`0.4.0`之前，整数的字面常量会被截断即`5/2=2`，但是之后是`2.5`)

**字面类型的运算还是字面常量，和非字面常量运算，就会隐式转化成普通类型**

### 字符串字面常量及类型

字符串字面常量（"foo"或者'bar'这样)，可以分段写(`"foo""bar"`等效为`"foobar"`)

字符串`"foo"`相当于3个字节，而不是4个字节，它不像C语言里以`\0`结尾。

字符串字面常量可以隐式的转换成`bytes1`,...`bytes32`。在合适的情况下，可以转换成`bytes`以及`string`

字符串字面常量只包含可打印的ASCII字符和下面的转义字符：

- `\<newline>` (转义实际换行)
- `\\` (反斜杠)
- `\'` (单引号)
- `\"` (双引号)
- `\b` (退格)
- `\f` (换页)
- `\n` (换行符)
- `\r` (回车)
- `\t` (标签 tab)
- `\v` (垂直标签)
- `\xNN` (十六进制转义，表示一个十六进制的值，)
- `\uNNNN` (unicode 转义，转换成UTF-8的序列)

### 十六进制的字面常量

十六进制字面常量以关键字 `hex` 打头，后面紧跟着用单引号或双引号引起来的字符串（例如，`hex"001122FF"` ）。 字符串的内容必须是一个十六进制的字符串，它们的值将使用二进制表示。

用空格分隔的多个十六进制字面常量被合并为一个字面常量： `hex"00112233" hex"44556677"` 等同于 `hex"0011223344556677"`

## 内置函数和变量

### 单位

- 币的单位默认是`wei`，也可以添加单位。

```
1 wei == 1;
1 gwei == 1e9;
1 ether == 1e18;
```

- 时间单位，默认是秒。但是需要注意闰秒和闰年的影响，这里的统计的时间并不是完全和日历上的时间相同。

```
1 == 1 seconds`
1 minutes == 60 seconds
1 hours == 60 minutes
1 days == 24 hours
1 weeks == 7 days
```

### 区块和交易的属性

括号内表示返回值类型

- `blockhash(uint blockNumber) returns (bytes32)`：指定区块的区块哈希，但是仅可用于最新的 256 个区块且不包括当前区块，否则返回0.
- `block.chainid` (`uint`): 当前链 id
- `block.coinbase` ( `address` ): 挖出当前区块的矿工地址
- `block.difficulty` ( `uint` ): 当前区块难度
- `block.gaslimit` ( `uint` ): 当前区块 gas 限额
- `block.number` ( `uint` ): 当前区块号
- `block.timestamp` ( `uint`): 自 unix epoch 起始当前区块以秒计的时间戳
- `gasleft() returns (uint256)` ：剩余的 gas
- `msg.data` ( `bytes` ): 完整的 calldata
- `msg.sender` ( `address` ): 消息发送者（当前调用）
- `msg.sig` ( `bytes4` ): calldata 的前 4 字节（也就是函数标识符）
- `msg.value` ( `uint` ): 随消息发送的 wei 的数量
- `tx.gasprice` (`uint`): 交易的 gas 价格
- `tx.origin` (`address payable`): 交易发起者（完全的调用链）

注意几大变化：

- `gasleft`原来是`msg.gas`
- `block.timestamp`原来是`now`
- `blockhash`原来是`block.blockhash`

### delete

`delete a`不是常规意义上的删除，而是给`a`赋默认值（即返回不带参数的声明的状态），比如`a`是整数，那么等同于`a=0`。对用动态数组是将数组长度变为0；对于静态数组是将每一个元素初始化；对于结构体就把每一个成员初始化；对于映射在原理上无效（不会影响映射关系），但是会删除其他的成员，如值。

```javascript
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.0 <0.9.0;

contract DeleteExample {
    uint data;
    uint[] dataArray;

    function f() public {
        uint x = data;
        delete x; // sets x to 0, does not affect data
        delete data; // sets data to 0, does not affect x
        uint[] storage y = dataArray;
        delete dataArray; // this sets dataArray.length to zero, but as uint[] is a complex object, also
        // y is affected which is an alias to the storage object
        // On the other hand: "delete y" is not valid, as assignments to local variables
        // referencing storage objects can only be made from existing storage objects.
        assert(y.length == 0);
    }
}
```

### ABI 编码及解码函数

详细原理和应用见下一节 应用二进制接口，需要明白 ABI 编码的含义才懂这些函数的用法。

- `abi.decode(bytes memory encodedData, (...)) returns (...)`: 对给定的数据进行ABI解码，而数据的类型在括号中第二个参数给出 。 例如: `(uint a, uint[2] memory b, bytes memory c) = abi.decode(data, (uint, uint[2], bytes))`
- `abi.encode(...) returns (bytes)`： 对给定参数进行编码

例如：

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.5.0 <0.9.0;
contract A {
    bytes public c = abi.encode(5,-1);
    uint public d;//5
    int public e;//-1
    
    constructor(){
        (d,e) = abi.decode(c,(uint,int));
    }
}
```



- `abi.encodePacked(...) returns (bytes)`：对给定参数执行 [紧打包编码](https://learnblockchain.cn/docs/solidity/abi-spec.html#abi-packed-mode) （即编码时不够 32 字节，不用补0了。
- `abi.encodeWithSelector(bytes4 selector, ...) returns (bytes)`： [ABI](https://learnblockchain.cn/docs/solidity/abi-spec.html#abi) - 对给定第二个开始的参数进行编码，并以给定的函数选择器作为起始的 4 字节数据一起返回
- `abi.encodeWithSignature(string signature, ...) returns (bytes)`：等价于 `abi.encodeWithSelector(bytes4(keccak256(signature), ...)`

用法如下：

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.5.0 <0.9.0;
contract A {
    uint public a;

    function add(uint b,uint c) public returns (uint) {
        a=a+b+c;
        return a;
    }

    bytes public encodedABI=abi.encodeWithSelector(this.add.selector,5,1);
    
    function callFunc()public returns(bool,bytes memory,uint) {
        bool flag=false;
        bytes memory result;
       (flag,result) = address(this).call(encodedABI);
       return (flag,result,a);
    }
    fallback() external{

    }
}
```



### 错误处理

`assert(bool condition)`,`require(bool condition)`,`require(bool condition, string memory message)`均是条件为假然后回滚。

`revert()`,`revert(string memory reason)` 立即回滚。

### 数学和密码学函数

- `addmod(uint x, uint y, uint k) returns (uint)`: 计算 `(x + y) % k`，加法会在任意精度下执行，并且加法的结果即使超过 `2**256` 也不会被截取。从 0.5.0 版本的编译器开始会加入对 `k != 0` 的校验（assert）。
- `mulmod(uint x, uint y, uint k) returns (uint)`: 计算 `(x * y) % k`，乘法会在任意精度下执行，并且乘法的结果即使超过 `2**256` 也不会被截取。从 0.5.0 版本的编译器开始会加入对 `k != 0` 的校验（assert）。
- `keccak256((bytes memory) returns (bytes32)`: 计算 Keccak-256 哈希。0.5.0以前有别名`sha3`.它一般用于：生成输入信息的独一无二的标识。

例如：函数选择器即函数签名

```js
pragma solidity >=0.5.0 <0.9.0;
contract A {
    uint public a;
    function add(uint b) public {
        a+=b;
    }
    bytes4 public selector = this.transfer.selector;
    bytes4 public genSelector = bytes4(keccak256("add(uint256)"));
    bool public isequal = (selector==genSelector);
}
```

- `sha256(bytes memory) returns (bytes32)`: 计算参数的 SHA-256 哈希。
- `ripemd160(bytes memory) returns (bytes20)`: 计算参数的 RIPEMD-160 哈希。
- `ecrecover(bytes32 hash, uint8 v, bytes32 r, bytes32 s) returns (address)`利用椭圆曲线签名恢复与公钥相关的地址。

函数参数对应于 ECDSA签名的值:

`r` = 签名的前 32 字节

`s` = 签名的第2个32 字节

`v` = 签名的最后一个字节

(以后还需要补充许多密码学知识)

### 合约相关

`this`：表示当前合约的实例。

`selfdestruct(address payable recipient)` ：在交易成功结束后，销毁合约，并且把余额转到指定地址。接受的合约不会运行。

### 类型信息

`type(X)` 返回`X`的类型，目前只支持整型和合约类型，未来计划会拓展。

**用于合约类型 `C` 支持以下属性:**

- `type(C).name`:

  获得合约名

- `type(C).creationCode`:

  获得包含创建合约的字节码的`memory  byte[]`数组。只能在内联汇编中使用。

- `type(C).runtimeCode`

  获得合同的运行时字节码的内存字节数组，通常在构造函数的内联汇编中使用。

在**接口类型**``I``下可使用:

- `type(I).interfaceId`:

  返回接口``I`` 的 `bytes4` 类型的接口 ID，接口 ID 参考： [EIP-165](https://learnblockchain.cn/docs/eips/eip-165.html) 定义的， 接口 ID 被定义为  接口内所有的函数的函数选择器（除继承的函数）的`XOR` （异或）。

对于**整型** `T` 有下面的属性可访问：

- `type(T).min`

  `T` 的最小值。

- `type(T).max`

  `T` 的最大值。

## 