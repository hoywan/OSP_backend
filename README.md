# OSP_backend

OSP_backend is a backend for an Online Survey Platform web application. It was developed using Golang and MongoDB.

## Installation

### Install GO

1.	Click link of the GO official website <a target="_blank" href="https://go.dev">https://go.dev</a> and download a compatible version based on your OS 

2.	Install it

### Install MongoDB

1.	Download a community version of MongoDB by following the instructions of <a target="_blank" href="https://www.mongodb.com/try/download/community">https://www.mongodb.com/try/download/community</a>
2.	Install the complete version with the MongoDB compass
3.	Click "+" icon

<img width="314" alt="image" src="https://github.com/user-attachments/assets/6d62fdf9-4950-400a-adb6-9d5bbded5007" />

4.	Enter "OSP" in the "Database Name" field and "Survey" in the "Collection Name"
5.	Click the "Create Database" button
6.	Click "Open MongoDB shell", type the following commands

```terminal
use admin
```
"switched to db admin" can be seen

You can use any user and pwd instead of the below
```terminal
db.createUser(
{
    user: "Andrew",
    pwd: "a1sdcvd23@",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
})
```

"{ok: 1}" can be seen

### Install Visual Studio Code, all dependencies and run the application

1.	Click link of the VS code official website <a target="_blank" href="https://code.visualstudio.com/download">https://code.visualstudio.com/download</a> and download a compatible version based on your OS
2.	Install it
3.	Donwload this repository and unzip it,
4.	Click "Open Folder" in VS code to open the folder, which should contains .env and main.go e.g. OSP_backend-main
5.	Change the "DB_USERNAME" and "DB_PASSWORD" in the .env file if you change the user and pwd in step 6 in "Install MongoDB"
6.	Click "Ctrl", "Shift" and "`" at the same time to open a new terminal
7.	Fix go module

```terminal
go env -w GO111MODULE=on
```
```terminal
go mod init OSP_backend-main{replace here with your folder name}
```
8.	Install dependencies by entering the commands below one by one with/without sudo

```teminal
sudo go get github.com/joho/godotenv
sudo go get go.mongodb.org/mongo-driver/mongo
sudo go get go.mongodb.org/mongo-driver/bson
sudo go get go.mongodb.org/mongo-driver/mongo/options
sudo go get github.com/gin-gonic/gin

```

9.	Enter "go run ." in the terminal, "Listening and serving HTTP on :8080" should be shown on the last output, which means the server is running
10.	For MacOS user, you can type curl command in the terminal app by following the doumentation below
For Windows user, you can install an extension on <a target="_blank" href="https://chromewebstore.google.com/detail/reqbin-http-client/gmmkjpcadciiokjpikmkkmapphbmdjok">https://chromewebstore.google.com/detail/reqbin-http-client/gmmkjpcadciiokjpikmkkmapphbmdjok</a>
and enter curl command in <a target="_blank" href="https://reqbin.com/curl">https://reqbin.com/curl</a>
11.	You can check view all data inside OSP>Survey in the MongoDB compass after you create one survey.

### Solution to MongoDB problem (can't be connected, ECONNREFUSED 127.0.0.1:27017 in Compass)
Enter the below in terminal to restart the service

```terminal
brew services restart mongodb-community@8.0
```

## API documentation

### Survey

### (1) Creating a survey
| POST | /surveys |
|------|----------|

Parameters: No parameters

Request body:

Schema:

|    Name   |  Type  |       Description      |
|:---------:|:------:|:----------------------:|
|   title   | string |      survey title      |
| questions |  array | contains all questions |

|    Name   |  Type  |       Description      |
|:---------:|:------:|:----------------------:|
|   question   | string |      question title      |
| question_format |  string | "Textbox" / "Multiple Choice" / "Likert Scale" |
| specification |  array | empty for "Textbox", > 2 elements for "Multiple Choice", > 3 elements for "Likert Scale" |


Example (creating a survey with 3 questions with different question formats):
```terminal
curl -X POST http://localhost:8080/surveys \
-H "Content-Type: application/json" \
-d '{
    "title": "Lecture Satisfaction Survey 2",
    "questions": [
        {
            "question": "What do you think about the difficulty of the lecture material?",
            "question_format": "Likert Scale",
            "specification": ["Very Easy", "Easy", "Neutral", "Difficult", "Very Difficult"]
        },
        {
            "question": "What do you think about my lecture style?",        
            "question_format": "Multiple Choice",
            "specification": ["Bad", "Good"]
        },
        {
            "question": "Type a comment about the lecture",
            "question_format": "Textbox",
            "specification": []
        }
    ]
}'
```

Responses:
| 200 | Survey successfully created. |
|------|----------|

Example:
{"message":"Survey successfully created. The token of this survey is: 5GXbe"}

