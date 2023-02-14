package fiber_server

import "github.com/wayla99/go_clean/src/use_case"

type staffListResponse struct {
	Data  []Staff `json:"data"`
	Total uint64  `json:"total"`
}

type Staff struct {
	Id        string `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (s Staff) toUseCase() use_case.Staff {
	return use_case.Staff{
		Id:        s.Id,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
	}
}

func newStaff(c use_case.Staff) Staff {
	return Staff{
		Id:        c.Id,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
	}
}
