package response

import "github.com/gin-gonic/gin"

type DataResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, DataResponse{Data: data, Error: nil})
}

func Error(c *gin.Context, status int, err string) {
	c.JSON(status, DataResponse{Data: nil, Error: err})
}

// response.error(c, http.StatusBadRequest, "error adding films please try again"})
// response.success(c, http.StatusCreated, gin.H{"message": "Film Added"})
