package main

import (
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
  _ "reflect"
  "time"
)

const (
  host = "<redacted>"
  port = 5432
  user = "<redacted>"
  password = "<redacted>"
  dbname = "<redacted>"
)

type AnswerDTO struct {
  ID          int
  UserID      int
  SurveyID    int
  QuestionID  int
  Answer      string
  Completed   bool
  Seq         int
  CreatedAt   time.Time
  UpdatedAt   time.Time
  Won         bool
  Data        string
  Lat         float64
  Lng         float64
}

// func main() {
//   db := connectToDb()
//
//   // fmt.Println(reflect.TypeOf(db))
//   // defers(*db)
//   defer db.Close()
//
//   err := db.Ping()
//   if err != nil {
//     panic(err)
//   }
//
//   fmt.Println("Successfully connected!")
//
//   // inserts(*db)
//   // updates(*db)
//   // selects(*db)
//   // deletes(*db)
//   queryMethod(*db)
// }

func connectToDb() *sql.DB {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " +
    "password=%s dbname=%s",
    host, port, user, password, dbname)

  db, err := sql.Open("postgres", psqlInfo)

  if err != nil {
    panic(err)
  }

  return db
}

func inserts(db sql.DB) {
  sqlStatement := `
  INSERT INTO answers (user_id, survey_id, question_id, answer, completed, int_col, bool_col, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  RETURNING id`
  id := 0
  // _, err := db.Exec(sqlStatement, 1, 52, 442, "foobar", false, 1, false, time.Now(), time.Now())
  err := db.QueryRow(sqlStatement, 1, 52, 442, "foobar", false, 1, false, time.Now(), time.Now()).Scan(&id)
  panicIfErr(err)
  fmt.Println("New record ID is:", id)
}

func updates(db sql.DB) {
  sqlStatement := `
  UPDATE answers
  SET answer = $1
  WHERE id = $2
  RETURNING question_id, answer;`
  var answer string
  var questionId string
  // res, err := db.Exec(sqlStatement, "Zoro", 18298)
  err := db.QueryRow(sqlStatement, "Zoro", 18298).Scan(&questionId, &answer)
  panicIfErr(err)

  // count, err := res.RowsAffected()
  // panicIfErr(err)
  // fmt.Println("rowsAffected:",count)
  fmt.Println("question_id:", questionId, ",answer:", answer)
}

func deletes(db sql.DB) {
  sqlStatement := `
  DELETE FROM answers
  WHERE id = $1
  `
  _, err := db.Exec(sqlStatement, 18297)
  panicIfErr(err)
}

func selects(db sql.DB) {
  sqlStatement := `
  SELECT id, user_id, survey_id, question_id, answer, created_at FROM answers WHERE id=$1;`
  var answer AnswerDTO
  // var questionId string
  // var answer string
  row := db.QueryRow(sqlStatement, 18299)
  switch err := row.Scan(&answer.ID, &answer.UserID, &answer.SurveyID,
    &answer.QuestionID, &answer.Answer, &answer.CreatedAt); err {
  case sql.ErrNoRows:
    fmt.Println("No rows were returned!")
    return
  case nil:
    // fmt.Println("question_id:", questionId, "answer:", answer)
    fmt.Println(answer)
  default:
    panic(err)
  }
}

func queryMethod(db sql.DB) {
  rows, err := db.Query("SELECT id, survey_id FROM answers LIMIT $1", 3)
  panicIfErr(err)
  defer rows.Close()
  // returns true when the next row is successfully prepared; false otherwise
  for rows.Next() {
    var id int
    var surveyId int
    err = rows.Scan(&id, &surveyId)
    panicIfErr(err)
    fmt.Println(id, surveyId)
  }
  // get any error encountered during iteration
  err = rows.Err()
  panicIfErr(err)
}


func panicIfErr(err error) {
  if err != nil {
    panic(err)
  }
}
