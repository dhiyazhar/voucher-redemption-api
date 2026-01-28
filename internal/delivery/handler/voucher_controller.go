package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
	"github.com/dhiyazhar/voucher-redemption-api/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type VoucherController struct {
	VoucherUsecase *usecase.VoucherUsecase
	Logger         *slog.Logger
	Validate       *validator.Validate
}

func NewVoucherController(voucherUsecase *usecase.VoucherUsecase, logger *slog.Logger, validator *validator.Validate) *VoucherController {
	return &VoucherController{
		VoucherUsecase: voucherUsecase,
		Logger:         logger,
		Validate:       validator,
	}
}

func (h *VoucherController) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateVoucherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON format")
		return
	}

	if err := h.Validate.Struct(&req); err != nil {
		Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}

	result, err := h.VoucherUsecase.CreateVoucher(r.Context(), &req)
	if err != nil {
		h.Logger.Warn("failed to create voucher", slog.Any("error", err))
		if errors.Is(err, model.ErrConflict) {
			Error(w, http.StatusConflict, "CONFLICT", "Voucher already exists")
			return
		}

		Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	Success(w, http.StatusCreated, "CREATED", result)
}

func (h *VoucherController) GetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.Logger.Warn("failed to get voucher by code", slog.Any("error", "Code parameter is not valid"))
		Error(w, http.StatusBadRequest, "BAD_REQUEST", "Code is required")
		return
	}

	result, err := h.VoucherUsecase.GetByCode(r.Context(), code)
	if err != nil {
		h.Logger.Warn("failed to get voucher by code", slog.Any("error", err))
		if errors.Is(err, model.ErrNotFound) {
			Error(w, http.StatusNotFound, "NOT_FOUND", "Your requested voucher is not exist")
			return
		}

		Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong")
		return
	}

	Success(w, http.StatusOK, "SUCCESS", result)
}

func (h *VoucherController) ClaimVoucher(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.Logger.Warn("failed to claim voucher", slog.Any("error", "code parameter is not valid"))
		Error(w, http.StatusBadRequest, "BAD_REQUEST", "Code is required")
		return
	}

	userIdStr := r.Header.Get("X-User-ID")
	if userIdStr == "" {
		Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "X-User-ID header required")
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid user ID")
		return
	}

	result, err := h.VoucherUsecase.ClaimVoucher(r.Context(), code, userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.Logger.Warn("failed to claim voucher", slog.String("code", code), slog.Any("error", err))
			Error(w, http.StatusNotFound, "NOT_FOUND", "Voucher is not found")
			return
		}

		if errors.Is(err, model.ErrExpiredVoucher) {
			h.Logger.Warn("failed to claim voucher", slog.String("code", code), slog.Any("error", err))

			Error(w, http.StatusConflict, "VOUCHER_EXPIRED", "Voucher is expired")
			return
		}

		if errors.Is(err, model.ErrQuotaExhausted) {
			h.Logger.Warn("failed to claim voucher", slog.String("code", code), slog.Any("error", err))
			Error(w, http.StatusConflict, "VOUCHER_QUOTA_EXCEEDED", "Voucher quota insufficient")
			return
		}

		if errors.Is(err, model.ErrAlreadyRedeemed) {
			Error(w, http.StatusConflict, "ALREADY_REDEEMED", "Voucher has already redeemed")
			return
		}

		h.Logger.Warn("failed to claim voucher", slog.String("code", code), slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong")
		return
	}

	Success(w, http.StatusOK, "SUCCESS", result)
}
