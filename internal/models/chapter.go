package models

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Chapter struct {
	Id            int       `json:"id"`
	MangaId       int       `json:"mangaId"`
	Title         string    `json:"title"`
	ChapterNumber float64   `json:"chapterNumber"`
	VolumeNumber  float64   `json:"volumeNumber"`
	Images        []string  `json:"images"`
	Date          time.Time `json:"date"`
	DateString    string    `json:"dateString"`
}

type ChapterModel struct {
	DB *pgxpool.Pool
}

func (ch *ChapterModel) Insert(mangaId int, title string, chapterNum float64, volumeNum float64) (int, error) {
	stmt := `insert into chapter(chapter_number, volume_number, title,date)
				values ($1,$2,$3,current_date) returning chapterid;`
	newChapter := Chapter{}
	result := ch.DB.QueryRow(context.Background(), stmt, chapterNum, volumeNum, title).Scan(&newChapter.Id)
	if result != nil {
		return 0, result
	}

	stmt = `insert into manga_chapter(chapterid, mangaid) values($1,$2);`
	_, err := ch.DB.Exec(context.Background(), stmt, newChapter.Id, mangaId)
	if err != nil {
		return 0, err
	}

	return newChapter.Id, nil
}

func (ch *ChapterModel) ChangeImages(id int, imagesPaths []string) error {
	stmt := `update chapter set images=$1 where chapterid=$2;`
	_, err := ch.DB.Exec(context.Background(), stmt, imagesPaths, id)
	if err != nil {
		return err
	}
	return nil
}
func (ch *ChapterModel) Get(mangaId int, chapterNum float64, volumeNum float64) (*Chapter, error) {
	stmt := `select ch.chapterid,m.mangaid,ch.title,ch.volume_number,ch.chapter_number,ch.images
			from chapter ch join manga_chapter mc on ch.chapterid = mc.chapterid 
						join manga m on m.mangaid = mc.mangaid
						where m.mangaid=$1 and ch.chapter_number=$2 and ch.volume_number=$3 `
	chapter := Chapter{}
	result := ch.DB.QueryRow(context.Background(), stmt, mangaId, chapterNum, volumeNum).Scan(&chapter.Id, &chapter.MangaId, &chapter.Title, &chapter.VolumeNumber, &chapter.ChapterNumber, &chapter.Images)
	if result != nil {
		return nil, result
	}

	return &chapter, nil
}

func (ch *ChapterModel) GetMangaChapters(mangaId int) ([]Chapter, error) {
	stmt := `select ch.chapterid,m.mangaid,ch.title,ch.volume_number,ch.chapter_number,ch.images,ch.date
			from chapter ch join manga_chapter mc on ch.chapterid = mc.chapterid 
						join manga m on m.mangaid = mc.mangaid
						where m.mangaid=$1`
	rows, err := ch.DB.Query(context.Background(), stmt, mangaId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	chapter := []Chapter{}
	for rows.Next() {
		chapter1 := Chapter{}
		result := rows.Scan(&chapter1.Id, &chapter1.MangaId, &chapter1.Title, &chapter1.VolumeNumber, &chapter1.ChapterNumber, &chapter1.Images, &chapter1.Date)
		chapter1.DateString = chapter1.Date.Format("2006-1-2")
		if result != nil {
			return nil, result
		}
		chapter = append(chapter, chapter1)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chapter, nil
}
