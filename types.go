package function

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Question struct {
	QuestionText string `yaml:"question_text"`
	Answer       string `yaml:"-"`
}

type Form struct {
	ID                   string     `yaml:"-"`
	UserID               int64      `yaml:"-"`
	Questions            []Question `yaml:"questions"`
	CurrentQuestionIndex int        `yaml:"-"`
	IsFilled             bool       `yaml:"-"`
}

func ReadForm(path string) (*Form, error) {
	var file_data []byte
	var err error
	if _, ok := os.LookupEnv("DEPLOYED"); ok {
		file_data, err = os.ReadFile(
			filepath.Join("./serverless_function_source_code/", path),
		)
	} else {
		file_data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}

	var form Form
	err = yaml.Unmarshal(file_data, &form)

	return &form, err
}

func NewForm(name string, user_id int64) (*Form, error) {
	form, err := ReadForm(name)
	if err != nil {
		return nil, err
	}

	form.CurrentQuestionIndex = 0
	form.UserID = user_id
	form.IsFilled = false

	return form, nil
}

func (form *Form) GetCurrentQuestion() *Question {
	return &form.Questions[form.CurrentQuestionIndex]
}

func (form *Form) IsLastQuestion() bool {
	return form.CurrentQuestionIndex == len(form.Questions)-1
}
