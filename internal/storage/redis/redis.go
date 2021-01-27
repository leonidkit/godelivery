package redis

import (
	"context"
	"fmt"
	"godelivery/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type DB struct {
	conn   redis.Conn
	expire time.Duration
}

func New(connUrl string, expireTime time.Duration) (*DB, error) {
	redisConn, err := redis.DialURL(connUrl)
	if err != nil {
		return &DB{}, fmt.Errorf("connecting to redis error: %v", err)
	}

	return &DB{
		conn:   redisConn,
		expire: expireTime,
	}, err
}

// Add item to redis db
func (d *DB) Create(ctx context.Context, item storage.Item) error {
	mkey := strings.ToUpper(item.FormatType) + strconv.FormatInt(item.ID, 10)
	res, err := redis.String(
		d.conn.Do("SET", mkey, item.Format, "EX", d.expire.Seconds()))
	if err != nil {
		return fmt.Errorf("inserting into redis error: %v", err)
	}

	if res != "OK" {
		return fmt.Errorf("result inserting not OK")
	}

	return err
}

// Remove item from redis db by ID
func (d *DB) Get(ctx context.Context, id int64, formatType string) (storage.Item, error) {
	mkey := strings.ToUpper(formatType) + strconv.FormatInt(id, 10)
	format, err := redis.Bytes(d.conn.Do("GET", mkey))
	if err != nil {
		return storage.Item{}, fmt.Errorf("get item from redis error: %v", err)
	}

	return storage.Item{
		ID:     id,
		Format: format,
	}, err
}

// Select item from redis db by ID
func (d *DB) Delete(ctx context.Context, id int64, formatType string) error {
	mkey := strings.ToUpper(formatType) + strconv.FormatInt(id, 10)
	_, err := redis.Int(d.conn.Do("DEL", mkey))
	if err != nil {
		return fmt.Errorf("delete item from redis error: %v", err)
	}

	return err
}

// Close opened connection
func (d *DB) Close() error {
	return d.conn.Close()
}
