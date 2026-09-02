package registros

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandlerRegistros(ctx context.Context, pool *pgxpool.Pool) *Handlers {
	return &Handlers{
		_service: NewService(NewStore(ctx, pool)),
	}
}
