package grades

func init() {
	students = []Student{
		{
			ID:        1,
			FirstName: "Mini",
			LastName:  "Mo",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 75,
				},
				{
					Title: "Homework 1",
					Type:  GradeHomework,
					Score: 79,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 88,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Tiny",
			LastName:  "So",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 68,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 73,
				},
				{
					Title: "Homework 1",
					Type:  GradeHomework,
					Score: 75,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 80,
				},
			},
		},
	}
}
