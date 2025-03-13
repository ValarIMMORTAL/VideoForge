package db

type Store interface {
	Querier
}

// 当前结构体用于支持事物
type SQLStore struct {
	*Queries
	db DBTX
}

func NewStore(db DBTX) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
