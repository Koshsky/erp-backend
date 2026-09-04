package audit

import (
	"encoding/json"
	"time"
)

// Event is one audit record shipped to the auditlog service. The JSON shape
// mirrors the auditlog store.Event schema (auditdb.audit_events).
type Event struct {
	TS           string          `json:"ts"`
	ActorUserID  *int64          `json:"actor_user_id,omitempty"`
	ActorEmail   string          `json:"actor_email,omitempty"`
	ActorRole    string          `json:"actor_role,omitempty"`
	ActorIP      string          `json:"actor_ip,omitempty"`
	Entity       string          `json:"entity"`
	Action       string          `json:"action"`
	EntityID     *int64          `json:"entity_id,omitempty"`
	Method       string          `json:"method"`
	Path         string          `json:"path"`
	Status       int             `json:"status"`
	DurationMS   *int            `json:"duration_ms,omitempty"`
	RequestBody  json.RawMessage `json:"request_body,omitempty"`
	ResponseBody json.RawMessage `json:"response_body,omitempty"`
}

// entity/action values used in route classification (see route.go).
const (
	entityAuth = "auth"
	entityUser = "user"

	actionLogin      = "login"
	actionRefresh    = "refresh"
	actionLogout     = "logout"
	actionChangePass = "change_password"
)

// nowRFC3339 returns the current time in the auditlog timestamp format.
func nowRFC3339() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")
}

// durationMS converts a duration to a millisecond pointer (min 1ms).
func durationMS(d time.Duration) *int {
	ms := max(int(d.Milliseconds()), 1)
	return &ms
}
