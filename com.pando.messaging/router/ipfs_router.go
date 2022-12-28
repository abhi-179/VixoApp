package router

import (
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	ipfs "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewIPFSController(e *echo.Echo, ipfsusecase ipfs.IPFSUsecase) {
	handler := &controller.IPFSController{
		Usecase: ipfsusecase,
	}

	e.POST("api/v1/ipfs/aws", handler.UploadFile)
	e.POST("api/v1/ipfs/save_backup_filehash", handler.SaveBackupFilehash)
	e.GET("api/v1/ipfs/get_backup_filehash", handler.GetBackupFilehash)
	e.POST("api/v1/ipfs", handler.UploadFileToIPFS)
}
