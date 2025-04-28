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

## API

<img width="690" alt="image" src="https://github.com/user-attachments/assets/ceaed8bb-ac51-4d1f-b083-1a2e67822d24" />
<img width="690" alt="image" src="https://github.com/user-attachments/assets/2cd51cda-fd18-4ad9-832c-682b6719748d" />
<img width="690" alt="image" src="https://github.com/user-attachments/assets/83d3d7ec-d93b-45e1-b984-f70ed18e9a3a" />


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

Explaination:
1.	combine the request body with an empty survey object
2.	use checkTitle function to ensure the title not already exists and between 2 to 300 characters
3.	use a for loop to generate a new tokens using generateRandomToken function ensure the token not duplicated, count=0 means no same token can be searched, stop until 100 iterations if still can't find an unique token
4.	check all questions using validateQuestions function which the question title,question format and specification will be checked
5.	set the time of the survey created and LastModifiedTime to be the present time
6.	set the responses array to be empty first
7.	insert the survey into the database
8.	Ouput successfully created and provide the token of the survey to the user

Responses:
| 200 | Survey successfully created. |
|------|----------|

Example:
{"message":"Survey successfully created. The token of this survey is: 5GXbe"}

The token is random generated, so it has very little chance to be "5GXbe" again. You should use the one generated when trying other APIs.

Responses

<img width="662" alt="image" src="https://github.com/user-attachments/assets/94d06121-662d-4f22-b26c-113546155e7a" />


|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid input"}        |      Binding input error      |
|      |     {"error":"The Survey Title already exists"}        |     The Survey Title already exists in DB     |
|      |     {"error":"The Survey title should be 2 to 300 characters"}        | The Survey title should be 2 to 300 characters |
|      |     {"error":"Cannot be an empty survey"}        | no questions |
|      |       {"error":"The question title should have at least 3 characters"}      | question title less than 3 characters |
|      |      {"error":"Invalid question format"}       | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |     {"error":"Textbox format should not have specification"}        | specification for "Textbox" is not empty  |
|      |     {"error":"Multiple Choice question should have at least 2 options"}        | specification for "Multiple Choice" has less than 2 elements |
|      |      {"error":"Likert Scale should have at least 3 options"}       | specification for "Likert Scale" has less than 3 elements |

| 500 | Internal Server Error |
|------|----------|
|   {"error":"Failed to check token"}   |Failed to check token|
|  {"error":"Failed to create survey"}    |Failed to create survey|
|  {"error":"Failed to generate a unique token"}    |Failed to generate a unique token|
|  {"error":"Failed to check the survey title"}    |Failed to check title|

### (2) Displaying a survey
| GET | /surveys/:token |
|------|----------|

Parameters: token(string)

Example

```terminal
curl -X GET http://localhost:8080/surveys/5GXbe
```

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	output the below if no error

Responses:
| 200 | Output a JSON containing survey title, questions array, time created, LastModifiedTime |
|------|----------|

Example:
```terminal
{
    "lastModifiedTime": "2025-04-28 16:04:53",
    "number_of_questions": 3,
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
    "time": "2025-04-28 16:04:53",
    "title": "Lecture Satisfaction Survey 7"
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

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	bind the request body with survey object
4.	perform checking on tiles and questions using checkTitle function and validateQuestions function respectively
5.	update the LastModifiedTime to present time
6.	only update Title, Questions, LastModifiedTime, Responses in the survey with that token in the database

Responses:
| 200 | Survey successfully updated |
|------|----------|

{"message":"Survey successfully updated"}

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |    {"error":"Invalid input"}         |      Binding input error      |
|      |     {"error":"The Survey Title already exists"}        |     The Survey Title already exists in DB     |
|      |     {"error":"The Survey title should be 2 to 300 characters"}        | The Survey title should be 2 to 300 characters |
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
|  {"error":"Failed to update survey"}    |Failed to update survey|
|  {"error":"Failed to check the survey title"}    |Failed to check the survey title|

### (4) Deleting a survey

| DELETE | /surveys/:token |
|------|----------|

Parameters: token(string)

Example

```terminal
curl -X DELETE http://localhost:8080/surveys/5GXbe
```

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	delete the survey in the database

repsonse:
| 200 | Survey successfully deleted |
|------|----------|

{"message":"Survey successfully deleted"}

|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|  {"error":"Failed to delete survey"}    |Failed to delete survey|

### Question

### (1) Inserting a question

| POST | /surveys/:token/:questionNo |
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
curl -X POST http://localhost:8080/surveys/5GXbe/2 \
-H "Content-Type: application/json" \
-d '{
	"question": "Do you want to know the topic in depth?",
	"question_format": "Likert Scale",
	"specification": ["Definitely No", "No", "No commnet", "Yes", "Definitely"]
}'
```

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	check the questionNo parameter should >= 1, question no. should not exceed total no. of questinos + 1(i.e. insert as the last question)
4.	bind the request body with newQuestion object
5.	perform checking on questions using validateQuestions function, need to change it into an array with only ine element for passing into the function
6.	append the question based on different situations and change the response of this question to be "No answer" for the responses array, the detailed, please see the comment in the code. By adding "No answer" can prevent affecting the correspondence between the questions array and response answer array e.g. {questions[0] should correspond to all answer[0],questions[1] should correspond to all answer[1] etc}

