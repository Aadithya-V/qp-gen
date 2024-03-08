package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"text/template"

	"github.com/Aadithya-V/qp-gen/database"
	"github.com/Aadithya-V/qp-gen/pkg/api/v1/models"
	"github.com/gin-gonic/gin"
)

func GenerateQpaperSetsFromDB(c *gin.Context, req *models.GenerateQpaperSetsFromDBRequest) error {

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
		return errors.New("unknown exam type")
	}

	questions, err := GetQuestionFromDB(req, units)
	if err != nil {
		return err
	}

	qStore := BuildQStore(questions, qsPerType)

	tmpl, err := template.ParseFiles("vec_template.tex")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

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

		c.Writer.WriteString("\n%-------------Question Paper Set Start--------------- \n")

		err = tmpl.Execute(c.Writer, tmplData)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return err
		}

		c.Writer.WriteString("\n%----------------Question Paper Set End ---------------- \n")

	}

	return nil
}

func GetQuestionFromDB(req *models.GenerateQpaperSetsFromDBRequest, units []string) ([]*Question, error) {

	inClause := "'" + strings.Join(units, "','") + "'"
	// sub_code = %s and
	query := fmt.Sprintf(`select unit, section, question from question_bank_data where unit in (%s);`, inClause)

	db := database.NewMySQLSession()

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

	for k, _ := range qsPerType {
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
