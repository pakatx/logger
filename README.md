# logger

GoLang logger package intended for back-end / server side usage.

## Usage

- `Init()` should be invoked to initialize the logger.
- Logging can be done with function logger.Log().

## Developer Notes

- `LogDispatcher()` is to be invoked as a part of waitgroup.
- A dispatcher go routine fetches from the channel, extracts the log message, and dumps the same in the log file.
- Logfile name: server.log.<no>, where "no" stands for logfile number.
- Current logfile has extension .1.
- Max allowed size of a logfile and rollover count can be configured.