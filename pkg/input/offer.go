package input
import (
		"github.com/go-playground/validator/v10"
)
type PawnRequest struct {
	Code   string `json:"code" validate:"required"`
	Offer  *int    `json:"offer" validate:"required"`
	Demand *int    `json:"demand" validate:"required"`
}

func ValidateRequest(pawnReq PawnRequest) error {
	validate := validator.New()
	return validate.Struct(pawnReq)
}