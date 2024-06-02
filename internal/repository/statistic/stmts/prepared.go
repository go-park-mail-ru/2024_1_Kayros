package stmts

import "database/sql"


func InitPrepareStatements(db *sql.DB) (map[string]*sql.Stmt, error) {
	statements := map[string]string {
		"addAnswer": `INSERT INTO quiz(question_id, user_id, rating, created_at) VALUES($1, $2, $3, $4)`,
		"getQuestionsOnFocus": `SELECT id, name, url, focus_id, param_type FROM question WHERE url=$1`,
		"getQuestions": `SELECT id, name, param_type FROM question`,
		"selectAnswerRatingMore8": `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`,
		"selectAnswerRatingLess8": `SELECT COUNT(*) FROM quiz WHERE rating<7 AND question_id=$1`,
		"getCountOfAnswers": `SELECT COUNT(*) FROM quiz WHERE question_id=$1`,
		"getAnswerCountRatingMore8": `SELECT COUNT(*) FROM quiz WHERE rating>8 AND question_id=$1`,
		"getAnswerCount": `SELECT COUNT(*) FROM quiz WHERE question_id=$1`,
	}	
	preparedStatements := make(map[string]*sql.Stmt, len(statements))
	for key, value := range statements {
		stmt, err := db.Prepare(value)
		if err != nil {
			return nil, err
		}
		preparedStatements[key] = stmt
	}

	return preparedStatements, nil
}
