package redisstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dracofs/Aranea/internal/crawler"
	goredis "github.com/redis/go-redis/v9"
)

func (f *Frontier) Enqueue(ctx context.Context, page crawler.Page) (bool, error) {
	payload, err := json.Marshal(page)
	if err != nil {
		return false, fmt.Errorf("encode page: %w", err)
	}

	result, err := enqueueScript.Run(
		ctx,
		f.client,
		[]string{f.seenKey, f.queueKey},
		page.URL,
		string(payload),
	).Int64()
	if err != nil {
		return false, fmt.Errorf("enqueue page: %w", err)
	}
	return result == 1, nil
}

func (f *Frontier) Claim(ctx context.Context) (crawler.Page, error) {
	for {
		payload, err := f.client.BLMove(
			ctx,
			f.queueKey,
			f.inFlightKey,
			"RIGHT",
			"LEFT",
			f.claimTimeout,
		).Result()
		if errors.Is(err, goredis.Nil) {
			if ctx.Err() != nil {
				return crawler.Page{}, ctx.Err()
			}
			continue
		}
		if err != nil {
			return crawler.Page{}, err
		}

		var page crawler.Page
		if err := json.Unmarshal([]byte(payload), &page); err != nil {
			return crawler.Page{}, fmt.Errorf("decode claimed page: %w", err)
		}
		page.Receipt = payload
		return page, nil
	}
}

func (f *Frontier) Ack(ctx context.Context, page crawler.Page) error {
	if page.Receipt == "" {
		return errors.New("acknowledge page: missing receipt")
	}
	removed, err := f.client.LRem(ctx, f.inFlightKey, 1, page.Receipt).Result()
	if err != nil {
		return fmt.Errorf("acknowledge page: %w", err)
	}
	if removed != 1 {
		return errors.New("acknowledge page: receipt not found")
	}
	return nil
}
