package function

import "testing"

func TestReadForm(t *testing.T) {
	form, err := ReadForm("test_form.yml")
	if err != nil {
		t.Fatalf("ReadForm(test_form.yml) = %v", err)
	}
	if len(form.Questions) != 2 {
		t.Fatalf("len(form.Questions) = %v", len(form.Questions))
	}

	if form.Questions[0].QuestionText != "Сколько тебе лет?" {
		t.Fatalf("form.Questions[0].QuestionText = %v", form.Questions[0].QuestionText)
	}

	if form.Questions[1].QuestionText != "В какой стране ты живешь?" {
		t.Fatalf("form.Questions[1].QuestionText = %v", form.Questions[1].QuestionText)
	}

}

func TestNewForm(t *testing.T) {
	form, err := NewForm("test_form.yml", 12345)
	if err != nil {
		t.Fatalf("NewForm(test_form.yaml, 12345) = %v", err)
	}

	if form.CurrentQuestionIndex != 0 {
		t.Fatalf("form.CurrentQuestionIndex = %v", form.CurrentQuestionIndex)
	}
	if form.UserID != 12345 {
		t.Fatalf("form.UserID = %v", form.UserID)
	}
	if form.IsFilled != false {
		t.Fatalf("form.IsFilled = %v", form.IsFilled)
	}

}
