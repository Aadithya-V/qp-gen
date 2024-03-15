package services

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"strings"
	"text/template"

	"github.com/Aadithya-V/qp-gen/database"
	"github.com/Aadithya-V/qp-gen/pkg/api/v1/models"
	"github.com/gin-gonic/gin"
)

func GenerateQpaperSetsFromDB(c *gin.Context, req *models.GenerateQpaperSetsFromDBRequest) (*bytes.Buffer, error) {

	var units []string

	// Unit 1:Section A: int
	var qsPerType map[string]int

	if req.ExamType == "CAT_1" {
		units = []string{"1", "2"}
		qsPerType = map[string]int{
			"1:A": 5,
			"1:B": 5,
			"1:C": 1,

			"2:A": 5,
			"2:B": 5,
			"2:C": 1,
		}
	} else if req.ExamType == "CAT_2" {
		units = []string{"3", "4"}
		qsPerType = map[string]int{
			"3:A": 5,
			"3:B": 5,
			"3:C": 1,

			"4:A": 5,
			"4:B": 5,
			"4:C": 1,
		}
	} else if req.ExamType == "SEM" {
		units = []string{"1", "2", "3", "4", "5"}
		qsPerType = map[string]int{
			"1:A": 2,
			"1:B": 2,
			"1:C": 1,

			"2:A": 2,
			"2:B": 2,
			"2:C": 0,

			"3:A": 2,
			"3:B": 2,
			"3:C": 0,

			"4:A": 2,
			"4:B": 2,
			"4:C": 1,

			"5:A": 2,
			"5:B": 2,
			"5:C": 0,
		}
	} else {
		return nil, errors.New("unknown exam type")
	}

	questions, err := GetQuestionFromDB(req, units)
	if err != nil {
		return nil, err
	}

	qStore := BuildQStore(questions, qsPerType)

	tmpl, err := template.ParseFiles("vec_template.tex")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return nil, err
	}

	//var overleafLinks []string = make([]string, 0)
	// Create a buffer to store the ZIP archive
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	for i := req.NumberOfSets; i > 0; i-- {

		qsbySection := qStore.PickQSet()

		type TemplatePipeline struct {
			// for section a
			SectionA []*Question
			// for b and c
			QuestionsBySection map[string][][2]*Question

			ExamDetails models.ExamDetails
		}

		var tmplData = TemplatePipeline{
			SectionA:           make([]*Question, 0),
			QuestionsBySection: make(map[string][][2]*Question),
			ExamDetails:        req.ExamDetails,
		}

		fmt.Println("--------------------------------------------------------")

		for k, v := range qsbySection {
			jsonD, _ := json.Marshal(v)
			fmt.Println(k, string(jsonD))

			var qpairs = make([][2]*Question, 0)

			if k == "A" {
				for i := 0; i < len(v); i++ {
					tmplData.SectionA = append(tmplData.SectionA, v[i])
				}
			} else {
				for j := 1; j < len(v); j += 2 {
					var qpair [2]*Question
					v[j-1].Choice = true
					qpair[0] = v[j-1]
					qpair[1] = v[j]
					qpairs = append(qpairs, qpair)
				}
			}
			tmplData.QuestionsBySection[k] = append(tmplData.QuestionsBySection[k], qpairs...)

		}

		delete(tmplData.QuestionsBySection, "A")

		// var buf bytes.Buffer

		// Create a new file in the archive
		fileWriter, err := zipWriter.Create(fmt.Sprintf("file-%v.tex", i))
		if err != nil {
			return nil, fmt.Errorf("Failed to create file in archive. Err: %w", err.Error())
		}

		err = tmpl.Execute(fileWriter, tmplData)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return nil, err
		}

		/* link := "https://www.overleaf.com/docs?encoded_snip=" + url.QueryEscape(buf.String())
		link = strings.ReplaceAll(link, "%5Cn", "%0A")
		link = strings.ReplaceAll(link, "%5C%5C", "%5C")

		overleafLinks = append(overleafLinks, link) */

	}

	if err = zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("Failed to close ZIP writer. Err: %w", err.Error())
	}

	return zipBuffer, nil
}

