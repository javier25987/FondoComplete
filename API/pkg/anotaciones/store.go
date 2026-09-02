package anotaciones

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

func (s *Store) CrearAnotacion(id int, nota *PostAnotacion) error {
	localCtx, cancel := context.WithTimeout(s.ctx, time.Second*5)
	defer cancel()

	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)

	q := `insert into anotaciones (id_user, anotacion, monto, tipo) values ($1, $2, $3, $4)`
	_, err = tx.Exec(localCtx, q, id, nota.Nota, nota.Monto, nota.Tipo)
	if err != nil {
		return err
	}

	if nota.Tipo == "MONETARIA" {
		q = `UPDATE usuarios SET multas_extras = multas_extras + $1 WHERE id = $2`
		_, err = tx.Exec(localCtx, q, nota.Monto, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit(s.ctx)
}

func (s *Store) ObtenerAnotaciones(id int) ([]*Anotacion, error) {
	localCtx, cancel := context.WithTimeout(s.ctx, time.Second*5)
	defer cancel()

	q := `SELECT codigo, anotacion, monto, tipo, fecha FROM anotaciones WHERE id_user = $1`
	rows, err := s.pool.Query(localCtx, q, id) // <<<<<-----------------------------------------------------------------------------------------------------
	if err != nil {
		return nil, err
	}

	var (
		codigo    int
		anotacion string
		monto     int
		tipo      string
		fecha     time.Time
	)

	result := []*Anotacion{}

	for rows.Next() {
		err := rows.Scan(&codigo, &anotacion, &monto, &tipo, &fecha)
		if err != nil {
			return nil, err
		}

		result = append(result, &Anotacion{
			Codigo: codigo,
			IdUser: id,
			Nota:   anotacion,
			Monto:  monto,
			Tipo:   tipo,
			Fecha:  fecha.Format("2006/01/02"),
		})
	}
	return result, nil
}
