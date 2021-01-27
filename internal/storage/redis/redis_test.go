package redis

import (
	"bytes"
	"context"
	"godelivery/internal/storage"
	"log"
	"os"
	"testing"
	"time"
)

var (
	db   = &DB{}
	ctx  = context.Background()
	item = storage.Item{
		ID:         1234,
		FormatType: "SF",
		Format:     []byte("<hello>world</hello>"),
	}
)

func TestMain(m *testing.M) {
	var err error

	db, err = New("redis://default@localhost:6379/0", time.Second*10)
	if err != nil {
		log.Fatal(err)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestCreate(t *testing.T) {
	err := db.Create(ctx, item)
	if err != nil {
		t.Fatalf("recieved error, but not expected: %v", err)
	}

	newItem, err := db.Get(ctx, item.ID, item.FormatType)
	if err != nil {
		t.Fatalf("recieved error, but not expected: %v", err)
	}

	if bytes.Compare(newItem.Format, item.Format) != 0 {
		t.Fatalf("want Format=%s, but get Fromat=%s", string(item.Format), string(newItem.Format))
	}
}

func TestGet(t *testing.T) {
	newItem, err := db.Get(ctx, item.ID, item.FormatType)
	if err != nil {
		t.Fatalf("recieved error, but not expected: %v", err)
	}

	if bytes.Compare(newItem.Format, item.Format) != 0 {
		t.Fatalf("want Format=%s, but get Fromat=%s", string(item.Format), string(newItem.Format))
	}
}

func TestDelete(t *testing.T) {
	err := db.Delete(ctx, item.ID, item.FormatType)
	if err != nil {
		t.Fatalf("recieved error, but not expected: %v", err)
	}

	newItem, err := db.Get(ctx, item.ID, item.FormatType)
	if err == nil {
		t.Fatalf("want error, but get: %v", newItem)
	}
}
