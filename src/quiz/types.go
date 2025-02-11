package quiz

type Quiz struct {
	Book            string          `json:"book"`
	Chapter         int             `json:"chapter"`
	QuizStructure   QuizStructure   `json:"quiz_structure"`
	KeyThemes       []string        `json:"key_themes"`
	KeyVerses       []KeyVerse      `json:"key_verses"`
	Characters      []Character     `json:"characters"`
	Events          []Event         `json:"events"`
	Concepts        []string        `json:"concepts"`
	CrossReferences CrossReferences `json:"cross_references"`
}

type QuizStructure struct {
	TotalQuestions int       `json:"total_questions"`
	Sections       []Section `json:"sections"`
}

type Section struct {
	Name          string     `json:"name"`
	Count         int        `json:"count"`
	QuestionTypes []string   `json:"question_types"`
	Questions     []Question `json:"questions"`
}

type Question struct {
	Question QuestionDetails `json:"question"`
}

type QuestionDetails struct {
	Text          string      `json:"text"`
	Options       interface{} `json:"options"`
	CorrectAnswer *int        `json:"correct_answer,omitempty"`
}

type KeyVerse struct {
	Reference string `json:"reference"`
}

type Character struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Event struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}

type CrossReferences struct {
	OldTestament []Reference `json:"Old Testament"`
	NewTestament []Reference `json:"New Testament"`
}

type Reference struct {
	Reference   string `json:"reference"`
	Description string `json:"description"`
}
