package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Stucker interface {
	Purchase()
	Counter() int
}

type Stuck struct {
	Count int
	mu    sync.RWMutex
}

func NewStuck(n int) *Stuck {
	return &Stuck{
		Count: n,
	}
}

func (s *Stuck) Purchase() {
	//模拟数据库延迟
	s.mu.Lock()
	defer s.mu.Unlock()
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(200)))
	s.Count--
}

func (s *Stuck) Counter() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Printf("left %d\n", s.Count)
	return s.Count
}

func setUpRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/purchase", Purchase)
	return r
}

func Purchase(c *gin.Context) {
	if TestStucker.Counter() > 0 {
		TestStucker.Purchase()
	}
}

var TestStucker = NewStuck(30000)

func main() {
	r := setUpRouter()

	r.Run(":9999")
}
