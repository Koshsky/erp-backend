// Package policies — единственное место с правилами доступа. Движок (rbac)
// остаётся чистым механизмом; правила разбиты на матрицу (роль × ресурс ×
// действие → зона владения) и маршрутные проверки (kind + параметры).
package policies

import (
	"sync/atomic"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

// Action — операция над сущностью.
type Action int

const (
	ActionView Action = iota
	ActionCreate
	ActionUpdate
	ActionDelete
)

// Scope — зона владения, требуемая для действия. Механизм один: «владелец
// родительского элемента». Отсутствие правила = ScopeNone (нет доступа).
type Scope int

const (
	ScopeNone Scope = iota
	ScopeAll
	ScopeOwn      // владелец самой строки (проект для rp; resource/worker для vp)
	ScopeParent   // владелец непосредственного родителя (управление на уровень вниз)
	ScopeAncestor // владелец предка на один уровень выше родителя (просмотр задач через проект)
)

// Строковые коды ресурсов (зеркалит V15 и kind-схемы).
const (
	resProject     = "project"
	resProcess     = "process"
	resTask        = "task"
	resMilestone   = "milestone"
	resAssignment  = "assignment"
	resState       = "state"
	resResource    = "resource"
	resWorker      = "worker"
	resComment     = "comment"
	resUserCatalog = "user_catalog"
	resRBACConfig  = "rbac_config"
)

// Строковые коды действий.
const (
	actView   = "view"
	actCreate = "create"
	actUpdate = "update"
	actDelete = "delete"
)

// Строковые коды зон владения.
const (
	scopeAll      = "all"
	scopeOwn      = "own"
	scopeParent   = "parent"
	scopeAncestor = "ancestor"
)

// Rule связывает роль и требуемую зону доступа.
type Rule struct {
	Role  string
	Scope Scope
}

// MatrixRule — строка матрицы (для сборки из БД и дефолтов).
type MatrixRule struct {
	Res   rbac.Resource
	Act   Action
	Role  string
	Scope Scope
}

// CurrentMatrix возвращает активную матрицу (для API matrix/explain).
func CurrentMatrix() Matrix { return snapshot() }

// DefaultMatrixRules возвращает правила встроенной матрицы (для reset).
func DefaultMatrixRules() []MatrixRule {
	var out []MatrixRule
	for res, byAction := range defaultMatrix.rules {
		for act, rules := range byAction {
			for _, r := range rules {
				out = append(out, MatrixRule{Res: res, Act: act, Role: r.Role, Scope: r.Scope})
			}
		}
	}
	return out
}

// Matrix — снимок матрицы прав (роль × ресурс × действие → правила).
type Matrix struct {
	rules map[rbac.Resource]map[Action][]Rule
}

// NewMatrix собирает матрицу из правил. Пары (role, resource, action)
// уникальны: последнее правило побеждает.
func NewMatrix(rules []MatrixRule) Matrix {
	m := Matrix{rules: make(map[rbac.Resource]map[Action][]Rule)}
	for _, r := range rules {
		byAction, ok := m.rules[r.Res]
		if !ok {
			byAction = make(map[Action][]Rule)
			m.rules[r.Res] = byAction
		}
		replaced := false
		for i, existing := range byAction[r.Act] {
			if existing.Role == r.Role {
				byAction[r.Act][i] = Rule{Role: r.Role, Scope: r.Scope}
				replaced = true
				break
			}
		}
		if !replaced {
			byAction[r.Act] = append(byAction[r.Act], Rule{Role: r.Role, Scope: r.Scope})
		}
	}
	return m
}

//nolint:gochecknoglobals // live snapshot; обновляется PolicyStore из БД, fallback — DefaultMatrix
var currentRules atomic.Pointer[Matrix]

// SetMatrix атомарно заменяет активную матрицу прав.
func SetMatrix(m Matrix) {
	currentRules.Store(&m)
}

// snapshot возвращает активную матрицу или встроенные дефолты (до первой
// загрузки из БД и в тестах).
func snapshot() Matrix {
	if m := currentRules.Load(); m != nil {
		return *m
	}
	return DefaultMatrix()
}

// DefaultMatrix — встроенная матрица по умолчанию: сериализация seed'а
// V15__rbac_policies.sql. Используется как fallback, источник reset и золотой
// тест эквивалентности. admin и worker не перечислены явно: admin — ScopeAll
// (инвариант), worker — ScopeNone (нет строк).
//
//nolint:gochecknoglobals // rule registry
var defaultMatrix = Matrix{rules: map[rbac.Resource]map[Action][]Rule{
	rbac.ResourceProject: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeOwn},
		},
		ActionCreate: {
			// rp создаёт проект в свою собственность (owner по умолчанию — сам).
			{userdomain.ProjectManager, ScopeOwn},
		},
		ActionUpdate: {
			// dp и admin редактируют любой проект, rp — свой.
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeOwn},
		},
		ActionDelete: {
			{userdomain.ProjectManager, ScopeOwn},
		},
	},
	rbac.ResourceProcess: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			// rp — процессы своих проектов (parent), vp — справочно все.
			{userdomain.ProjectManager, ScopeParent},
			{userdomain.ProcessOwner, ScopeAll},
		},
		ActionCreate: {
			{userdomain.ProjectManager, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProjectManager, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProjectManager, ScopeParent},
		},
	},
	rbac.ResourceTask: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAncestor},
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeParent},
		},
	},
	rbac.ResourceMilestone: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAncestor},
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeParent},
		},
	},
	rbac.ResourceAssignment: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAncestor},
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeParent},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeParent},
		},
	},
	// === Timesheet ===
	// States: справочник без владельца; vp видит (для табеля), управляет только admin.
	rbac.ResourceState: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeAll},
		},
		ActionCreate: {},
		ActionUpdate: {},
		ActionDelete: {},
	},
	// Resource categories: admin — все, vp — свои (own); vp создаёт в свою собственность.
	rbac.ResourceResource: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
	},
	// Workers: admin — все, vp — свои подчинённые (manager_id); vp создаёт в свою команду.
	rbac.ResourceWorker: {
		ActionView: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionCreate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionUpdate: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
		ActionDelete: {
			{userdomain.ProcessOwner, ScopeOwn},
		},
	},
	// Comments: доступа нет в общей матрице — права считаются по родительской
	// задаче (см. route-проверки task.comment.*: list/create по task.view,
	// delete — автор или право обновления задачи).
	rbac.ResourceComment: {},
	// Виртуальные ресурсы: user_catalog — каталог пользователей для пикеров
	// (dp/rp/vp + admin); rbac_config — управление автосозданием/RBAC (только
	// admin, bypass) — строк в матрице нет.
	rbac.ResourceUserCatalog: {
		ActionView: {
			{userdomain.ProjectDirector, ScopeAll},
			{userdomain.ProjectManager, ScopeAll},
			{userdomain.ProcessOwner, ScopeAll},
		},
	},
	rbac.ResourceRBACConfig: {},
}}

