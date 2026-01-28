package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
	"github.com/dhiyazhar/voucher-redemption-api/internal/usecase"
)

type RedemptionController struct {
	RedemptionUsecase usecase.RedemptionUsecase
	Logger            *slog.Logger
}

func NewRedemptionController(redemptionUsecase usecase.RedemptionUsecase, logger *slog.Logger) *RedemptionController {
	return &RedemptionController{
		RedemptionUsecase: redemptionUsecase,
		Logger:            logger,
	}
}

func (h *RedemptionController) GetByVoucherCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.Logger.Error("failed to get redemptions", slog.Any("error", "code parameter is not valid"))
		Error(w, http.StatusBadRequest, "BAD_REQUEST", "Code is required")
		return
	}

	result, err := h.RedemptionUsecase.GetByVoucherCode(r.Context(), code)
	if err != nil {
		h.Logger.Error("failed to get redemptions by voucher code",
			slog.String("code", code),
			slog.Any("error", err),
		)

		if errors.Is(err, model.ErrNotFound) {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Voucher not found")
			return
		}

		Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong")
		return
	}

	Success(w, http.StatusOK, "SUCCESS", result)
}
