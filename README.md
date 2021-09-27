# Tracker

This repository contains `tracker`, a simple tool to keep track of scheduled backups.

## Clone the project

```
$ git clone https://github.com/10Pines/tracker
$ cd tracker
$ make build
```

## [Rationale](rationale/)
Tracker follows Cloudwatch's evaluation approach:

When you create a task, you specify three settings to enable `tracker` to evaluate when to change the task state:
* **Name** is a symbolic string that describes the task  
* **Datapoints** is the number of data points that will conform the evaluation period for the task. The data points will be consecutive and treated as 1 per day.
* **Tolerance** is the number of days without acks or `jobs` from a given task that before triggering an error state.

| System got | Datapoints | Tolerance | Status       | Why?                                                                                            |
|------------|------------|-----------|--------------|-------------------------------------------------------------------------------------------------|
| X X X X X  | 5          | 0         | OK           | Over the last 5 days, the task tracked 5 jobs (x),5 (5 jobs) >= 5 (5 datapoints - 0 tolerance)  |
| _ X X X _  | 5          | 2         | OK           | 3 jobs >= 5 datapoints - 2 tolerance                                                            |
| X X X      | 5          | 0         | INSUFFICIENT | Task needs 5 datapoints to determine status. We'll have to wait 2 more days                     |
| X X X _ _  | 5          | 1         | ERROR        | This assertion failed: 3 jobs >= 5 datapoints - 1 tolerance                                     |

## [pkg](gopkg/)

The `tracker` package exposes a customizable tracker client that wraps the rest API.
or manipulate Go programs.

```
t := tracker.New(apiKey)
err := t.TrackJob(taskID)
```


## [CLI](cli/)

Exposes a simple CLI to track jobs. Refer to [releases](https://github.com/10Pines/tracker/releases).
`API_KEY` is required to be defined as env var.

```
> API_KEY='...' tracker track $TASK_ID
```