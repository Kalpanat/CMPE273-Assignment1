# CMPE273-Assignment1

This is a JSON-RPC Stock Trading System running HTTP Server and Client using Gorilla json/rpc library

## Usage

### Install

```
go get github.com/kalpanat/CMPE273-Assignment1
```

### Start the  server:

```
cd CMPE273-Assignment1/json-rpc
go run server.go

```

### Start the  client:

```
cd CMPE273-Assignment1/json-rpc
go run client.go

```

### Make RPC calls for buying stock using curl
```
curl  -H "Content-Type: application/json"  -d '{"method":"FinanceService.Response","params":[{"StockMessage":"FB:20%;GOOG:40%,NFLX:30%,AAPL:10%", "Budget":10000}], "id": 0}' http://localhost:8080/rpc

```

### Result
```
{"result":{"Message":"TradeId : 1812 StockString: GOOG:6:$626.91,NFLX:28:$106.11,AAPL:9:$110.38,FB:21:$92.07 UnvestedAmount : 340.570000\n"},"error":null,"id":0}

```
### Make RPC calls for checking profile using curl

```
curl  -H "Content-Type: application/json"  -d '{"method":"FinanceService.TradeRequest","params":[{"TradeId":9029}], "id": 0}' http://localhost:8080/rpc

```

### Result
```
{"result":{"Message":"StockString: FB:21:$92.07 ,GOOG:6:$626.91 ,NFLX:28:$106.11 ,AAPL:9:$110.38  CurrentMarketValue: 9659.429688 UnvestedAmount: 340.570007"},"error":null,"id":0}

```

### Also supports compressed format

```
curl -sv --compressed -H "Content-Type: application/json"  -d '{"method":"FinanceService.Response","params":
[{"StockMessage":"FB:20%;GOOG:40%,NFLX:30%,AAPL:10%", "Budget":10000}], "id": 0}' http://localhost:8080/rpc
```

```
{"result":{"Message":"TradeId : 9029 StockString: FB:21:$92.07,GOOG:6:$626.91,NFLX:28:$106.11,AAPL:9:$110.38 UnvestedAmount : 340.570000\n"},"error":null,"id":0}
* timeout on name lookup is not supported
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /rpc HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.44.0
> Accept: */*
> Accept-Encoding: deflate, gzip
> Content-Type: application/json
> Content-Length: 127
>
} [127 bytes data]
* upload completely sent off: 127 out of 127 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Sun, 04 Oct 2015 06:10:41 GMT
< Content-Length: 162
<
{ [162 bytes data]
* Connection #0 to host localhost left intact
```