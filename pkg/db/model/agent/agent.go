package agent

import (
	"context"
	"gorm.io/gorm"
	"serverApi/pkg/tools/tx"
)

func NewAgent(db *gorm.DB) *Agent {
	return &Agent{
		db: db,
	}
}

type Agent struct {
	db *gorm.DB
	tx tx.Tx
}

func (o *Agent) GetAgentInfo(ctx context.Context, platform int, page int32, size int32) (int64, error) {
	return 0, nil
}