func GetQuestionFromDB(req *models.GenerateQpaperSetsFromDBRequest, units []string) ([]*Question, error) {

	inClause := "'" + strings.Join(units, "','") + "'"
	// sub_code = %s and
	query := fmt.Sprintf(`select unit, section, question from question_bank_data where unit in (%s);`, inClause)

	db := database.NewMySQLSession()
	defer db.Close()

	results, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error while executing DB query")
	}

	questions := make([]*Question, 0)

	for results.Next() {
		question := Question{}

		var jsonData string

		err := results.Scan(&question.Unit, &question.Section, &jsonData)
		if err != nil {
			fmt.Println(err.Error())
			if err == sql.ErrNoRows {
				break
			}
			return nil, errors.New("error scanning from db")
		}

		if err := json.Unmarshal([]byte(jsonData), &question); err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("error unmarshalling question json from db. probably wrong format")
		}

		question.Type = question.Unit + ":" + question.Section

		fmt.Println(question)

		questions = append(questions, &question)
	}
	return questions, nil
}

type QStore struct {
	ByTypes   map[string][]*Question
	QsPerType map[string]int
}

func BuildQStore(questions []*Question, qsPerType map[string]int) QStore {

	ret := QStore{
		ByTypes:   make(map[string][]*Question),
		QsPerType: qsPerType,
	}

	for k := range qsPerType {
		ret.ByTypes[k] = make([]*Question, 0)
	}

	for _, question := range questions {
		fmt.Println(question)

		ret.ByTypes[question.Type] = append(ret.ByTypes[question.Type], question)
	}

	return ret
}

type Question struct {
	Unit         string
	Section      string
	Type         string
	PickedCount  int
	Choice       bool
	Marks        string     `json:"mark"`
	Text         string     `json:"text"`          // can be empty
	SubQuestions []Question `json:"sub_questions"` // sub-divisions
}

// qs by section in result
func (QS *QStore) PickQSet() map[string][]*Question {

	ret := make(map[string][]*Question)

	for qType, qs := range QS.ByTypes {
		if len(qs) == 0 {
			continue
		}
		questions := pickQs(qs, QS.QsPerType[qType])

		for _, question := range questions {
			if _, ok := ret[question.Section]; !ok {
				ret[question.Section] = make([]*Question, 0)
			}
			ret[question.Section] = append(ret[question.Section], question)
		}
		// fmt.Printf("\ntype: %v marks: %v question: %s", qType, ret[qType].Marks, ret[qType].Qs[0])
	}

	/* for k, v := range ret {
		fmt.Printf("\ntype: %v marks: %v question: %s", k, v.Marks, v.Qs[0])
	} */

	return ret
}

func pickQs(qs []*Question, nums int) []*Question {
	fmt.Println("MARKER pickQs()")
	ret := make([]*Question, 0)

	r := rand.Intn(len(qs) - 1)

	for i := r + 1; i < len(qs); i = (i + 1) % len(qs) {
		if nums == 0 {
			return ret
		}
		if i == r {
			i = rand.Intn(len(qs))
			q := qs[i]
			ret = append(ret, q)
			q.PickedCount++
			nums--

			continue
		}

		if qs[i].PickedCount == 0 {
			ret = append(ret, qs[i])
			qs[i].PickedCount++
			nums--
		}
	}
	fmt.Println("MARKER NOW")
	for _, s := range ret {

		fmt.Println(s)
	}

	return ret
}

// unit,section,external_ref,question,sub-question,marks,sub-question,marks,sub-question,marks,sub-question,marks
var CSV_HEADER_ROW = []string{"unit", "section", "external_ref", "question", "sub_question", "marks", "sub_question", "marks", "sub_question", "marks", "sub_question", "marks"}

func ParseAndSaveQuestionsFromCSV(c *gin.Context, fh *multipart.FileHeader, subCode, year string) error {
	file, err := fh.Open()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer file.Close()

	// Parse the CSV file
	records, err := parseCSVFile(file)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	j, _ := json.Marshal(records)
	fmt.Println(string(j))

	if len(records) <= 1 { // including only header row and no data
		err := errors.New("no data in uploaded file")
		log.Println(err.Error())
		return err
	}

	headerRow := records[0]
	records = records[1:]

	if !stringArraysEqual(headerRow, CSV_HEADER_ROW) {
		fmt.Println(headerRow)
		err := fmt.Errorf("wrong header format. required- %s", CSV_HEADER_ROW)
		log.Println(err.Error())
		return err
	}

	dbRows, err := prepareDataForBatchInsert(records)
	if err != nil {
		err := fmt.Errorf("error preparing data for db batch insert- error: %w", err)
		log.Println(err.Error())
		return err
	}

	err = BatchInsertCsvData(dbRows, year, subCode)
	if err != nil {
		err := fmt.Errorf("error during db batch insert- error: %w", err)
		log.Println(err.Error())
		return err
	}

	return nil
}

