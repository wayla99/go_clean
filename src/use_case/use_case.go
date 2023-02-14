package use_case

import (
	"context"
	"errors"

	"github.com/wayla99/go_clean/src/entity/staff"
)

var (
	ErrStaffNotFound = errors.New("staff not found")
)

type UseCase struct {
	staffRepository StaffRepository
}

type StaffRepository interface {
	CreateStaff(ctx context.Context, staff staff.Staff) (string, error)
	GetStaffs(ctx context.Context, offset, limit int64, search string) ([]staff.Staff, uint64, error)
	GetStaffById(ctx context.Context, staffId string) (staff.Staff, error)
	UpdateStaffById(ctx context.Context, staffId string, staff staff.Staff) error
	DeleteStaffById(ctx context.Context, staffId string) error
	Health(ctx context.Context) error
}

func New(staffRepo StaffRepository) *UseCase {
	return &UseCase{staffRepository: staffRepo}
}
