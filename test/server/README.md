Test Server for Humble
----------------------
### About

This is the backend server for testing humble. In particular, it is used to
test code in the model package which communiciates with a REST API.

The server is written in go and runs on port 3000. It accepts Content-Types of
application/json, application/x-www-form-urlencoded, or multipart/form-data and
responds with JSON. It does validations and returns 422 errors when validations
fail.

This is a test server specifically designed for testing the humble framework.
As such, it is designed to be completely idempotent. That means nothing you do will
actually change the data on the server, and sending the same request will always
give you the same response. However, when possible the responses are designed to mimic
that of a real server that does hold state.


### Getting Up and Running

- Install dependencies with `go get`
- Run the server with `go run main.go`

### Endpoints

#### GET /todos

List all existing todos. Since this server is idempotent, the response is always exactly the same.

**Parameters**: none

**Example Responses**:

Success:

```json
[
  {
    "Id": 0,
    "Title": "Todo 0",
    "IsCompleted": false
  },
  {
    "Id": 1,
    "Title": "Todo 1",
    "IsCompleted": false
  },
  {
    "Id": 2,
    "Title": "Todo 2",
    "IsCompleted": true
  }
]
```

#### POST /todos

Simulate creation of a new todo item. Since this server is idempotent, the state never changes
and the list of todos always stays the same. However, the server will respond exactly as if the
todo were created, and even will assign it an id. The id assigned is always 3. The rest of the
response can vary, as the Title and IsCompleted field of the response will match the form data
that was submitted for the todo. The server will also validate the submitted data by requiring
the Title and IsCompleted fields. If either is not provided, it returns a validation error.

**Parameters**:

| Field       | Type    | Description     |
| ----------- | ------- | --------------- |
| Title       | string  | The title of the new todo. |
| IsCompleted | bool    | Whether or not the todo is completed. |


**Example Responses**:

Success:

```json
{
  "Id": 3,
  "Title": "New Test Todo",
  "IsCompleted": false
}
```

Validation error:

```json
{
  "Title": [
    "Title is required."
  ]
}
```

#### PUT /todos/{id}

Simulate editing an existing todo item. Since this server is idempotent, the state never changes
and the list of todos always stays the same. However, the server will respond exactly as if the
todo were updated, and will respond with json data representing the updated todo. The server will
return an error if the id is not an integer between 0 and 2.

**Parameters**:

| Field       | Type    | Description     |
| ----------- | ------- | --------------- |
| Title       | string  | The title of the todo. |
| IsCompleted | bool    | Whether or not the todo has been completed. |

**Example Responses**:

Success:

```json
{
  "Id": 3,
  "Title": "Updated Title",
  "IsCompleted": true
}
```

Error:

```json
{
  "error": "Could not find todo with id = 5"
}
```

#### DELETE /todos/{id}

Simulate deletion of an existing todo item. Since this server is idempotent, the state never changes
and the list of todos always stays the same. However, the server will respond exactly as if the
todo were deleted. A successful response is always just an empty json object. The server will
return an error if the id is not an integer between 0 and 2.

**Parameters**: none

**Example Responses**:

Success:

```json
{}
```
