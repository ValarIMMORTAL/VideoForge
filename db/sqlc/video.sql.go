// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: video.sql

package db

import (
	"context"
)

const getVideosByUid = `-- name: GetVideosByUid :many
select id, title, url, duration, user_id, created_at, delete_at from videos
where user_id = $1
and delete_at is null
`

func (q *Queries) GetVideosByUid(ctx context.Context, userID int64) ([]Video, error) {
	rows, err := q.db.Query(ctx, getVideosByUid, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Video{}
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Url,
			&i.Duration,
			&i.UserID,
			&i.CreatedAt,
			&i.DeleteAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertVideo = `-- name: InsertVideo :one
insert into videos(
                   title,
                   url,
                   duration,
                   user_id
) values (
          $1,$2,$3,$4
         )  RETURNING id, title, url, duration, user_id, created_at, delete_at
`

type InsertVideoParams struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	Duration int32  `json:"duration"`
	UserID   int64  `json:"user_id"`
}

func (q *Queries) InsertVideo(ctx context.Context, arg InsertVideoParams) (Video, error) {
	row := q.db.QueryRow(ctx, insertVideo,
		arg.Title,
		arg.Url,
		arg.Duration,
		arg.UserID,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.Duration,
		&i.UserID,
		&i.CreatedAt,
		&i.DeleteAt,
	)
	return i, err
}
