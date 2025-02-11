package quiz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetQuiz(book string, chapter int) (Quiz, error) {
	data, err := getQuizJson(book, chapter)
	if err != nil {
		return Quiz{}, err
	}
	var quiz Quiz

	err = json.Unmarshal([]byte(data), &quiz)
	if err != nil {
		return Quiz{}, err
	}

	return quiz, nil
}

func getQuizJson(book string, chapter int) (string, error) {
	filename := filepath.Join(
		"data", "quizzes", book,
		fmt.Sprintf("%s-%d.json", strings.ToLower(book), chapter),
	)

	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(file), nil
}
