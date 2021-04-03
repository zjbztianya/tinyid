package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
	pb "tinyid/api/idgen/v1"
	"tinyid/internal/conf"
	"tinyid/pkg/ip"
)

const (
	reportInterval = 3 * time.Second
	minBackwards   = 3 * time.Second
	workerIDPath   = "/id/snowflake"
)

func NewEtcd(c *conf.Data_Etcd) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: c.DialTimeout.AsDuration()})
	if err != nil {
		panic(errors.Wrap(err, "init etcd error"))
	}
	return cli
}

// InitEtcd
func InitEtcd(ctx context.Context, data *Data, rpcAddr string) error {
	lease := clientv3.NewLease(data.etcd)
	kv := data.kv
	grant, err := lease.Grant(ctx, int64(data.ttl.Seconds()))
	if err != nil {
		return err
	}

	_, port, _ := net.SplitHostPort(rpcAddr)
	nodeKey := workerIDPath + "/" + ip.InternalIP() + ":" + port
	resp, err := kv.Get(ctx, workerIDPath, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	//不存在根节点，机器首次启动,直接上报时间
	if resp.Count == 0 {
		_, err := kv.Put(ctx, nodeKey, strconv.FormatInt(time.Now().UnixNano(), 10), clientv3.WithLease(grant.ID))
		if err != nil {
			return err
		}
	} else {
		var workerID int64 = -1
		for i, kv := range resp.Kvs {
			if string(kv.Key) == nodeKey {
				workerID = int64(i)
				timestamp, _ := strconv.ParseInt(string(kv.Value), 10, 64)
				nodeTime, now := time.Unix(0, timestamp), time.Now()
				if nodeTime.After(now) {
					//时间回拨过大,启动失败
					interval := nodeTime.Sub(now)
					if interval > minBackwards {
						return errors.Errorf("clock backwards to large, interval:%v", interval)
					}
					// 回拨时间短，sleep一会
					time.Sleep(interval + 100*time.Millisecond)
				}
				break
			}
		}

		if workerID == -1 {
			workerID = resp.Count
		}
		data.workerID = workerID
		//给正在提供服务的其他机器发送PRC，得到平均的系统时间,超过阈值,启动失败
		avgTime := avgSysTime(ctx, workerID, resp)
		if math.Abs(float64(time.Since(time.Unix(0, avgTime)))) > float64(minBackwards) {
			return errors.Errorf("sys time inaccuracy, avg:%v now:%v", time.Unix(0, avgTime), time.Now())
		}

		_, err = kv.Put(ctx, nodeKey, strconv.FormatInt(time.Now().UnixNano(), 10), clientv3.WithLease(grant.ID))
		if err != nil {
			return err
		}
	}

	hb, err := data.etcd.KeepAlive(ctx, grant.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case _, ok := <-hb:
				if !ok {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	//定时上报自己的时间
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			time.Sleep(reportInterval)
			kv.Put(ctx, nodeKey, strconv.FormatInt(time.Now().UnixNano(), 10), clientv3.WithLease(grant.ID))
		}
	}()
	return nil
}

func avgSysTime(ctx context.Context, workerID int64, resp *clientv3.GetResponse) int64 {
	var (
		total    time.Duration
		nodeSize int64
	)

	for i, kv := range resp.Kvs {
		if workerID == int64(i) {
			continue
		}
		key := string(kv.Key)
		pos := strings.LastIndexByte(key, '/')
		endpoint := key[pos:]
		go func() {
			sctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
			defer cancel()
			conn, err := grpc.Dial(sctx, grpc.WithEndpoint(endpoint))
			if err != nil {
				return
			}
			c := pb.NewIdgenClient(conn)
			reply, err := c.CurrentTime(sctx, &pb.CurrentTimeRequest{})
			if err != nil {
				return
			}
			nodeSize++
			total += reply.Time.AsDuration()
		}()
	}
	if nodeSize == 0 {
		nodeSize++
		total += time.Duration(time.Now().UnixNano())
	}

	return int64(float64(total.Nanoseconds()) / float64(nodeSize))
}
