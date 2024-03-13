package models

type TriviaBundle struct {
	ID         int      `json:"id,omitempty"`
	Question   string   `json:"question"`
	Category   string   `json:"category"`
	Answers    []Answer `json:"answers"`
	ShowAnswer bool     `json:"show_answer,omitempty"`
}

type Answer struct {
	ID             int    `json:"id,omitempty"`
	TriviaBundleID int    `json:"trivia_bundle_id,omitempty"`
	AnswerText     string `json:"answer_text"`
	IsCorrect      bool   `json:"is_correct"`
}

func (t *TriviaBundle) GetID() int {
	return t.ID
}

func (t *TriviaBundle) GetQuestion() string {
	return t.Question
}

func (t *TriviaBundle) GetCategory() string {
	return t.Category
}

func (t *TriviaBundle) GetAnswers() []Answer {
	return t.Answers
}

func (t *TriviaBundle) GetShowAnswer() bool {
	return t.ShowAnswer
}

func (a *Answer) GetID() int {
	return a.ID
}

func (a *Answer) GetTriviaBundleID() int {
	return a.TriviaBundleID
}

func (a *Answer) GetAnswerText() string {
	return a.AnswerText
}

func (a *Answer) GetIsCorrect() bool {
	return a.IsCorrect
}

func (a *Answer) SetID(id int) {
	a.ID = id
}

func (a *Answer) SetTriviaBundleID(triviaBundleID int) {
	a.TriviaBundleID = triviaBundleID
}

func (a *Answer) SetAnswerText(answerText string) {
	a.AnswerText = answerText
}

func (a *Answer) SetIsCorrect(isCorrect bool) {
	a.IsCorrect = isCorrect
}

func (t *TriviaBundle) SetID(id int) {
	t.ID = id
}

func (t *TriviaBundle) SetQuestion(question string) {
	t.Question = question
}

func (t *TriviaBundle) SetCategory(category string) {
	t.Category = category
}

func (t *TriviaBundle) SetAnswers(answers []Answer) {
	t.Answers = answers
}

func (t *TriviaBundle) SetShowAnswer(showAnswer bool) {
	t.ShowAnswer = showAnswer
}

func NewTriviaBundle(id int, question string, category string, answers []Answer, showAnswer bool, nextBundles []int, prevBundles []int, gameOver bool) *TriviaBundle {
	return &TriviaBundle{
		ID:         id,
		Question:   question,
		Category:   category,
		Answers:    answers,
		ShowAnswer: showAnswer,
	}
}

func NewAnswer(id int, triviaBundleID int, answerText string, isCorrect bool) *Answer {
	return &Answer{
		ID:             id,
		TriviaBundleID: triviaBundleID,
		AnswerText:     answerText,
		IsCorrect:      isCorrect,
	}
}
