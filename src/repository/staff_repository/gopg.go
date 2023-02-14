package staff_repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-pg/pg/v10/orm"

	"github.com/wayla99/go_clean/src/use_case"

	"github.com/wayla99/go_clean/src/entity/staff"

	"github.com/go-pg/pg/v10"
	lop "github.com/samber/lo/parallel"
)

type logger struct {
}

func NewGoPg(ctx context.Context, db *pg.DB) (use_case.StaffRepository, error) {
	g := &goPg{db: db}
	err := g.initTable(ctx)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g goPg) initTable(ctx context.Context) error {
	tOpts := orm.CreateTableOptions{
		Temp:          false,
		IfNotExists:   true,
		FKConstraints: true,
	}

	err := g.db.Model((*goPgStaff)(nil)).Context(ctx).CreateTable(&tOpts)
	if err != nil {
		return err
	}

	return nil
}

func (l logger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	qs, _ := q.FormattedQuery()
	fmt.Println(string(qs))
	return ctx, nil
}

func (l logger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	return nil
}

type goPg struct {
	db *pg.DB
}

func (g goPg) Health(ctx context.Context) error {
	return g.db.Ping(ctx)
}

func (g goPg) CreateStaff(ctx context.Context, staff staff.Staff) (string, error) {
	d, err := toGoPgStaff(staff)
	if err != nil {
		return "", err
	}

	err = g.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		_, err := tx.ModelContext(ctx, &d).Insert()
		if err != nil {
			return err
		}

		return nil
	})

	return strconv.Itoa(d.ID), err
}

func (g goPg) GetStaffs(ctx context.Context, offset, limiit int64, search string) ([]staff.Staff, uint64, error) {
	var defaults []goPgStaff
	count, err := g.db.ModelContext(ctx, &defaults).Offset(int(offset)).Limit(int(limiit)).SelectAndCount()

	if err != nil {
		return nil, 0, err
	}

	return lop.Map(defaults, func(item goPgStaff, _ int) staff.Staff {
		return item.toEntity()
	}), uint64(count), nil
}

func (g goPg) GetStaffById(ctx context.Context, staffId string) (staff.Staff, error) {
	var s goPgStaff
	err := g.db.ModelContext(ctx, &s).Where("staff_id = ?", staffId).First()

	if err != nil {
		return staff.Staff{}, err
	}
	return s.toEntity(), nil
}

func (g goPg) UpdateStaffById(ctx context.Context, staffId string, staff staff.Staff) error {
	err := g.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		s, err := toGoPgStaff(staff)
		if err != nil {
			return err
		}

		_, err = tx.ModelContext(ctx, &s).Where("staff_id = ?", staffId).UpdateNotZero()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (g goPg) DeleteStaffById(ctx context.Context, staffId string) error {
	r, err := g.db.ModelContext(ctx, (*goPgStaff)(nil)).Where("staff_id = ?", staffId).Delete()
	if err != nil {
		return err
	}

	if r.RowsAffected() == 0 {
		return use_case.ErrStaffNotFound
	}
	return nil
}
