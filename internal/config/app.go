package config

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/dhiyazhar/voucher-redemption-api/internal/delivery/handler"
	"github.com/dhiyazhar/voucher-redemption-api/internal/delivery/handler/route"
	"github.com/dhiyazhar/voucher-redemption-api/internal/repository"
	"github.com/dhiyazhar/voucher-redemption-api/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB       *sql.DB
	Config   *viper.Viper
	Logger   *slog.Logger
	Mux      *http.ServeMux
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	voucherRepository := repository.NewVoucherRepository(config.DB)
	redemptionRepository := repository.NewRedemptionRepository(config.DB)

	voucherUsecase := usecase.NewVoucherUsecase(config.DB, voucherRepository, redemptionRepository, config.Logger)
	redemptionUsecase := usecase.NewRedemptionUsecase(redemptionRepository, voucherRepository, config.Logger)

	voucherController := handler.NewVoucherController(voucherUsecase, config.Logger, config.Validate)
	redemptionController := handler.NewRedemptionController(*redemptionUsecase, config.Logger)

	routeConfig := route.RouteConfig{
		VoucherController:    voucherController,
		RedemptionController: redemptionController,
		Logger:               config.Logger,
	}

	routeConfig.SetupRoutes(config.Mux)
}
