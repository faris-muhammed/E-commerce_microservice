package handlers

import (
	sellerpb "github.com/faris-muhammed/e-protofiles/sellerlogin"
	"github.com/gin-gonic/gin"
)

type SellerLoginHandler struct {
	sellerpb.SellerServiceClient
}

func (s *SellerLoginHandler) SellerLoginHttp(c *gin.Context) {

}
