package server

import (
	"bytes"
	"github.com/bradfitz/gomemcache/memcache"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

var mc *memcache.Client
var illegalChars = regexp.MustCompile("[+ :-]")

func init() {
	mc = memcache.New("127.0.0.1:11211")
}

func cleanKey(bits []string) string {
	return illegalChars.ReplaceAllString(strings.Join(bits, "/"), "-")
}

func (s *server) cache(w io.Writer, keys []string, gen func() io.Reader) {
	start := time.Now()

	var item *memcache.Item
	var err error
	key := cleanKey(keys)

	if s.config.CacheEnabled {
		item, err = mc.Get(key)

		if err == nil && item != nil {
			log.Println("Cache hit for", key)
			w.Write(item.Value)
			return
		}
		log.Println("Cache miss for", key)
	}

	content := gen()

	if s.config.CacheEnabled {
		var buffer bytes.Buffer
		output := io.MultiWriter(w, &buffer)

		io.Copy(output, content)

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
	} else {
		io.Copy(w, content)
	}

	log.Println("Completed rendering", key, "in", time.Now().Sub(start))
}
