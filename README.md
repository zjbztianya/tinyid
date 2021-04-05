### Introduction
Tinyid是用Golang开发的一个分布式ID生成服务，参考了[美团Leaf](https://tech.meituan.com/2017/04/21/mt-leaf.html)
的设计，同时支持**DB号段模式**和**snowflake模式**。

### Features
- 趋势递增
- HTTP和GRPC两种访问方式
- DB段号模式
	- 双buffer缓存优化
	- 自适应步长
- snowflake模式
	- 时钟回拨处理
	- 多节点时钟校验
    - etcd分配workerID

### Benchmark
机器配置为8C16G的MBP
- DB段号
   
```
wrk -t16 -c400 -d1m --latency 'http://127.0.0.1:8000/api/segment/tinyid-segment-test'
Running 1m test @ http://127.0.0.1:8000/api/segment/tinyid-segment-test
  16 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.75ms    2.21ms  91.17ms   78.65%
    Req/Sec     5.28k   627.41     8.08k    68.73%
  Latency Distribution
     50%    4.56ms
     75%    5.67ms
     90%    7.03ms
     99%   11.86ms
  5048175 requests in 1.00m, 596.98MB read
  Socket errors: connect 0, read 381, write 0, timeout 0
Requests/sec:  84099.87
Transfer/sec:      9.95MB
```

- Snowflake

```
wrk -t16 -c400 -d1m --latency 'http://127.0.0.1:8000/api/snowflake'
Running 1m test @ http://127.0.0.1:8000/api/snowflake
  16 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.17ms    2.18ms 108.75ms   80.78%
    Req/Sec     4.83k   747.23     7.47k    66.00%
  Latency Distribution
     50%    5.18ms
     75%    6.21ms
     90%    7.13ms
     99%   10.91ms
  4610760 requests in 1.00m, 589.22MB read
  Socket errors: connect 0, read 384, write 0, timeout 0
Requests/sec:  76808.95
Transfer/sec:      9.82MB
```

### Getting Started
本服务基于[Kratos V2](https://github.com/go-kratos/kratos) 框架

#### Required
- go
- mysql
- etcd

#### Build and Run
```
cd tinyid
go mod download
go build -o ./bin/ ./...
./bin/tinyid  -conf ./configs
```
#### Test

```sql
insert into tinyid_alloc(biz_tag, max_id, step, description) values('tinyid-segment-test', 1, 2000, 'Test Tinyid Segment Mode Get ID')
```
- DB段号模式api调用

```
curl -i 'http://127.0.0.1:8000/api/segment/tinyid-segment-test'
```
返回如下:

```
HTTP/1.1 200 OK
Content-Length: 15
Connection: keep-alive
Content-Type: application/json
Date: Sun, 04 Apr 2021 09:47:16 GMT
Keep-Alive: timeout=4
Proxy-Connection: keep-alive

{"id":2}
```
- snowflake模式api调用

```
curl -i 'http://127.0.0.1:8000/api/snowflake'
```
返回如下:

```
HTTP/1.1 200 OK
Content-Length: 26
Connection: keep-alive
Content-Type: application/json
Date: Sun, 04 Apr 2021 10:00:07 GMT
Keep-Alive: timeout=4
Proxy-Connection: keep-alive

{"id":1378648567383588864}
```