package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

)

//structs

type Response struct { //e.g. "name":"Ho Yeung", "answer":["Tutorial", "red", "satisfied", "good"]
	Name string `json:"name"`
	Answer []string `json:"answer"` // store the answer of all questions
	Time string `json:"time"` // store the time of responses
}

type Question struct {
	Question string `json:"question"` //question title
	QuestionFormat string `json:"question_format"` // "Textbox" / "Multiple Choice" / "Likert Scale"
	Specification []string `json:"specification"`
}

type Survey struct {
	ID string `bson:"_id,omitempty"`
	Title string `json:"title"` //survey title
	Token string `json:"token"`
	Questions []Question `json:"questions"`
	//Responses []Response `json:"responses"` //not create the responses field first
}

type SurveyWithResponses struct { //for displaying surveys that hv responses
	ID string `bson:"_id,omitempty"`
	Title string `json:"title"` //survey title
	Token string `json:"token"`
	Questions []Question `json:"questions"`
	Responses []Response `json:"responses"`
}

type EditQuestion struct { //for editing surveys
	Question string `json:"question"`
	QuestionFormat string `json:"question_format"`
	Specification []string `json:"specification"`
}

// global variable for the mongodb collection
var surveysCollection *mongo.Collection

//generate random token when creating a new survey
func generateRandomToken() string {
	//array containing letters and numbers
	const combination = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const tokenLength = 5

	token := make([]byte, tokenLength)
	for i := range token {
		randIndex := rand.Intn(len(combination)) //make a random index
		token[i] = combination[randIndex] //choose it from the combination arr
	} //e.g. {"d","w","2","a","3"}

	//transform the byte array into a string arr and join all elements
	return strings.Join([]string{string(token)}, "") //e.g. dw2a3
}

//validate questions
func validateQuestions(questions []Question, c*gin.Context) error {
	//check no empty survey
	if len(questions)< 1 {
		c.JSON(400, gin.H{"error": "Cannot be an empty survey"})
		return fmt.Errorf("Cannot be an empty survey") //return an error, preventing return nil
	}		

	for i := range questions {
		//check question title
		if (len(questions[i].Question) < 3) {
			c.JSON(400, gin.H{"error": "The question title should have at least 3 characters"})
			return fmt.Errorf("The question title should have at least 3 characters")
		}
		//check if the question format is valid, question format should be "Textbox", "Multiple Choice" or "Likert Scale"
		if (questions[i].QuestionFormat != "Textbox" && questions[i].QuestionFormat != "Multiple Choice" && questions[i].QuestionFormat != "Likert Scale") {
			c.JSON(400, gin.H{"error": "Invalid question format"})
			return fmt.Errorf("Invalid question format")
		}
		if questions[i].QuestionFormat == "Textbox" {
			if (len(questions[i].Specification) > 0) {
				c.JSON(400, gin.H{"error": "Textbox format should not have specification"})
				return fmt.Errorf("Textbox format should not have specification")
			}
		} else if (questions[i].QuestionFormat == "Multiple Choice") {
			//check the array have at least 2 elements
			if len(questions[i].Specification) < 2 {
				c.JSON(400, gin.H{"error": "Multiple Choice question should have at least 2 options"})
				return fmt.Errorf("Multiple Choice question should have at least 2 options")
			}
		} else { //Likert Scale
			//check the array have at least 3 elements
			if len(questions[i].Specification) < 3 {
				c.JSON(400, gin.H{"error": "Likert Scale should have at least 3 options"})
				return fmt.Errorf("Likert Scale should have at least 3 options")
			}
		}
	}

	return nil //no error
}

