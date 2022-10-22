package db

import (
	"ModeFlex/api"
	"ModeFlex/data"
	"context"
	"time"
)

func GetStats() (w *api.Stats, e error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	defer cancel()
	row := data.Database.QueryRowContext(ctx, "SELECT * FROM stats")

	err := row.Scan(&w.TotalRenders)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func IncrTotalRenders() error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	defer cancel()
	_, err := data.Database.ExecContext(ctx, "UPDATE stats SET total_renders = total_renders + 1")
	if err != nil {
		return err
	}
	return nil
}
