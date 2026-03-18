package scaffoldtmpl

// CacheTmpl 缓存层模板
var CacheTmpl = `package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"{{.Name}}/internal/domain"
)

// {{.NameUpper}}Cache {{.Name}} 缓存接口
type {{.NameUpper}}Cache interface {
	Get(ctx context.Context, id string) (*domain.{{.NameUpper}}, error)
	Set(ctx context.Context, u *domain.{{.NameUpper}}) error
	Delete(ctx context.Context, id string) error
}

// Redis{{.NameUpper}}Cache Redis 缓存实现
type Redis{{.NameUpper}}Cache struct {
	client     redis.Cmdable
	expiration time.Duration
}

// New{{.NameUpper}}Cache 创建缓存
func New{{.NameUpper}}Cache(client redis.Cmdable) {{.NameUpper}}Cache {
	return &Redis{{.NameUpper}}Cache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// key 生成缓存 key
// 遵循格式：webook:{{.Name}}:profile:id
func (c *Redis{{.NameUpper}}Cache) key(id string) string {
	return fmt.Sprintf("webook:{{.Name}}:profile:%s", id)
}

func (c *Redis{{.NameUpper}}Cache) Get(ctx context.Context, id string) (*domain.{{.NameUpper}}, error) {
	key := c.key(id)
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var u domain.{{.NameUpper}}
	err = json.Unmarshal(val, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Redis{{.NameUpper}}Cache) Set(ctx context.Context, u *domain.{{.NameUpper}}) error {
	key := c.key(u.ID)
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, val, c.expiration).Err()
}

func (c *Redis{{.NameUpper}}Cache) Delete(ctx context.Context, id string) error {
	key := c.key(id)
	return c.client.Del(ctx, key).Err()
}
`
