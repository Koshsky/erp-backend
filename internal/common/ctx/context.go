package ctx

type UserContext struct {
	ID       int64  `json:"id"`
	Role     string `json:"role"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	TenantID int64  `json:"tenant_id,omitempty"`
}

type RequestContext struct {
	User      UserContext
	RequestID string
	ClientIP  string
	UserAgent string
}

const (
	KeyUser    = "user"
	KeyRequest = "request"
)
