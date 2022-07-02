## 表达式和控制结构

Solidity 支持 `if`, `else`, `while`, `do`, `for`, `break`, `continue`, `return`这些和C语言一样的关键字。

Solidity还支持 `try`/`catch` 语句形式的异常处理， 但仅用于 外部函数调用 和 合约创建调用．

由于不支持非布尔类型值转换成布尔类型，因此`if(1){}`是不合法的。

### 函数调用

#### 内部调用

内部调用再EVM中只是简单的跳转（不会产生实际的消息调用），传递当前的内存的引用，效率很高。但是仍然要避免过多的递归，因为每次进入内部函数都会占用一个堆栈槽，而最多只有1024个堆栈槽。

#### 外部调用

- 只有`external`或者`public`的函数才可以通过消息调用而不是单纯的跳转调用，外部函数的参数会暂时复制在内存中。

- **`this`不可以出现在构造函数里，因为此时合约还没有完成**。  

- **调用时可以指定 value 和 gas** 。这里导入合约使用的时初始化合约实例然后赋予地址。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.6.2 <0.9.0;

contract InfoFeed {
    function info() public payable returns (uint ret) { return 42; }
}

contract Consumer {
    InfoFeed feed;
    function setFeed(InfoFeed addr) public { feed = addr; }
    function callFeed() public { feed.info{value: 10, gas: 800}(); }
}
```

需要注意到，`function callFeed() public { feed.info{value: 10, gas: 800}(); }`，花括号`{ feed.info{value: 10, gas: 800}`里的只是修饰，实际调用的时圆括号`()`里的内容。再0.7.0前，使用的时`f.value(x).gas(g)()`。

一般我们不推荐使用call调用除了`fallback`函数之外的函数，但是在考虑节省gas和保证安全性的前提下可以尝试。

#### 函数参数写法

调用函数时参数还有一种写法：与函数声明的的名字对应。当然，最常见的还是按照顺序，忽略函数参数的名字。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.0 <0.9.0;

contract C {
    mapping(uint => uint) data;

    function f() public {
        set({value: 2, key: 3});
    }

    function set(uint key, uint value) public {
        data[key] = value;
    }

}
```

### 用`new`创建合约实例

在已知一个合约完整的代码的前提下（比如写在同一个文件内），就可以使用`contractName newContractInstance{value:initial value}(constructor para)` ，（注意无法限定gas，但是可以写明发送多少以太币过去)。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
contract D {
    uint public x;
    constructor(uint a) payable {
        x = a;
    }
}

contract C {
    D d = new D(4); // will be executed as part of C's constructor

    function createD(uint arg) public {
        D newD = new D(arg);
        newD.x();
    }

    function createAndEndowD(uint arg, uint amount) public payable {
        // Send ether along with the creation
        D newD = new D{value: amount}(arg);
        newD.x();
    }
}
```

### 合约创建的新合约地址

合约的地址时由创建时交易的nonce和创建者的地址决定，但是还可以选择一个32个字节的`salt` 来改变生成合约地址的方式，合约地址将会由创建者的地址、给定`salt`、被创建合约的字节码及参数共同决定。下面是计算方法：

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.7.0;

contract D {
    uint public x;
    constructor(uint a) {
        x = a;
    }
}

contract C {
    function createDSalted(bytes32 salt, uint arg) public {
        /// 这个复杂的表达式只是告诉我们，如何预先计算合约地址。
        /// 这里仅仅用来说明。
        /// 实际上，你仅仅需要 ``new D{salt: salt}(arg)``.
        address predictedAddress = address(uint160(uint(keccak256(abi.encodePacked(
            bytes1(0xff),
            address(this),
            salt,
            keccak256(abi.encodePacked(
                type(D).creationCode,
                arg
            ))
        )))));

        D d = new D{salt: salt}(arg);
        require(address(d) == predictedAddress);
    }
}
```

这一特性使得在销毁合约之后在重新在同一地址生成代码相同的合约。但是，尽管创建的字节码相同，但是由于编译器会检查外部的状态变化，`deploy bytecode`可能会不一样。

