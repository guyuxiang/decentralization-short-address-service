# SAS (Short Address Service)

去中心化短链接服务 - 基于 Cosmos SDK 的区块链应用

## 简介

SAS 是一个运行在 Cosmos 区块链上的短链接/短链接服务。用户可以购买唯一的短链接（如 `abc123`），将其映射到任意长 URL，实现类似 `bit.ly` 的短链接功能，但具有去中心化特性。

## 功能特性

- **短链接生成**: 1-6位 Base62 编码（可自定义长度）
- **自定义长度**: 支持1-6位自由设置
- **直接访问**: 无需 `/s` 前缀，直接使用 `/{sUrl}` 访问
- **域名式交易**: 支持购买、出售、定价短链接
- **快速查询**: Bloom Filter + LRU 缓存优化
- **自动跳转**: 访问短链接自动重定向到目标网站
- **安全性**: URL 协议校验，防止钓鱼攻击
- **过期机制**: 默认1年租期，7天宽限期后自动回收
- **资金托管**: Escrow 托管机制保障交易安全
- **黑名单**: 屏蔽恶意/非法 URL 和域名
- **批量操作**: 支持批量设置长 URL（A/B测试）
- **访问统计**: 记录点击次数，展示热门链接
- **数据持久化**: BloomFilter/LRU/黑名单/统计自动持久化

## 技术架构

```
sas/
├── types.go           # 数据结构定义
├── msgs.go            # 区块链消息定义
├── keeper.go          # 状态持久化
├── handler.go         # 消息处理逻辑
├── querier.go         # 查询处理
├── urlManager.go      # 短 URL 生成管理
├── base62.go          # Base62 编解码
├── bloomFilter.go     # 布隆过滤器
├── lru.go             # LRU 缓存
├── codec.go           # 编解码注册
│
└── client/
    ├── cli/           # CLI 命令
    │   ├── tx.go      # 交易命令
    │   └── query.go   # 查询命令
    └── rest/          # REST API
        └── rest.go    # HTTP 接口
```

## 快速开始

### 1. 安装依赖

```bash
# 确保已安装 Go 1.13+
go version

# 安装 Cosmos SDK 依赖
go mod tidy

# 编译项目
go build -o sasd ./cmd/sasd
go build -o sascli ./cmd/sascli
```

### 2. 初始化并启动节点

```bash
# 初始化节点
rm -rf ~/.sasd
sasd init --chain-id test-chain --home ~/.sasd

# 添加测试账户（替换为你的地址）
sasd add-genesis-account cosmos1rxh8pl8k3gea67t4uhw2387v9hqpgz3u2awk2g 1000000stake --home ~/.sasd

# 启动节点
nohup sasd start --home ~/.sasd > /tmp/sasd.log 2>&1 &
```

### 3. 配置CLI并创建钱包

```bash
# 配置CLI连接节点
sascli config --home ~/.sascli --chain-id test-chain node http://localhost:26657

# 创建钱包
echo -e "password\npassword" | sascli keys add test --home ~/.sascli
```

### 4. 启动REST服务器

```bash
# 启动REST服务器（监听0.0.0.0:80，可从公网访问）
nohup sascli rest-server --home ~/.sascli --node http://localhost:26657 --chain-id test-chain --trust-node --laddr tcp://0.0.0.0:80 > /tmp/rest.log 2>&1 &
```

### 5. 购买短链接

```bash
# 购买短链接（自动生成）
echo -e "password\n" | sascli tx sas buy-sUrl "" 1stake --from test --home ~/.sascli --chain-id test-chain --node http://localhost:26657 --yes

# 购买指定短链接
echo -e "password\n" | sascli tx sas buy-sUrl vpn 1stake --from test --home ~/.sascli --chain-id test-chain --node http://localhost:26657 --yes

# 购买指定长度的短链接（1-6位）
echo -e "password\n" | sascli tx sas buy-sUrl "" 1stake --from test --home ~/.sascli --chain-id test-chain --node http://localhost:26657 --yes --length 3
```

### 6. 设置长链接

```bash
# 将短链接映射到长 URL
echo -e "password\n" | sascli tx sas set_lUrl vpn https://example.com --from test --home ~/.sascli --chain-id test-chain --node http://localhost:26657 --yes
```

#### 设置长链接

```bash
# 将短链接映射到长 URL
sascli tx sas set_lUrl abc123 https://google.com --from <your-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### 设置价格

```bash
# 设置短链接售价
sascli tx sas set_price abc123 200stake --from <your-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### 设置出售

```bash
# 设置是否可出售
sascli tx sas set_sell abc123 true --from <your-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### 续期

```bash
# 续期短链接（按天）
sascli tx sas renew abc123 30 --from <your-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### Escrow 托管交易

```bash
# 创建托管购买（资金暂存链上）
sascli tx sas buy-escrow abc123 100stake --from buyer --chain-id <chain-id> --yes --gas-prices 0.00001stake

# 确认交易（卖方确认后资金转给卖方）
sascli tx sas confirm-escrow abc123 --from seller --chain-id <chain-id> --yes --gas-prices 0.00001stake

# 取消交易（资金退回买方）
sascli tx sas cancel-escrow abc123 --from buyer --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### 批量操作

```bash
# 批量设置长 URL（用于 A/B 测试）
sascli tx sas batch-set-lurl abc123 "https://a.com,https://b.com" --from <your-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### 黑名单管理

