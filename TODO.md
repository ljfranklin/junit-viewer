## TODO

- summary of last X runs (human-readable)
- summary of last X runs filtered by metadata (e.g. azure)
- summary of last X runs filtered by stdout/stderr text (e.g. "502")
- top X most frequent failing tests over last X runs
- graphs?

## Summary

- unmarshal all given junit files
- merge results
- print given summary

## UX

junit-viewer --input-files results/*.xml --output-type pass-fail

```
## Last 10 runs

| TESTS |    PASS     |    FAIL    | TIME  |           WHEN            |
|-------|-------------|------------|-------|-------------------------- |
|    30 | 30 (100.0%) | 0 (0.0%)   | 9.837 | 2006-01-02T15:04:05Z07:00 |
|     8 | 0 (0.0%)    | 8 (100.0%) | 0.003 | 2006-01-02T15:04:05Z07:00 |
...
```

junit-viewer --input-files results/*.xml --output-type frequent-failures

```
## Most frequent failures

| FAILED  |   TEST    |         LAST RAN          |
|---------|-----------|-------------------------- |
| 5 (50%) | TestS3Get | 2006-01-02T15:04:05Z07:00 |
| 1 (10%) | TestS3Put | 2006-01-02T15:04:05Z07:00 |
...
```

junit-viewer --input-files results/*.xml --output-type failure-regex --output-args patterns.txt

patterns.txt
```
\s502\s
InsufficientResources
```

```
## Failures matching regex

| MATCHES |         REGEX         |       LAST FAILURE        |
|---------|-----------------------|---------------------------|
|      10 | \s502\s               | 2006-01-02T15:04:05Z07:00 |
|       3 | InsufficientResources | 2006-01-02T15:04:05Z07:00 |
...
```
