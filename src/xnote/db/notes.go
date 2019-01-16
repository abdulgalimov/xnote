package db

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
	"xnote/models"
)

var notesScheme = `
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    text varchar(4000) NULL DEFAULT NULL,
    file_id varchar(128) NULL DEFAULT '',
    file_type SMALLINT NULL DEFAULT 0,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);
`

type noteInner struct {
	Id        int64          `db:"id"`
	UserId    int            `db:"user_id"`
	Text      string         `db:"text"`
	FileId    string         `db:"file_id"`
	FileType  byte           `db:"file_type"`
	CreatedAt mysql.NullTime `db:"created_at"`
}

func (n *noteInner) note() *models.Note {
	return &models.Note{
		Id:        n.Id,
		UserId:    n.UserId,
		Text:      n.Text,
		FileId:    n.FileId,
		FileType:  n.FileType,
		CreatedAt: &n.CreatedAt.Time,
	}
}

type dbNotes struct {
}

func (n *dbNotes) FindAll(userId int, countOnPage int, pageNum int) ([]*models.Note, int, error) {
	notesInner := make([]*noteInner, 0)
	var err error
	var count int
	if countOnPage == 0 {
		err = dbInstance.Select(&notesInner, `SELECT * FROM notes WHERE user_id=?;`, userId)
		if err != nil {
			return nil, 0, err
		}
		count = len(notesInner)
	} else {
		err = dbInstance.Select(&notesInner, `SELECT * FROM notes WHERE user_id=? LIMIT ? OFFSET ?;`, userId, countOnPage, pageNum*countOnPage)
		if err != nil {
			return nil, 0, err
		}
		err := dbInstance.QueryRow(`SELECT COUNT(*) as count FROM notes WHERE user_id=? ;`, userId).Scan(&count)
		fmt.Println("res", count, err)
	}
	//
	var notes []*models.Note
	for _, noteIn := range notesInner {
		notes = append(notes, noteIn.note())
	}
	return notes, count, nil
}

func (n *dbNotes) Find(noteId int) (*models.Note, error) {
	var note noteInner
	err := dbInstance.Get(&note, `SELECT * FROM notes WHERE id=? LIMIT 1;`, noteId)
	if err != nil {
		return nil, err
	}
	return note.note(), nil
}

func (n *dbNotes) Create(userId int, text string) (*models.Note, error) {
	query := `INSERT INTO notes (user_id, text) VALUES(?, ?);`
	res, err := dbInstance.Exec(query, userId, text)
	if err != nil {
		return nil, err
	}
	//
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	//
	now := time.Now()
	note := models.Note{
		Id:        id,
		UserId:    userId,
		Text:      text,
		FileId:    "",
		FileType:  0,
		CreatedAt: &now,
	}
	return &note, nil
}

func (n *dbNotes) Delete(noteId int) error {
	query := `DELETE FROM notes WHERE id=?;`
	_, err := dbInstance.Exec(query, noteId)
	return err
}
