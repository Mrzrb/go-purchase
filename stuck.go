package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Stucker interface {
	Purchase()
	Counter() int
}

const redisKey = "structcnt"

var ctx = context.Background()

type Stuck struct {
	Count    int
	cache    *redis.Client
	resultCh chan struct{}
	mu       sync.RWMutex
}

func NewStuck(n int) *Stuck {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	rdb.Set(ctx, redisKey, n, time.Hour)
	s := &Stuck{
		Count:    n,
		cache:    rdb,
		resultCh: make(chan struct{}),
	}
	go func() {
		for {
			<-s.resultCh
			s.Count--
			log.Printf("success purchase one")
		}
	}()

	go func() {
		timer := time.NewTicker(time.Second)
		for range timer.C {
			log.Printf("monitor %d left", s.Count)
		}
	}()
	return s
}

// redis CAS
var purchaseLua = redis.NewScript(`
    local cnt = tonumber(redis.call("GET", KEYS[1]))
    if cnt <= 0 then
        return false
    else
        redis.call("SET", KEYS[1], cnt-1)
        return true
    end

`)

func (s *Stuck) Purchase() {
	//模拟数据库延迟
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := purchaseLua.Run(ctx, s.cache, []string{fmt.Sprintf("%s", redisKey)}).Result()

	if err != nil {
		log.Printf("err purchase %s", err.Error())
		return
	}
	succ, ok := res.(int64)
	log.Println(succ, ok)
	if ok && succ == 1 {
		s.resultCh <- struct{}{}
	}
}

func (s *Stuck) Counter() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cnt, err := s.cache.Get(ctx, fmt.Sprintf("%s", redisKey)).Result()
	if err != nil {
		log.Printf("error : %s", err)
	}
	cntn, _ := strconv.Atoi(cnt)
	log.Printf("left %d\n", cntn)
	return s.Count
}
