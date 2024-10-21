package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) uploadFile(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://"+h.cfg.UploaderURL+"/upload", c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,  err)
		return
	}
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,  err)
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func (h *Handler) deleteFile(c *gin.Context) {
	req, err := http.NewRequest("POST", "http://"+h.cfg.UploaderURL+"/delete", c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,  err)
		return
	}
	req.Header = c.Request.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError,  err)
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
