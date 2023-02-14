package staff_repository

import (
	"strconv"

	"github.com/wayla99/go_clean/src/entity/staff"
)

type goPgStaff struct {
	tableName struct{} `pg:"staff,alias:c"`
	ID        int      `pg:"staff_id,pk,notnull"`
	FirstName string   `pg:"first_name"`
	LastName  string   `pg:"last_name"`
	Email     string   `pg:"email"`
}

func toGoPgStaff(s staff.Staff) (goPgStaff, error) {
	var err error
	id := 0

	if s.Id != "" {
		id, err = strconv.Atoi(s.Id)
		if err != nil {
			return goPgStaff{}, err
		}
	}

	return goPgStaff{
		ID:        id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}, nil
}

func (s goPgStaff) toEntity() staff.Staff {
	return staff.Staff{
		Id:        strconv.Itoa(s.ID),
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}
}