The token is random generated, so it has less chance to be "5GXbe" again. You should use the one generated when trying other APIs.

|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid input"}        |      Binding input error      |
|      |     {"error":"The survey title must have at least 3 characters"}        | survey title less than 3 characters |
|      |     {"error":"Cannot be an empty survey"}        | no questions |
|      |       {"error":"The question title should have at least 3 characters"}      | question title less than 3 characters |
|      |      {"error":"Invalid question format"}       | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |     {"error":"Textbox format should not have specification"}        | specification for "Textbox" is not empty  |
|      |     {"error":"Multiple Choice question should have at least 2 options"}        | specification for "Multiple Choice" has less than 2 elements |
|      |      {"error":"Likert Scale should have at least 3 options"}       | specification for "Likert Scale" has less than 3 elements |

| 500 | Internal Server Error |
|------|----------|
|      |Failed to check token|
|      |Failed to create survey|

Example:
{"error":"Failed to check token"}

{"error":"Failed to create survey"}

### (2) Displaying a survey
| GET | /surveys/:token |
|------|----------|

Parameters: token(string)

Example

```terminal
curl -X GET http://localhost:8080/surveys/5GXbe
```

Responses:
| 200 | Output a JSON containing survey title and questions array |
|------|----------|

Example:
```terminal
{
    "questions": [
        {
            "question": "What do you think about the difficulty of the lecture material?",
            "question_format": "Likert Scale",
            "specification": [
                "Very Easy",
                "Easy",
                "Neutral",
                "Difficult",
                "Very Difficult"
            ]
        },
        {
            "question": "What do you think about my lecture style?",
            "question_format": "Multiple Choice",
            "specification": [
                "Bad",
                "Good"
            ]
        },
        {
            "question": "Type a comment about the lecture",
            "question_format": "Textbox",
            "specification": []
        }
    ],
    "title": "Lecture Satisfaction Survey 2"
}
```

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             |      Invalid token e.g. not equal to 5 characters or containing any special characters    |

{"error":"Invalid token"}

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

### (3) Editing a survey in all fields

| PUT | /surveys/:token |
|------|----------|

Parameters: token(string)

Request body:

Schema:

|    Name   |  Type  |       Description      |
|:---------:|:------:|:----------------------:|
|   title   | string |      survey title      |
| questions |  array | contains all questions |

Example (Change from 3 questions to 2 new question, and update the survey title):
```terminal
curl -X PUT http://localhost:8080/surveys/5GXbe \
-H "Content-Type: application/json" \
-d '{
    "title": "Satisfaction Survey 3",
		"questions": [
			{
				"question": "How satisfied are you with our service?",
				"question_format": "Likert Scale",
				"specification": ["Very Unsatisfied", "Unsatisfied", "Neutral", "Satisfied", "Very Satisfied", "No Comment"]
			},
			{
				"question": "What is the color of the sky?",
				"question_format": "Multiple Choice",
				"specification": ["Red", "Blue", "Green"]
			}

    ]
}'
```

Responses:
| 200 | Survey successfully updated |
|------|----------|

{"message":"Survey successfully updated"}

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |    {"error":"Invalid input"}         |      Binding input error      |
|      |       {"error":"The survey title must have at least 3 characters"}      | survey title less than 3 characters |
|      |      {"error":"Cannot be an empty survey"}       | no questions |
|      | {"error":"The question title should have at least 3 characters"} | question title less than 3 characters |
|      |      {"error":"Invalid question format"}       | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |     {"error":"Textbox format should not have specification"}        | specification for "Textbox" is not empty  |
|      |      {"error":"Multiple Choice question should have at least 2 options"}       | specification for "Multiple Choice" has less than 2 elements |
|      |       {"error":"Likert Scale should have at least 3 options"}      | specification for "Likert Scale" has less than 3 elements |

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|      |Failed to update survey|

{"error":"Failed to update survey"}

### (4) Deleting a survey

| DELETE | /surveys/:token |
|------|----------|

Parameters: token(string)

Example

```terminal
curl -X DELETE http://localhost:8080/surveys/5GXbe
```

repsonse:
| 200 | Survey successfully deleted |
|------|----------|

{"message":"Survey successfully deleted"}

|    400   |  Bad Request  |       Description      |
|      |             | Invalid token e.g. not equal to 5 characters or containing any special characters    |

{"error":"Invalid token"}

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|      |Failed to delete survey|

{"error":"Failed to delete survey"}

### Question

### (1) Editing a question

| PUT | /surveys/:token/:questionNo |
|------|----------|

Parameters: token(string), questionNo(string)[-> Integer later]

Request body:

Schema:

|    Name   |  Type  |       Description      |
|:---------:|:------:|:----------------------:|
|   question   | string |      question title      |
| question_format |  string | "Textbox" / "Multiple Choice" / "Likert Scale" |
| specification |  array | empty for "Textbox", > 2 elements for "Multiple Choice", > 3 elements for "Likert Scale" |

