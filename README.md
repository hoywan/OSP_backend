# OSP_backend

OSP_backend is a backend for an Online Survey Platform web application. It was developed using Golang and MongoDB.

## Installation

### Install GO

### Install MongoDB

### Install Visual Studio Code

### Install all dependencies

```bash
pip install foobar
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

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             |      Binding input error      |
|      |             | survey title less than 3 characters |
|      |             | no questions |
|      |             | question title less than 3 characters |
|      |             | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |             | specification for "Textbox" is not empty  |
|      |             | specification for "Multiple Choice" has less than 2 elements |
|      |             | specification for "Likert Scale" has less than 3 elements |

Example:
{"error":"Invalid input"}

{"error":"The survey title must have at least 3 characters"}

{"error":"Cannot be an empty survey"}

{"error":"The question title should have at least 3 characters"}

{"error":"Invalid question format"}

{"error":"Textbox format should not have specification"}

{"error":"Multiple Choice question should have at least 2 options"}

{"error":"Likert Scale should have at least 3 options"}

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

{"message":"Survey successfully updated"

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |             |      Binding input error      |
|      |             | survey title less than 3 characters |
|      |             | no questions |
|      |             | question title less than 3 characters |
|      |             | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |             | specification for "Textbox" is not empty  |
|      |             | specification for "Multiple Choice" has less than 2 elements |
|      |             | specification for "Likert Scale" has less than 3 elements |

Example:
{"error":"Invalid token"}

{"error":"Invalid input"}

{"error":"The survey title must have at least 3 characters"}

{"error":"Cannot be an empty survey"}

{"error":"The question title should have at least 3 characters"}

{"error":"Invalid question format"}

{"error":"Textbox format should not have specification"}

{"error":"Multiple Choice question should have at least 2 options"}

{"error":"Likert Scale should have at least 3 options"}

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

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |             |      Binding input error      |
|      |             |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |
|      |             | survey title less than 3 characters |
|      |             | no questions |
|      |             | question title less than 3 characters |
|      |             | not either "Textbox" / "Multiple Choice" / "Likert Scale" |
|      |             | specification for "Textbox" is not empty  |
|      |             | specification for "Multiple Choice" has less than 2 elements |
|      |             | specification for "Likert Scale" has less than 3 elements |

{"error":"Invalid token"}

{"error":"Invalid input"}

{"error":"Invalid question number"}

{"error":"The survey title must have at least 3 characters"}

{"error":"Cannot be an empty survey"}

{"error":"The question title should have at least 3 characters"}

{"error":"Invalid question format"}

{"error":"Textbox format should not have specification"}

{"error":"Multiple Choice question should have at least 2 options"}

{"error":"Likert Scale should have at least 3 options"}

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
|      |Failed to delete the question|

{"error":"Failed to delete the question"}

### Response

### (1) Submitting a response

| POST | /surveys/:token/responses |
|------|----------|

Parameters: token(string)

Request body:

Schema:

|    Name   |  Type  |       Description      |
|:---------:|:------:|:----------------------:|
|   name   | string |      survey title      |
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

|    400   |  Bad Request  |       Description      |
|:---------:|:------:|:----------------------:|
|      |             | Invalid token e.g. not equal to 5 characters or containing any special characters    |
|      |             |      Binding input error      |
|      |             |      Invalid question number e.g. less or equal than 0 or exceed the total no. of questions     |
|      |             |      Name less than 3 characters     |
|      |             |     response not contain any answers     |
|      |             |      the no. of elements in the answer array not match to the total no. of survey question     |
|      |             |      Answer for Textbox question < 3 characters     |
|      |             |      Answer is not the option that included in the specifcation of survey questions for MC & Likert Scale    |


{"error":"Invalid token"}

{"error":"Invalid input"}

{"error":"Invalid question number"}

{"error":"Your Name must have at least 3 characters"}

{"error":"Not allow empty response"}

{"error":"Please answer the exact number of questions"}

{"error":"Answer must have at least 3 characters"}

{"error":""Answer is not an option for question 2"}

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