//check token
func checkToken(token string, c*gin.Context) error {
	// check whether the token is valid i.e. length is 5 and not include any special characters
	if ((len(token) != 5) || (strings.ContainsAny(token, "!@#$%^&*()_+{}|:<>?~`-=[]\\;',./"))) {
		c.JSON(400, gin.H{"error": "Invalid token"})
		return fmt.Errorf("Invalid token")
	}
	return nil
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf(".env file cannot be loaded")
	}

	// connect to mongodb
	opt := options.Client().ApplyURI("mongodb://localhost:27017/?connect=direct").SetAuth(options.Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		AuthSource: "admin",
		AuthMechanism: "SCRAM-SHA-256",
	})

	// create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		log.Fatal(err)
	}

	surveysCollection = client.Database("OSP").Collection("Survey")

	// check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	router := gin.Default()

// create a new survey 
	router.POST("/surveys", func(c *gin.Context) {
		var survey Survey
		if err := c.BindJSON(&survey); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		//check survey's title
		if len(survey.Title) < 3 {
			c.JSON(400, gin.H{"error": "The survey title must have at least 3 characters"})
			return
		}

		//check token is not repeated, else gen a new one
		for (true) {
			survey.Token = generateRandomToken()

			// check whether the token already exists
			count, err := surveysCollection.CountDocuments(context.TODO(), bson.M{"token": survey.Token})
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to check token"})
				return
			}

			if count == 0 { // token is unique, stop iterating
				break
			}
		}

		//check all questions
		if (validateQuestions(survey.Questions, c)) != nil {
			return
		}

		// insert the survey into the database
		_, err := surveysCollection.InsertOne(context.TODO(), survey)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create survey"})
			return
		}
		
		c.JSON(200, gin.H{"message": fmt.Sprintf("Survey successfully created. The token of this survey is: %s", survey.Token)}) //show the survey token

	})

	// test for adding 1Q
	// curl -X POST http://localhost:8080/surveys \
	// -H "Content-Type: application/json" \
	// -d '{
	// 	"title": "Satisfaction Survey",
	// 	"questions": [
	// 		{
	// 			"question": "How satisfied are you with our service?",
	// 			"question_format": "Likert Scale",
	// 			"specification": ["Very Unsatisfied", "Unsatisfied", "Neutral", "Satisfied", "Very Satisfied"]
	// 		}
	// 	]
	// }' 

	// test for adding 3Q
	// curl -X POST http://localhost:8080/surveys \
	// -H "Content-Type: application/json" \
	// -d '{
	// 		"title": "Lecture Satisfaction Survey 2",
	// 		"questions": [
	// 				{
	// 						"question": "What do you think about the difficulty of the lecture material?",
	// 						"question_format": "Likert Scale",
	// 						"specification": ["Very Easy", "Easy", "Neutral", "Difficult", "Very Difficult"]
	// 				},
	// 				{
	// 						"question": "What do you think about my lecture style?",        
	// 						"question_format": "Multiple Choice",
	// 						"specification": ["Bad", "Good"]
	// 				},
	// 				{
	// 						"question": "Type a comment about the lecture",
	// 						"question_format": "Textbox",
	// 						"specification": []
	// 				}
	// 		]
	// }'

	// display a new survey using the input token
	router.GET("/surveys/:token", func(c *gin.Context) {
		token := c.Param("token") // Get the input token
		if checkToken(token,c) != nil {
			return
		}

		var survey Survey
		
		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey) // search survey
		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}
		// return the survey title & questions
		c.IndentedJSON(200, gin.H{"title": survey.Title, "questions": survey.Questions}) //IndentedJSON is for prettier output
	})

	//test: curl -X GET http://localhost:8080/surveys/aPdlq

