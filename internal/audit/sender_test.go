//nolint:testpackage // tests the unexported Sender internals (same pattern as middleware_test.go)
package audit

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"
)

// stubService records delivered events and can fail on demand.
type stubService struct {
	mu     sync.Mutex
	events []Event
	err    error
}

func (s *stubService) Send(_ context.Context, ev Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return s.err
	}
	s.events = append(s.events, ev)
	return nil
}

func (s *stubService) sent() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Event(nil), s.events...)
}

func newTestSender(client eventSender, syncMode bool) *Sender {
	return newSender(slog.New(slog.DiscardHandler), client, syncMode)
}

// stopSender stops the sender with a bounded context: with a failing client
// the drain would otherwise retry every buffered event (tens of seconds).
func stopSender(s *Sender) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	s.Stop(ctx)
}

// TestCriticalEventBypassesFullQueue checks that a critical event (login) is
// delivered synchronously with retries even when the client is failing, while
// a non-critical event under the same conditions is dropped by the full
// buffer — i.e. a full queue never loses the authentication trail (M2).
func TestCriticalEventBypassesFullQueue(t *testing.T) {
	t.Parallel()
	client := &stubService{err: errors.New("loki down")}
	s := newTestSender(client, false)
	defer stopSender(s)

	// Fill the buffer: every buffered event fails after its retry budget, so
	// the worker cannot drain the queue while we enqueue.
	for range senderQueueCapacity + 50 {
		s.Enqueue(Event{Entity: entityUser, Action: "update", Path: "/x"})
	}
	if dropped := s.Dropped(); dropped == 0 {
		t.Fatal("Dropped() = 0, want non-zero after overflowing the buffer")
	}

	// A critical event must not be counted as dropped: it is sent inline
	// (synchronously) rather than enqueued.
	droppedBefore := s.Dropped()
	start := time.Now()
	s.Enqueue(Event{Entity: entityAuth, Action: actionLogin, Path: "/auth/login", Critical: true})
	elapsed := time.Since(start)

	if s.Dropped() != droppedBefore {
		t.Errorf("Dropped() = %d, want %d: the critical event was dropped instead of sent inline",
			s.Dropped(), droppedBefore)
	}
	if elapsed < senderBaseBackoff {
		t.Errorf("elapsed = %v, want the inline retry budget (%v+)", elapsed, senderBaseBackoff)
	}
}

// TestNonCriticalDropIsCounted checks that a non-critical event dropped by a
// full buffer is counted (M2): with a failing client the worker cannot drain
// the queue, so further buffered events overflow.
func TestNonCriticalDropIsCounted(t *testing.T) {
	t.Parallel()
	client := &stubService{err: errors.New("loki down")}
	s := newTestSender(client, false)
	defer stopSender(s)

	for range senderQueueCapacity + 50 {
		s.Enqueue(Event{Entity: entityUser, Action: "update", Path: "/x"})
	}

	if dropped := s.Dropped(); dropped == 0 {
		t.Fatal("Dropped() = 0, want non-zero after overflowing the buffer")
	}
}

// TestSyncModeSendsInline checks sync mode delivers every event inline.
func TestSyncModeSendsInline(t *testing.T) {
	t.Parallel()
	client := &stubService{}
	s := newTestSender(client, true)
	s.Enqueue(Event{Entity: entityAuth, Action: actionLogout, Path: "/auth/logout"})

	sent := client.sent()
	if len(sent) != 1 || sent[0].Action != actionLogout {
		t.Fatalf("sent = %+v, want one logout event", sent)
	}
}

// TestSendRetriesThenFails checks the retry budget is exhausted before giving
// up (failures are only logged, never propagated to the request).
func TestSendRetriesThenFails(t *testing.T) {
	t.Parallel()
	client := &stubService{err: errors.New("loki down")}
	s := newTestSender(client, true)

	start := time.Now()
	s.Enqueue(Event{Entity: entityAuth, Action: actionLogin, Critical: true})
	if elapsed := time.Since(start); elapsed < senderBaseBackoff {
		t.Errorf("elapsed = %v, want at least one backoff (%v)", elapsed, senderBaseBackoff)
	}
}

// TestIsCriticalAction checks the critical set covers auth events only.
func TestIsCriticalAction(t *testing.T) {
	t.Parallel()
	if !isCriticalAction(actionLogin) || !isCriticalAction(actionLogout) {
		t.Error("login/logout must be critical")
	}
	if isCriticalAction("update") || isCriticalAction(actionRefresh) {
		t.Error("non-auth events must not be critical")
	}
}
