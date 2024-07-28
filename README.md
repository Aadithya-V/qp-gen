Active Branch - feat/v2


http://localhost:8001/qp-gen/api/docs/index.html#/  - swagger

curl -X POST http://localhost:8001/qp-gen/api/v1/upload/654af/adf5   -F "file=@csv_example.csv"   -H "Content-Type: multipart/form-data"


{
  "exam_details": {
    "academic_year": "2023-2024",
    "college_name": "SRM Valliammai Engineering College",
    "department_name": "CSE",
    "exam_date": "12-12-12",
    "exam_duration": "3 hrs",
    "exam_end_time": "12:00",
    "exam_instructions": [
      ""
    ],
    "exam_session": "FN",
    "exam_start_time": "9:00",
    "exam_total_marks": "100",
    "examination_name": "CAT 2",
    "semester": "VII",
    "subject_code": "CS-101",
    "subject_name": "Automata Theory"
  },
  "number_of_sets": 3,
  "q_paper_codes": [
    ""
  ],
  "questions_by_type": [
    {
      "choice_allowed": false,
      "marks": 2,
      "questions": [
         "Define Green computing",
        "What are the 3Rs of Green IT?",
        "Distinguish between EI and BI.",
        "Define carbon foot print.",
        "What is ERBS?",
"Classify the challenges of carbon economy.",
"Illustrate the concepts of Business intelligence.",
"Describe the impact of BI to EI.",
"Give elements of an ERBS forming the Green Strategies Mix.",
"What are the steps in developing an ERBS?",
"Generalize about Green organizational goals to be achieved through policy development.",
"Evaluate about Lean Impact on Green Computing.",
"Define green sustainable policy.",
"Interpret the need for green computing.",
"Categorize the green IT drivers?",
"Classify the 5 M’s of Carbon metrics",
"Predict the types of carbon emissions under scope?",
"Define Green Assets.",
"Illustrate the type of assets.",
"Discover the idea of Green Data Centers.",
"Define Carbon Emitting Bit.",
"Analyze the factors influencing Green data center.",
"Summarize the list of Green Process Categories.",
"Interpret the factors of Green BPM.",
"Distinguish between coupling and cohesion.",
"When to use Patterns?",
"Analyze Green Enterprise Architecture."

      ],
      "section_name": "",
      "total_questions": 5,
      "type_number": 1
    }
  ]
}



{{if $question.Marks}}
\part[{{$question.Marks}}] 
{{else}}
\part[] 
{{end}}
{{ $question.Text }} 
{{if $question.Marks}}
\droppoints
{{end}}
