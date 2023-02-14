package use_case

import (
	"context"

	lop "github.com/samber/lo/parallel"

	"github.com/wayla99/go_clean/src/entity/staff"
)

type Staff struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
}

func (s Staff) toEntity() staff.Staff {
	return staff.Staff{
		Id:        s.Id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}
}

func (uc UseCase) enrichStaff(ctx context.Context, s staff.Staff) (Staff, error) {
	var err error
	return Staff{
		Id:        s.Id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}, err
}

func (uc UseCase) CreateStaff(ctx context.Context, s Staff) (string, error) {
	sf := s.toEntity()
	if err := sf.Validate(); err != nil {
		return "", err
	}
	sf.Id = ""
	insertId, err := uc.staffRepository.CreateStaff(ctx, sf)
	if err != nil {
		return "", err
	}
	return insertId, nil
}

func (uc UseCase) GetStaffs(ctx context.Context, offset, limit int64, search string) ([]Staff, uint64, error) {
	sf, total, err := uc.staffRepository.GetStaffs(ctx, offset, limit, search)
	if err != nil {
		return nil, 0, err
	}

	return lop.Map(sf, func(item staff.Staff, _ int) Staff {
		s, err := uc.enrichStaff(ctx, item)
		if err != nil {
			return Staff{}
		}

		return s
	}), total, err
}

func (uc UseCase) GetStaffById(ctx context.Context, staffId string) (Staff, error) {
	s, err := uc.staffRepository.GetStaffById(ctx, staffId)
	if err != nil {
		return Staff{}, err
	}

	return uc.enrichStaff(ctx, s)
}

func (uc UseCase) UpdateStaffById(ctx context.Context, staffId string, s Staff) error {
	sf := s.toEntity()

	if err := sf.Validate(); err != nil {
		return err
	}

	err := uc.staffRepository.UpdateStaffById(ctx, staffId, sf)
	if err != nil {
		return err
	}

	return nil
}

func (uc UseCase) DeleteStaffById(ctx context.Context, staffId string) error {
	err := uc.staffRepository.DeleteStaffById(ctx, staffId)
	if err != nil {
		return err
	}

	return nil
}
