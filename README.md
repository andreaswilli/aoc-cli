# Advent of Code CLI

A CLI to make solving [Advent of Code](https://adventofcode.com) puzzles more
convenient.

## Features

- **Run multiple days** with a single command - even if they are implemented in
  **different languages**.
- **Watch** the filesystem and re-run whenever changes are detected.
- **Verify results** against the content of a file so you immediately know if
  a refactoring broke something.

## Installation

Clone this repo and run `cd cmd/aoc && go install` from its root. This will
make the `aoc` command available in your path.

> ![Note]
> Go needs to be installed on your system for this to work.

## Project Structure

Your project should be structured something like this:

```text
advent-of-code/
├─ 2024/
│  ├─ 01/
│  │  ├─ expected.txt
│  │  ├─ solution.js
│  ├─ 02/
│  │  ├─ expected.txt
│  │  ├─ solution.py
│  ├─ ...
├─ aoc-cli.json
```

Year and day folders can be named differently, it just needs to match the
commands you run (see [Usage](#usage)).

The names of the files containing the source code need to be configured (see
[Config File](#config-file)).

The file `expected.txt` can be omitted if the output should not be checked.
Otherwise it should contain the expected output written to stdout by your code.

## Config File

The config file is in JSON format and defines engines to run your code. Find
some examples [here](cmd/aoc/2024/aoc-cli.json).

```json
{
  "engines": [
    {
      "name": "engine name (e.g. node)",
      "cmd": "command to run, use placeholder {{entryFile}}",
      "entryFile": "file name of the source file (e.g. solution.js)"
    },
    {
      ...
    }
  ]
}
```

## Usage

There are two main commands, `run` and `watch`.

### Run

Run once and exit.

#### Single Day

Run day `01` of year `2024`:

```sh
aoc run 2024/01
```

#### All Days of Year

Run all days of year `2024`:

```sh
aoc run 2024
```

### Watch

Watch re-runs a day every time one of the files in its folder changes.

#### Single Day

Watch day `01` of year `2024`:

```sh
aoc watch 2024/01
```

#### All Days of Year

Watch all days of year `2024`:

```sh
aoc watch 2024
```