下面是创建多个合约的例子：

```js
// SPDX-License-Identifier: MIT
pragma solidity ^0.7.6;

contract Car {
    address public owner;
    string public model;

    constructor(address _owner, string memory _model) payable {
        owner = _owner;
        model = _model;
    }
}

contract CarFactory {
    Car[] public cars;

    function create(address _owner, string memory _model) public {
        Car car = new Car(_owner, _model);
        cars.push(car);
    }

    function createAndSendEther(address _owner, string memory _model)
        public
        payable
    {
        Car car = (new Car){value: msg.value}(_owner, _model);
        cars.push(car);
    }

    function getCar(uint _index)
        public
        view
        returns (address owner, string memory model, uint balance)
    {
        Car car = cars[_index];

        return (car.owner(), car.model(), address(car).balance);
    }
}
```

**特别提到，调用已部署的合约，应当先引入合约的接口（或者源代码），然后`合约名 a=合约名(地址)`**

### 元组的赋值行为

函数的返回值可以是元组，因此就可以用元组的形式接收，但是必须按照顺序排列。在0.5.0之后，两个元组的大小必须相同，用逗号表示间隔，可以**空着省略元素**。注意，不允许赋值和声明都出现在元组里，比如`(x, uint y) = (1, 2);`不合法。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.5.0 <0.9.0;

contract C {
    uint index;

    function f() public pure returns (uint, bool, uint) {
        return (7, true, 2);
    }

    function g() public {
        // Variables declared with type and assigned from the returned tuple,
        // not all elements have to be specified (but the number must match).
        (uint x, , uint y) = f();
        // Common trick to swap values -- does not work for non-value storage types.
        (x, y) = (y, x);
        // Components can be left out (also for variable declarations).
        (index, , ) = f(); // Sets the index to 7
    }
}
```

**注意：元组的赋值行为，它仍然保留了引用类型。**

### 错误处理

调用和因这次调用而形成的调用链出现异常就会回滚所有更改，但是可以使用`try`或者`catch` 只回滚到这一层（回滚不会到底，如 A 调用B, B 调用 C, 如果B 调用 C 时出错导致回滚，不会消除 A 调用 B 造成的影响）。

底层函数错误是不会回滚的，而是返回 `false` 和 `error instance`.

有两种错误类型，一种是`error`，表示常规的错误。而`Panic`则表示代码没有问题，

`assert`函数，用于检查内部错误，返回`Panic(uint256)`，错误代码分别表示：

1. 0x00: 由编译器本身导致的Panic.
2. 0x01:  `assert` 的参数（表达式）结果为 false 。
3. 0x11: 在``unchecked { … }``外，算术运算结果向上或向下溢出。
4. 0x12: 除以0或者模0.
5. 0x21: 不合适的枚举类型转换。
6. 0x22: 访问一个没有正确编码的`storage`byte数组.
7. 0x31: 对空数组 `.pop()` 。
8. 0x32: 数组的索引越界或为负数。
9. 0x41: 分配了太多的内存或创建的数组过大。
10. 0x51: 如果你调用了零初始化内部函数类型变量。

`Error(string)`的异常（错误提示信息）由编译器产生，有以下情况：

1.  `require` 的参数为 `false` 。
2.  触发`revert`或者`revert("discription")`
3.  执行外部函数调用合约没有代码。
4.  `payable` 修饰的函数（包括构造函数和 fallback 函数），接收以太币。
5.  合约通过 getter 函数接收以太币 。

以下即可能是`Panic`也可能是`error`

1.  `.transfer()` 失败。
2.  通过消息调用调用某个函数，但该函数没有正确结束（例如, 它耗尽了 gas，没有对应的函数，或者本身抛出一个异常）。低级操作不会抛出异常，而通过返回 `false` 来指示失败。
3.  如果你使用 `new` 关键字创建未完成的合约 。

 **注意：** `Panic`异常使用的是`invalid`操作码，会消耗所有可用gas. 在 都会 版本之前，require 也是这样。

**注意：**`revert errorInstance` 其中的`errorInstance`是自定义的错误实例，用`errorInstance`的名字来表示错误，在编码的时候只占4个字节（如果带参数的话可能不一样），因此，远比`Error(string)`的方式节省gas。错误实例和函数调用与错误实例同名且同参数的函数的函数的ABI编码相同，也就是说错误实例的数据是由ABI编码后的4个字节的选择器组成的。而这个选择器是错误实例的签名的keccak256-hash 的前4个字节。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.4;

contract VendingMachine {
    address owner;
    error Unauthorized();
    function buy(uint amount) public payable {
        if (amount > msg.value / 2 ether)
            revert("Not enough Ether provided.");
        // Alternative way to do it:
        require(
            amount <= msg.value / 2 ether,
            "Not enough Ether provided."
        );
        // Perform the purchase.
    }
    function withdraw() public {
        if (msg.sender != owner)
            revert Unauthorized();

        payable(msg.sender).transfer(address(this).balance);
    }
}
```

