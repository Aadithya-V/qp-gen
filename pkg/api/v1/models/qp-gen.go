package models

type ExamDetails struct {
	College                string   `json:"college_name"`
	Department             string   `json:"department_name"`
	DepartmentAbbreviation string   `json:"department_abbreviation"`
	ExamName               string   `json:"examination_name"`
	Regulation             string   `json:"regulation"`
	CourseYear             string   `json:"course_year"`
	SubjectName            string   `json:"subject_name"`
	SubjectCode            string   `json:"subject_code"`
	CommonTo               string   `json:"common_to"`
	Semester               string   `json:"semester"`
	Year                   string   `json:"year"`
	AcademicYear           string   `json:"academic_year"`
	Date                   string   `json:"exam_date"`
	Duration               string   `json:"exam_duration"`
	StartTime              string   `json:"exam_start_time"`
	EndTime                string   `json:"exam_end_time"`
	Session                string   `json:"exam_session"`
	TotalMarks             string   `json:"exam_total_marks"`
	Instructions           []string `json:"exam_instructions"`
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

	NumberOfSets int8     `json:"number_of_sets"`
	QpaperCodes  []string `json:"q_paper_codes"`

	QuestionsByType []*QuestionsByType `json:"questions_by_type"`
}

type GenerateQpaperSetsFromDBRequest struct {
	ExamDetails ExamDetails `json:"exam_details"`

	SubjectCode string `json:"subject_code"`

	ExamType string `json:"exam_type"`

	NumberOfSets int8     `json:"number_of_sets"`
	QpaperCodes  []string `json:"q_paper_codes"`
}