```bash
# 添加 URL 到黑名单
sascli tx sas add-blacklist https://evil.com url --from <admin-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake

# 添加域名到黑名单
sascli tx sas add-blacklist evil.com domain --from <admin-key> --chain-id <chain-id> --yes --gas-prices 0.00001stake
```

#### 查询操作

```bash
# 查询长链接
sascli query sas lurl abc123 --chain-id <chain-id>

# 查询地址详情
sascli query sas laddress abc123 --chain-id <chain-id>

# 查询所有短链接（分页）
sascli query sas surls --chain-id <chain-id>

# 按所有者查询
sascli query sas owner-surls <owner-address> --chain-id <chain-id>

# 查询访问统计
sascli query sas stats --chain-id <chain-id>
```

## REST API

启动节点后，默认端口 80（或 1317）：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/{sUrl}` | 访问短链接自动跳转到长URL |
| GET | `/sas/adress/{sUrl}/lUrl` | 查询长链接 |
| GET | `/sas/adress/{sUrl}/lAddress` | 查询地址详情 |
| GET | `/sas/adress/sUrls` | 查询所有短链接 |
| GET | `/sas/stats` | 查询访问统计 |
| POST | `/sas/adress` | 购买短链接 |
| PUT | `/sas/adress/lUrl` | 设置长链接 |
| PUT | `/sas/adress/price` | 设置价格 |
| PUT | `/sas/adress/sell` | 设置是否出售 |

### 自动跳转

访问短链接自动重定向（无需 `/s` 前缀）：

```
http://localhost/vpn
→ 重定向到 https://example.com

http://openshort.cloud/code
→ 重定向到 https://github.com/guyuxiang/decentralization-short-address-service
```

## 消息类型

| 消息 | 说明 |
|------|------|
| `MsgBuySUrl` | 购买短链接 |
| `MsgSetLUrl` | 设置/更新长链接 |
| `MsgSetPrice` | 设置售价 |
| `MsgSetSell` | 设置是否可出售 |
| `MsgRenew` | 续期 |
| `MsgBuySUrlEscrow` | 托管购买 |
| `MsgConfirmEscrow` | 确认托管交易 |
| `MsgCancelEscrow` | 取消托管交易 |
| `MsgBatchSetLUrl` | 批量设置长链接 |
| `MsgAddBlackList` | 添加黑名单 |

## 技术细节

### Base62 编码

- 字符集: `0-9a-zA-Z` (62个字符)
- 短链接长度: 1-6位（可自定义）
- 1位: 62 个
- 2位: 3,844 个
- 3位: 238,328 个
- 4位: 14,776,336 个
- 5位: 916,132,832 个
- 6位: 56,800,235,584 个（约568亿）

### 缓存优化

- **Bloom Filter**: 快速判断短链接是否存在
- **LRU Cache**: 缓存热点查询结果
- **持久化**: 缓存数据自动保存到磁盘

### 过期机制

- 默认租期: 365 天
- 宽限期: 7 天（过期后7天内可续期）
- 自动清理: 每个区块检查并清理超过宽限期的短链接

### 手续费

- 购买/交易: 5% 手续费
- 续期: 1 代币/天

### 安全校验

- LUrl 必须以 `http://` 或 `https://` 开头
- LUrl 长度限制 2048 字符
- 重定向时二次校验协议
- 黑名单拦截恶意 URL/域名

## 开发说明

### 初始化全局组件

组件在 `init()` 中自动初始化：

```go
// urlManager.go
var gC = &globeCounter{
    number: new(uint32),
}

// bloomFilter.go  
GlobalBloomFilter = NewBloomFilter(1024*1024, 16)

// lru.go
LruCache = New(10000)
```

### 模块注册

在 `app.go` 中注册 SAS 模块：

```go
app.ModuleBasics = module.NewBasicManager(
    // ... other modules
    sas.AppModuleBasic{},
)
```

### Keeper 初始化

```go
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, dataDir string) Keeper {
    k := Keeper{
        coinKeeper: coinKeeper,
        storeKey:   storeKey,
        cdc:        cdc,
        dataDir:    dataDir,  // 用于持久化缓存数据
    }
    k.loadBloomFilter()
    k.loadLRUCache()
    k.loadBlackList()
    k.loadStats()
    return k
}
```

## 注意事项

1. 购买短链接需要支付代币 + 5% 手续费
2. 设置长链接前必须先拥有该短链接
3. 只有设置为可出售(sell=true)才能被购买
4. 出价必须高于当前售价
5. 短链接过期后有7天宽限期可续期
6. Escrow 交易更安全，建议大额交易使用
7. 黑名单功能需要管理员权限
8. 短链接长度必须为 1-6 位字符
9. 访问短链接无需 `/s` 前缀，直接使用 `/{sUrl}` 访问
10. 所有交易命令需要添加 `--chain-id`、`--yes` 参数
11. CLI 查询命令需要添加 `--chain-id` 参数

## License

MIT License
