package services

import (
	"fmt"
	"math/rand"
	"text/template"

	"github.com/Aadithya-V/qp-gen/pkg/api/v1/models"
	"github.com/gin-gonic/gin"
)

func GenerateQpaperSetsInLatex(c *gin.Context, req *models.GenerateQpaperSetsInLatexRequest) error {

	// Read template from file
	tmpl, err := template.ParseFiles("latex_template.tex")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	QS := BuildQStore(req.QuestionsByType)

	/* resp := GenerateQpaperSetsInLatexResponse{
		QSets: make([][]byte, req.NumberOfSets),
	} */

	for i := req.NumberOfSets; i > 0; i-- {
		templData := TemplateData{
			ExamDetails:     *req.ExamDetails,
			QuestionsByType: QS.PickQSet(),
		}

		c.Writer.WriteString("\n-------------Question Paper Set Start--------------- \n")

		err = tmpl.Execute(c.Writer, templData.QuestionsByType[1])
		if err != nil {
			fmt.Println("Error executing template:", err)
			return err
		}

		c.Writer.WriteString("\n----------------Question Paper Set End ---------------- \n")

	}

	return nil
}

type GenerateQpaperSetsInLatexResponse struct {
	QSets [][]byte `json:"question_paper_sets_array"`
}

type TemplateData struct {
	ExamDetails     models.ExamDetails
	QuestionsByType map[int][]string
}

type Q struct {
	Type int
	Id   int
	Text string
	//Qq          *Qq
	PickedCount int
}

type Qq struct {
	Marks float32
	Text  string // can be empty
	Qq    []*Qq  // sub-divisions
}

type QStore struct {
	ByTypes   map[int][]*Q
	QsPerType map[int]int
}

func BuildQStore(QsByType []*models.QuestionsByType) QStore {

	ret := QStore{
		ByTypes:   make(map[int][]*Q),
		QsPerType: make(map[int]int),
	}

	for _, qsType := range QsByType {
		fmt.Println(qsType)
		if _, ok := ret.ByTypes[qsType.TypeNumber]; !ok {
			ret.ByTypes[qsType.TypeNumber] = make([]*Q, 0)
		}

		for _, qStr := range qsType.Questions {
			var q = Q{
				Type:        qsType.TypeNumber,
				Text:        qStr,
				PickedCount: 0,
				// Id
			}
			ret.ByTypes[qsType.TypeNumber] = append(ret.ByTypes[qsType.TypeNumber], &q)
		}
		ret.QsPerType[qsType.TypeNumber] = qsType.TotalQuestions
	}

	return ret
}

// number of qs per type in opt
func (QS *QStore) PickQSet() map[int][]string {

	ret := make(map[int][]string)

	for qType, qs := range QS.ByTypes {
		if _, ok := ret[qType]; !ok { // if not requrired here
			ret[qType] = make([]string, 0)
		}

		ret[qType] = append(ret[qType], pickQs(qs, QS.QsPerType[qType])...)
	}

	return ret
}

func pickQs(qs []*Q, nums int) []string {

	ret := make([]string, 0)

	r := rand.Intn(len(qs) - 1)

	for i := r + 1; i < len(qs); i = (i + 1) % len(qs) {
		if nums == 0 {
			return ret
		}
		if i == r {
			i = rand.Intn(len(qs))
			q := qs[i]
			ret = append(ret, q.Text)
			q.PickedCount++
			nums--

			continue
		}

		if qs[i].PickedCount == 0 {
			ret = append(ret, qs[i].Text)
			qs[i].PickedCount++
			nums--
		}
	}
	return ret
}

/* var qs = []Q{
	Q{
		Qq: &Qq{
			Text:  "Sample text 1",
			Marks: 2,
		},
		Id:   1,
		Type: 1,
	},
	Q{
		Qq: &Qq{
			Text:  "Sample text 2",
			Marks: 2,
		},
		Id:   2,
		Type: 1,
	},
	Q{
		Qq: &Qq{
			Text:  "Sample text",
			Marks: 2,
		},
		Id:   2,
		Type: 1,
	},
} */
