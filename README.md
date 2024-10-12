# wanikani-progress-plotter-tool

Simple tool to fetch Level progression data, and then plot a progress chart

## Memos

- We use the default `slog` logger. We don't create a customer slog logger instance.
- a secret `.envrc` file is used to store the Wanikani API Key, which is then handled with [direnv](https://direnv.net).
- Output is written to a fixed location
- I won't use full dependency injection, interfaces for clients etc, as this is just a quick and dirty tool.
