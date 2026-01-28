package route

import (
	"log/slog"
	"net/http"

	"github.com/dhiyazhar/voucher-redemption-api/internal/delivery/handler"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type RouteConfig struct {
	VoucherController    *handler.VoucherController
	RedemptionController *handler.RedemptionController
	Logger               *slog.Logger
}

func (r *RouteConfig) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", r.healthCheck)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("POST /api/vouchers", r.VoucherController.Create)
	mux.HandleFunc("GET /api/vouchers/{code}", r.VoucherController.GetByCode)
	mux.HandleFunc("POST /api/vouchers/{code}/claim", r.VoucherController.ClaimVoucher)
	mux.HandleFunc("GET /api/vouchers/{code}/redemptions", r.RedemptionController.GetByVoucherCode)
}

func (r *RouteConfig) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
