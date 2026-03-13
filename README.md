# SAS (Short Address Service)

去中心化短地址服务 - 基于 Cosmos SDK 的区块链应用

## 简介

SAS 是一个运行在 Cosmos 区块链上的短地址/短链接服务。用户可以购买唯一的短地址（如 `abc123`），将其映射到任意长 URL，实现类似 `bit.ly` 的短链接功能，但具有去中心化特性。

## 功能特性

- **短地址生成**: 6位 Base62 编码（568亿+ 组合）
- **域名式交易**: 支持购买、出售、定价短地址
- **快速查询**: Bloom Filter + LRU 缓存优化
- **自动跳转**: 访问短地址自动重定向到目标网站
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
```

### 2. 启动区块链节点

```bash
# 初始化节点
星辰cosmos init <节点名称>

# 启动节点
cosmos start
```

### 3. 使用 CLI 操作

#### 购买短地址

```bash
# 购买新的短地址（不指定则自动生成）
sascli tx sas buy-sUrl "" 100sastoken --from <your-key>

# 购买指定短地址
sascli tx sas buy-sUrl abc123 100sastoken --from <your-key>
```

#### 设置长链接

```bash
# 将短地址映射到长 URL
sascli tx sas set_lUrl abc123 https://google.com --from <your-key>
```

#### 设置价格

```bash
# 设置短地址售价
sascli tx sas set_price abc123 200sastoken --from <your-key>
```

#### 设置出售

```bash
# 设置是否可出售
sascli tx sas set_sell abc123 true --from <your-key>
```

#### 续期

```bash
# 续期短地址（按天）
sascli tx sas renew abc123 30 --from <your-key>
```

#### Escrow 托管交易

```bash
# 创建托管购买（资金暂存链上）
sascli tx sas buy-escrow abc123 100sastoken --from buyer

# 确认交易（卖方确认后资金转给卖方）
sascli tx sas confirm-escrow abc123 --from seller

# 取消交易（资金退回买方）
sascli tx sas cancel-escrow abc123 --from buyer
```

#### 批量操作

```bash
# 批量设置长 URL（用于 A/B 测试）
sascli tx sas batch-set-lurl abc123 "https://a.com,https://b.com" --from <your-key>
```

#### 黑名单管理

```bash
# 添加 URL 到黑名单
sascli tx sas add-blacklist https://evil.com url --from <admin-key>

# 添加域名到黑名单
sascli tx sas add-blacklist evil.com domain --from <admin-key>
```

#### 查询操作

```bash
# 查询长链接
sascli query sas lurl abc123

# 查询地址详情
sascli query sas laddress abc123

# 查询所有短地址（分页）
sascli query sas surls

# 按所有者查询
sascli query sas owner-surls <owner-address>

# 查询访问统计
sascli query sas stats
```

## REST API

启动节点后，默认端口 1317：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/sas/adress/{sUrl}/lUrl` | 查询长链接 |
| GET | `/sas/adress/{sUrl}/lAddress` | 查询地址详情 |
| GET | `/sas/adress/sUrls` | 查询所有短地址 |
| GET | `/sas/stats` | 查询访问统计 |
| POST | `/sas/adress` | 购买短地址 |
| PUT | `/sas/adress/lUrl` | 设置长链接 |
| PUT | `/sas/adress/price` | 设置价格 |
| PUT | `/sas/adress/sell` | 设置是否出售 |

### 自动跳转

访问短地址自动重定向：

```
http://localhost:1317/s/abc123
→ 重定向到 https://google.com
```

## 消息类型

| 消息 | 说明 |
|------|------|
| `MsgBuySUrl` | 购买短地址 |
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
- 短地址长度: 6位
- 最大组合: 62^6 ≈ 568亿

### 缓存优化

- **Bloom Filter**: 快速判断短地址是否存在
- **LRU Cache**: 缓存热点查询结果
- **持久化**: 缓存数据自动保存到磁盘

### 过期机制

- 默认租期: 365 天
- 宽限期: 7 天（过期后7天内可续期）
- 自动清理: 每个区块检查并清理超过宽限期的短地址

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

1. 购买短地址需要支付代币 + 5% 手续费
2. 设置长链接前必须先拥有该短地址
3. 只有设置为可出售(sell=true)才能被购买
4. 出价必须高于当前售价
5. 短地址过期后有7天宽限期可续期
6. Escrow 交易更安全，建议大额交易使用
7. 黑名单功能需要管理员权限

## License

MIT License
