package postgres

import (
	"context"
	"fmt"
	"github.com/alimzhanoff/property-finder/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, cfg config.DatabaseConfig) *Postgres {
	const op = "storage.property.postgres.NewPG"
	var err error
	pgOnce.Do(func() {
		pgInstance = &Postgres{}
		pgInstance.db, err = new(ctx, cfg)
		if err != nil {
			log.Fatalf("%s: cannot connect to postgres db: %v", op, err)
		}
	})

	return pgInstance
}
func new(ctx context.Context, config config.DatabaseConfig) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName)

	log.Println(dsn)
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	log.Println(err)
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
}
