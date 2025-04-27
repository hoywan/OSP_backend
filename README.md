# OSP_backend

OSP_backend is a backend for an Online Survey Platform web application. It was developed using Golang and MongoDB.

## Installation

Use the package manager [pip](https://pip.pypa.io/en/stable/) to install foobar.

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

### (3) Editing a whole survey

| PUT | /surveys/:token |
|------|----------|

Parameters: token(string)


Example (Change from 3 questions to 1 new question, and update the survey title):
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
			}
    ]
}'
```
