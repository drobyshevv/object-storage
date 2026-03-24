package handler

import (
	"strconv"

	"github.com/drobyshevv/object-storage/internal/service"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *service.FileService
}

func NewFileHandler(s *service.FileService) *FileHandler {
	return &FileHandler{service: s}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Upload(c, file)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)
}

func (h *FileHandler) List(c *gin.Context) {
	files, err := h.service.GetAll(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, files)
}

func (h *FileHandler) Download(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	body, filename, err := h.service.GetFile(c, id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	defer body.Close()

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.DataFromReader(200, -1, "application/octet-stream", body, nil)
}

func (h *FileHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.Delete(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
