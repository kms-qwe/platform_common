package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
	pgv1 "github.com/kms-qwe/platform_common/pkg/client/postgres/pg"
)

type manager struct {
	db pgClient.Transactor
}

// NewTransactionManager создает новый менеджер транзакций, который удовлетворяет интерфейсу db.TxManager
func NewTransactionManager(db pgClient.Transactor) pgClient.TxManager {
	return &manager{
		db: db,
	}
}

// ReadCommitted реализация Read Committed транзакцию
func (m *manager) ReadCommitted(ctx context.Context, f pgClient.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

// transaction основная функция, которая выполняет указанный пользователем обработчик в транзакции
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, f pgClient.Handler) (err error) {
	// Если это вложенная транзакция, пропускаем инициацию новой транзакции и выполняем обработчик.
	_, ok := ctx.Value(pgv1.TxKey).(pgx.Tx)
	if ok {
		return f(ctx)
	}

	// Стартуем новую транзакцию.
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	// Кладем транзакцию в контекст.
	ctx = pgv1.MakeContextTx(ctx, tx)

	// Настраиваем функцию отсрочки для отката или коммита транзакции.
	defer func() {
		// восстанавливаемся после паники
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		// откатываем транзакцию, если произошла ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = fmt.Errorf("errRollback: %v: %w", errRollback, err)
			}

			return
		}

		// если ошибок не было, коммитим транзакцию
		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("tx commit falied: %w", err)
			}
		}
	}()

	// Выполните код внутри транзакции.
	// Если функция терпит неудачу, возвращаем ошибку, и функция отсрочки выполняет откат
	// или в противном случае транзакция коммитится.
	if err = f(ctx); err != nil {
		err = fmt.Errorf("failed executing code inside transaction: %w", err)
	}

	return err
}
