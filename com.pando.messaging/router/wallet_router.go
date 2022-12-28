package router

import (
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	wallet "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewWalletController(e *echo.Echo, walletUsecase wallet.WalletUsecase) {
	handler := &controller.WalletController{
		Usecase: walletUsecase,
	}
	e.POST("api/v1/wallet/create_wallet", handler.CreateWallet)
	e.GET("api/v1/wallet/get_balance", handler.GetBalance)
	e.POST("api/v1/wallet/add_token", handler.AddToken)
	e.POST("api/v1/wallet/request_token", handler.RequestToken)
	e.GET("api/v1/wallet/reject_or_accept_request", handler.RejectRequest)
	e.POST("api/v1/wallet/send_token", handler.SendToken)
	e.GET("api/v1/wallet/get_transactions", handler.GetTransactions)
	e.GET("api/v1/wallet/view_spend_analytics", handler.ViewSpendAnalytics)
	e.POST("api/v1/wallet/wallet_statement", handler.WalletStatement)
	e.GET("api/v1/wallet/get_wallet_id", handler.GetWalletId)
	e.GET("api/v1/wallet/get_recent_transactions", handler.RecentTransactions)
	e.GET("api/v1/wallet/show_pending_requests", handler.ShowPendingRequests)
	e.POST("api/v1/wallet/send_token_to_admin", handler.SendTokenToAdmin)
	e.GET("api/v1/wallet/show_own_token_requests", handler.ShowOwnTokenRequests)

}
