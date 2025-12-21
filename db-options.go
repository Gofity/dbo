package dbo

import (
	"context"
	"time"

	"github.com/Gofity/dbo/scopes"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type O struct {
	DB      *gorm.DB
	Scopes  []scopes.Scope
	Clauses []clause.Expression
	Timeout time.Duration
}

func (x *O) newContext() (ctx context.Context, cancel context.CancelFunc) {
	if x.Timeout > 0 {
		return context.WithTimeout(context.Background(), x.Timeout)
	}

	return context.WithCancel(context.Background())
}
