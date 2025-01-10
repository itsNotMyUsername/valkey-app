package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/valkey-io/valkey-go"
)

var (
	client       valkey.Client
	ctx          = context.Background()
	redisAddress = "localhost:6379"
)

type keysStore struct {
	mu   sync.Mutex
	data map[string]struct{}
}

func (ks *keysStore) Add(key string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.data[key] = struct{}{}
}

func (ks *keysStore) Get() []string {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	keys := make([]string, 0)
	for key := range ks.data {
		keys = append(keys, key)
	}
	return keys
}

var keys *keysStore

func main() {
	keys = &keysStore{
		mu:   sync.Mutex{},
		data: make(map[string]struct{}),
	}

	if os.Getenv("REDIS_ADDRESS") != "" {
		redisAddress = os.Getenv("REDIS_ADDRESS")
	}

	var err error
	client, err = valkey.NewClient(valkey.ClientOption{
		TLSConfig: nil,
		//Username:  "default",
		//Password:              "",
		InitAddress: []string{redisAddress},
	})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	for {
		now := time.Now()
		if err := write(); err != nil {
			log.Println(err)
			return
		}
		value, err := read()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("value", value, "elapsed", time.Since(now).Seconds())
		time.Sleep(time.Second)
	}
}

func write() error {
	value, err := read()
	if err != nil {
		if !strings.Contains(err.Error(), "valkey nil message") {
			log.Fatal(err)
		}
	}
	if value != "" {
		fmt.Println("value", value)
		return nil
	}
	fmt.Println("value not found")

	time.Sleep(2 * time.Second)

	key := gofakeit.Adjective()
	keys.Add(key)

	err = client.Do(ctx, client.B().Set().Key(key).Value("value").Build()).Error()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func read() (string, error) {
	k := keys.Get()
	for _, key := range k {
		value, err := client.Do(ctx, client.B().Get().Key(key).Build()).ToString()
		if err != nil {
			return "", err
		}
		return value, nil
	}

	return "", nil
}
