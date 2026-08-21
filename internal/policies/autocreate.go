package policies

import (
	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoglobals // rule registry
var autoCreatePolicies = []rbac.Policy{
	{Name: "autocreate.list", Check: adminOnly},
	{Name: "autocreate.update", Check: adminOnly},
}

// adminOnly lets only the admin through (конфигурация автосоздания).
func adminOnly(rc *rbac.CheckCtx) error {
	if rc.User.Role != userdomain.Admin {
		return errors.ErrForbidden
	}
	return nil
}
