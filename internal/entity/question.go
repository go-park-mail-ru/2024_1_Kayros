package entity

import "database/sql"

type Question struct {
	Id        uint64
	Name      string
	Url       string
	FocusId   string
	ParamType string
}

type QuestionSql struct {
	Id        uint64         `json:"id"`
	Name      string         `json:"name"`
	Url       string         `json:"url"`
	FocusId   sql.NullString `json:"focus_id"`
	ParamType string         `json:"param_type"`
}

func QuestionFromDB(sqlRow *QuestionSql) *Question {
	return &Question{
		Id:        sqlRow.Id,
		Name:      sqlRow.Name,
		Url:       sqlRow.Url,
		FocusId:   String(sqlRow.FocusId),
		ParamType: sqlRow.ParamType,
	}
}
