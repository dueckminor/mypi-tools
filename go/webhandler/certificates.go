package webhandler

import (
	"net/http"

	"github.com/dueckminor/mypi-tools/go/mypi/setup"
	"github.com/gin-gonic/gin"
)

type CertHandler struct {
	certificates *setup.Certificates
}

func NewCertHandler(certs ...*setup.Certificates) *CertHandler {
	if len(certs) == 1 {
		return &CertHandler{
			certificates: certs[0],
		}
	}
	return &CertHandler{
		certificates: setup.NewCertificates(),
	}
}

func (h *CertHandler) RegisterEndpoints(r gin.IRoutes) {
	r.GET("/api/certificates", h.GetCertificates)
	r.POST("/api/certificates", h.PostCertificates)
}

// GetCertificates return list of hosts
func (h *CertHandler) GetCertificates(c *gin.Context) {
	result, err := h.certificates.Get()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *CertHandler) PostCertificates(c *gin.Context) {
	_, err := h.certificates.CreatePKI()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	result, err := h.certificates.Get()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, result)
}
