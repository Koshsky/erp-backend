package policies

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// RouteSpec — configurable definition of a route policy (kind + parameters).
// The mechanisms (builders) stay in code; the kind choice and parameters come from the DB.
type RouteSpec struct {
	Name   string
	Kind   string
	Params map[string]any
}

// KindBuilder builds a route check from parameters (validates them on write
// and when loading from the DB).
type KindBuilder func(params map[string]any) (func(*rbac.CheckCtx) error, error)

// ParamInfo describes a kind parameter (for the API and UI reference).
type ParamInfo struct {
	Key      string `json:"key"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// KindInfo describes a kind for the API reference.
type KindInfo struct {
	Name   string      `json:"name"`
	Params []ParamInfo `json:"params"`
}

// Route policy kind codes.
const (
	kindList         = "list"
	kindEntity       = "entity"
	kindCreate       = "create"
	kindOwnerMatch   = "owner_match"
	kindAuthorOr     = "author_or"
	kindParentAction = "parent_action"
)

// Kind parameter keys (JSONB).
const (
	paramResource        = "resource"
	paramQueryKey        = "query_key"
	paramAction          = "action"
	paramOwner           = "owner"
	paramParentResource  = "parent_resource"
	paramParentFrom      = "parent_from"
	paramOwnerKey        = "owner_key"
	paramDefaultSelf     = "default_self"
	paramPrimaryResource = "primary_resource"
	paramPrimaryFrom     = "primary_from"
	paramCompareResource = "compare_resource"
	paramCompareFrom     = "compare_from"
	paramExemptRoles     = "exempt_roles"
	paramAuthorResource  = "author_resource"
	paramAuthorIDParam   = "author_id_param"
	paramRightResource   = "right_resource"
	paramRightAction     = "right_action"
)

// Parameter value types (API reference).
const (
	paramTypeString     = "string"
	paramTypeBool       = "bool"
	paramTypeStringList = "string_list"
)

// Owner resolution modes for the entity kind.
const (
	ownerModeID   = "id"   // owner by id from the URL (default)
	ownerModeNone = "none" // resource without an owner chain (state, virtual resources)
)

// body/query keys appearing in the default specifications.
const (
	bodyKeyProjectID  = "project_id"
	bodyKeyProcessID  = "process_id"
	bodyKeyOwnerID    = "owner_id"
	bodyKeyManagerID  = "manager_id"
	bodyKeyTaskID     = "task_id"
	bodyKeyResourceID = "resource_id"
	bodyKeyCommentID  = "comment_id"
)

//nolint:gochecknoglobals // kind registry (mechanisms; static)
var kindRegistry = map[string]KindBuilder{
	kindList:         buildList,
	kindEntity:       buildEntity,
	kindCreate:       buildCreate,
	kindOwnerMatch:   buildOwnerMatch,
	kindAuthorOr:     buildAuthorOr,
	kindParentAction: buildParentAction,
}

//nolint:gochecknoglobals // kind parameter schemas (mirror the builders)
var kindParamSchemas = map[string][]ParamInfo{
	kindList: {
		{Key: paramResource, Type: paramTypeString, Required: true},
		{Key: paramQueryKey, Type: paramTypeString, Required: true},
	},
	kindEntity: {
		{Key: paramResource, Type: paramTypeString, Required: true},
		{Key: paramAction, Type: paramTypeString, Required: true},
		{Key: paramOwner, Type: paramTypeString, Required: false}, // "id" (default) | "none"
	},
	kindCreate: {
		{Key: paramResource, Type: paramTypeString, Required: true},
		{Key: paramParentResource, Type: paramTypeString, Required: false},
		{Key: paramParentFrom, Type: paramTypeString, Required: false},
		{Key: paramOwnerKey, Type: paramTypeString, Required: false},
		{Key: paramDefaultSelf, Type: paramTypeBool, Required: false},
	},
	kindOwnerMatch: {
		{Key: paramResource, Type: paramTypeString, Required: true},
		{Key: paramAction, Type: paramTypeString, Required: false},
		{Key: paramPrimaryResource, Type: paramTypeString, Required: true},
		{Key: paramPrimaryFrom, Type: paramTypeString, Required: true},
		{Key: paramCompareResource, Type: paramTypeString, Required: true},
		{Key: paramCompareFrom, Type: paramTypeString, Required: true},
		{Key: paramExemptRoles, Type: paramTypeStringList, Required: false},
	},
	kindAuthorOr: {
		{Key: paramAuthorResource, Type: paramTypeString, Required: true},
		{Key: paramAuthorIDParam, Type: paramTypeString, Required: true},
		{Key: paramRightResource, Type: paramTypeString, Required: true},
		{Key: paramRightAction, Type: paramTypeString, Required: true},
	},
	kindParentAction: {
		{Key: paramResource, Type: paramTypeString, Required: true},
		{Key: paramAction, Type: paramTypeString, Required: true},
		{Key: paramParentResource, Type: paramTypeString, Required: true},
		{Key: paramParentFrom, Type: paramTypeString, Required: true},
	},
}

// =============================================
// Parameter parsing and validation (JSONB → typed values)
// =============================================

func strParam(params map[string]any, key string, required bool) (string, error) {
	v, ok := params[key]
	if !ok || v == nil {
		if required {
			return "", fmt.Errorf("%s: обязательный параметр отсутствует", key)
		}
		return "", nil
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("%s: ожидается непустая строка", key)
	}
	return s, nil
}

func boolParam(params map[string]any, key string, def bool) (bool, error) {
	v, ok := params[key]
	if !ok || v == nil {
		return def, nil
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%s: ожидается boolean", key)
	}
	return b, nil
}

func strListParam(params map[string]any, key string) ([]string, error) {
	v, ok := params[key]
	if !ok || v == nil {
		return nil, nil
	}
	switch items := v.(type) {
	case []string:
		return items, nil
	case []any:
		out := make([]string, 0, len(items))
		for _, item := range items {
			str, okItem := item.(string)
			if !okItem || str == "" {
				return nil, fmt.Errorf("%s: ожидается список строк", key)
			}
			out = append(out, str)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("%s: ожидается список строк", key)
	}
}

// kindError returns a validation error carrying the kind name.
func kindError(kind, key string, err error) error {
	return fmt.Errorf("params[%s].%s: %w", kind, key, err)
}

// =============================================
// Builders (mechanisms — code)
// =============================================

// viewOrForbidden returns 404 for view (do not disclose whether the object
// exists) and 403 for the other actions.
func viewOrForbidden(act Action) error {
	if act == ActionView {
		return errors.ErrNotFound
	}
	return errors.ErrForbidden
}

// buildList — the standard listing rule: matrix view check + validation of
// the owner query parameter (ScopeAll — any id, otherwise one's own).
func buildList(params map[string]any) (func(*rbac.CheckCtx) error, error) {
	resName, err := strParam(params, paramResource, true)
	if err != nil {
		return nil, kindError(kindList, paramResource, err)
	}
	res, ok := ParseResource(resName)
	if !ok {
		return nil, kindError(kindList, paramResource, fmt.Errorf("неизвестный ресурс %q", resName))
	}
	key, err := strParam(params, paramQueryKey, true)
	if err != nil {
		return nil, kindError(kindList, paramQueryKey, err)
	}
	return ListCheck(res, key), nil
}

// buildEntity — the standard entity-by-id-from-URL rule; owner:none is for
// resources without an owner chain (state, virtual resources).
func buildEntity(params map[string]any) (func(*rbac.CheckCtx) error, error) {
	resName, err := strParam(params, paramResource, true)
	if err != nil {
		return nil, kindError(kindEntity, paramResource, err)
	}
	res, ok := ParseResource(resName)
	if !ok {
		return nil, kindError(kindEntity, paramResource, fmt.Errorf("неизвестный ресурс %q", resName))
	}
	actName, err := strParam(params, paramAction, true)
	if err != nil {
		return nil, kindError(kindEntity, paramAction, err)
	}
	act, ok := ParseAction(actName)
	if !ok {
		return nil, kindError(kindEntity, paramAction, fmt.Errorf("неизвестное действие %q", actName))
	}
	owner, err := strParam(params, paramOwner, false)
	if err != nil {
		return nil, kindError(kindEntity, paramOwner, err)
	}
	if owner == ownerModeNone {
		return noOwnerCheck(res, act), nil
	}
	return EntityCheck(res, act), nil
}

// noOwnerCheck — matrix check for a resource without owners (state, virtual ones).
func noOwnerCheck(res rbac.Resource, act Action) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		if !Authorize(rc.User.Role, res, act, rbac.Owners{}, rc.User.ID) {
			return viewOrForbidden(act)
		}
		return nil
	}
}

// buildCreate — the standard creation rule: by the parent owner from the body
// (parent_resource + parent_from) or into one's own ownership (owner_key +
// default_self: 0 → the current user).
func buildCreate(params map[string]any) (func(*rbac.CheckCtx) error, error) {
	resName, err := strParam(params, paramResource, true)
	if err != nil {
		return nil, kindError(kindCreate, paramResource, err)
	}
	res, ok := ParseResource(resName)
	if !ok {
		return nil, kindError(kindCreate, paramResource, fmt.Errorf("неизвестный ресурс %q", resName))
	}
	parentResource, err := strParam(params, paramParentResource, false)
	if err != nil {
		return nil, kindError(kindCreate, paramParentResource, err)
	}
	parentFrom, err := strParam(params, paramParentFrom, false)
	if err != nil {
		return nil, kindError(kindCreate, paramParentFrom, err)
	}
	ownerKey, err := strParam(params, paramOwnerKey, false)
	if err != nil {
		return nil, kindError(kindCreate, paramOwnerKey, err)
	}
	defaultSelf, err := boolParam(params, paramDefaultSelf, true)
	if err != nil {
		return nil, kindError(kindCreate, paramDefaultSelf, err)
	}

	if parentFrom != "" {
		if parentResource == "" {
			return nil, kindError(
				kindCreate,
				paramParentResource,
				fmt.Errorf("parent_from задан, а parent_resource — нет"),
			)
		}
		parentRes, okParent := ParseResource(parentResource)
		if !okParent {
			return nil, kindError(kindCreate, paramParentResource, fmt.Errorf("неизвестный ресурс %q", parentResource))
		}
		return CreateCheck(res, parentByID(parentRes, parentFrom)), nil
	}

	if ownerKey == "" {
		return nil, kindError(kindCreate, paramOwnerKey, fmt.Errorf("требуется parent_from или owner_key"))
	}
	makeOwners, okOwner := createOwners(res)
	if !okOwner {
		return nil, kindError(kindCreate, paramOwnerKey,
			fmt.Errorf("создание с owner_key неприменимо к ресурсу %q", resName))
	}
	return func(rc *rbac.CheckCtx) error {
		ownerID, bodyErr := rc.BodyID(ownerKey)
		if bodyErr != nil {
			return bodyErr
		}
		if ownerID == 0 && defaultSelf {
			ownerID = rc.User.ID
		}
		if !Authorize(rc.User.Role, res, ActionCreate, makeOwners(ownerID), rc.User.ID) {
			return errors.ErrForbidden
		}
		return nil
	}, nil
}

// createOwners returns an owners constructor for creating "into one's own
// ownership" (the field depends on the resource).
func createOwners(res rbac.Resource) (func(int64) rbac.Owners, bool) {
	switch res {
	case rbac.ResourceProject:
		return func(id int64) rbac.Owners { return rbac.Owners{ProjectOwner: id} }, true
	case rbac.ResourceWorker, rbac.ResourceResource:
		return func(id int64) rbac.Owners { return rbac.Owners{Owner: id} }, true
	case rbac.ResourceProcess, rbac.ResourceTask, rbac.ResourceMilestone,
		rbac.ResourceAssignment, rbac.ResourceState, rbac.ResourceComment,
		rbac.ResourceUserCatalog, rbac.ResourceRBACConfig,
		rbac.ResourceUserAdmin, rbac.ResourceStateAdmin, rbac.ResourceOrgStructure,
		rbac.ResourceAudit:
		return nil, false
	}
	return nil, false
}

// ownerMatchSpec — parsed parameters of the owner_match kind.
type ownerMatchSpec struct {
	res         rbac.Resource
	act         Action
	primaryRes  rbac.Resource
	primaryFrom string
	compareRes  rbac.Resource
	compareFrom string
	exemptRoles []string
}

// parseOwnerMatch parses and validates the owner_match kind parameters.
func parseOwnerMatch(params map[string]any) (ownerMatchSpec, error) {
	var spec ownerMatchSpec
	resName, err := strParam(params, paramResource, true)
	if err != nil {
		return spec, kindError(kindOwnerMatch, paramResource, err)
	}
	spec.res, err = parseResourceCode(resName, kindOwnerMatch, paramResource)
	if err != nil {
		return spec, err
	}
	actName, err := strParam(params, paramAction, false)
	if err != nil {
		return spec, kindError(kindOwnerMatch, paramAction, err)
	}
	spec.act = ActionCreate
	if actName != "" {
		act, ok := ParseAction(actName)
		if !ok {
			return spec, kindError(kindOwnerMatch, paramAction, fmt.Errorf("неизвестное действие %q", actName))
		}
		spec.act = act
	}
	if spec.primaryRes, err = parseResourceParam(params, kindOwnerMatch, paramPrimaryResource); err != nil {
		return spec, err
	}
	if spec.primaryFrom, err = strParam(params, paramPrimaryFrom, true); err != nil {
		return spec, kindError(kindOwnerMatch, paramPrimaryFrom, err)
	}
	if spec.compareRes, err = parseResourceParam(params, kindOwnerMatch, paramCompareResource); err != nil {
		return spec, err
	}
	if spec.compareFrom, err = strParam(params, paramCompareFrom, true); err != nil {
		return spec, kindError(kindOwnerMatch, paramCompareFrom, err)
	}
	if spec.exemptRoles, err = strListParam(params, paramExemptRoles); err != nil {
		return spec, kindError(kindOwnerMatch, paramExemptRoles, err)
	}
	return spec, nil
}

// parseResourceParam reads a resource code from kind parameters.
func parseResourceParam(params map[string]any, kind, key string) (rbac.Resource, error) {
	name, err := strParam(params, key, true)
	if err != nil {
		return 0, kindError(kind, key, err)
	}
	return parseResourceCode(name, kind, key)
}

// parseResourceCode parses a string resource code, reporting failures as kind errors.
func parseResourceCode(name, kind, key string) (rbac.Resource, error) {
	res, ok := ParseResource(name)
	if !ok {
		return 0, kindError(kind, key, fmt.Errorf("неизвестный ресурс %q", name))
	}
	return res, nil
}

// buildOwnerMatch — a cross-entity rule: matrix check against the primary
// entity owner + a shared-owner requirement with the compare entity
// (admin and the exempt_roles are exempt from this business rule).
func buildOwnerMatch(params map[string]any) (func(*rbac.CheckCtx) error, error) {
	spec, err := parseOwnerMatch(params)
	if err != nil {
		return nil, err
	}
	return func(rc *rbac.CheckCtx) error {
		primaryID, bodyErr := rc.BodyID(spec.primaryFrom)
		if bodyErr != nil {
			return bodyErr
		}
		primaryOwners, ownerErr := rc.Owners(spec.primaryRes, primaryID)
		if ownerErr != nil {
			return ownerErr
		}
		if !Authorize(rc.User.Role, spec.res, spec.act, primaryOwners, rc.User.ID) {
			return errors.Forbidden(
				"недостаточно прав: действие доступно владельцу родительского элемента (или администратору)",
			)
		}
		if slices.Contains(spec.exemptRoles, rc.User.Role) {
			return nil
		}
		compareID, bodyErr := rc.BodyID(spec.compareFrom)
		if bodyErr != nil {
			return bodyErr
		}
		compareOwners, ownerErr := rc.Owners(spec.compareRes, compareID)
		if ownerErr != nil {
			return ownerErr
		}
		if !primaryOwners.SharesOwner(compareOwners) {
			return errors.Forbidden("ресурс не принадлежит владельцу")
		}
		return nil
	}, nil
}

// buildAuthorOr — a disjunction: the entity author (via the author_resource
// owner chain) OR the right_action on right_resource (the parent entity).
func buildAuthorOr(params map[string]any) (func(*rbac.CheckCtx) error, error) {
	authorResource, err := strParam(params, paramAuthorResource, true)
	if err != nil {
		return nil, kindError(kindAuthorOr, paramAuthorResource, err)
	}
	authorRes, ok := ParseResource(authorResource)
	if !ok {
		return nil, kindError(kindAuthorOr, paramAuthorResource, fmt.Errorf("неизвестный ресурс %q", authorResource))
	}
	authorIDParam, err := strParam(params, paramAuthorIDParam, true)
	if err != nil {
		return nil, kindError(kindAuthorOr, paramAuthorIDParam, err)
	}
	rightResource, err := strParam(params, paramRightResource, true)
	if err != nil {
		return nil, kindError(kindAuthorOr, paramRightResource, err)
	}
	rightRes, okRight := ParseResource(rightResource)
	if !okRight {
		return nil, kindError(kindAuthorOr, paramRightResource, fmt.Errorf("неизвестный ресурс %q", rightResource))
	}
	rightAction, err := strParam(params, paramRightAction, true)
	if err != nil {
		return nil, kindError(kindAuthorOr, paramRightAction, err)
	}
	rightAct, okAct := ParseAction(rightAction)
	if !okAct {
		return nil, kindError(kindAuthorOr, paramRightAction, fmt.Errorf("неизвестное действие %q", rightAction))
	}

	return func(rc *rbac.CheckCtx) error {
		authorID, parseErr := strconv.ParseInt(rc.C.Param(authorIDParam), 10, 64)
		if parseErr != nil {
			return errors.BadRequest("invalid id")
		}
		owners, ownerErr := rc.Owners(authorRes, authorID)
		if ownerErr != nil {
			return ownerErr
		}
		if owners.Owner != 0 && owners.Owner == rc.User.ID {
			return nil // the author deletes their own
		}
		if !Authorize(rc.User.Role, rightRes, rightAct, owners, rc.User.ID) {
			return errors.ErrForbidden
		}
		return nil
	}, nil
}

// buildParentAction — an action on a child entity authorized by the parent
// entity from the body: used by bulk endpoints without a :id path parameter
// (/task/order, /process/order), whose request carries the parent id
// (process_id / project_id). Matches the matrix for the child entity against
// the parent's owner chain.
func buildParentAction(params map[string]any) (func(*rbac.CheckCtx) error, error) {
	resName, err := strParam(params, paramResource, true)
	if err != nil {
		return nil, kindError(kindParentAction, paramResource, err)
	}
	res, ok := ParseResource(resName)
	if !ok {
		return nil, kindError(kindParentAction, paramResource, fmt.Errorf("неизвестный ресурс %q", resName))
	}
	actName, err := strParam(params, paramAction, true)
	if err != nil {
		return nil, kindError(kindParentAction, paramAction, err)
	}
	act, ok := ParseAction(actName)
	if !ok {
		return nil, kindError(kindParentAction, paramAction, fmt.Errorf("неизвестное действие %q", actName))
	}
	parentResName, err := strParam(params, paramParentResource, true)
	if err != nil {
		return nil, kindError(kindParentAction, paramParentResource, err)
	}
	parentRes, ok := ParseResource(parentResName)
	if !ok {
		return nil, kindError(kindParentAction, paramParentResource, fmt.Errorf("неизвестный ресурс %q", parentResName))
	}
	parentFrom, err := strParam(params, paramParentFrom, true)
	if err != nil {
		return nil, kindError(kindParentAction, paramParentFrom, err)
	}

	return func(rc *rbac.CheckCtx) error {
		parentID, bodyErr := rc.BodyID(parentFrom)
		if bodyErr != nil {
			return bodyErr
		}
		owners, ownerErr := rc.Owners(parentRes, parentID)
		if ownerErr != nil {
			return ownerErr
		}
		if !Authorize(rc.User.Role, res, act, owners, rc.User.ID) {
			return errors.ErrForbidden
		}
		return nil
	}, nil
}

// CommentDeleteCheck — comment deletion: the author always, others by task
// update right (a wrapper of the author_or kind; kept for tests).
func CommentDeleteCheck() func(*rbac.CheckCtx) error {
	check, err := buildAuthorOr(map[string]any{
		paramAuthorResource: resComment,
		paramAuthorIDParam:  bodyKeyCommentID,
		paramRightResource:  resTask,
		paramRightAction:    actUpdate,
	})
	if err != nil {
		panic("policies: invalid built-in author_or params: " + err.Error())
	}
	return check
}

// =============================================
// Shared mechanisms (used by builders and tests)
// =============================================

// ListCheck — the standard listing rule (see buildList).
func ListCheck(rsrc rbac.Resource, key string) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		scope := scopeFor(rc.User.Role, rsrc, ActionView)
		if scope == ScopeNone {
			return errors.ErrForbidden
		}
		raw := rc.C.Query(key)
		if raw == "" {
			return nil
		}
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 0 {
			return errors.BadRequest("invalid " + key)
		}
		if scope != ScopeAll && id != rc.User.ID {
			return errors.ErrForbidden
		}
		return nil
	}
}

// parentByID resolves the parent owner from the body by key (for CreateCheck).
func parentByID(rsrc rbac.Resource, key string) func(*rbac.CheckCtx) (rbac.Owners, error) {
	return func(rc *rbac.CheckCtx) (rbac.Owners, error) {
		id, err := rc.BodyID(key)
		if err != nil {
			return rbac.Owners{}, err
		}
		return rc.Owners(rsrc, id)
	}
}

// EntityCheck — the standard entity-by-id-from-URL rule: matrix against
// owners; a denied view does not disclose existence (404).
func EntityCheck(rsrc rbac.Resource, act Action) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		id, err := rc.ParamID()
		if err != nil {
			return err
		}
		owners, err := rc.Owners(rsrc, id)
		if err != nil {
			return err
		}
		if !Authorize(rc.User.Role, rsrc, act, owners, rc.User.ID) {
			return viewOrForbidden(act)
		}
		return nil
	}
}

// CreateCheck — the standard creation rule by the parent owner from the body.
func CreateCheck(rsrc rbac.Resource, parent func(*rbac.CheckCtx) (rbac.Owners, error)) func(*rbac.CheckCtx) error {
	return func(rc *rbac.CheckCtx) error {
		owners, err := parent(rc)
		if err != nil {
			return err
		}
		if !Authorize(rc.User.Role, rsrc, ActionCreate, owners, rc.User.ID) {
			return errors.ErrForbidden
		}
		return nil
	}
}

// =============================================
// Default route specifications (mirror seed V15) and assembly
// =============================================

//nolint:gochecknoglobals // default route specs (mechanism; static)
var defaultRouteSpecs = []RouteSpec{
	{
		Name:   "project.list",
		Kind:   kindList,
		Params: map[string]any{paramResource: resProject, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "project.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resProject, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name:   "project.create",
		Kind:   kindCreate,
		Params: map[string]any{paramResource: resProject, paramOwnerKey: bodyKeyOwnerID, paramDefaultSelf: true},
	},
	{
		Name:   "project.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resProject, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "project.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resProject, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		Name:   "process.list",
		Kind:   kindList,
		Params: map[string]any{paramResource: resProcess, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "process.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resProcess, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name: "process.create",
		Kind: kindCreate,
		Params: map[string]any{
			paramResource:       resProcess,
			paramParentResource: resProject,
			paramParentFrom:     bodyKeyProjectID,
		},
	},
	{
		Name:   "process.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resProcess, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "process.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resProcess, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		// Bulk reorder: the request carries the parent project id in the body
		// (no :id path param) — authorized via the project owner chain.
		Name: "process.order",
		Kind: kindParentAction,
		Params: map[string]any{
			paramResource:       resProcess,
			paramAction:         actUpdate,
			paramParentResource: resProject,
			paramParentFrom:     bodyKeyProjectID,
		},
	},
	{Name: "task.list", Kind: kindList, Params: map[string]any{paramResource: resTask, paramQueryKey: bodyKeyOwnerID}},
	{
		Name:   "task.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resTask, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name: "task.create",
		Kind: kindCreate,
		Params: map[string]any{
			paramResource:       resTask,
			paramParentResource: resProcess,
			paramParentFrom:     bodyKeyProcessID,
		},
	},
	{
		Name:   "task.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resTask, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "task.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resTask, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		// Bulk reorder: the request carries the parent process id in the body
		// (no :id path param) — authorized via the process owner chain.
		Name: "task.order",
		Kind: kindParentAction,
		Params: map[string]any{
			paramResource:       resTask,
			paramAction:         actUpdate,
			paramParentResource: resProcess,
			paramParentFrom:     bodyKeyProcessID,
		},
	},
	{
		Name:   "milestone.list",
		Kind:   kindList,
		Params: map[string]any{paramResource: resMilestone, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "milestone.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resMilestone, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name: "milestone.create",
		Kind: kindCreate,
		Params: map[string]any{
			paramResource:       resMilestone,
			paramParentResource: resProcess,
			paramParentFrom:     bodyKeyProcessID,
		},
	},
	{
		Name:   "milestone.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resMilestone, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "milestone.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resMilestone, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		Name:   "assignment.list",
		Kind:   kindList,
		Params: map[string]any{paramResource: resAssignment, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "assignment.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resAssignment, paramAction: actView, paramOwner: ownerModeID},
	},
	{Name: "assignment.create", Kind: kindOwnerMatch, Params: map[string]any{
		paramResource: resAssignment, paramAction: actCreate,
		paramPrimaryResource: resTask, paramPrimaryFrom: bodyKeyTaskID,
		paramCompareResource: resResource, paramCompareFrom: bodyKeyResourceID,
		paramExemptRoles: []string{userdomain.Admin},
	}},
	{
		Name:   "assignment.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resAssignment, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "assignment.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resAssignment, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		Name:   "task.comment.list",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resTask, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name:   "task.comment.create",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resTask, paramAction: actView, paramOwner: ownerModeID},
	},
	{Name: "task.comment.delete", Kind: kindAuthorOr, Params: map[string]any{
		paramAuthorResource: resComment, paramAuthorIDParam: bodyKeyCommentID,
		paramRightResource: resTask, paramRightAction: actUpdate,
	}},
	{
		Name:   "worker.list",
		Kind:   kindList,
		Params: map[string]any{paramResource: resWorker, paramQueryKey: bodyKeyManagerID},
	},
	{
		Name:   "worker.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resWorker, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name:   "worker.create",
		Kind:   kindCreate,
		Params: map[string]any{paramResource: resWorker, paramOwnerKey: bodyKeyManagerID, paramDefaultSelf: true},
	},
	{
		Name:   "worker.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resWorker, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "worker.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resWorker, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		Name:   "user.picker",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resUserCatalog, paramAction: actView, paramOwner: ownerModeNone},
	},
	{
		// User profile mutations (an employee IS a system user): gated by the
		// user_admin virtual resource — admin by default, grantable via the matrix.
		Name:   "user_admin.create",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resUserAdmin, paramAction: actCreate, paramOwner: ownerModeNone},
	},
	{
		Name:   "user_admin.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resUserAdmin, paramAction: actUpdate, paramOwner: ownerModeNone},
	},
	{
		Name:   "user_admin.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resUserAdmin, paramAction: actDelete, paramOwner: ownerModeNone},
	},
	{
		Name:   "resource.list",
		Kind:   kindList,
		Params: map[string]any{paramResource: resResource, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "resource.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resResource, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name:   "resource.create",
		Kind:   kindCreate,
		Params: map[string]any{paramResource: resResource, paramOwnerKey: bodyKeyOwnerID, paramDefaultSelf: true},
	},
	{
		Name:   "resource.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resResource, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "resource.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resResource, paramAction: actDelete, paramOwner: ownerModeID},
	},
	{
		Name:   "resource.member-list",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resResource, paramAction: actView, paramOwner: ownerModeID},
	},
	{
		Name:   "resource.member-add",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resResource, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "resource.member-remove",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resResource, paramAction: actUpdate, paramOwner: ownerModeID},
	},
	{
		Name:   "calendar.view",
		Kind:   kindList,
		Params: map[string]any{paramResource: resResource, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "state.list",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resState, paramAction: actView, paramOwner: ownerModeNone},
	},
	{
		Name:   "state.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resState, paramAction: actView, paramOwner: ownerModeNone},
	},
	{
		Name:   "state.create",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resState, paramAction: actCreate, paramOwner: ownerModeNone},
	},
	{
		Name:   "state.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resState, paramAction: actUpdate, paramOwner: ownerModeNone},
	},
	{
		Name:   "state.delete",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resState, paramAction: actDelete, paramOwner: ownerModeNone},
	},
	// Planning aggregates (/planning/*): gated by the view matrix of the
	// underlying domain. The list kind checks the view right; row scoping
	// stays in the SQL via ViewScopeCode (no owner query param is exposed).
	{
		Name:   "planning.projects",
		Kind:   kindList,
		Params: map[string]any{paramResource: resProject, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "planning.processes",
		Kind:   kindList,
		Params: map[string]any{paramResource: resProcess, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "planning.tasks",
		Kind:   kindList,
		Params: map[string]any{paramResource: resTask, paramQueryKey: bodyKeyOwnerID},
	},
	{
		Name:   "autocreate.list",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resRBACConfig, paramAction: actView, paramOwner: ownerModeNone},
	},
	{
		Name:   "autocreate.update",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resRBACConfig, paramAction: actUpdate, paramOwner: ownerModeNone},
	},
	{
		Name:   "rbac.manage",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resRBACConfig, paramAction: actView, paramOwner: ownerModeNone},
	},
	// Audit-log admin section: virtual resource, admin only (bypass). No rows
	// in the matrix — the page and the /audit/* routes are admin-only.
	{
		Name:   "audit.view",
		Kind:   kindEntity,
		Params: map[string]any{paramResource: resAudit, paramAction: actView, paramOwner: ownerModeNone},
	},
}

// DefaultRouteSpecs returns the default route specifications.
func DefaultRouteSpecs() []RouteSpec {
	out := make([]RouteSpec, 0, len(defaultRouteSpecs))
	return append(out, defaultRouteSpecs...)
}

// Kinds returns the catalog of kinds and their parameters.
func Kinds() []KindInfo {
	infos := make([]KindInfo, 0, len(kindRegistry))
	for _, name := range []string{kindList, kindEntity, kindCreate, kindOwnerMatch, kindAuthorOr, kindParentAction} {
		infos = append(infos, KindInfo{Name: name, Params: kindParamSchemas[name]})
	}
	return infos
}

// BuildPolicies assembles checks from specifications (validates kind + parameters).
// An error in one specification cancels the whole build — the caller decides (fail-closed).
func BuildPolicies(specs []RouteSpec) ([]rbac.Policy, error) {
	out := make([]rbac.Policy, 0, len(specs))
	for _, spec := range specs {
		builder, ok := kindRegistry[spec.Kind]
		if !ok {
			return nil, fmt.Errorf("route %q: неизвестный kind %q", spec.Name, spec.Kind)
		}
		check, err := builder(spec.Params)
		if err != nil {
			return nil, fmt.Errorf("route %q: %w", spec.Name, err)
		}
		out = append(out, rbac.Policy{Name: spec.Name, Check: check})
	}
	return out, nil
}

// ValidateSpec validates a specification (kind + parameters) without building the check.
func ValidateSpec(spec RouteSpec) error {
	builder, ok := kindRegistry[spec.Kind]
	if !ok {
		return fmt.Errorf("неизвестный kind %q", spec.Kind)
	}
	if spec.Name == "" {
		return fmt.Errorf("пустое имя маршрутной проверки")
	}
	if strings.TrimSpace(spec.Name) != spec.Name {
		return fmt.Errorf("имя маршрутной проверки не должно содержать пробелы по краям")
	}
	_, err := builder(spec.Params)
	return err
}
