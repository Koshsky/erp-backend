package policies

import "github.com/Koshsky/erp-backend/internal/middleware/rbac"

// ProvideAll возвращает политики по умолчанию — стартовый снапшот и fallback
// при недоступной/пустой БД. На рантайме реальный набор поставляется
// PolicyStore через rbac.Middleware.Refresh.
func ProvideAll() []rbac.Policy {
	policies, err := BuildPolicies(DefaultRouteSpecs())
	if err != nil {
		panic("policies: invalid default spec: " + err.Error())
	}
	return policies
}