Example (change the question 1 of survey with token 5GXbe):

```terminal
curl -X PUT http://localhost:8080/surveys/5GXbe/1 \
-H "Content-Type: application/json" \
-d '{
	"question": "How satisfied are you with our service?",
	"question_format": "Multiple Choice",
	"specification": ["Unsatisfied", "Neutral", "Satisfied"]
}'
```

Responses:
| 200 | The survey question successfully updated |
|------|----------|

{"message":"The survey question successfully updated"

|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |     {"error":"Invalid input"}        |      Binding input error      |
|      |     {"error":"Invalid question number"}       |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |
|      |     {"error":"The survey title must have at least 3 characters"}        | survey title less than 3 characters |
|      |      {"error":"Cannot be an empty survey"}       | no questions |
|      |     {"error":"The question title should have at least 3 characters"}        | question title less than 3 characters |
|      |     {"error":"Invalid question format"}        | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |     {"error":"Textbox format should not have specification"}        | specification for "Textbox" is not empty  |
|      |      {"error":"Multiple Choice question should have at least 2 options"}       | specification for "Multiple Choice" has less than 2 elements |
|      |      {"error":"Likert Scale should have at least 3 options"}       | specification for "Likert Scale" has less than 3 elements |

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|      |Failed to update a question|

{"error":"Failed to update a question"}

### (2) Deleting a question

| DELETE | /surveys/:token/:questionNo |
|------|----------|

Parameters: token(string), questionNo(string)[-> Integer later]

Example (delete the question 2 of survey with token 5GXbe):

```terminal
curl -X DELETE http://localhost:8080/surveys/5GXbe/2
```

Response:
| 200 | The survey question successfully updated |
|------|----------|

{"message":"The question is successfully deleted"}

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |             |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |

{"error":"Invalid token"}

{"error":"Invalid question number"}

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|   {"error":"Failed to delete the question"}   |Failed to delete the question|
|   {"error":"Failed to delete null value"}    |Failed to delete the null value|
|   {"error":"Failed to update the modified time"}   |Failed to update the modified time|

### Response

### (1) Submitting a response

| POST | /surveys/:token/responses |
|------|----------|

Parameters: token(string)

Request body:

Schema:

|    Name   |  Type  |       Description      |
|:---------:|:------:|:----------------------:|
|   name   | string |      response name      |
| answer |  array | contains all answer to each question in the corresponding survey |

Example (submit one response to the survey with token 5GXbe):

```terminal
curl -X POST http://localhost:8080/surveys/5GXbe/responses \
-H "Content-Type: application/json" \
-d '{
	"name": "WAN Ho Yeung",
	"answer": ["Satisfied","Blue"]
}'
```

Response:

| 200 | Reponse successfully submitted |
|------|----------|

{"message":"Reponse successfully submitted"}

<img width="506" alt="image" src="https://github.com/user-attachments/assets/4d4456e3-7020-4f3a-8827-ff03eaa13877" />

|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |     {"error":"Invalid input"}        |      Binding input error      |
|      |     {"error":"Invalid question number"}       |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |
|      |     {"error":"Your Name should be 2 to 100 characters"}        |      response name should be 2 to 100 char     |
|      |      {"error":"Not allow empty response"}       |     response not contain any answers     |
|      |{"error":"Please answer the exact number of questions"}|      the no. of elements in the answer array not match to the total no. of survey question     |
|      |     {"error":"Answer should be 1 to 300 characters"}        |      Answer for Textbox question < 1 or > 300 characters      |
|      |      {"error":""Answer is not an option for question X"}       |      Answer is not the option that included in the specifcation of questions for MC & Likert Scale    |

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|      |Failed to submit reponse|

{"error":"Failed to submit reponse"}

### (2) Displaying all responses in a survey

| POST | /surveys/:token/responses |
|------|----------|

Parameters: token(string)

Example (Display all responses of the survey with token 5GXbe):

```terminal
curl -X GET http://localhost:8080/surveys/5GXbe/responses
```

Response:

| 200 | Output a JSON containing response array |
|------|----------|

```terminal
{
    "response": [
        {
            "name": "WAN Ho Yeung",
            "answer": [
                "Satisfied",
                "Blue"
            ],
            "time": "2025-04-28 00:52:10"
        },
        {
            "name": "Andrew",
            "answer": [
                "Unsatisfied",
                "Red"
            ],
            "time": "2025-04-28 01:05:57"
        }
    ]
}
```
|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             | Invalid token e.g. not equal to 5 characters or containing any special characters    |

{"error":"Invalid token"}

| 404 | Survey not found with the input token |
|------|----------|
|      |No response e.g. the response array does not exists or empty|


{"error":"Survey not found"}

{"error":"No response"}