**注意：**`require` 是可执行的函数，在 `require(condition, f())` 里，函数 `f` 会被执行，即便 `condition` 为 True .

**注意：**`Error(string)`函数会返回16进制的 错误提示信息。

**注意：**`throw`等同于`reverse()` 但是在0.5.0废除了。

#### `try`/`catch`

**`try`后面只能接外部函数调用或者是创建合约`new ContractName`的表达式**，并且花括号里面的错误会立即回滚，当花括号调用合约以外的函数（或者以外部调用的形式调用函数，如用 `this`）出现错不会造成当前合约回滚。用 `try` 尝试调用的外部函数如果需要返回参数，就要在 `returns` 后面声明返回参数的类型，如果外部调用执行成功就可以获取返回值，继续执行花括号内的语句，花括号的语句都完全成功了，就会跳过后面的 `catch`；但是如果失败就会根据错误类型跳转到对应的 `catch` 里面。如下面的代码：

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >0.8.0;
//接口类型，后面会介绍，如果熟悉 Golang 的接口则很容易理解。
interface DataFeed { function getData(address token) external returns (uint value); }

contract FeedConsumer {
    DataFeed feed;//从接口创建合约
    uint errorCount;//记录错误次数
    function rate(address token) public returns (uint value, bool success) {
        // 如果有十个及以上的错误就回滚
        require(errorCount < 10);
        try feed.getData(token) returns (uint v) {//尝试调用 外部的 getData 函数，执行成功就获得返回值，然后继续执行花括号内的内容
            return (v, true);
        } catch Error(string memory /*reason*/) {
            // 执行 revert 语句造成的回滚，返回错误提示信息
            errorCount++;
            return (0, false);
        } catch Panic(uint /*errorCode*/) {
            // Panic类型错误。
            errorCount++;
            return (0, false);
        } catch (bytes memory /*lowLevelData*/) {
            // 无返回提示的底层错误。
            errorCount++;
            return (0, false);
        }
    }
}
```

Solidity支持不同的`catch`代码块：

- `catch Error(string memory reason) { ... }`: 对应的执行条件是 `revert("reasonString")` or `require(false, "reasonString")` 或者是执行时内部的错误.
- `catch Panic(uint errorCode) { ... }`: 用于接收 Panic 类型的错误，比如用了 `assert`，数组下标越界，除以0，这些语言层面的错误。
- `catch (bytes memory lowLevelData) { ... }`:  如果发送错误类型不是前两种，比如无错误提示信息，或者是返回的错误提示信息无法解码（比如由编译器版本变迁造成），这个语句就会提供底层的编码后的错误提示信息。
- `catch { ... }`: 接收所有错误类型，但是不能出现前面的判断错误类型的分句。

注意：为了接收所有方式的错误，最后要使用 `catch { ...}`  或者 `catch (bytes memory lowLevelData) { ... }`.

注意：调用失败的原因多种多样，错误消息可能是来自调用链中的某一环，不一定来自被调用的合约。比如gas不足。在调用时会保留1/64的gas，以保证当前合约顺利执行。