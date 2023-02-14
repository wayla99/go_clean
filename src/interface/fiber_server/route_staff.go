package fiber_server

import (
	"errors"
	"strconv"

	"github.com/samber/lo"
	"github.com/wayla99/go_clean/src/use_case"

	"github.com/gofiber/fiber/v2"
	"github.com/wayla99/go_clean/src/interface/fiber_server/docs"
)

func (f *FiberServer) addRouteStaff(base fiber.Router) {
	r := base.Group(docs.SwaggerInfo.BasePath)
	r.Post("/staffs", f.createStaff)
	r.Get("/staffs", f.getStaffs)
	r.Get("/staffs/:staff_id", f.getStaffById)
	r.Put("/staffs/:staff_id", f.updateStaffById)
	r.Delete("/staffs/:staff_id", f.deleteStaffById)
}

// createStaff godoc
// @Summary create staffs
// @Description return array of created id
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param data body Staff true "The input staff struct"
// @Success 201 {string} string "IDs of created staff"
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs [post]
func (f *FiberServer) createStaff(ctx *fiber.Ctx) error {
	var staff Staff
	err := ctx.BodyParser(&staff)
	if err != nil {
		return f.errorHandler(ctx, ErrInvalidPayload)
	}

	insertId, err := f.useCase.CreateStaff(getSpanContext(ctx), staff.toUseCase())
	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return ctx.SendString(insertId)
}

// getStaffs godoc
// @Summary get staffs
// @Description return rows of staff
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param offset query number false "offset number"
// @Param limit query number false "limit number"
// @Param search query string false "search string"
// @Success 200 {object} staffListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs [get]
func (f *FiberServer) getStaffs(ctx *fiber.Ctx) error {
	offset, err := strconv.ParseInt(ctx.Query("offset", "0"), 10, 64)
	if err != nil {
		return f.errorHandler(ctx, errors.New("invalid offset"))
	}

	limit, err := strconv.ParseInt(ctx.Query("limit", "10"), 10, 64)
	if err != nil {
		return f.errorHandler(ctx, errors.New("invalid limit"))
	}

	search := ctx.Query("search")

	staff, total, err := f.useCase.GetStaffs(getSpanContext(ctx), offset, limit, search)
	if err != nil {
		return f.errorHandler(ctx, err)
	}

	resp := staffListResponse{
		Data: lo.Map(staff, func(item use_case.Staff, _ int) Staff {
			return newStaff(item)
		}),
		Total: total,
	}
	return ctx.JSON(resp)

}

// getStaffById godoc
// @Summary get staff by id
// @Description return a row of staff
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param staff_id path string true "staff id of staff to be fetched"
// @Success 200 {object} Staff
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs/{staff_id} [get]
func (f *FiberServer) getStaffById(ctx *fiber.Ctx) error {
	sf, err := f.useCase.GetStaffById(
		getSpanContext(ctx),
		ctx.Params("staff_id"),
	)

	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return ctx.JSON(newStaff(sf))
}

// updateStaffById godoc
// @Summary update staff
// @Description return OK
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param staff_id path string true "staff id of staff to be updated"
// @Param data body Staff true "The input staff struct"
// @Success 200 {string} string OK
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs/{staff_id} [put]
func (f *FiberServer) updateStaffById(ctx *fiber.Ctx) error {
	var sf Staff
	err := ctx.BodyParser(&sf)
	if err != nil {
		return f.errorHandler(ctx, ErrInvalidPayload)
	}

	err = f.useCase.UpdateStaffById(
		getSpanContext(ctx),
		ctx.Params("staff_id"),
		sf.toUseCase(),
	)

	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return ctx.Send([]byte(OK))

}

// deleteStaffById godoc
// @Summary delete staff
// @Description return OK
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param staff_id path string true "staff id of staff to be deleted"
// @Success 200 {string} string OK
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs/{staff_id} [delete]
func (f *FiberServer) deleteStaffById(ctx *fiber.Ctx) error {
	err := f.useCase.DeleteStaffById(getSpanContext(ctx), ctx.Params("staff_id"))
	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return ctx.Send([]byte(OK))
}