// DefaultMatrix возвращает встроенную матрицу по умолчанию (копию).
func DefaultMatrix() Matrix {
	return defaultMatrix
}

// ScopeFor возвращает требуемую зону доступа для (role, resource, action).
// admin получает ScopeAll (защитный инвариант, не хранится в БД).
func (m Matrix) ScopeFor(role string, res rbac.Resource, act Action) Scope {
	if role == userdomain.Admin {
		return ScopeAll
	}
	rules, ok := m.rules[res][act]
	if !ok {
		return ScopeNone
	}
	for _, r := range rules {
		if r.Role == role {
			return r.Scope
		}
	}
	return ScopeNone
}

// ownField возвращает владельца самой строки для ресурса (0 — ресурс нельзя
// проверять зоной ScopeOwn: у сущности нет собственного владельца).
func ownField(res rbac.Resource, owners rbac.Owners) int64 {
	switch res {
	case rbac.ResourceProject:
		return owners.ProjectOwner
	case rbac.ResourceProcess:
		return owners.ProcessOwner
	case rbac.ResourceResource, rbac.ResourceWorker:
		return owners.Owner
	case rbac.ResourceTask, rbac.ResourceMilestone, rbac.ResourceAssignment,
		rbac.ResourceState, rbac.ResourceComment,
		rbac.ResourceUserCatalog, rbac.ResourceRBACConfig:
		return 0
	}
	return 0
}