// edit a survey in all fields
	router.PUT("/surveys/:token", func(c *gin.Context) {
		token := c.Param("token") // get the input token
		if checkToken(token,c) != nil {
			return
		}

		var survey Survey

		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey) // search survey
		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}

		if err := c.BindJSON(&survey); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		//check survey title
		if len(survey.Title) < 3 {
			c.JSON(400, gin.H{"error": "Title must have at least 3 characters"})
			return
		}	

		//check all questions
		if (validateQuestions(survey.Questions, c)) != nil {
			return
		}

		// update this survey in the database

		_, err = surveysCollection.UpdateOne(context.TODO(), 
			bson.M{"token": token}, 
			bson.M{"$set": bson.M{"title": survey.Title, "questions": survey.Questions}})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update survey"})
			return
		}
		c.JSON(200, gin.H{"message": "Survey successfully updated"})
		
	})

	//test --update token(aPdlq) change the title & change the specification
	// curl -X PUT http://localhost:8080/surveys/aPdlq \
	// -H "Content-Type: application/json" \
	// -d '{
	// 	"title": "CS111 Satisfaction Survey",
		// "questions": [
		// 	{
		// 		"question": "How satisfied are you with our service?",
		// 		"question_format": "Likert Scale",
		// 		"specification": ["Very Unsatisfied", "Unsatisfied", "Neutral", "Satisfied", "Very Satisfied", "No Comment"]
		// 	}
	// 	]
	// }'

//edit a question
	router.PUT("/surveys/:token/:questionNo", func(c *gin.Context) {
		token := c.Param("token") // get the input token
		if checkToken(token,c) != nil {
			return
		}

		var survey Survey

		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey) // search survey
		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}

		var editQuestion EditQuestion
		if err := c.BindJSON(&editQuestion); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		questionNo, err := strconv.Atoi(c.Param("questionNo")) // get the question number string -> int
		if err != nil || questionNo <= 0 || questionNo > len(survey.Questions) {
			c.JSON(400, gin.H{"error": "Invalid question number"})
			return
		}

		arrIndex := questionNo - 1 //e.g. Q1 is at Questions[0]

		survey.Questions[arrIndex].Question = editQuestion.Question
		survey.Questions[arrIndex].QuestionFormat = editQuestion.QuestionFormat
		survey.Questions[arrIndex].Specification = editQuestion.Specification

		//check all questions
		if (validateQuestions(survey.Questions, c)) != nil {
			return
		}
		// update in the database
		_, err = surveysCollection.UpdateOne(context.TODO(),
			bson.M{"token": token}, 
			bson.M{"$set": bson.M{"questions": survey.Questions}})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update a question"})
			return
		}
		c.JSON(200, gin.H{"message": "The survey question successfully updated"})
	})

	// test --update question 1 of aPdlq
	// curl -X PUT http://localhost:8080/surveys/aPdlq/1 \
	// -H "Content-Type: application/json" \
	// -d '{
	// 	"question": "How satisfied are you with our service?",
	// 	"question_format": "Multiple Choice",
	// 	"specification": ["Unsatisfied", "Neutral", "Satisfied"]
	// }'

//delete a survey
	router.DELETE("/surveys/:token", func(c *gin.Context) {
		token := c.Param("token") // get the input token
		if checkToken(token,c) != nil {
			return
		}

		var survey Survey

		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey) // search survey

		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}

		// Delete the survey from the database
		_, err = surveysCollection.DeleteOne(context.TODO(),bson.M{"token": token})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete survey"})
			return
		}
		c.JSON(200, gin.H{"message": "Survey successfully deleted"})
	})

	//curl -X DELETE http://localhost:8080/surveys/aPdlq

//delete a question
	router.DELETE("/surveys/:token/:questionNo", func(c *gin.Context) {
		token := c.Param("token") // get the input token
		if checkToken(token,c) != nil {
			return
		}

		var survey Survey
		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey) // search survey
		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}

		questionNo, err := strconv.Atoi(c.Param("questionNo")) // get the question number string -> int
		if err != nil || questionNo < 0 || questionNo > len(survey.Questions) {
			c.JSON(400, gin.H{"error": "Invalid question number"})
			return
		}

		arrIndex := questionNo - 1 //e.g. Q1 is at Questions[0]

		//delete Questions[arrIndex]
		_,err = surveysCollection.UpdateOne(context.TODO(), 
			bson.M{"token": token}, 
			bson.M{"$unset": bson.M{fmt.Sprintf("questions.%d", arrIndex): ""}}) // delete the question based on index			

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete the question"})
			return
		}

		//delete the null value in the array
		_, err = surveysCollection.UpdateOne(context.TODO(),
			bson.M{"token": token},
			bson.M{"$pull": bson.M{"questions": nil}})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete question"})
			return
		}
		c.JSON(200, gin.H{"message": "The question is successfully deleted"})
	})

	// test --delete question 2 of 1BZHv
	// curl -X DELETE http://localhost:8080/surveys/1BZHv/2

