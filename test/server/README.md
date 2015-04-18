Test Server for Humble
----------------------
### About

This is the backend server for testing humble. In particular, it is used to
test code in the model package which communiciates with a REST API.

The server is written in go and runs on port 3000. It accepts Content-Types of
application/json, application/x-www-form-urlencoded, or multipart/form-data and
responds with JSON. It does validations and returns 422 errors when validations
fail.

This backend simply stores todos in memory so you don't have to worry about
setting up a database. This means there is no persistence, so if you restart
the server all your todos will be gone. It's perfect for testing and building
new things, but not recommended for use in production.


### Getting Up and Running

- Install dependencies with `go get`
- Run the server with `go run main.go`

### Endpoints

#### GET /todos

List all existing todos, ordered by time of creation.

**Parameters**: none

**Example Responses**:

Success:

```json
[
  {
    "id": 0,
    "title": "Write a frontend framework in Go",
    "isCompleted": false
  },
  {
    "id": 1,
    "title": "???",
    "isCompleted": false
  },
  {
    "id": 2,
    "title": "Profit!",
    "isCompleted": false
  }
]
```

#### POST /todos

Create a new todo item.

**Parameters**:

| Field    | Type    | Description     |
| ---------| ------- | --------------- |
| title    | string  | The title of the new todo. |


**Example Responses**:

Success:

```json
{
  "id": 3,
  "title": "Take out the trash",
  "isCompleted": false
}
```

Validation error:

```json
{
  "title": [
    "title is required."
  ]
}
```

#### PUT /todos/{id}

Edit an existing todo item.

**Parameters**:

| Field       | Type    | Description     |
| ----------- | ------- | --------------- |
| title       | string  | The title of the todo. |
| isCompleted | bool    | Whether or not the todo has been completed. |

**Example Responses**:

Success:

```json
{
  "id": 3,
  "title": "Handle the garbage",
  "isCompleted": false
}
```

#### DELETE /todos/{id}

Delete an existing todo item.

**Parameters**: none

**Example Responses**:

Success:

```json
{}
```