// parentField возвращает владельца непосредственного родителя для ресурса
// (0 — у ресурса нет родителя в иерархии project → process → task/…).
func parentField(res rbac.Resource, owners rbac.Owners) int64 {
	switch res {
	case rbac.ResourceProcess:
		return owners.ProjectOwner
	case rbac.ResourceTask, rbac.ResourceMilestone, rbac.ResourceAssignment:
		return owners.ProcessOwner
	case rbac.ResourceProject, rbac.ResourceState, rbac.ResourceResource,
		rbac.ResourceWorker, rbac.ResourceComment,
		rbac.ResourceUserCatalog, rbac.ResourceRBACConfig:
		return 0
	}
	return 0
}

// ancestorField возвращает владельца предка на один уровень выше родителя
// (для задачи/вехи/назначения — владелец проекта; 0 — предка нет).
func ancestorField(res rbac.Resource, owners rbac.Owners) int64 {
	switch res {
	case rbac.ResourceTask, rbac.ResourceMilestone, rbac.ResourceAssignment:
		return owners.ProjectOwner
	case rbac.ResourceProject, rbac.ResourceProcess, rbac.ResourceState,
		rbac.ResourceResource, rbac.ResourceWorker, rbac.ResourceComment,
		rbac.ResourceUserCatalog, rbac.ResourceRBACConfig:
		return 0
	}
	return 0
}

// Authorize сообщает, может ли пользователь выполнить действие над сущностью
// с её владельцами.
func Authorize(role string, res rbac.Resource, act Action, owners rbac.Owners, userID int64) bool {
	switch snapshot().ScopeFor(role, res, act) {
	case ScopeNone:
		return false
	case ScopeAll:
		return true
	case ScopeOwn:
		owner := ownField(res, owners)
		return userID != 0 && owner != 0 && owner == userID
	case ScopeParent:
		parent := parentField(res, owners)
		return userID != 0 && parent != 0 && parent == userID
	case ScopeAncestor:
		ancestor := ancestorField(res, owners)
		return userID != 0 && ancestor != 0 && ancestor == userID
	default:
		return false
	}
}

// Can сообщает, может ли роль выполнить действие в принципе
// (грубая проверка перед загрузкой списков).
func Can(role string, res rbac.Resource, act Action) bool {
	return scopeFor(role, res, act) != ScopeNone
}

// ViewScopeCode возвращает строковый код зоны просмотра для листинг-запросов
// (all|own|parent|ancestor). SQL применяет ровно эту зону к owner-цепочке.
func ViewScopeCode(role string, res rbac.Resource) string {
	return ScopeName(scopeFor(role, res, ActionView))
}

// scopeFor — пакетная обёртка для билдеров проверок.
func scopeFor(role string, res rbac.Resource, act Action) Scope {
	return snapshot().ScopeFor(role, res, act)
}

//nolint:gochecknoglobals // resource codex (стабильный словарь, зеркалит V15)
var resourceNames = map[rbac.Resource]string{
	rbac.ResourceProject:     resProject,
	rbac.ResourceProcess:     resProcess,
	rbac.ResourceTask:        resTask,
	rbac.ResourceMilestone:   resMilestone,
	rbac.ResourceAssignment:  resAssignment,
	rbac.ResourceState:       resState,
	rbac.ResourceResource:    resResource,
	rbac.ResourceWorker:      resWorker,
	rbac.ResourceComment:     resComment,
	rbac.ResourceUserCatalog: resUserCatalog,
	rbac.ResourceRBACConfig:  resRBACConfig,
}

