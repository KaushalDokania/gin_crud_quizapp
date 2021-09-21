package main

import (

	// Why do we need this package?

	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // If you want to use mysql or any other db, replace this line
)

var db *gorm.DB // declaring the db globally
var err error

type Quiz struct {
	ID          uint   `json:"id";gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Questions   []Question
}

type Question struct {
	ID   uint   `json:"id";gorm:"primary_key"`
	Name string `json:"name"`
	// Question string `json:"question"`
	QuizID  uint `json:"quizid"`
	Points  int  `json:points`
	Options []Option
}

type Option struct {
	ID         uint   `json:"optid";gorm:"primary_key"`
	QuestionID uint   `json:"qid"`
	OptionName string `json:"optname"`
	Correct    bool   `json:"iscorrect"`
}

func main() {
	os.Remove("gorm.db") // delete the file to avoid duplicated records.
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Option{}, &Question{}, &Quiz{})

	// routes
	r := gin.Default()
	r.GET("/api/quiz/", GetQuizzes)
	r.GET("/api/quiz/:quiz_id", GetQuiz)
	r.POST("/api/quiz/", CreateQuiz)

	r.GET("/api/questions", GetQuestions)
	r.GET("/api/questions/:question_id", GetQuestion)
	r.POST("/api/questions", CreateQuestion)
	r.GET("/api/quiz-questions/:quiz_id", GetQuizQuestions)

	quiz1 := Quiz{Name: "quiz1", Description: "Sample quiz 1", Questions: []Question{
		{Name: "ques1", QuizID: 1, Points: 10},
		{Name: "ques2", QuizID: 1, Points: 20},
	}}
	quiz2 := Quiz{Name: "quiz2", Description: "Sample quiz 2", Questions: []Question{
		{Name: "ques3", QuizID: 2, Points: 5},
		{Name: "ques4", QuizID: 2, Points: 15},
	}}
	db.Create(&quiz1)
	db.Create(&quiz2)

	// ques1 := Question{Name: "ques1", QuizID: 1, Points: 10}
	// append(quiz1, ques1)
	log.Println("inserted dummy data")

	r.Use((cors.Default()))
	r.Run(":8080") // Run on port 8080
}

func GetQuizzes(c *gin.Context) {
	var quizzes []Quiz
	// var genres []Genre
	if err := db.Find(&quizzes).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.Header("access-control-allow-origin", "*")
		c.JSON(200, quizzes)
	}
}

func GetQuiz(c *gin.Context) {
	quizId := c.Params.ByName("quiz_id")
	fmt.Println(quizId)
	var quiz Quiz
	var questions []Question

	if err := db.Where("id = ?", quizId).First(&quiz).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		log.Println("inside else")
		db.Model(&quiz).Related(&questions)
		for i, q := range questions {
			var options []Option
			db.Model(&q).Related(&options)
			questions[i].Options = options
		}
		c.Header("access-control-allow-origin", "*")
		c.JSON(200, quiz)
	}
}

func GetQuestions(c *gin.Context) {
	var questions []Question
	// var genres []Genre
	if err := db.Find(&questions).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.Header("access-control-allow-origin", "*")
		c.JSON(200, questions)
	}
}

func GetQuestion(c *gin.Context) {
	questionId := c.Params.ByName("question_id")
	fmt.Println(questionId)
	var question Question
	var options []Option

	if err := db.Where("id = ?", questionId).First(&question).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		log.Println("inside else")
		db.Model(&question).Related(&options)
		c.Header("access-control-allow-origin", "*")
		c.JSON(200, question)
	}
}

func CreateQuiz(c *gin.Context) {
	var quiz Quiz
	c.BindJSON(&quiz)
	db.Create(&quiz)
	c.Header("access-control-allow-origin", "*")
	c.JSON(200, quiz)
}

func CreateQuestion(c *gin.Context) {
	quizId := c.Params.ByName("quiz_id")
	var quiz Quiz
	if err := db.Where("id = ?", quizId).First(&quiz).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		var ques Question
		c.BindJSON(&ques)
		db.Create(&ques)
		db.Model(&ques).UpdateColumn("quiz_id", quiz.ID)
		c.Header("access-control-allow-origin", "*")
		c.JSON(200, ques)
	}
}

func GetQuizQuestions(c *gin.Context) {
	quizId := c.Params.ByName("quiz_id")
	fmt.Println(quizId)
	var quiz Quiz
	var questions []Question

	if err := db.Where("id = ?", quizId).First(&quiz).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		log.Println("inside else")
		db.Model(&quiz).Related(&questions)
		log.Println(quiz)
		for i, q := range questions {
			var options []Option
			db.Model(&q).Related(&options)
			questions[i].Options = options
		}
		c.Header("access-control-allow-origin", "*")
		c.JSON(200, quiz)
	}
}
