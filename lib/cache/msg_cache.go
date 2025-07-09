/*
*

	@author:
	@date : 2025/5/16
*/
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/link1st/gowebsocket/v2/lib/redislib"
	"github.com/link1st/gowebsocket/v2/models/message"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	MsgOfflinePrefix = "acc:msg:offline:" // 数据不重复提交
	MsgCacheTime     = 24 * 60 * 60 * 30 * 12
)

// getMsgOfflineKey 获取离线信息
func getMsgOfflineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", MsgOfflinePrefix, userKey)
	return
}

// PopOfflineMsg 获取用户离线数据
func PopOfflineMsg(userKey string) (msg *message.Resp, err error) {
	redisClient := redislib.GetClient()
	key := getMsgOfflineKey(userKey)
	ctx := context.Background()
	// 从队列左边弹出一个元素（FIFO）
	data, err := redisClient.LPop(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("队列为空", userKey, err)
			return nil, nil
		}
		fmt.Println("获取离线消息失败", userKey, err)
		return nil, err
	}
	msg = &message.Resp{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println("获取用户离线数据 json Unmarshal", userKey, err)
		return
	}
	fmt.Printf("获取用户离线数据 %v", msg)
	return
}

// PushOfflineMsg 设置用户在线数据
func PushOfflineMsg(userKey string, msg string) (err error) {
	redisClient := redislib.GetClient()
	key := getMsgOfflineKey(userKey)
	//string 转结构体
	//valueByte, err := json.Marshal(msg)
	//if err != nil {
	//	fmt.Println("添加用户离线数据 json Marshal", key, err)
	//	return
	//}
	_, err = redisClient.RPush(context.Background(), key, msg).Result()
	if err != nil {
		fmt.Println("保存用户离线数据", key, err)
		return
	}
	// 设置队列整体过期时间
	redisClient.Expire(context.Background(), key, time.Second*MsgCacheTime)
	return
}

// PopAllOfflineMsgs 取出所有离线消息并自动销毁
func PopAllOfflineMsgs(userKey string) ([]*message.Resp, error) {
	redisClient := redislib.GetClient()
	key := getMsgOfflineKey(userKey)
	ctx := context.Background()

	// 获取所有消息
	dataList, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println("获取离线消息列表失败", key, err)
		return nil, err
	}

	// 删除整个 key（取出后清空）
	_ = redisClient.Del(ctx, key).Err()

	msgs := make([]*message.Resp, 0, len(dataList))
	for _, data := range dataList {
		msg := &message.Resp{}
		err = json.Unmarshal([]byte(data), msg)
		if err != nil {
			fmt.Println("解析离线消息失败", key, err)
			continue
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
