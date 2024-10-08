package http

import (
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/nbisso/storicard-challenge/domain"
	"github.com/nbisso/storicard-challenge/registry"
)

// Ping
//
//	@Summary		ping
//	@Description	ping
//	@Tags			ping
//	@Produce		json
//	@Success		200	{object}	string
//	@Router			/ping [get]
func RegisterRoutes(r *gin.Engine, reg *registry.Registry) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/migrations", func(c *gin.Context) {
		postMigration(c, reg)
	})

	r.GET("/users/:id/balance", func(c *gin.Context) {
		getBalance(c, reg)
	})
}

// Transactions
//
//	@Summary		Transactions
//	@Description	Transactions
//	@Tags			Transactions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"id"
//	@Param			from	query		string	false	"from"
//	@Param			to		query		string	false	"to"
//	@Success		200		{object}	domain.TransactionResult
//	@Router			/users/{id}/transactions [get]
func getBalance(c *gin.Context, reg *registry.Registry) {
	filter := domain.TransactionFilter{}
	userID := c.Param("id")

	if userID == "" {
		c.JSON(400, gin.H{
			"error": "user id is required",
		})
	}

	c.ShouldBindQuery(&filter)

	filter.UserID = userID

	result, err := reg.MigrationUsecases.GetUserBalance(c.Request.Context(), filter)

	if err != nil {

		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, result)
}

// Migrations
//
//	@Summary		Migrations
//	@Description	Migrations
//	@Tags			Migrations
//	@Accept			json
//	@Produce		json
//	@Param			file	formData	file	true	"file"
//	@Success		200		{object}	domain.Migration
//	@Router			/migrations [post]
func postMigration(c *gin.Context, reg *registry.Registry) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(400, gin.H{
			"error": "file is required",
		})
	}

	bytefile, err := file.Open()

	if err != nil {
		c.JSON(400, gin.H{
			"error": "file is required",
		})
	}

	fileb, _ := io.ReadAll(io.Reader(bytefile))

	defer bytefile.Close()

	req := domain.MigrationRequest{
		CsvFile: fileb,
	}

	result, err := reg.MigrationUsecases.NewMigration(c.Request.Context(), req)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, result)
}
