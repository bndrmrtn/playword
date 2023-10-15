# PlayWord API with Go Fiber

This is my first project with Go and Fiber. It's an API for a simple app.
It contains 2 endpoints, a `GET` and a `POST`.

| Method | Resource  |
|--------|-----------|
| `GET`  | /api/game |
| `POST` | /api/game |

Every endpoint returns a `JSON` response.

### `GET` endpoint

This endpoint returns a `JSON` response like this:

```json
{
  "data": {
    "letters":["w","r","p","e","o"],
    "length":5,
    "max_trials":3,
    "trials":0,
    "guessed":0,
    "expiration":"2023-10-15T14:01:55.2868383+02:00"
  },
  "status":{
    "code":200,
    "message":"OK"
  }
}
```

So, it contains a `status` for the status, and `data`
for the game data, it contains the `letters` in an array
and the letters are mixed. The `length` is the length of the word.
The `max_trials` is the maximum number of trials for the user to guess the word.
The `trials` is the current trials, and the `guessed` is how many words the user guessed.
An `expiration` is the expiration time to guess the word (currently not implemented).

### `POST` endpoint

For the `POST` endpoint you have to send a
`JSON` request body with the following parameter:
```json
{
  "guess": "power"
}
```

Then the server will return a `JSON` response.

#### When the word is successfully guessed

```json
{
    "valid": true,
    "trials": 1,
    "endgame": false,
    "msg": {
        "type": "success",
        "text": "Congratulation you\nguessed the word! :D"
    }
}
```

#### Otherwise

```json
{
    "valid": false,
    "trials": 2,
    "endgame": false,
    "msg": {
        "type": "error",
        "text": "No problem! You have 1 more trial"
    }
}
```

### The projects structure

```tree
📦 Project
    📂 database
    📂 handlers
        📂 types
    📂 helpers
    📂 http
        📂 session
    📂 models
    📂 routes
    📜 main.go
```