<img width="608" alt="image" src="https://github.com/user-attachments/assets/81e59f19-7be9-40d2-adb2-240d5ef9d61f" />
<img width="268" alt="image" src="https://github.com/user-attachments/assets/ee89cfa6-0c1f-48ef-bb83-861e34aed5bd" />

8.	update the LastModifiedTime to present time
9.	only update Questions, LastModifiedTime, Responses in the survey with that token in the database

Responses:
| 200 | The question is successfully inserted with a preview of all questins |
|------|----------|

```terminal
{
    "All Questions": [
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
            "question": "Do you want to know the topic in depth?",
            "question_format": "Likert Scale",
            "specification": [
                "Definitely No",
                "No",
                "No commnet",
                "Yes",
                "Definitely"
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
    "message": "The question is successfully inserted"
}
```

|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |     {"error":"Invalid input"}        |      Binding input error      |
|      |     {"error":"Invalid question number"}       |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |
|      |     {"error":"The Survey title should be 2 to 300 characters"}        | The Survey title should be 2 to 300 characters |
|      |     {"error":"The Survey Title already exists"}        | The Survey Title already exists in DB |
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
|   {"error":"Failed to insert a question"}   |Failed to insert a question|

### (2) Editing a question

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

Example (Edit the question 2 of survey with token 5GXbe):

```terminal
curl -X PUT http://localhost:8080/surveys/5GXbe/2 \
-H "Content-Type: application/json" \
-d '{
	"question": "How satisfied are you with our service?",
	"question_format": "Multiple Choice",
	"specification": ["Unsatisfied", "Neutral", "Satisfied"]
}'
```

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	check the questionNo parameter should >= 1, question no. should not exceed total no. of questinos + 1(i.e. insert as the last question)
4.	bind the request body with EditQuestion object
5.	assign the Question,QuestionFormat and Specification into the specific question with index (questionNo - 1)
6.	perform validation in Questions array using validateQuestions function
7.	replace all existed response of this question with "Deleted", preventing unmatched responses
8.	update the LastModifiedTime to present time
9.	only update Questions, LastModifiedTime, Responses in the survey with that token in the database

Responses:
| 200 | The survey question successfully updated |
|------|----------|

{"message":"The survey question successfully updated"}

|    400   |  Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |     {"error":"Invalid input"}        |      Binding input error      |
|      |     {"error":"Invalid question number"}       |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |
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

### (3) Deleting a question

| DELETE | /surveys/:token/:questionNo |
|------|----------|

Parameters: token(string), questionNo(string)[-> Integer later]

Example (delete the question 2 of survey with token 5GXbe):

```terminal
curl -X DELETE http://localhost:8080/surveys/5GXbe/2
```

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	check the questionNo parameter should >= 1, question no. should not exceed total no. of questinos + 1(i.e. insert as the last question)
5.	delete the corresponding question with index (questionNo - 1) in the survey with the token parameter
6.	update the LastModifiedTime to present time
7.	if the question is deleted, the answer of this question should also be deleted,and the rest of the answers should be shifted to the left or just delete all last elements in each answer if the question is also the last
9.	only update lastModifiedTime and Responses in the survey with that token in the database

Response:
| 200 | The survey question successfully updated |
|------|----------|

{"message":"The question is successfully deleted"}

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |     {"error":"Invalid token"}        | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |       {"error":"Invalid question number"}      |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |

| 404 | Survey not found with the input token |
|------|----------|

{"error":"Survey not found"}

