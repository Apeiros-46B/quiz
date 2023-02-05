# PocketBase backend

## Setup

```bash
go build
./pocketbase migrate up
```

## Configuration

```bash
./pocketbase serve
```

- Visit `http://localhost:8090/_` and create an admin account
- Set up collection records:

<details>
  <summary><code>questions</code> collection</summary><br>

  A collection of questions in the quiz.  
  (Visible by anyone)

  - `index` (number): number used to determine order of questions in the quiz
    - higher = later
    - must start at 0 and increase
  - `question` (text): question text
  - `choices` (JSON): JSON array of choices (in text form)
</details>

<details>
  <summary><code>correct</code> collection</summary><br>

  A collection of correct answers, each corresponding to a question.  
  (Only visible by admins)

  - `question` (relation): the corresponding question
  - `correct` (number): the choice index which is correct
    - the choices array supplied in the `questions` collection is 0-indexed
    - 0 means the first choice in the choices array is correct
</details>

<details>
  <summary><code>settings</code> collection</summary><br>

  A collection of settings consisting of key-value pairs.  
  (Visible by anyone, except for the `successKW` setting)

  - `key` (text):
    <details>
      <summary>valid keys</summary>

      - `startHTML` (HTML): the HTML on the start screen (excluding the start button)
      - `startButtonTxt` (plain text): the text on the start button
      - `nextButtonTxt` (plain text): the text on the button to advance to the next question
      - `submitButtonTxt` (plain text): the text on the button to advance past the last question (submit answers)
      - `proceedConfirmTxt` (plain text): the text in the confirmation dialog shown when the user advances to the next question without answering the current one
      - `failureHTML` (HTML): the HTML on the failure screen (excluding the keyword)
      - `failureKW` (HTML): the HTML keyword on the failure screen
      - `successHTML` (HTML): the HTML on the success screen (excluding the keyword)
      - `successKW` (HTML): the HTML keyword on the success screen
        <details>
          <summary>available placeholders for embedding stats</summary>

          - `%userID%`: the user's ID (as reported in the `attempts` collection)
          - `%attempts%`: how many attempts this run took (previous unsuccessful attempts + this successful attempt) for the user
          - `%attemptTime%`: how much time, in milliseconds, the user took for this successful attempt
          - `%totalTime%`: how much time, in milliseconds, the user took over all attempts in this run
        </details>
    </details>
  - `value` (JSON): the value of the setting (either a JSON string or a JSON array of strings)
    - if the value is an array of strings, then they will be concatenated into one string (delimited by newline)

## Running

See [main README](../README.md).
