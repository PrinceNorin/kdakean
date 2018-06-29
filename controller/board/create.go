package boardController

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kdakean/kdakean/model"
	"github.com/kdakean/kdakean/util/messages"
	"github.com/kdakean/kdakean/util/validator"
)

func PostBoardHandler(c *gin.Context) {
	var b model.Board
	c.Bind(&b)

	f := &model.CreateBoardJSON{
		Name:        b.Name,
		Description: b.Description.String,
		UserId:      b.UserId,
	}
	msg := messages.GetMessages(c)
	if err := validator.Validate(f, msg); err != nil {
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	board, err := model.CreateBoard(&b)
	if err != nil {
		msg.AddErrorT("user_id", err.Error())
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"board": board,
	})
}

func PutBoardHandler(c *gin.Context) {
	var b model.Board
	c.Bind(&b)
	if id, err := strconv.ParseUint(c.Param("id"), 10, 0); err == nil {
		b.Id = uint(id)
	}

	f := &model.UpdateBoardJSON{
		Name:        b.Name,
		Description: b.Description.String,
	}
	msg := messages.GetMessages(c)
	if err := validator.Validate(f, msg); err != nil {
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	board, err := model.UpdateBoard(&b)
	if err != nil {
		msg.ErrorT(err)
		c.JSON(http.StatusBadRequest, msg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"board": board,
	})
}
