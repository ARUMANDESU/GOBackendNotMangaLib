package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Manga struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Author          string    `json:"author"`
	Type            string    `json:"type"`
	LastUpdatedTime time.Time `json:"lastUpdatedTime"`
	Status          string    `json:"status"`
	MangaImg        string    `json:"mangaImg"`
}
type MangaCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Type        string `json:"type"`
	Status      string `json:"status"`
}

func NewManga() *Manga {
	return &Manga{Id: 0, Name: "chainsaw man", Description: "something", Type: "Manga"}
}

type MangaModel struct {
	DB *pgxpool.Pool
}

func (m *MangaModel) Insert(title string, description string, author string, mangaType string, status string) (int, error) {

	stmt := `insert into manga(name,description,author,type,last_updated_time,status,rating)
			values ($1,$2,$3,$4,current_timestamp,$5,5) returning mangaid`
	manga := &Manga{}
	result := m.DB.QueryRow(context.Background(), stmt, title, description, author, mangaType, status).Scan(&manga.Id)
	if result != nil {
		return 0, result
	}

	return manga.Id, nil
}

func (m *MangaModel) Get(id int) (*Manga, error) {
	stmt := `select mangaid,name,description,author,type,last_updated_time,status,mangaimg from manga where mangaid=$1;`
	manga := &Manga{}
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&manga.Id, &manga.Name, &manga.Description, &manga.Author, &manga.Type, &manga.LastUpdatedTime, &manga.Status, &manga.MangaImg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return manga, nil
}

func (m *MangaModel) Latest() ([]*Manga, error) {
	stmt := `select mangaid,name,description, author, type, last_updated_time, status,mangaimg  from manga order by last_updated_time desc limit 50;`
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	manga := []*Manga{}

	for rows.Next() {

		manga1 := &Manga{}

		err = rows.Scan(&manga1.Id, &manga1.Name, &manga1.Description, &manga1.Author, &manga1.Type, &manga1.LastUpdatedTime, &manga1.Status, &manga1.MangaImg)
		if err != nil {
			return nil, err
		}

		manga = append(manga, manga1)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return manga, nil
}

func (m *MangaModel) ChangeImg(id int, imgPath string) error {
	imgPath = "http://localhost:5000" + imgPath
	stmt := `update manga set mangaimg=$1 where mangaid=$2;`
	_, err := m.DB.Exec(context.Background(), stmt, imgPath, id)
	if err != nil {
		return err
	}
	return nil
}
