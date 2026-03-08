# 支付机构与银行账单API调研报告

**调研日期**: 2026-03-08
**项目**: account 个人记账应用

---

## 目录
1. [概述](#概述)
2. [支付机构](#支付机构)
   - [支付宝](#支付宝)
   - [微信支付](#微信支付)
   - [京东支付](#京东支付)
3. [银行](#银行)
   - [招商银行](#招商银行)
   - [工商银行](#工商银行)
   - [农业银行](#农业银行)
   - [中信银行](#中信银行)
   - [浦发银行](#浦发银行)
   - [中国银行](#中国银行)
4. [总结与建议](#总结与建议)

---

## 概述

本报告对国内主流支付机构和银行的账单获取接口进行了全面调研，旨在为 account 个人记账应用的"自动记账"功能提供技术参考。

**核心结论**: **绝大多数支付机构和银行都不直接向个人开发者开放账单查询API**，仅对企业/商户提供相关接口。

---

## 支付机构

### 支付宝

#### 1. API 提供情况
- **个人账户账单API**: ❌ **不对外开放**
- **商户账单API**: ✅ 提供

**主要API接口**:
- `alipay.data.bill.downloadurl.get` - 查询对账单下载地址接口（商户使用）

**官方文档**:
- https://opendocs.alipay.com/apis/api_15/alipay.data.bill.downloadurl.get
- https://open.alipay.com/

#### 2. 前提条件
- 必须注册**企业开发者账号**（个人无法使用）
- 完成企业资质认证
- 创建应用并获取 APPID
- 配置应用密钥（RSA2签名）
- 申请"查询对账单下载地址"接口权限
- 签约相关产品（如"支付宝收单"）

#### 3. 日期指定方式
- **仅支持按单日或单月查询**，不支持任意起止日期范围
- 参数: `bill_date`
  - 日账单格式: `yyyy-MM-dd`（如 "2025-03-07"）
  - 月账单格式: `yyyy-MM`（如 "2025-03"）
- 参数: `bill_type`
  - `trade`: 商户基于支付宝交易收单的业务账单
  - `signcustomer`: 基于商户支付宝余额收入及支出等资金变动的账务账单

**限制**:
- 对账单通常在每日上午9点后生成
- 通常只能查询最近1年内的账单
- 沙箱环境不支持此接口

#### 4. 返回字段（商户账单CSV）

| 字段名 | 说明 | 示例 |
|--------|------|------|
| 交易时间 | 交易发生的具体时间 | 2025-03-07 14:30:00 |
| 交易号 | 支付宝交易号 | 202503072200143XXXXXX |
| 商户订单号 | 商户系统订单号 | 202503070001 |
| 业务类型 | 交易业务类型 | 即时到账交易 |
| 交易名称 | 商品名称或交易描述 | 商品购买 |
| 交易对方 | 交易对方信息 | 某某商户 |
| 金额（元） | 交易金额 | 100.50 |
| 收入/支出 | 收入或支出标识 | 支出 |
| 交易状态 | 交易当前状态 | 交易成功 |
| 服务费（元） | 支付宝收取的手续费 | 0.60 |
| 成功退款（元） | 已退款金额 | 0.00 |
| 备注 | 交易备注信息 | 用户备注 |
| 支付渠道 | 支付所用渠道 | 支付宝余额 |

**个人账户手动导出CSV字段**:
- 交易时间、交易分类、交易对方、商品名称、金额、收/付款方式、交易状态、备注

#### 5. API调用示例

```python
from alipay import AliPay
import datetime
import requests

# 初始化支付宝客户端
alipay = AliPay(
    appid="your_app_id",
    app_notify_url=None,
    app_private_key_string="""-----BEGIN RSA PRIVATE KEY-----
您的应用私钥
-----END RSA PRIVATE KEY-----""",
    alipay_public_key_string="""-----BEGIN PUBLIC KEY-----
支付宝公钥
-----END PUBLIC KEY-----""",
    sign_type="RSA2",
    debug=False
)

# 设置账单日期
bill_date = (datetime.datetime.now() - datetime.timedelta(days=1)).strftime("%Y-%m-%d")

# 调用账单下载地址查询接口
result = alipay.api_alipay_data_bill_downloadurl_get(
    bill_type="trade",
    bill_date=bill_date
)

if result.get("code") == "10000":
    download_url = result.get("bill_download_url")
    # 下载账单文件（注意：下载地址有效期仅30秒）
    response = requests.get(download_url)
```

**HTTP请求示例**:
```
POST https://openapi.alipay.com/gateway.do
Content-Type: application/x-www-form-urlencoded

app_id=2021000111666666
&method=alipay.data.bill.downloadurl.get
&format=json
&charset=utf-8
&sign_type=RSA2
&sign=计算后的签名值
&timestamp=2025-03-08 10:00:00
&version=1.0
&biz_content={"bill_type":"trade","bill_date":"2025-03-07"}
```

---

### 微信支付

#### 1. API 提供情况
- **个人账户账单API**: ❌ **不对外开放**
- **商户账单API**: ✅ 提供

**主要API接口**:
- `GET https://api.mch.weixin.qq.com/v3/bill/tradebill` - 申请交易账单
- `GET https://api.mch.weixin.qq.com/v3/billdownload/file` - 下载账单

**官方文档**:
- https://pay.weixin.qq.com/wiki/doc/api/index.html
- https://pay.weixin.qq.com/docs/merchant/apis/download-bill/download-bill.html

#### 2. 前提条件
- 必须拥有**企业资质**（个人无法使用）
- 在微信支付商户平台注册并完成企业认证
- 拥有商户号
- 创建应用并获取API密钥
- 申请商户API证书用于接口签名认证
- 确保商户号已开通账单下载权限

#### 3. 日期指定方式
- **仅支持按单日查询**，不支持直接查询日期范围
- 参数: `bill_date`，格式为 `YYYY-MM-DD` 或 `YYYYMMDD`
- 仅支持2023年1月1日之后的账单
- 账单生成时间: T+1日上午10点后可下载前一日账单

**若要查询日期范围**: 需要循环调用API，逐日下载账单。

#### 4. 返回字段（交易账单CSV）

| 字段名 | 描述 |
|--------|------|
| 交易时间 | 订单支付时间 |
| 公众账号ID | 商户的公众号ID |
| 商户号 | 商户的商户号 |
| 子商户号 | 服务商模式下的子商户号 |
| 设备号 | 微信支付分配的终端设备号 |
| 微信订单号 | 微信支付订单号 |
| 商户订单号 | 商户系统内部订单号 |
| 用户标识 | 微信用户的OpenID |
| 交易类型 | 如JSAPI、NATIVE等 |
| 交易状态 | 如SUCCESS、REFUND等 |
| 付款银行 | 付款银行类型 |
| 货币种类 | 如CNY |
| 总金额 | 订单总金额 |
| 代金券或立减优惠金额 | 优惠金额 |
| 微信退款单号 | 微信退款单号 |
| 商户退款单号 | 商户退款单号 |
| 退款金额 | 退款金额 |
| 商品名称 | 商品描述 |
| 商户数据包 | 商户附加数据 |

**资金账单返回字段**:
- 记账时间、账户类型、资金变动方向、收支类型、金额、账户余额、业务类型、业务单据号

#### 5. API调用示例

```python
import requests

# 配置参数
merchant_id = "your_merchant_id"
serial_no = "your_serial_no"

# 申请账单
bill_date = "2024-01-01"
bill_type = "ALL"
url = f"https://api.mch.weixin.qq.com/v3/bill/tradebill?bill_date={bill_date}&bill_type={bill_type}"

# 注意：实际使用时需要实现完整的签名逻辑
headers = {
    "Authorization": "WECHATPAY2-SHA256-RSA2048 mchid='...',serial_no='...',nonce_str='...',timestamp='...',signature='...'",
    "Accept": "application/json"
}

response = requests.get(url, headers=headers)
result = response.json()
download_url = result["download_url"]

# 下载账单
download_response = requests.get(download_url, headers=headers)
with open(f"bill_{bill_date}.csv.gz", "wb") as f:
    f.write(download_response.content)
```

---

### 京东支付

#### 1. API 提供情况
- **个人账户账单API**: ❌ **没有公开的个人账单查询API**
- **商户账单API**: ✅ 提供

**官方文档入口**:
- 京东支付开发者中心: https://payapi.jd.com
- 京东商户平台: https://merchant.jd.com

#### 2. 前提条件
- 仅支持**企业用户**注册，个人无法申请
- 在京东开放平台完成企业认证
- 创建应用获取 AppKey 和 AppSecret
- 在商户后台申请账单/对账相关API权限
- 配置服务器IP白名单
- 完成线上签约，选择相应产品套餐

#### 3. 日期指定方式（商户API）
- 主要参数:
  - `merchantNo`: 商户号
  - `billDate` 或 `startDate`/`endDate`: 账单日期或起止日期
  - 日期格式通常为 `yyyy-MM-dd`

#### 4. 返回字段（商户API）

| 字段类别 | 可能包含的字段 |
|---------|--------------|
| **基础信息** | 交易流水号、商户订单号、交易时间 |
| **金额信息** | 交易金额、支付金额、退款金额、手续费 |
| **交易信息** | 交易类型（支付/退款）、交易状态、支付方式 |
| **账户信息** | 商户号、门店信息、终端信息 |
| **对方信息** | 用户标识、支付账号信息（脱敏） |
| **附加信息** | 商品描述、备注信息、附加数据 |

#### 5. API调用示例（结构示例）

```
// 请求URL
POST https://payapi.jd.com/queryBill

// 请求参数
{
  "merchantNo": "M123456789",
  "startDate": "2025-01-01",
  "endDate": "2025-01-31",
  "signType": "MD5",
  "timestamp": "20250201120000",
  "sign": "a1b2c3d4e5f6..."
}

// 响应数据
{
  "code": "200",
  "message": "success",
  "data": {
    "totalCount": 100,
    "billList": [
      {
        "tradeNo": "JD20250101123456789",
        "orderNo": "ORDER_001",
        "tradeTime": "2025-01-01 10:30:00",
        "tradeType": "PAY",
        "amount": 100.00,
        "status": "SUCCESS",
        "remark": "商品购买"
      }
    ]
  }
}
```

---

## 银行

### 招商银行

#### 1. API 提供情况
- **个人账户账单API**: ❌ **暂不支持个人开发者申请**
- **企业账单API**: ✅ 提供

**官方平台**:
- 招商银行一网通开放平台: https://open.cmbchina.com/
- 开发者文档中心: https://open.cmbchina.com/devDocuments

#### 2. 前提条件
- 必须是**企业法人、个体工商户或其他组织**
- 需要提交的证件:
  - 有效的营业执照（三证合一）
  - 组织机构代码证、税务登记证（或三证合一证件）
  - 法人身份证件
  - 企业营业执照扫描件（需加盖公章）

**申请流程**:
1. 访问招商银行开放平台注册账号
2. 登录后进行企业认证
3. 提交企业认证申请，等待审核（通常1-3个工作日）
4. 审核通过后创建应用，获取API密钥（AppID/Secret）
5. 申请"账户查询"等相关接口权限

**应用场景限制**: API仅适用于企业自身经营场景，禁止用于外包服务、转售或变相为第三方提供服务。

#### 3. 日期指定方式
根据银行业API通用模式，通常通过以下参数指定时间范围:

```json
{
  "acctNo": "账号",
  "startDate": "2026-01-01",
  "endDate": "2026-01-31",
  "currentPage": 1,
  "pageSize": 20
}
```

日期格式通常使用 `YYYY-MM-DD` 格式。

#### 4. 返回字段（推测）

| 字段名 | 说明 | 示例 |
|--------|------|------|
| transDate | 交易日期 | 2026-01-05 |
| transTime | 交易时间 | 10:30:00 |
| transAmt | 交易金额 | 1000.00 |
| currency | 币种 | CNY |
| transType | 交易类型 | 转账支出/转账收入/消费/还款等 |
| counterparty | 交易对手方 | 张三/XX公司 |
| remark | 备注 | 报销款/工资/货款等 |
| balance | 账户余额 | 50000.00 |
| serialNo | 交易流水号 | 202601051000001234 |

#### 5. API调用示例（推测结构）

```bash
POST https://api.cmbchina.com/obp/v1/account/transaction
Content-Type: application/json
Authorization: Bearer {access_token}
X-API-Key: {your_api_key}
X-Timestamp: 1741420800000
X-Signature: {sign_value}

{
  "acctNo": "6226091234567890",
  "startDate": "2026-01-01",
  "endDate": "2026-01-31",
  "currentPage": 1,
  "pageSize": 20
}
```

---

### 工商银行

#### 1. API 提供情况
- **个人账户账单API**: ❌ **通常不直接对外公开**
- **企业/商户账单API**: ✅ 提供

**官方网站**:
- 工商银行开放平台: https://open.icbc.com.cn
- 工商银行开发者社区: https://developer.icbc.com.cn

#### 2. 前提条件

**注册开发者账号**:
- 访问工商银行开放平台注册账号
- 完成邮箱/手机验证

**企业资质认证**（需要提交）:
- 企业营业执照
- 组织机构代码证（或三证合一证件）
- 法人身份证照片
- 企业基本信息

**审核周期**: 通常1-3个工作日

**创建应用并申请权限**:
- 在控制台创建应用
- 选择需要使用的API产品
- 获取API密钥（App ID、App Secret）
- 部分接口需要申请API证书

#### 3. 日期指定方式

**常见日期参数**:
- `BeginDate` / `start_date`: 开始日期（格式: `YYYYMMDD`）
- `EndDate` / `end_date`: 结束日期（格式: `YYYYMMDD`）
- `WorkDate`: 工作日期（部分接口需要）

**示例（对账单文件查询）**:
```xml
<BeginDate>20230501</BeginDate>
<EndDate>20230515</EndDate>
```

#### 4. 返回字段

| 字段类别 | 字段名（常见） | 说明 | 格式示例 |
|---------|--------------|------|---------|
| **交易时间** | `txn_date` / `tran_date` | 交易日期 | 20231027 |
| | `txn_time` | 交易时间 | 143022 |
| **金额信息** | `txn_amt` | 交易金额 | |
| | `cur_code` | 币种 | CNY（人民币）|
| **账户信息** | `acct_no` | 本方账号 | |
| | `acct_name` | 本方账户名 | |
| | `opp_acct_no` | 对方账号 | |
| | `opp_acct_name` | 对方账户名 | |
| **交易信息** | `txn_type` | 交易类型 | 转账/消费/存款等 |
| | `postscript` / `summary` | 摘要/备注 | |
| | `balance` | 交易后余额 | |
| **其他信息** | `seq_no` | 交易流水号 | |
| | `corp_no` | 商户客户号 | |
| | `mer_code` | 商户编号 | |

#### 5. API调用示例

**XML请求示例（对账文件查询）**:
```xml
<?xml version="1.0" encoding="GBK"?>
<ap>
    <CCTransCode>ICBC_DCP_AP_ReconFileDownQry</CCTransCode>
    <WorkDate>20230515</WorkDate>
    <UserID>商户操作员号</UserID>
    <UserName>操作员姓名</UserName>
    <ReqSeqNo>唯一请求流水号</ReqSeqNo>
    <CorpNo>商户客户号</CorpNo>
    <MerCode>商户编号</MerCode>
    <BeginDate>20230501</BeginDate>
    <EndDate>20230515</EndDate>
    <FileType>1</FileType>
</ap>
```

**Python SDK查询示例**:
```python
from icbcbc import ICBCBClient

client = ICBCBClient(
    app_id='你的APP_ID',
    merchant_id='你的商户号',
    private_key_path='private_key.pem',
    icbc_public_key_path='icbc_public_key.pem',
    gateway_url='https://apipc.es.icbc.com.cn'
)

response = client.query_paid_bills(
    begin_date='20230901',
    end_date='20230910'
)
```

---

### 农业银行

#### 1. API 提供情况
- **个人账户账单API**: ❌ **需要企业认证后获取详细信息**
- **企业账单API**: ✅ 提供

**官方平台**:
- 中国农业银行开放银行: https://openbank.abchina.com.cn/Portal/

#### 2. 前提条件

**企业资质**:
- 必须是企业用户，个人开发者通常无法申请
- 提供企业营业执照、组织机构代码证等资质证明

**注册认证**:
- 在开放银行平台注册企业账号
- 完成企业身份认证和验证

**申请权限**:
- 提交API使用申请，说明业务场景
- 申请具体的产品权限（如账户查询、交易明细等）
- 可能需要线下商务洽谈

**技术对接**:
- 配置服务器IP白名单
- 申请数字证书用于接口签名
- 完成联调测试

#### 3. 日期指定方式
具体参数需查阅官方文档（需企业认证后查看）。根据银行业通用模式，通常使用 `startDate`/`endDate` 参数，格式为 `YYYY-MM-DD` 或 `YYYYMMDD`。

#### 4. 返回字段
具体字段需查阅官方文档。基于银行业通用模式，可能包含: 交易流水号、交易时间、金额、币种、余额、对方账户信息、交易类型、备注等。

#### 5. API调用示例
需登录农业银行开放银行平台获取官方文档和示例代码。

---

### 中信银行

#### 1. API 提供情况
- **个人账户账单API**: ❌ **暂不支持个人开发者**
- **企业账单API**: ✅ 提供

**官方开放平台**:
- https://open.ecitic.com （或 https://open.citicbank.com）

#### 2. 前提条件
- 必须是**企业/机构客户**，个人开发者暂无法申请
- 在开放平台注册企业账号
- 提交企业认证资料:
  - 企业基本信息
  - 企业营业执照
  - 组织机构代码证
  - 税务登记证
- 创建应用，获取 AppID 和 AppSecret
- 申请"账户服务"类API权限，包括交易明细查询

#### 3. 日期指定方式
根据银行API通用规范，通常通过以下参数指定时间范围:

```json
{
  "accountNo": "6226xxxxxxxxx1234",
  "startDate": "20260101",
  "endDate": "20260308",
  "pageNo": 1,
  "pageSize": 50
}
```

日期格式通常为 `YYYYMMDD`。

#### 4. 返回字段

| 字段名称 | 说明 | 示例 |
|---------|------|------|
| `transactionId` | 交易流水号 | 202603081234567890 |
| `transactionDate` | 交易日期 | 2026-03-08 |
| `transactionTime` | 交易时间 | 14:30:25 |
| `amount` | 交易金额 | 128.50 |
| `currency` | 币种 | CNY |
| `balance` | 账户余额 | 10000.00 |
| `debitCreditFlag` | 借贷标志 | D(支出)/C(收入) |
| `counterpartyAccountName` | 对方账户名称 | 某某超市 |
| `counterpartyAccountNo` | 对方账号 | 6225xxxxxxxxx5678 |
| `transactionType` | 交易类型 | 消费/转账/理财等 |
| `channel` | 交易渠道 | 手机银行/网银/POS等 |
| `summary` | 交易摘要/备注 | 购物消费 |
| `status` | 交易状态 | 成功/失败/处理中 |

#### 5. API调用示例（通用模式）

```python
import requests
import json
from datetime import datetime

# 配置信息
APP_ID = "your_app_id"
APP_SECRET = "your_app_secret"
API_URL = "https://api.open.citicbank.com/v1/account/transaction/query"

def query_transactions(account_no, start_date, end_date):
    payload = {
        "appId": APP_ID,
        "timestamp": datetime.now().strftime("%Y%m%d%H%M%S"),
        "accountNo": account_no,
        "startDate": start_date,
        "endDate": end_date,
        "pageNo": 1,
        "pageSize": 100
    }

    # 生成签名（实际需按中信银行规范）
    # signature = generate_signature(payload, APP_SECRET)
    # payload["signature"] = signature

    response = requests.post(API_URL, json=payload, timeout=30)
    result = response.json()

    if result.get("code") == "0000":
        return result.get("data", {})
    else:
        print(f"API调用失败: {result.get('message')}")
        return None
```

---

### 浦发银行

#### 1. API 提供情况
- **个人账户账单API**: ❌ **需要企业认证**
- **企业账单API**: ✅ 提供

**官方平台**:
- 浦发银行开放平台: https://open.spdb.com.cn/

#### 2. 前提条件
- 注册开发者账号
- **企业认证**:
  - 营业执照
  - 法人身份证
  - 其他企业相关证明材料
- 申请API权限: 在控制台申请具体的API产品权限
- 安全认证: 生产环境通常需要使用CFCA证书或OAuth2.0令牌

#### 3. 日期指定方式
根据银行业API通用设计模式:

```json
{
  "accountNo": "6225XXXXXXXXXXXX",
  "startDate": "2026-01-01",
  "endDate": "2026-01-31",
  "pageNum": 1,
  "pageSize": 50
}
```

日期格式通常为 `YYYY-MM-DD`。

#### 4. 返回字段

| 字段类别 | 常见字段名 | 说明 |
|---------|-----------|------|
| **基本信息** | `transactionId` / `flowId` | 交易流水号 |
| | `transactionTime` / `txnTime` | 交易时间 |
| | `transactionDate` | 交易日期 |
| **账户信息** | `accountNo` | 本方账号 |
| | `accountName` | 本方户名 |
| | `counterpartyAccountNo` | 对方账号 |
| | `counterpartyAccountName` | 对方户名 |
| | `counterpartyBankName` | 对方开户行 |
| **金额信息** | `amount` / `transactionAmount` | 交易金额 |
| | `currency` | 币种 |
| | `balance` | 交易后余额 |
| | `creditAmount` | 贷方金额（收入）|
| | `debitAmount` | 借方金额（支出）|
| **交易信息** | `transactionType` / `bizType` | 交易类型 |
| | `remark` / `memo` / `summary` | 备注/摘要/附言 |
| | `channel` | 交易渠道 |
| | `status` | 交易状态 |

#### 5. API调用示例

**请求示例**:
```http
POST /api/v1/account/transaction/query
Host: open.spdb.com.cn
Content-Type: application/json
Authorization: Bearer {access_token}
X-API-Key: {your_api_key}
X-Timestamp: 20260308120000
X-Signature: {request_signature}

{
  "accountNo": "6225123456789012",
  "startDate": "2026-01-01",
  "endDate": "2026-01-31",
  "pageNum": 1,
  "pageSize": 20
}
```

---

### 中国银行

#### 1. API 提供情况
- **个人账户账单API**: ❌ **需要企业资质**
- **企业账单API**: ✅ 提供

**相关链接**:
- 中银开放平台: https://open.boc.cn

#### 2. 前提条件
根据银行业通用规范，通常需要:
- **企业资质**: 通常需要是企业法人，个人开发者一般无法申请
- **营业执照**: 需提供有效期内的企业营业执照
- **开发者账号**: 在中银开放平台注册企业开发者账号
- **应用场景审核**: 需说明API使用场景、业务范围等
- **合同签署**: 可能需要签署合作协议和保密协议
- **技术对接**: 按照银行安全规范进行开发（通常涉及签名、加密等）

**建议**: 直接联系中国银行客户经理或开放平台客服获取准确的申请流程。

#### 3. 日期指定方式
根据银行业通用设计，查询账单API通常包含以下与时间相关的参数:

| 参数名 | 格式 | 说明 |
|--------|------|------|
| `startDate` 或 `beginDate` | `YYYY-MM-DD` | 查询起始日期 |
| `endDate` | `YYYY-MM-DD` | 查询结束日期 |
| `accountNo` | 字符串 | 需查询的账户/卡号（可能需要脱敏） |
| `pageSize` | 数字 | 每页返回的记录数（分页用） |
| `pageIndex` 或 `pageNo` | 数字 | 页码（分页用） |

#### 4. 返回字段（基于行业通用规范）

**必选字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| `transactionId` | String | 交易唯一标识/流水号 |
| `transactionDate` | DateTime | 交易日期时间（ISO 8601格式） |
| `amount` | Decimal | 交易金额 |
| `currency` | String | 币种（如CNY、USD） |
| `direction` | String | 交易方向（收入/支出 或 DEBIT/CREDIT） |
| `balance` | Decimal | 交易后账户余额 |
| `counterpartyName` | String | 交易对手名称/商户名称 |
| `transactionType` | String | 交易类型编码（如消费、转账、取现等） |

**可选字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| `counterpartyAccountNo` | String | 对手账号（通常脱敏） |
| `counterpartyBankName` | String | 对手开户行 |
| `remark` 或 `summary` | String | 交易备注/摘要 |
| `terminalInfo` | String | 终端信息（ATM编号、POS终端号等） |
| `status` | String | 交易状态（成功/失败/处理中） |
| `fee` | Decimal | 手续费金额 |
| `exchangeRate` | Decimal | 汇率（外币交易时） |

#### 5. API调用示例（基于通用金融API模式）

**步骤1: 获取Access Token**
```http
POST /oauth2/token HTTP/1.1
Host: open.boc.cn（示例域名）
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials
&client_id=YOUR_CLIENT_ID
&client_secret=YOUR_CLIENT_SECRET
```

**步骤2: 查询账单**
```http
GET /api/v1/accounts/transactions?accountNo=622202********1234&startDate=2026-01-01&endDate=2026-01-31&pageSize=20&pageIndex=1 HTTP/1.1
Host: open.boc.cn（示例域名）
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 总结与建议

### 核心发现

| 机构类型 | 机构名称 | 个人账单API | 企业账单API | 备注 |
|---------|---------|-----------|-----------|------|
| **支付机构** | 支付宝 | ❌ | ✅ | 仅对商户开放 |
| | 微信支付 | ❌ | ✅ | 仅对商户开放 |
| | 京东支付 | ❌ | ✅ | 仅对商户开放 |
| **银行** | 招商银行 | ❌ | ✅ | 需企业资质 |
| | 工商银行 | ❌ | ✅ | 需企业资质 |
| | 农业银行 | ❌ | ✅ | 需企业资质 |
| | 中信银行 | ❌ | ✅ | 需企业资质 |
| | 浦发银行 | ❌ | ✅ | 需企业资质 |
| | 中国银行 | ❌ | ✅ | 需企业资质 |

### 对 account 项目的建议

#### 方案一：继续使用现有的CSV导入方式（推荐）✅

**现状**: 项目已实现完善的CSV导入功能，支持:
- 支付宝账单CSV解析
- 微信支付账单CSV解析
- 银行对账单CSV解析
- 通用CSV解析

**优点**:
- ✅ 已实现，无需额外开发
- ✅ 合规性高，不涉及API授权问题
- ✅ 用户操作路径清晰

**建议优化**:
1. **完善导出指南**: 在UI中提供各平台详细的账单导出步骤指南
2. **增加京东支付支持**: 添加京东账单CSV解析器
3. **批量导入优化**: 支持一次性导入多个CSV文件

#### 方案二：OAuth授权方案（可行性低）

调研发现，目前国内主流支付机构和银行**均未向个人开发者开放个人账户账单的OAuth授权接口**。此方案暂时不可行。

#### 方案三：浏览器自动化（不推荐）⚠️

可以考虑使用Selenium等工具模拟用户登录网银/支付宝导出账单，但存在以下问题:
- ❌ **法律风险**: 可能违反平台服务条款
- ❌ **技术风险**: 页面结构变化会导致功能失效
- ❌ **安全风险**: 需要用户提供账号密码
- ❌ **维护成本**: 需要持续跟进页面变化

**此方案不推荐使用**。

### 最终建议

**继续使用并完善现有的CSV导入方式**，这是目前最合规、最稳定、最可行的方案。

对于"自动记账"功能，可以重新定义为:
- **半自动化**: 引导用户快速导出账单并导入应用
- **智能解析**: 优化CSV解析的准确度和覆盖率
- **导入预览**: 提供更友好的导入预览和确认界面

---

**报告结束**

*本报告基于2026年3月的公开信息调研，API政策可能随时变化，建议以各平台官方文档为准。*
