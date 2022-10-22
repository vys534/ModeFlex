package data

import (
	"ModeFlex/api"
	"image"
	"net/http"
	"sync"
	"time"
)

type lastUpdatedStruct struct {
	Map map[int64]int64
	mu  sync.Mutex
}

var LastUpdatedMap = &lastUpdatedStruct{Map: make(map[int64]int64, 0)}

var Configuration *api.Configuration

var ClientTokenExpires int64

var ClientToken string

var OsuHTTPClient = &http.Client{Timeout: 5 * time.Second}

var ModeArray = []string{"osu", "taiko", "fruits", "mania"}

var IconOsu *image.RGBA
var IconTaiko *image.RGBA
var IconFruits *image.RGBA
var IconMania *image.RGBA

func (c *lastUpdatedStruct) Set(k int64, v int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	LastUpdatedMap.Map[k] = v
}