//submit a repsonse to a survey
	router.POST("/surveys/:token/responses", func(c *gin.Context) {
		token := c.Param("token") // get the input token
		if checkToken(token,c) != nil {
			return
		}

		var survey Survey

		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey) // search survey
		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}

		var response Response
		if err := c.BindJSON(&response); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		//get the current time
		response.Time = time.Now().Format("2006-01-02 15:04:05")

		//check response name
		if len(response.Name) < 3 {
			c.JSON(400, gin.H{"error": "Your Name must have at least 3 characters"})
			return
		}

		//check no empty response
		if len(response.Answer) < 1 {
			c.JSON(400, gin.H{"error": "Not allow empty response"})
			return
		}

		//check the no of answers == no of questions
		if len(response.Answer) != len(survey.Questions) {
			c.JSON(400, gin.H{"error": "Please answer the exact number of questions"})
			return
		}

		//check whether the answer algin with its specification?
		for i := range survey.Questions {
			if (survey.Questions[i].QuestionFormat == "Textbox" ) {
				if len(response.Answer[i]) < 3 {
					c.JSON(400, gin.H{"error": "Answer must have at least 3 characters"})
					return
				}
			} else { //MC or Likert Scale, check the answer whether is one of the options?

				for j := range survey.Questions[i].Specification {
					if ((survey.Questions[i].Specification[j]) == (response.Answer[i])) {
						break
					}
					if (j == len(survey.Questions[i].Specification)-1) { //if the answer is not in the options for last iteration
						currentQuestionNo := i + 1
						c.JSON(400, gin.H{"error": "Answer is not an option for question " + fmt.Sprint(currentQuestionNo)})
						return
					}
				}	
			} 
		}

		// appending the response to the survey
		_, err = surveysCollection.UpdateOne(context.TODO(), 
			bson.M{"token": token},
			bson.M{"$push": bson.M{"responses": response}})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to submit reponse"})
			return
		}
		c.JSON(200, gin.H{"message": "Reponse successfully submitted"})
	})

	// test: curl -X POST http://localhost:8080/surveys/1BZHv/responses \
	// -H "Content-Type: application/json" \
	// -d '{
	// 	"name": "WAN Ho Yeung",
	// 	"answer": ["Very Satisfied"]
	// }'

	// curl -X POST http://localhost:8080/surveys/1BZHv/responses \
	// -H "Content-Type: application/json" \
	// -d '{
	// 		"name": "Ho Yeung",
	// 		"answer": ["Very Easy","Good","can be faster"]
	// }'

//display responses of a survey
	router.GET("/surveys/:token/responses", func(c *gin.Context) {
		token := c.Param("token")
		if checkToken(token,c) != nil {
			return
		}

		var survey SurveyWithResponses

		err := surveysCollection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&survey)
		if err != nil {
			c.JSON(404, gin.H{"error": "Survey not found"})
			return
		}

		if (len(survey.Responses) < 1) {
			c.JSON(404, gin.H{"error": "No response"})
			return
		}

		c.IndentedJSON(200, gin.H{"response": survey.Responses})

	})

	//curl -X GET http://localhost:8080/surveys/1tRpL/responses

	router.Run() // run on port 8080

	// close the mongodb connection when the application stop
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to database is closed.")
	}()
}