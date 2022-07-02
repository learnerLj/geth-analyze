## 代码结构

直观理解代码结构，下面是铸造，生成代币的代码。

```solidity
pragma solidity ^0.4;
contract Coin{
    //set the "address" type variable minter
    address public minter; 
    mapping (address =>uint) public balances; 
    // set an event so as to be seen publicly
    event Sent(address from,address to,uint amount); 
    //constructor only run once when creating contract,unable to invoke
    //"msg" is the address of creator."msg.sender"  is 
    constructor()public{
        minter=msg.sender;
    }
    //铸币
    //can only be called by creator
    function mint(address receiver,uint amount)public{
        require(msg.sender ==minter);
        balances[receiver]+=amount;
    }
    //转账 
    function send(address receiver,uint amount)public{
        require(balances[msg.sender]>= amount);
        balances[msg.sender]-=amount;
        balances[receiver]+=amount;
        emit Sent(msg.sender,receiver,amount);
    }

}
```



### 版本标识

 `pragma`

版本标识，是pragmatic information的简称，用于启动编译器检查，避免因为solidity更新后造成的不兼容和语法变动的错误。**只对本文件有效，如果导入其他文件，版本标识不会被导入，而是采用工作的文件自身的版本标识**

```
pragma solidity ^0.5.2;
```

这里`^`表示从0.5.2到0.6（不含）的版本



### 导入其他文件

`import "filename";` 这种导入方式会把导入文件的所有全局符号都导入到工作文件的全局作用域，会污染命名空间，不建议这么使用。

```
import * as symbolName from "filename";
//等价于
import "filename" as symbolName;
```

这样所有的全局符号都以`symbolName.symbol`的格式提供。

我们还可以设置别名，别名和重定义的符号名，都可以表示导入的文件里的全局符号。

```
import {symbol1 as alias, symbol2} from "filename";
```

支持从网络中导入，如：`import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/release-v3.3/contracts/cryptography/ECDSA.sol";`

### 路径

路径的形式和Linux下的完全一致，但是要避免使用`..`。我们可以引入指定路径的文件，如`import "./filename" as simbolName`，是当前目录下的文件。引用的文件除了本地文件，也可以是网络资源。

**实际solc编译器使用的时候可以指定路径的重映射，编译器可以从重映射的位置读取文件。尤其是使用网络文件的时候** 例如，可以使`github.com/ethereum/dapp-bin/library` 会被重映射到 `/usr/local/dapp-bin/library` ,格式如下。

```
solc github.com/ethereum/dapp-bin/=/usr/local/dapp-bin/ source.sol
```

更具体地会在solc编译器地部分说明。而truffle框架和remix就相对智能，可以通过网络获取文件。

### 注释

单行注释`//`,多行注释`/*......*/`

一种natspec注释，他是用`///`或者`/**......*/`，它里面可以使用Doxygen样式来给出相关地信息。

Doxygen样式地注释可以使特殊地注释形式变得可识别，方便读取和自动提取信息。主要有

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.21 <0.9.0;

/** @title Shape calculator.
 * @file（文件名）
 * @author John Doe <jdoe@example.com>（作者）
 * @version 1.0 （版本）
 * @details （细节）
 * @date （年-月-日）
 * @license （版权协议）
 * @brief （类的简单概述）
 * @section LICENSE（这一段的主要内容）
 * @param  Description of method's or function's input parameter（形式参数说明）
 * @return Description of the return value（返回说明）
 * @retval （返回值说明）
 * @attention（注意）
 * @warning（警告）
 * @var（变量声明）
 * @bug（代码缺陷）
 * @exception（异常）
 */
contract ShapeCalculator {
    /// @dev Calculates a rectangle's surface and perimeter.
    /// @param w Width of the rectangle.
    /// @param h Height of the rectangle.
    /// @return s The calculated surface.
    /// @return p The calculated perimeter.
    function rectangle(uint w, uint h) public pure returns (uint s, uint p) {
        s = w * h;
        p = 2 * (w + h);
    }
}
```

特别地，可以使用 `pragma abicoder v1` 或者 `pragma abicoder v1` 指定 ABI 的编码器和解码器版本，一般而言 0.8.0 以后，默认使用 v2 版本。

### 全局变量

状态变量是永久地存储在合约存储中的值，它具有数据的类型，也有可见性的属性。**在函数外的都是`storage`全局变量**。

```solidity
pragma solidity >=0.4.0 <0.9.0;

contract TinyStorage {
    uint storedXlbData; // 状态变量
    // ...
}
```

### 函数

函数是代码的可执行单元。函数通常在合约内部定义，但也可以在合约外定义。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >0.7.0 <0.9.0;

contract TinyAuction {
    function Mybid() public payable { // 定义函数
        // ...
    }
}

// Helper function defined outside of a contract
function helper(uint x) pure returns (uint) {
    return x * 2;
}
```

