// Package maintenance запускает фоновую нормализацию данных БД.
package maintenance

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/config"
)

// Normalizer периодически вызывает fn_normalize_employee_states(): сливает
// смежные диапазоны состояний сотрудников (employee_states) в непрерывные.
type Normalizer struct {
	cfg    config.MaintenanceConfig
	logger *slog.Logger
	pool   *pgxpool.Pool
	cancel context.CancelFunc
	done   chan struct{}
}

// New builds the normalizer.
func New(cfg config.MaintenanceConfig, pool *pgxpool.Pool, logger *slog.Logger) *Normalizer {
	return &Normalizer{cfg: cfg, logger: logger, pool: pool}
}

// Start launches the background loop if maintenance is enabled.
func (n *Normalizer) Start() {
	if !n.cfg.Enabled {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	n.cancel = cancel
	n.done = make(chan struct{})
	go n.loop(ctx)
	n.logger.Info("maintenance normalizer started", "interval", n.cfg.Interval)
}

func (n *Normalizer) loop(ctx context.Context) {
	defer close(n.done)
	// Первый прогон сразу после старта (могли накопиться несведённые диапазоны),
	// затем — по таймеру.
	n.normalize(ctx)
	ticker := time.NewTicker(time.Duration(n.cfg.Interval))
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n.normalize(ctx)
		}
	}
}

func (n *Normalizer) normalize(ctx context.Context) {
	if _, err := n.pool.Exec(ctx, "SELECT fn_normalize_employee_states()"); err != nil {
		n.logger.ErrorContext(ctx, "employee states normalization failed", "error", err)
		return
	}
	n.logger.InfoContext(ctx, "employee states normalized")
}

// Stop cancels the background loop and waits for it to finish.
func (n *Normalizer) Stop(ctx context.Context) {
	if n.cancel == nil {
		return
	}
	n.cancel()
	select {
	case <-n.done:
	case <-ctx.Done():
	}
}
