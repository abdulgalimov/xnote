package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/xnoteapp/app/common"
	"github.com/go-sql-driver/mysql"
	"time"
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
	ID        int64          `db:"id"`
	UserID    int            `db:"user_id"`
	Text      string         `db:"text"`
	FileID    string         `db:"file_id"`
	FileType  byte           `db:"file_type"`
	CreatedAt mysql.NullTime `db:"created_at"`
}

func (n *noteInner) note() *common.Note {
	return &common.Note{
		ID:        n.ID,
		UserID:    n.UserID,
		Text:      n.Text,
		FileID:    n.FileID,
		FileType:  n.FileType,
		CreatedAt: &n.CreatedAt.Time,
	}
}

type dbNotes struct {
	instance *sqlx.DB
}

func (n *dbNotes) FindAll(userID int, countOnPage int, pageNum int) ([]*common.Note, int, error) {
	notesInner := make([]*noteInner, 0)
	var err error
	var count int
	if countOnPage == 0 {
		err = n.instance.Select(&notesInner, `SELECT * FROM notes WHERE user_id=?;`, userID)
		if err != nil {
			return nil, 0, err
		}
		count = len(notesInner)
	} else {
		err = n.instance.Select(&notesInner, `SELECT * FROM notes WHERE user_id=? LIMIT ? OFFSET ?;`, userID, countOnPage, pageNum*countOnPage)
		if err != nil {
			return nil, 0, err
		}
		err := n.instance.QueryRow(`SELECT COUNT(*) as count FROM notes WHERE user_id=? ;`, userID).Scan(&count)
		fmt.Println("res", count, err)
	}
	//
	var notes []*common.Note
	for _, noteIn := range notesInner {
		notes = append(notes, noteIn.note())
	}
	return notes, count, nil
}

func (n *dbNotes) Find(noteID int) (*common.Note, error) {
	var note noteInner
	err := n.instance.Get(&note, `SELECT * FROM notes WHERE id=? LIMIT 1;`, noteID)
	if err != nil {
		return nil, err
	}
	return note.note(), nil
}

func (n *dbNotes) Create(userID int, text string) (*common.Note, error) {
	query := `INSERT INTO notes (user_id, text) VALUES(?, ?);`
	res, err := n.instance.Exec(query, userID, text)
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
	note := common.Note{
		ID:        id,
		UserID:    userID,
		Text:      text,
		FileID:    "",
		FileType:  0,
		CreatedAt: &now,
	}
	return &note, nil
}

func (n *dbNotes) Delete(noteID int) error {
	query := `DELETE FROM notes WHERE id=?;`
	_, err := n.instance.Exec(query, noteID)
	return err
}
