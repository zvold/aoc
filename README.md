# Solutions for [Advent of Code](https://adventofcode.com/) puzzles

> [!NOTE]
> This repository doesn't include the "real" puzzle inputs, which AoC generates
> for each user individually.

> [!NOTE]
> For some days, the code doesn't output a complete answer, because I solved it
> on a piece of paper (examples are 2023-06, 2023-24, 2024-24) or with some
> external tooling (example is 2023-25) and didn't bother converting it to code.

## 2023

The code is written in Go, and (provided you have Go runtime installed) can be
run with:

```
# Run on a small, embedded, test input.
go run github.com/zvold/aoc/2023/go/day01@latest
```

To run the code on your own puzzle input, use the `--input` command-line flag:

```
go run github.com/zvold/aoc/2023/go/day01@latest --input=<input_file>
```

## 2024

Similar to 2023, the solutions are written in Go:

```
# Run on a small example input.
go run github.com/zvold/aoc/2024/go/day01@latest

# Run on your own input.
go run github.com/zvold/aoc/2024/go/day01@latest --input=<input_file>
```
