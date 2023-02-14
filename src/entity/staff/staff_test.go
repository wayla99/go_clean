package staff

import (
	"testing"
)

func TestStaff_Validate(t *testing.T) {
	tests := []struct {
		name    string
		staff   Staff
		wantErr bool
	}{
		{
			name: "invalid staff firstname min1",
			staff: Staff{
				Id:        "1",
				FirstName: "1",
				LastName:  "2023",
				Email:     "test1@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid staff email",
			staff: Staff{
				Id:        "2",
				FirstName: "test",
				LastName:  "2023",
				Email:     "ddsf",
			},
			wantErr: true,
		},
		{
			name: "valid staff",
			staff: Staff{
				Id:        "3",
				FirstName: "dee",
				LastName:  "dee",
				Email:     "test3@gmail.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.staff.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
