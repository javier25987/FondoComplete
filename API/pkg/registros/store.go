package registros

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, pool *pgxpool.Pool) *Store {
	return &Store{ctx: ctx, pool: pool}
}

func (s *Store) GetMes(mes int) (*TablaMeses, error) {
	localCtx, cancel := context.WithTimeout(s.ctx, time.Second*5)
	defer cancel()

	var result TablaMeses

	q := "SELECT fecha, ingreso, egreso FROM registros WHERE mes = ?"

	rows, err := s.pool.Query(localCtx, q, mes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		fecha      string
		ingreso    int
		egreso     int
		diferencia int
		muestra    string
	)

	for rows.Next() {
		if err := rows.Scan(&fecha, &ingreso, &egreso); err != nil {
			return nil, err
		}

		diferencia = ingreso - egreso
		if diferencia > 0 {
			muestra = "✅"
		} else if diferencia == 0 {
			muestra = "⭐"
		} else {
			muestra = "🚨"
		}

		result.Fecha = append(result.Fecha, fecha)
		result.Ingreso = append(result.Ingreso, ingreso)
		result.Egreso = append(result.Egreso, egreso)
		result.Diferencia = append(result.Diferencia, diferencia)
		result.Muestra = append(result.Muestra, muestra)
	}

	return &result, nil
}
