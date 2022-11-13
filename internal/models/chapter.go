package models

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Chapter struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	ChapterNumber string   `json:"chapterNumber"`
	VolumeNumber  string   `json:"volumeNumber"`
	Images        []string `json:"images"`
}

type ChapterModel struct {
	DB *pgxpool.Pool
}

func (ch *ChapterModel) Insert(title string, images []string, chapterNum float64, volumeNum float64) (int, error) {
	stmt := `insert into chapter(images, chapter_number, volume_number, title)
				values ($1,$2,$3,$4) returning chapterid`
	newChapter := Chapter{}
	result := ch.DB.QueryRow(context.Background(), stmt, images, chapterNum, volumeNum, title).Scan(&newChapter.Id)
	if result != nil {
		return 0, result
	}

	return newChapter.Id, nil
}
