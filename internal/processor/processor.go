package processor

import db "github.com/pule1234/VideoForge/db/sqlc"

type Processor struct {
	store *db.Queries
}

func NewProcessor(store *db.Queries) *Processor {
	return &Processor{
		store: store,
	}
}
