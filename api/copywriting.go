package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"net/http"
)

func (server *Server) getCopyWriting(c *gin.Context) {
	var req getCopyWritingRequest

	arg := db.ListCopiesParams{
		Date: pgtype.Timestamp{
			Time:  req.date,
			Valid: true,
		},
		Limit:  req.num,
		Offset: req.num * (req.page - 1),
	}
	copies, err := server.store.ListCopies(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败" + err.Error()})
		return
	}

	var resp getCopyWritingResponse
	for _, copy := range copies {
		resp.items = append(resp.items, Copywriting{
			ID:      copy.ID,
			Date:    copy.Date.Time,
			Title:   copy.Title,
			Source:  copy.Source,
			Content: copy.Content,
		})
	}

	c.JSON(http.StatusOK, resp)
}
