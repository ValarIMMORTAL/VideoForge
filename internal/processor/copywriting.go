package processor

import (
	"encoding/json"
	"fmt"
	"github.com/pule1234/VideoForge/internal/models"
)

func CreateCopyWriting(item models.TrendingItem) error {
	res, _ := json.Marshal(item)
	fmt.Println(res)
	return nil
}
