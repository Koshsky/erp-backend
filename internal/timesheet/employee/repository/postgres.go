//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/common/nullable"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/repository/sqlc"
)

type EmployeeRepository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	db     *sqlc.Queries
}

func NewEmployeeRepository(logger *slog.Logger, pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{
		logger: logger,
		pool:   pool,
		db:     sqlc.New(pool),
	}
}

func (r *EmployeeRepository) IsResourceActive(ctx context.Context, resourceID int64) (bool, error) {
	return r.db.IsResourceActive(ctx, resourceID)
}

func (r *EmployeeRepository) ListEmployeesByResourceID(
	ctx context.Context,
	resourceID int64,
) ([]domain.Employee, error) {
	rows, err := r.db.ListEmployeesByResourceID(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	employees := make([]domain.Employee, 0, len(rows))
	for _, row := range rows {
		employees = append(employees, mapEmployeeByResourceRow(row))
	}
	return employees, nil
}

func (r *EmployeeRepository) ListEmployees(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.db.ListEmployees(ctx)
	if err != nil {
		return nil, err
	}
	employees := make([]domain.Employee, 0, len(rows))
	for _, row := range rows {
		employees = append(employees, mapEmployeeListRow(row))
	}
	return employees, nil
}

func (r *EmployeeRepository) ListEmployeesByManagerID(ctx context.Context, managerID int64) ([]domain.Employee, error) {
	rows, err := r.db.ListEmployeesByManagerID(ctx, managerID)
	if err != nil {
		return nil, err
	}
	employees := make([]domain.Employee, 0, len(rows))
	for _, row := range rows {
		employees = append(employees, mapEmployeeByManagerRow(row))
	}
	return employees, nil
}

func (r *EmployeeRepository) FindEmployee(ctx context.Context, id int64) (*domain.Employee, error) {
	row, err := r.db.FindEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped := mapEmployeeRow(row)
	return &mapped, nil
}

func (r *EmployeeRepository) CreateEmployee(ctx context.Context, employee domain.Employee) (*domain.Employee, error) {
	row, err := r.db.CreateEmployee(ctx, sqlc.CreateEmployeeParams{
		ResourceID:      employee.ResourceID,
		Name:            employee.Name,
		Position:        employee.Position,
		ManagerID:       nullable.ToInt8(employee.ManagerID),
		HireDate:        toDate(employee.HireDate),
		TerminationDate: toDate(employee.TerminationDate),
	})
	if err != nil {
		return nil, err
	}

	return r.findFull(ctx, row.ID)
}

func (r *EmployeeRepository) UpdateEmployee(ctx context.Context, employee domain.Employee) (*domain.Employee, error) {
	row, err := r.db.UpdateEmployee(ctx, sqlc.UpdateEmployeeParams{
		EmployeeID:      employee.ID,
		ResourceID:      employee.ResourceID,
		Name:            employee.Name,
		Position:        employee.Position,
		ManagerID:       nullable.ToInt8(employee.ManagerID),
		HireDate:        toDate(employee.HireDate),
		TerminationDate: toDate(employee.TerminationDate),
	})
	if err != nil {
		return nil, err
	}

	return r.findFull(ctx, row.ID)
}

func (r *EmployeeRepository) DeleteEmployee(ctx context.Context, id int64) error {
	return r.db.DeleteEmployee(ctx, id)
}

func (r *EmployeeRepository) ListStates(
	ctx context.Context,
	employeeID int64,
	start, end time.Time,
) ([]domain.EmployeeState, error) {
	rows, err := r.db.ListStatesByEmployeeRange(ctx, sqlc.ListStatesByEmployeeRangeParams{
		EmployeeID: employeeID,
		StartDate:  start,
		EndDate:    end,
	})
	if err != nil {
		return nil, err
	}
	states := make([]domain.EmployeeState, 0, len(rows))
	for _, row := range rows {
		states = append(states, mapStateRow(row))
	}
	return states, nil
}

// overlapState — перекрывающий интервал с его состоянием.
type overlapState struct {
	StateID   int64
	StartDate time.Time
	EndDate   time.Time
}

// SetStateRange перезаписывает диапазон [start, end] заданным состоянием:
// вычитает [start, end] из пересекающихся интервалов (сохраняя их остатки),
// удаляет покрытые и вставляет новый.
func (r *EmployeeRepository) SetStateRange(
	ctx context.Context,
	employeeID, stateID int64,
	start, end time.Time,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	rows, err := q.ListOverlappingStates(ctx, sqlc.ListOverlappingStatesParams{
		EmployeeID: employeeID,
		StartDate:  start,
		EndDate:    end,
	})
	if err != nil {
		return err
	}
	if err = q.DeleteOverlapping(ctx, sqlc.DeleteOverlappingParams{
		EmployeeID: employeeID,
		StartDate:  start,
		EndDate:    end,
	}); err != nil {
		return err
	}
	if err = insertResidues(ctx, q, employeeID, toOverlapStates(rows), start, end); err != nil {
		return err
	}
	if _, err = q.InsertStateRange(ctx, sqlc.InsertStateRangeParams{
		EmployeeID: employeeID,
		StateID:    stateID,
		StartDate:  start,
		EndDate:    end,
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeleteStateRange очищает диапазон [start, end] (сохраняя остатки пересекающихся интервалов);
// при stateID != nil — только состояния этого типа.
func (r *EmployeeRepository) DeleteStateRange(
	ctx context.Context,
	employeeID int64,
	start, end time.Time,
	stateID *int64,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	var overlaps []overlapState
	if stateID == nil {
		overlaps, err = loadAndDeleteAll(ctx, q, employeeID, start, end)
	} else {
		overlaps, err = loadAndDeleteByState(ctx, q, employeeID, *stateID, start, end)
	}
	if err != nil {
		return err
	}
	if err = insertResidues(ctx, q, employeeID, overlaps, start, end); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// loadAndDeleteAll загружает и удаляет все интервалы, пересекающие [start, end].
func loadAndDeleteAll(
	ctx context.Context,
	q *sqlc.Queries,
	employeeID int64,
	start, end time.Time,
) ([]overlapState, error) {
	rows, err := q.ListOverlappingStates(ctx, sqlc.ListOverlappingStatesParams{
		EmployeeID: employeeID,
		StartDate:  start,
		EndDate:    end,
	})
	if err != nil {
		return nil, err
	}
	return toOverlapStates(rows), q.DeleteOverlapping(ctx, sqlc.DeleteOverlappingParams{
		EmployeeID: employeeID,
		StartDate:  start,
		EndDate:    end,
	})
}

// loadAndDeleteByState загружает и удаляет интервалы заданного состояния, пересекающие [start, end].
func loadAndDeleteByState(
	ctx context.Context,
	q *sqlc.Queries,
	employeeID, stateID int64,
	start, end time.Time,
) ([]overlapState, error) {
	rows, err := q.ListOverlappingStatesByState(ctx, sqlc.ListOverlappingStatesByStateParams{
		EmployeeID: employeeID,
		StateID:    stateID,
		StartDate:  start,
		EndDate:    end,
	})
	if err != nil {
		return nil, err
	}
	return toOverlapStatesByState(rows), q.DeleteOverlappingByState(ctx, sqlc.DeleteOverlappingByStateParams{
		EmployeeID: employeeID,
		StateID:    stateID,
		StartDate:  start,
		EndDate:    end,
	})
}

// insertResidues вставляет части пересекающихся интервалов, не входящие в [start, end]
// (левые/правые остатки), сохраняя их исходное состояние.
func insertResidues(
	ctx context.Context,
	q *sqlc.Queries,
	employeeID int64,
	overlaps []overlapState,
	start, end time.Time,
) error {
	for _, o := range overlaps {
		if o.StartDate.Before(start) {
			if _, err := q.InsertStateRange(ctx, sqlc.InsertStateRangeParams{
				EmployeeID: employeeID,
				StateID:    o.StateID,
				StartDate:  o.StartDate,
				EndDate:    start.AddDate(0, 0, -1),
			}); err != nil {
				return err
			}
		}
		if o.EndDate.After(end) {
			if _, err := q.InsertStateRange(ctx, sqlc.InsertStateRangeParams{
				EmployeeID: employeeID,
				StateID:    o.StateID,
				StartDate:  end.AddDate(0, 0, 1),
				EndDate:    o.EndDate,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func toOverlapStates(rows []sqlc.ListOverlappingStatesRow) []overlapState {
	result := make([]overlapState, 0, len(rows))
	for _, row := range rows {
		result = append(result, overlapState{
			StateID:   row.StateID,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		})
	}
	return result
}

func toOverlapStatesByState(rows []sqlc.ListOverlappingStatesByStateRow) []overlapState {
	result := make([]overlapState, 0, len(rows))
	for _, row := range rows {
		result = append(result, overlapState{
			StateID:   row.StateID,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		})
	}
	return result
}

// findFull возвращает сотрудника с названием категории по id.
func (r *EmployeeRepository) findFull(ctx context.Context, id int64) (*domain.Employee, error) {
	row, err := r.db.FindEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	mapped := mapEmployeeRow(row)
	return &mapped, nil
}

func mapEmployeeByResourceRow(row sqlc.ListEmployeesByResourceIDRow) domain.Employee {
	return domain.Employee{
		ID:              row.ID,
		ResourceID:      row.ResourceID,
		ResourceTitle:   row.ResourceTitle,
		Position:        row.Position,
		ManagerID:       nullable.Int64Ptr(row.ManagerID),
		HireDate:        fromDate(row.HireDate),
		TerminationDate: fromDate(row.TerminationDate),
	}
}

func mapEmployeeRow(row sqlc.FindEmployeeRow) domain.Employee {
	return domain.Employee{
		ID:              row.ID,
		ResourceID:      row.ResourceID,
		ResourceTitle:   row.ResourceTitle,
		Position:        row.Position,
		ManagerID:       nullable.Int64Ptr(row.ManagerID),
		HireDate:        fromDate(row.HireDate),
		TerminationDate: fromDate(row.TerminationDate),
	}
}

func mapEmployeeListRow(row sqlc.ListEmployeesRow) domain.Employee {
	return domain.Employee{
		ID:              row.ID,
		ResourceID:      row.ResourceID,
		ResourceTitle:   row.ResourceTitle,
		Position:        row.Position,
		ManagerID:       nullable.Int64Ptr(row.ManagerID),
		HireDate:        fromDate(row.HireDate),
		TerminationDate: fromDate(row.TerminationDate),
	}
}

func mapEmployeeByManagerRow(row sqlc.ListEmployeesByManagerIDRow) domain.Employee {
	return domain.Employee{
		ID:              row.ID,
		ResourceID:      row.ResourceID,
		ResourceTitle:   row.ResourceTitle,
		Position:        row.Position,
		ManagerID:       nullable.Int64Ptr(row.ManagerID),
		HireDate:        fromDate(row.HireDate),
		TerminationDate: fromDate(row.TerminationDate),
	}
}

func mapStateRow(row sqlc.ListStatesByEmployeeRangeRow) domain.EmployeeState {
	return domain.EmployeeState{
		ID:          row.ID,
		EmployeeID:  row.EmployeeID,
		StateID:     row.StateID,
		StateCode:   row.StateCode,
		StateName:   row.StateName,
		IsAvailable: row.IsAvailable,
		StartDate:   row.StartDate,
		EndDate:     row.EndDate,
	}
}

// fromDate разворачивает nullable-дату (pgtype.Date) в [time.Time].
func fromDate(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

// toDate упаковывает [time.Time] в nullable-дату (pgtype.Date).
func toDate(v *time.Time) pgtype.Date {
	if v == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *v, Valid: true}
}
