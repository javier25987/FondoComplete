package anotaciones

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandlerAnotaciones(ctx context.Context, pool *pgxpool.Pool) *Handlers {
	return &Handlers{
		_service: NewService(NewStore(ctx, pool)),
	}
}