函数调用可发生在合约内部或外部，且函数有严格的可见性限制，对于谁可以调用它有着明确的规定（ [可见性和 getter 函数](https://learnblockchain.cn/docs/solidity/contracts.html#visibility-and-getters)）。

函数的返回值可以是元组，接收时需要一一对应。

### 函数修饰

函数修饰符用来修饰函数，比如添加函数执行前必须的先决条件。这样可以方便地实现代码复用。

```solidity
contract Owner {
   modifier onlyOwner {
      require(msg.sender == owner);
      _;
   }
   modifier costs(uint price) {
      if (msg.value >= price) {
         _;
      }
   }
}
```

函数体会插入在修饰函数的下划线`_`的位置。所以只有当修饰条件满足之后才能执行这个函数，否则报错。

注意下面的用法。实际上常常会被继承，作为模块复用。

可以看到，使用的格式

`function funcName(params) 可见性修饰 函数属性修饰 函数修饰器 returns(params)`

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Test{
    uint public a;
    uint public b;
    function set(uint _a,uint _b) public{
        a=_a;
        b=_b;
    }
    modifier Func(uint _a)
    {
        require(a>_a,"error:a is so small.");
        _;
    }
    function f(uint _a) public view Func(_a) returns(uint) {
        return _a;
    }
}
```

### 事件

事件是能方便地调用以太坊虚拟机日志功能的接口，分为设置事件和触发事件。

设置事件只需要 `event 事件名(params)` 。

触发事件 `emit 事件名(实参)`，注意触发事件和设置事件的参数类型需要匹配。

```solidity
pragma solidity >=0.4.21 <0.9.0;
contract TinyAuction {
    event HighestBidIncreased(address bidder, uint amount); // 事件

    function bid() public payable {
        // ...
        emit HighestBidIncreased(msg.sender, msg.value); // 触发事件
    }
}
```

## 合约

合约的构造函数至多一个，只在部署执行一次。

创建合约的方式可以是：Remix 这样的IDE、合约创建合约、用web3.js API.

部署的在区块链上的代码包括了所有可调用的函数或者是被其他函数调用的函数，但是不包括构造函数代码和只被构造函数调用的内部函数的代码。

构造函数的参数的ABI编码在合约的代码之后传递，web3.js可以略过这个细节。

支持合约类型和地址类型的强制转换。

### 函数和变量的可见性

可见性标识符在类型标识的后面。

`external`: 外部函数作为合约接口的一部分，可以被交易或者其他合约调用。 外部函数 `f` 不能以内部调用的方式调用（即 `f` 不起作用，但 `this.f()` 可以）。 

`public`: public 函数是合约接口的一部分，可以在内部或通过消息调用。对于 public 状态变量， 会自动生成一个 getter 函数。

`internal` : 只能在当前合约内部或它派生合约中访问，不使用 `this` 调用。

`private`: private 函数和状态变量仅在当前定义它们的合约中使用，并且不能被派生合约使用（如继承）。

**注意：**区块链所有信息都是透明的，这里的可见性只是针对其他合约或者调用者的是否有权限，访问或者修改状态。

**getter函数**：`public` 的状态变量会自动生成一个 getter 函数，内部调用时相当于状态变量，外部调用时相当于一个函数。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.0 <0.9.0;

contract C {
    uint public data;
    function x() public returns (uint) {
        data = 3; // internal access
        return this.data(); // external access
    }
}
```

如果这个 `public` 的全局变量是一个数组，那么 getter 函数就只能通过下标访问单个元素，但是结构体中的数组或者是映射不能够返回。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.0 <0.9.0;

contract Complex {
    struct Data {
        uint a;
        bytes3 b;
        mapping (uint => uint) map;
        uint[3] c;
        uint[] d;
        bytes e;
    }
    mapping (uint => mapping(bool => Data[])) public data;
}
```

等效为

```solidity
function data(uint arg1, bool arg2, uint arg3)
    public
    returns (uint a, bytes3 b, bytes memory e)
{
    a = data[arg1][arg2][arg3].a;
    b = data[arg1][arg2][arg3].b;
    e = data[arg1][arg2][arg3].e;
}
```

### 函数修饰器

函数修饰器会在函数执行前见擦汗条件，只有标记为`virtual`的情况下，才会被继承的合约覆盖。使用方法看下面的例子。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >0.7.0 <0.9.0;

contract owned {
    constructor() { owner = payable(msg.sender); }

    address owner;

    // 函数修饰器通过继承在派生合约中起作用。
    // 函数体会被插入到特殊符号 _; 的位置。
       modifier onlyOwner {
        require(
            msg.sender == owner,
            "Only owner can call this function."
        );
        _;
    }
}

contract destructible is owned {
    //调用格式是在 可见性修饰符（或者view(payable)权限修饰符） 之后，returns之前 
    function destroy() public onlyOwner {
        selfdestruct(owner);
    }
}

contract priced {
    // 修改器可以接收参数：
    modifier costs(uint price) {
        if (msg.value >= price) {
            _;
        }
    }
}

contract Register is priced, destructible {
    mapping (address => bool) registeredAddresses;
    uint price;

    constructor(uint initialPrice) { price = initialPrice; }

    function register() public payable costs(price) {
        registeredAddresses[msg.sender] = true;
    }

    function changePrice(uint _price) public onlyOwner {
        price = _price;
    }
}

contract Mutex {
    bool locked;
    modifier noReentrancy() {
        require(
            !locked,
            "Reentrant call."
        );
        locked = true;
        _;
        locked = false;
    }

    // 这个函数受互斥量保护，这意味着 `msg.sender.call` 中的重入调用不能再次调用  `f`。
    function f() public noReentrancy returns (uint) {
        (bool success,) = msg.sender.call("");
        require(success);
        return 7;
    }
}
```

函数修饰器只能在当前合约或者是继承的合约中使用。库合约内的函数修饰器只能在库合约中定义及使用。

**如果一个函数中有许多修饰器**，写法上以空格隔开，执行时依次执行：首先进入第一个函数修饰器，然后一直执行到`_;`接着跳转回函数体，进入第二个修饰器，以此类推。到达最后一层时，一次返回到上一层修饰器的`_;`后。

**修饰器不能够隐式地访问或者修改函数的变量，也不能够给函数提供返回值，只有规定的给修饰器的传入的参数才能够被修饰器使用**。

**显式地在修饰器中使用 `return` 不会影响函数地返回值，但是可能提前结束，就不会执行`_;`**处地函数体了。修饰器和函数中的 `return` 都只会跳出当前的代码块，进入上一层的堆栈。

`_` 可以在修饰器中多次出现，每一处都会执行函数体（注意包括函数地其他修饰器）。

修饰器的参数可以是任意表达式，函数中可见的函数外的变量，在修饰器中都是可见的。但是修饰器中的变量对函数不可见。

### 构造函数

如果没有构造函数，就等同于有默认的构造函数`constructor() {}`.

在继承中，构造函数有两种写法，一种是继承时直接给参数，形如`is Base(7)`；另外一种是在子合约的构造函数中定义，这很适用于依赖子合约状态给父合约的构造函数赋值的情况。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Base {
    uint x;
    constructor(uint _x) { x = _x; }
}

contract Derived1 is Base(7) {
    constructor() {}
}

contract Derived2 is Base {
    constructor(uint _y) Base(_y * _y) {}
}
```



### 常量和不可变量

全局变量如果有 `constant` 或者 `immutable` 标识，表示他们在合约创建后不可改变（但是可以在创建时可以使用使用 `constructor` 修饰。他们的区别在于：

- `constant`的值必须是全局变量，且声明时就要确定，且不可在构造函数中修改，因为它是**在编译时就确定且固定的**。而且在构造函数中，给 `constant` 赋值的表达式必须是返回固定的值，不能是运行时才确定的值。

- `immutable` 既可以在全局变量声明时确定（此后不可用构造函数修改），也可以在构造函数中确定（但只能赋值一次），因为**在构建时才确定并且固定的**。创建 `immutable` 变量发生在返回合约的 `creation code` 之前，编译器会发生值替换，修改合约的 `runtime code` ，这会造成区块链上实际存储的代码和 `runtime code` 有些差异。

在编译时，编译器不会给这些变量留储存位置，而是把常量和不可变量当作常量表达式，因此相比于常规的全局变量，消耗的gas少得多。

`constant` 的常量将会把赋值给它的表达式复制到所有访问它的位置，然后再进行求值的运算，类似于 C 语言的 `#define  a (7*5)`。`immutable` 的不变量则是将表达式的值复制到访问它的位置，但是占用固定的32个字节，类似于 `#define  a (35)` 。因此，不可变量占用空间较多，而且实际计算表达式时会优化，`constant` 的常量可能更加省gas

只有值类型或者常量字符串 `string` 才支持 `constant` 和 `immutable` 的标识。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.4;

uint constant X = 32**22 + 8;

contract C {
    string constant TEXT = "abc";
    bytes32 constant MY_HASH = keccak256("abc");
    uint immutable decimals;
    uint immutable maxBalance;
    address immutable owner = msg.sender;

    constructor(uint _decimals, address _reference) {
        decimals = _decimals;
        // Assignments to immutables can even access the environment.
        maxBalance = _reference.balance;
    }

    function isBalanceTooHigh(address _other) public view returns (bool) {
        return _other.balance > maxBalance;
    }
}
```

### 函数

#### 自由函数

函数既可以定义在合约内，也可以定义在合约外。

定义在合约外的函数叫做自由函数，一定是`internal`类型，就像一个内部函数库一样，会包含在所有调用他们的合约内，就像写在对应位置一样。但是自由函数不能直接访问全局变量和其他不在作用域下的函数（比如，需要通过地址引入合约，再使用合约内的函数）。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >0.7.0 <0.9.0;

function sum(uint[] memory _arr) pure returns (uint s) {
    for (uint i = 0; i < _arr.length; i++)
        s += _arr[i];
}

contract ArrayExample {
    bool found;
    function f(uint[] memory _arr) public {
        // This calls the free function internally.
        // The compiler will add its code to the contract.
        uint s = sum(_arr);
        require(s >= 10);
        found = true;
    }
}
```

#### 参数和返回值

外部函数 不可以接受多维数组作为参数，除非原文件加入 `pragma abicoder v2;`，以启用启用ABI v2版编码功能。 （注：在 0.7.0 之前是使用``pragma experimental ABIEncoderV2;``）

非内部函数无法返回多维动态数组、结构体、映射。如果添加 `pragma abicoder v2;` 启用 ABI V2 编码器，则是可以的返回更多类型，不过 `mapping` 仍然是受限的。

内部函数默认可以接受多维数组作为参数。

返回值的变量名可以出现，也可以省略。当变量名出现时，可以不写明`return`，但是如果和全局变量重名，就会局部覆盖。

#### `view` 函数

`view` 函数不能产生任何修改。由于操作码的原因，`view ` 库函数不会在运行时阻止状态改变，不过编译时静态检查器会发现这个问题。

以下行为都视为修改状态：

1. 修改状态变量。
2. 触发事件。
3. 创建其它合约。
4. 使用 `selfdestruct`。
5. 通过调用发送以太币。
6. 调用任何没有标记为 `view` 或者 `pure` 的函数。
7. 使用低级调用。
8. 使用包含特定操作码的内联汇编。

注意：`constant` 之前是 `view` 的别名，在0.5.0之后移除。

注意：`getter`方法会自动标记为`view`。

注意：在0.5.0前，`view`函数仍然可能产生状态修改。

#### `pure`函数

`pure` 函数不会读取状态，也不会改变状态。但是由于EVM的更新，也可能读取状态，而且无法在虚拟机水平上强制不读取状态。

以下行为视为读取状态：

1. 读取状态变量。
2. 访问 `address(this).balance` 或者 `<address>.balance`。
3. 访问 `block`，`tx`， `msg` 中任意成员 （除 `msg.sig` 和 `msg.data` 之外）。
4. 调用任何未标记为 `pure` 的函数。
5. 使用包含某些操作码的内联汇编。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.5.0 <0.9.0;

contract C {
    function f(uint a, uint b) public pure returns (uint) {
        return a * (b + 42);
    }
}
```

在`try/catch`中的回滚，不会视作状态改变。



### 事件

事件是对EVM日志的简短总结，可以通过RPC接口监听。触发事件时，设置好的参数就会记录在区块链的交易日志中，永久的保存，但是合约本身是不可以访问这些日志的。可以通过带有日志的Merkle证明的合约，来检查日志是否存在于区块链上。由于合约中仅能访问最近的 256 个区块哈希，所以还需要提供区块头信息。

也可以对事件中至多三个参数附加 `indexed` 属性，他们就会成为 `topics` 数据结构的一部分（详细请查看 ABI 部分编码的方式）。一个`topic`只可以容纳32个字节，对于`indexed`的引用类型会把值的 Keccak-256 hash 储存在一个`topic`。`topic` 允许通过过滤器来搜索事件，比如出发事件的合约地址。

没有 `indexed` 的参数就会被ABI编码后存储在日志。

如果没有使用 `anonymous` 标识符的话，事件的签名的哈希值就会是一个`topic`，如果使用了的话就无法通过除了触发它的合约地址之外的方式（如：事件的参数）去筛选事件。但是匿名事件在部署和调用时更节省gas，而且可以使用四个`index`(虽然没啥用了)。

```solidity
pragma solidity  >=0.4.21 <0.9.0;

contract ClientReceipt {
    event Deposit(
        address indexed _from,
        bytes32 indexed _id,
        uint _value
    );

    function deposit(bytes32 _id) public payable {
        // 事件使用 emit 触发事件。
        // 我们可以过滤对 `Deposit` 的调用，从而用 Javascript API 来查明对这个函数的任何调用（甚至是深度嵌套调用）。
        emit Deposit(msg.sender, _id, msg.value);
    }
}
```

## 