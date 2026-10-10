package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// TrainPercent is the share of the corpus assigned to the training set.
const TrainPercent = 80

// splitCutoffSQL is the single definition of the train/test split (TE §5.2.2):
// the cut-off is the fecha of the comment at 0-based position
// floor(TrainPercent% of N) in chronological order. Train is fecha < cut-off and
// test is fecha >= cut-off, so every training comment is strictly earlier than
// every test comment, across all categories. When the 80% boundary falls inside
// a day, that whole day goes to test and train holds slightly less than 80%.
const splitCutoffSQL = `
	SELECT fecha FROM comentarios
	ORDER BY fecha
	LIMIT 1 OFFSET (SELECT COUNT(*) * ? / 100 FROM comentarios)`

// SplitCutoff returns the first fecha of the test set. Training comments are
// those with fecha < cut-off; test comments those with fecha >= cut-off.
func SplitCutoff(d *sql.DB) (string, error) {
	var cut string
	err := d.QueryRow(splitCutoffSQL, TrainPercent).Scan(&cut)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("db: split cut-off of an empty corpus")
	}
	if err != nil {
		return "", fmt.Errorf("db: split cut-off: %w", err)
	}
	return cut, nil
}
