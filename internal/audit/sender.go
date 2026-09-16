package audit

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Koshsky/erp-backend/internal/config"
)

// Buffered sender constants.
const (
	// senderQueueCapacity bounds the in-memory event buffer (drop + log when
	// full — audit must never block or fail the user request).
	senderQueueCapacity = 1024
	// senderMaxAttempts is the retry budget for one event.
	senderMaxAttempts = 3
	// senderBaseBackoff is the initial retry delay (doubled per attempt).
	senderBaseBackoff = 50 * time.Millisecond
	// clientTimeoutForRetry bounds each retry attempt context (the client has
	// its own timeout; this is a safety net for the retry context).
	clientTimeoutForRetry = 5 * time.Second
)

// eventSender ships one event to the audit store (implemented by *Client).
type eventSender interface {
	Send(ctx context.Context, ev Event) error
}

// Sender drains audit events to the auditlog service. In the default async
// mode events are buffered and retried without blocking the request; with
// sync=true each event is sent synchronously (strict durability, slower).
// Critical events (login/logout, M2) are always delivered synchronously so a
// buffer overflow can never lose the authentication trail.
type Sender struct {
	logger  *slog.Logger
	client  eventSender
	sync    bool
	queue   chan Event
	done    chan struct{}
	stopped atomic.Bool
	dropped atomic.Int64
	wg      sync.WaitGroup
}

// NewSender builds the sender and, in async mode, starts its worker.
func NewSender(logger *slog.Logger, client *Client, cfg config.AuditConfig) *Sender {
	return newSender(logger, client, cfg.Sync)
}

// newSender builds the sender around the given event sender (test seam).
func newSender(logger *slog.Logger, client eventSender, syncMode bool) *Sender {
	s := &Sender{
		logger: logger,
		client: client,
		sync:   syncMode,
	}
	if !s.sync {
		s.queue = make(chan Event, senderQueueCapacity)
		s.done = make(chan struct{})
		s.wg.Add(1)
		go s.run()
	}
	return s
}

// Enqueue schedules an event for delivery (async) or sends it immediately
// (sync mode). Critical events are always sent synchronously; in async mode a
// full queue drops a non-critical event with an error log and a running
// dropped counter (M2).
func (s *Sender) Enqueue(ev Event) {
	if s.sync || ev.Critical {
		if err := s.sendWithRetry(ev); err != nil {
			s.logger.Error("audit send failed",
				"error", err, "entity", ev.Entity, "action", ev.Action, "path", ev.Path,
				"critical", ev.Critical)
		}
		return
	}
	if s.stopped.Load() {
		return
	}
	select {
	case s.queue <- ev:
	default:
		dropped := s.dropped.Add(1)
		s.logger.Error("audit queue full, dropping event",
			"entity", ev.Entity, "action", ev.Action, "path", ev.Path,
			"dropped_total", dropped)
	}
}

// Dropped returns how many buffered events were dropped since startup.
func (s *Sender) Dropped() int64 { return s.dropped.Load() }

// Stop signals the worker and waits (bounded by ctx) for the buffer to drain.
func (s *Sender) Stop(ctx context.Context) {
	if s.sync {
		return
	}
	s.stopped.Store(true)
	close(s.done)
	done := make(chan struct{})
	go func() { s.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-ctx.Done():
		s.logger.WarnContext(ctx, "audit sender stop timed out, buffered events may be lost")
	}
}

// run drains the queue until Stop.
func (s *Sender) run() {
	defer s.wg.Done()
	for {
		select {
		case ev := <-s.queue:
			s.deliver(ev)
		case <-s.done:
			// Drain the remaining buffered events (best effort, bounded by the
			// client timeout) before exiting.
			for {
				select {
				case ev := <-s.queue:
					s.deliver(ev)
				default:
					return
				}
			}
		}
	}
}

// deliver sends one buffered event and logs failures.
func (s *Sender) deliver(ev Event) {
	if err := s.sendWithRetry(ev); err != nil {
		s.logger.Error("audit send failed",
			"error", err, "entity", ev.Entity, "action", ev.Action, "path", ev.Path)
	}
}

// sendWithRetry sends one event with a bounded retry loop.
func (s *Sender) sendWithRetry(ev Event) error {
	var lastErr error
	backoff := senderBaseBackoff
	for attempt := 1; attempt <= senderMaxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), clientTimeoutForRetry)
		err := s.client.Send(ctx, ev)
		cancel()
		if err == nil {
			return nil
		}
		lastErr = err
		if attempt < senderMaxAttempts {
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return lastErr
}
