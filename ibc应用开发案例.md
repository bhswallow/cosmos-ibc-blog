# IBC应用开发案例

## 1.生成链planet

用ignite生成一条新的区块链名字叫planet。

```shell
ignite scaffold chain planet --no-module
```

解析：

- ignite：是快速构建区块链的脚手架工具
- planet：表示启动链的名字为planet
- --no-module：表示不指定任何module

## 2.在链planet中加载blog模块

使用ignite生成一个Blog的模块，并且集成IBC。

```shell
cd planet
ignite scaffold module blog --ibc
```

解析：

- blog：进入planet项目目录，使用ignite scaffold脚手架创建一个blog的module

- --ibc：blog module需要集成ibc模块

  

## 3. 对模块blog的操作

### 3.1 给blog模块添加针对博文（post）的增删改查。

```shell
ignite scaffold list post title content creator --no-message --module blog

```

解析：

- post：创建一个结构体叫post
- list：这个post是个list，表示在创建数据库时，用数组形式存储
- title、content、creator：post结构体的三个字段，没设定类型，表示都是string类型。
- --no-message：不需要创建message。
  - message：我可以发送一笔交易，但是交易的message可以触发新增的增删改查的操作。
- --module blog：表示实在blog这个module下生成。



### 3.2添加已发成功博文（sentPost）的增删改查。

```shell
ignite scaffold list sentPost postID title chain creator --no-message --module blog
```

解析：

- sentPost：创建一个sentPost结构体
- postID：对手链收到了上面的post消息，它会返回对应postId返回给发送者链，发送者链需要把postID包含再sentPost里面，并保存到数据库里面。
- postID、title、chain、creator：sentPost结构体里面的四个字段，没设定类型，默认string。
- --no-message：不需要创建message。
  - message：我可以发送一笔交易，但是交易的message可以触发新增的增删改查的操作。



### 3.3 添加发送超时博文（timeoutPost）的增删改查。

```shell
ignite scaffold list timedoutPost title chain creator --no-message --module blog
```

解析：

- timedoutPost：如果一条博客发送出去了，但是可能因为网络的原因发生超时了，这时候需要把一条timedoutPost记录到数据库里面。
- title、chain、creator：这条timedoutPost需要包含title、chain、creator。
- --no-message：这条timedoutPost相关的增删改查操作不需要响应的message。

### 3.4 添加已发成功博文（updatePost）增删改查

```shell
ignite scaffold list updatePost postID title content --no-message --module blog
```



## 4.数据包

### 4.1 添加IBC发送数据包和确认数据包的结构。

```shell
ignite scaffold packet ibcPost title content --ack postID --module blog

```

解析：

- packet：数据包

- title、content：数据包包含博客的title和content。并在确认数据包里面包含postID。

- postID：postID是由数据库生成的

  

## 5.修改对应应用代码

### 5.1 添加creator

 在proto/blog/packet.proto目录下修改`IbcPostPacketData`，添加创建人`Creator`， 并重新编译proto文件。在x/blog/keeper/msg_server_ibc_post.go。编译完成后在x/blog/keeper/msg_server_ibc_post.go中发送数据包前更新`Creator`。

```shell
ignite chain build
```

### 5.2 修改keeper方法中的`OnRecvIbcPostPacket `。

```go
id := k.AppendPost(
        ctx,
        types.Post{
            Creator: packet.SourcePort + "-" + packet.SourceChannel + "-" + data.Creator,
            Title:   data.Title,
            Content: data.Content,
        },
    )

    packetAck.PostID = strconv.FormatUint(id, 10)
```

### 5.3 修改keeper方法中的`OnAcknowledgementIbcPostPacket `。

```
k.AppendSentPost(
            ctx,
            types.SentPost{
                Creator: data.Creator,
                PostID:  packetAck.PostID,
                Title:   data.Title,
                Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
            },
        )
```

### 5.4  修改keeper方法中的`OnTimeoutIbcPostPacket `。

```go
k.AppendTimedoutPost(
        ctx,
        types.TimedoutPost{
            Creator: data.Creator,
            Title:   data.Title,
            Chain:   packet.DestinationPort + "-" + packet.DestinationChannel,
        },
    )
```

## 6.配置两条链

### 6.1 添加链启动的配置文件。

```yaml
# earth.yml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
  - name: bob
    coins: ["500token", "100000000stake"]
validator:
  name: alice
  staked: "100000000stake"
faucet:
  name: bob
  coins: ["5token", "100000stake"]
genesis:
  chain_id: "earth"
init:
  home: "$HOME/.earth"
  
# mars.yml
accounts:
  - name: alice
    coins: ["1000token", "1000000000stake"]
  - name: bob
    coins: ["500token", "100000000stake"]
validator:
  name: alice
  staked: "100000000stake"
faucet:
  host: ":4501"
  name: bob
  coins: ["5token", "100000stake"]
host:
  rpc: ":26659"
  p2p: ":26658"
  prof: ":6061"
  grpc: ":9092"
  grpc-web: ":9093"
  api: ":1318"
genesis:
  chain_id: "mars"
init:
  home: "$HOME/.mars"
```

### 6.2 编译cmd的planet命令

```shell
cd path/planet/cmd/planetd
# 将cmd下的planetd的main.go编译成plannetd
go build
```

### 6.3 分别启动两条链

```shell
ignite chain serve -c earth.yml

ignite chain serve -c mars.yml
```

## 7. 启动中继器

### 7.1 配置并启动relayer

```shell
rm -rf ~/.ignite/relayer  # 如果之前启动过，要删掉

ignite relayer configure -a \
  --source-rpc "http://0.0.0.0:26657" \
  --source-faucet "http://0.0.0.0:4500" \
  --source-port "blog" \
  --source-version "blog-1" \
  --source-gasprice "0.0000025stake" \
  --source-prefix "cosmos" \
  --source-gaslimit 300000 \
  --target-rpc "http://0.0.0.0:26659" \
  --target-faucet "http://0.0.0.0:4501" \
  --target-port "blog" \
  --target-version "blog-1" \
  --target-gasprice "0.0000025stake" \
  --target-prefix "cosmos" \
  --target-gaslimit 300000

ignite relayer connect
```

## 8.两条链之间的互操作

### 8.1 从earth链向mars链发送博文数据包（注意修改channel id）

channelid从启动的relayer日志中可以看到，我这里是channel-0

```shell
cd path/planet/cmd/planetd
./planetd tx blog send-ibc-post blog channel-0 "Hello" "Hello Mars, I'm Alice from Earth" --from alice --chain-id earth --home ~/.earth
```

### 8.2 通过rpc查询验证结果。

```shell
./planetd q blog list-post --node tcp://localhost:26659

./planetd q blog list-sent-post
```
