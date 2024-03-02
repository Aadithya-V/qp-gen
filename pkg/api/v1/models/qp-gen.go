package models

type ExamDetails struct {
	Name         string   `json:"exam_name"`
	Code         string   `json:"exam_code"`
	Date         string   `json:"exam_date"`
	Duration     string   `json:"exam_duration"`
	StartTime    string   `json:"exam_start_time"`
	EndTime      string   `json:"exam_end_time"`
	Session      string   `json:"exam_session"`
	TotalMarks   string   `json:"exam_total_marks"`
	Instructions []string `json:"exam_instructions"`
}

type QuestionsByType struct {
	TypeNumber     int     `json:"type_number"`
	SectionName    string  `json:"section_name"`
	TotalQuestions int     `json:"total_questions"`
	Marks          float32 `json:"marks"`
	ChoiceAllowed  bool    `json:"choice_allowed"`

	Questions []string `json:"questions"`
}

type GenerateQpaperSetsInLatexRequest struct {
	ExamDetails *ExamDetails `json:"exam_details"`

	NumberOfSets int8 `json:"number_of_sets"`

	QuestionsByType []*QuestionsByType `json:"questions_by_type"`
}
