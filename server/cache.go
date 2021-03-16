package server

import (
	"bytes"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"licenseplate.wtf/util"
)

var mc *memcache.Client
var illegalChars = regexp.MustCompile("[+ :-]")

func init() {
	mc = memcache.New("127.0.0.1:11211")
}

func cleanKey(bits []string) string {
	return strings.ReplaceAll(
		illegalChars.ReplaceAllString(strings.Join(bits, "/"), "-"),
		"--0000-UTC", "")
}

func (s *server) cache(w io.Writer, keys []string, gen func(io.Writer)) {
	start := time.Now()

	var item *memcache.Item
	var err error
	key := cleanKey(keys)

	// We want to directly write the output without the need for io.Copy
	var output io.Writer

	// We may need to save the bytes though to divert to cache
	var buffer bytes.Buffer

	if s.config.CacheEnabled {
		util.LogTime("fetching from cache", func() {
			item, err = mc.Get(key)
		})

		if err == nil && item != nil {
			log.Println("Cache hit for", key)
			w.Write(item.Value)
			return
		}
		log.Println("Cache miss for", key)

		// If we are caching, also save to buffer
		output = io.MultiWriter(w, &buffer)
	} else {
		output = w
	}

	// Call the user code to generate the output to be cached.
	gen(output)

	if s.config.CacheEnabled {
		util.LogTime("saving to cache", func() {
			item = &memcache.Item{
				Key:        key,
				Value:      buffer.Bytes(),
				Flags:      0,
				Expiration: 0,
			}
			err := mc.Set(item)
			if err != nil {
				log.Println("Problem setting cache:", err)
			}
		})
	}

	log.Println("Completed rendering", key, "in", time.Now().Sub(start))
}