| 500 | Internal Server Error |
|------|----------|
|   {"error":"Failed to delete the question"}   |Failed to delete the question|
|   {"error":"Failed to delete null value"}    |Failed to delete the null value|
|   {"error":"Failed to update the modified time/responses"}   |Failed to update the modified time/responses|

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
	"answer": ["Easy","Bad","I love lecture"]
}'
```

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	bind the request body with empty response object
4.	use the current time as response.Time
5.	check no empty response, the number of answers is equal to the number of questions and align with specification
6.	append the response to the survey with that token in the database

Response:

| 200 | Reponse successfully submitted |
|------|----------|

{"message":"Reponse successfully submitted"}

<img width="662" alt="image" src="https://github.com/user-attachments/assets/034fa2fc-b425-4f02-939d-ec4ad1256349" />

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

| POST | /surveys/:token/responses/:displayMode |
|------|----------|

Parameters: token(string),displayMode(string individual/overview)

Explaination:
1.	use checkToken function to check the token parameter, ensuring it has exactly 5 characters and not include any special characters
2.	find the survey using the token as a filter and decode it into the empty survey object
3.	check if there any responses, output "No response" if no
4.	check the displayMode parameter should only be either individual or overview

Example (Display the overview of all responses of the survey with token 5GXbe, hide names, show statistics):

```terminal
curl -X GET http://localhost:8080/surveys/5GXbe/responses/overview
```
Explaination:

1.	the questions array stores all Q&A pairs, the answer will be accepted form the response array when the corresponding index in the Questions array is matched
2.	the counts[answer] store all count for each answer like counts[blue] = 2 if 2 of answer blue are stored in the response array
3.	calculate the percentage of single answer in all answers

Response:

| 200 | Output a JSON containing title, number of responses, questions array storing (question + all answers for that question) |
|------|----------|

```terminal
{
    "title": "Lecture Satisfaction Survey 7",
    "number_of_responses": 3,
    "questions": [
        {
            "question": "What do you think about the difficulty of the lecture material?",
            "answer": [
                "Easy (1, 33.33%)",
                "Difficult (1, 33.33%)",
                "Very Easy (1, 33.33%)"
            ]
        },
        {
            "question": "What do you think about my lecture style?",
            "answer": [
                "Bad (3, 100.00%)"
            ]
        },
        {
            "question": "Type a comment about the lecture",
            "answer": [
                "I love lecture (1, 33.33%)",
                "Can speak faster (2, 66.67%)"
            ]
        }
    ]
}
```

Example (Display individual responses of the survey with token 5GXbe):

```terminal
curl -X GET http://localhost:8080/surveys/5GXbe/responses/individual
```

Explaination:

1.	the qa, Q&A pairs stores question from the questions array and answer in answer array, like a matching
2.	every response in responses array of ouput will contain the name, qa, time of response submission

Response:

| 200 | Output a JSON containing survey title and responses array storing name, Q&A pairs, time |
|------|----------|

```terminal
{
    "title": "Lecture Satisfaction Survey 7",
    "responses": [
        {
            "name": "WAN Ho Yeung",
            "qa": [
                {
                    "question": "What do you think about the difficulty of the lecture material?",
                    "answer": "Easy"
                },
                {
                    "question": "What do you think about my lecture style?",
                    "answer": "Bad"
                },
                {
                    "question": "Type a comment about the lecture",
                    "answer": "I love lecture"
                }
            ],
            "time": "2025-04-28 16:16:32"
        },
        {
            "name": "Andrew",
            "qa": [
                {
                    "question": "What do you think about the difficulty of the lecture material?",
                    "answer": "Difficult"
                },
                {
                    "question": "What do you think about my lecture style?",
                    "answer": "Bad"
                },
                {
                    "question": "Type a comment about the lecture",
                    "answer": "Can speak faster"
                }
            ],
            "time": "2025-04-28 16:17:24"
        },
        {
            "name": "Candy",
            "qa": [
                {
                    "question": "What do you think about the difficulty of the lecture material?",
                    "answer": "Very Easy"
                },
                {
                    "question": "What do you think about my lecture style?",
                    "answer": "Bad"
                },
                {
                    "question": "Type a comment about the lecture",
                    "answer": "Can speak faster"
                }
            ],
            "time": "2025-04-28 16:19:09"
        }
    ]
}
```
|    400   | Output  |       Description      |
|:---------:|:------:|:----------------------:|
|      |      {"error":"Invalid token"}       | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |      {"error":"Invalid display mode"}       | the input display mode is nither individual nor overview    |



| 404 | Description |
|------|----------|
|   {"error":"Survey not found"}   |Survey not found with the input token |
|   {"error":"No response"}   |No response e.g. the response array does not exists or empty|