//nolint:gochecknoglobals // action codex
var actionNames = map[Action]string{
	ActionView:   actView,
	ActionCreate: actCreate,
	ActionUpdate: actUpdate,
	ActionDelete: actDelete,
}

//nolint:gochecknoglobals // scope codex («none» не хранится: отсутствие строки = нет доступа)
var scopeNames = map[Scope]string{
	ScopeNone:     "",
	ScopeAll:      scopeAll,
	ScopeOwn:      scopeOwn,
	ScopeParent:   scopeParent,
	ScopeAncestor: scopeAncestor,
}

// ResourceName возвращает строковый код ресурса ("" — неизвестен).
func ResourceName(res rbac.Resource) string { return resourceNames[res] }

// ParseResource разбирает строковый код ресурса.
func ParseResource(s string) (rbac.Resource, bool) {
	for res, name := range resourceNames {
		if name == s {
			return res, true
		}
	}
	return 0, false
}

// ActionName возвращает строковый код действия ("" — неизвестно).
func ActionName(act Action) string { return actionNames[act] }

// ParseAction разбирает строковый код действия.
func ParseAction(s string) (Action, bool) {
	for act, name := range actionNames {
		if name == s {
			return act, true
		}
	}
	return 0, false
}

// ScopeName возвращает строковый код зоны ("" — нет доступа, не хранится).
func ScopeName(scope Scope) string { return scopeNames[scope] }

// ParseScope разбирает строковый код зоны.
func ParseScope(s string) (Scope, bool) {
	for scope, name := range scopeNames {
		if name == s {
			return scope, true
		}
	}
	return 0, false
}

//nolint:gochecknoglobals // applicability maps зон (полные: все ресурсы перечислены)
var ownApplicable = map[rbac.Resource]bool{
	rbac.ResourceProject:     true,
	rbac.ResourceProcess:     true,
	rbac.ResourceTask:        false,
	rbac.ResourceMilestone:   false,
	rbac.ResourceAssignment:  false,
	rbac.ResourceState:       false,
	rbac.ResourceResource:    true,
	rbac.ResourceWorker:      true,
	rbac.ResourceComment:     false,
	rbac.ResourceUserCatalog: false,
	rbac.ResourceRBACConfig:  false,
}

//nolint:gochecknoglobals // applicability maps зон (полные: все ресурсы перечислены)
var parentApplicable = map[rbac.Resource]bool{
	rbac.ResourceProject:     false,
	rbac.ResourceProcess:     true,
	rbac.ResourceTask:        true,
	rbac.ResourceMilestone:   true,
	rbac.ResourceAssignment:  true,
	rbac.ResourceState:       false,
	rbac.ResourceResource:    false,
	rbac.ResourceWorker:      false,
	rbac.ResourceComment:     false,
	rbac.ResourceUserCatalog: false,
	rbac.ResourceRBACConfig:  false,
}

//nolint:gochecknoglobals // applicability maps зон (полные: все ресурсы перечислены)
var ancestorApplicable = map[rbac.Resource]bool{
	rbac.ResourceProject:     false,
	rbac.ResourceProcess:     false,
	rbac.ResourceTask:        true,
	rbac.ResourceMilestone:   true,
	rbac.ResourceAssignment:  true,
	rbac.ResourceState:       false,
	rbac.ResourceResource:    false,
	rbac.ResourceWorker:      false,
	rbac.ResourceComment:     false,
	rbac.ResourceUserCatalog: false,
	rbac.ResourceRBACConfig:  false,
}

// ScopeApplicable сообщает, применима ли зона к ресурсу (для валидации правил).
func ScopeApplicable(res rbac.Resource, scope Scope) bool {
	switch scope {
	case ScopeAll:
		return true
	case ScopeOwn:
		return ownApplicable[res]
	case ScopeParent:
		return parentApplicable[res]
	case ScopeAncestor:
		return ancestorApplicable[res]
	case ScopeNone:
		return false
	}
	return false
}