type DbRow struct {
	Unit        string
	Section     string
	Marks       string
	ExternalRef string // q number in question bank / external-ref
	Question    DbRowQuestion
}

type DbRowQuestion struct {
	Marks        string          `json:"mark,omitempty"`
	Text         string          `json:"text,omitempty"`          // can be empty
	SubQuestions []DbRowQuestion `json:"sub_questions,omitempty"` // sub-divisions
}

func BatchInsertCsvData(dbRows []*DbRow, year, subCode string) error {
	db := database.NewMySQLSession()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Prepare the insert statement
	//select unit, section, question from question_bank_data where
	stmt, err := tx.Prepare("INSERT INTO question_bank_data (academic_year, sub_code, external_ref, unit, section, mark, question) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, row := range dbRows {
		j, err := json.Marshal(row.Question)
		if err != nil {
			tx.Rollback() // Rollback if executing the statement fails
			log.Println(err)
			return err
		}

		_, err = stmt.Exec(year, subCode, row.ExternalRef, row.Unit, row.Section, row.Marks, string(j))
		if err != nil {
			tx.Rollback() // Rollback if executing the statement fails
			log.Println(err)
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback if committing the transaction fails
		log.Println(err)
		return err
	}

	log.Println("Batch insert completed successfully")

	return nil
}

func validateCsvRecord(record []string) bool {
	return true
}

func prepareDataForBatchInsert(records [][]string) ([]*DbRow, error) {
	dbRows := make([]*DbRow, 0)
	for _, record := range records {
		if !validateCsvRecord(record) || len(record) < len(CSV_HEADER_ROW) {
			continue // continue if empty
			/* err := fmt.Errorf("csv record not valid- %s", record)
			log.Printf(err.Error())
			return nil, err */
		}
		// unit,section,question,sub-question,marks,sub-question,marks,sub-question,marks,sub-question,marks
		dbRow := &DbRow{
			Unit:        record[0],
			Section:     record[1],
			ExternalRef: record[2],
		}

		switch dbRow.Section {
		case "A":
			dbRow.Marks = "2"
		case "B":
			dbRow.Marks = "13"
		case "C":
			dbRow.Marks = "15"
		}

		dbRow.Question.Marks = dbRow.Marks
		dbRow.Question.Text = record[3]

		dbRow.Question.SubQuestions = make([]DbRowQuestion, 0)

		if len(record[4]) != 0 {
			dbRow.Question.SubQuestions = append(dbRow.Question.SubQuestions, DbRowQuestion{
				Text:  record[4],
				Marks: record[5],
			})
		}

		if len(record[6]) != 0 {
			dbRow.Question.SubQuestions = append(dbRow.Question.SubQuestions, DbRowQuestion{
				Text:  record[6],
				Marks: record[7],
			})
		}

		if len(record[8]) != 0 {
			dbRow.Question.SubQuestions = append(dbRow.Question.SubQuestions, DbRowQuestion{
				Text:  record[8],
				Marks: record[9],
			})
		}

		if len(record[10]) != 0 {
			dbRow.Question.SubQuestions = append(dbRow.Question.SubQuestions, DbRowQuestion{
				Text:  record[10],
				Marks: record[11],
			})
		}

		dbRows = append(dbRows, dbRow)
	}
	return dbRows, nil
}

func parseCSVFile(file io.Reader) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func stringArraysEqual(arr1, arr2 []string) bool {
	fmt.Printf("\n len arr1: %v   len arr2: %v \n", len(arr1), len(arr2))
	if len(arr1) != len(arr2) {
		fmt.Println("MARKER within len not equal")
		return false
	}

	for i := range arr1 {
		fmt.Printf("\n arr1 elem: %x   arr2 elem: %x \n", removeBOM(arr1[i]), removeBOM(arr2[i]))
		if removeBOM(arr1[i]) != removeBOM(arr2[i]) {
			fmt.Println("MARKER elements not equal")
			return false
		}
	}

	return true
}

// Function to remove UTF-8 BOM from a string
func removeBOM(s string) string {
	if len(s) >= 3 && s[0] == '\xEF' && s[1] == '\xBB' && s[2] == '\xBF' {
		return s[3:]
	}
	return s
}
