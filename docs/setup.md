### 启动redis-server
```text
默认启动redis server即可；
redis-server &
```

### 编译
```text
1. 进入项目根目录，执行make编译；
2. 会在 ./bin/ 目录下生成一个 memo 可执行文件；
```

### 启动 
```text
 1. 关注配置文件, 目前的配置文件支持 从10.0.1.102机器sub 比特币全节点数据；
 2. sh ./start.sh # 此步骤会一次启动 daemon api web三个服务；
```

### 配置文件
```toml

# This example config contains all avaliable configuration options.
# Options are assigned their default or left blank with a FIXME: mark.


# ---- needed for both the API and the daemon ----
[redis]
host = ""                 # FIXME:
port = "6379"             # FIXME:
passwd = ""               # FIXME:
connection = "redis"


# ---- needed only for the API ----
[api]
port = "23485" # this needs to be a string
production = false


# ---- needed only for the deamon ----

# Needed when fetching via REST interface
[bitcoind.rest]
protocol = "http"
host = "10.0.1.102"
port = "8332"
responseTimeout = 30      # seconds

# Needed when fetching via JSON-RPC interface
[bitcoind.jsonrpc]
protocol = "http"
host = "10.0.1.102"
port = "8332" # regtest 18443 # testnet 18332
username = "user"
password ="password"
responseTimeout = 30 # seconds

[mempool]
enable = true            # enable mempool fetching
fetchInterface = "JSON-RPC"   # fetch mempool via REST or JSON-RPC
fetchInterval = 60        # seconds
callSaveMempool = false    # calls the savemempool RPC on average after every fourth fetch
[mempool.processing]
processCurrentMempool = true
processHistoricalMempool = true
processTransactionStats = true

[feeratefetcher]
enable = true
fetchInterval = 180

[log]
enableTrace = true
colorizeOutput = true

[zmq]
enable = true            # enable the zmq interface
host = "10.0.1.102"        # currently only tcp connections are supported
port = "28832"
[zmq.subscribeTo]
rawTx     = true
rawBlock  = true        # needed to write recentBlocks to database 
hashTx    = false
hashBlock = true
rawTx2    = true       # this needs a custom bitcoind patch https://github.com/0xB10C/bitcoin/tree/2019-06-rawtx2-zmq-for-memod to work
[zmq.saveMempoolEntries]
enable = true
dbPath = "./sqlite.db"

[mysql]
host = "localhost"
port = "3306"
user = "user"
password = "password"
db = "bitcoin"
```

### 比特币全节点
```shell
1. git clone github.com/charlist/bitcoin # commit id:3cccbe442f84433321b4b0c2544fe79e1f8ab712

2. 按照docs目录下的编译手册，进行编译；

3. 启动命令如下：
./bitcoind  -rpcuser=user -rpcpassword=password -server -txindex -zmqpubrawtx2=tcp://10.0.1.102:28832 -zmqpubrawblock=tcp://10.0.1.102:28832 -rpcport=8332 -rpcallowip=10.0.1.102/255.255.255.0 -rpcbind=10.0.1.102:8332
```