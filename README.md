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
wrk -t4 -c1000 -d30s --latency 'http://127.0.0.1:8000/api/segment/tinyid-segment-test'
Running 30s test @ http://127.0.0.1:8000/api/segment/tinyid-segment-test
  4 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.55ms    6.61ms 254.60ms   97.56%
    Req/Sec    13.78k     2.12k   18.90k    79.67%
  Latency Distribution
     50%    4.74ms
     75%    6.06ms
     90%    7.37ms
     99%   13.42ms
  1639373 requests in 30.01s, 192.30MB read
  Socket errors: connect 751, read 143, write 0, timeout 0
  Non-2xx or 3xx responses: 1
Requests/sec:  54629.04
Transfer/sec:      6.41MB
```

- Snowflake

```
wrk -t4 -c1000 -d30s --latency 'http://127.0.0.1:8000/api/snowflake'
Running 30s test @ http://127.0.0.1:8000/api/snowflake
  4 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.01ms    3.09ms 112.92ms   75.53%
    Req/Sec    14.97k     1.70k   19.56k    81.25%
  Latency Distribution
     50%    4.14ms
     75%    5.40ms
     90%    6.68ms
     99%   11.85ms
  1787329 requests in 30.01s, 228.41MB read
  Socket errors: connect 751, read 142, write 0, timeout 0
Requests/sec:  59563.64
Transfer/sec:      7.61MB
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