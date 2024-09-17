# Kojin

- [Kojin](#kojin)
  - [Description](#description)
  - [How does it work](#how-does-it-work)
  - [Why I built Kojin](#why-i-built-kojin)
  - [Example of usage (with Gleam)](#example-of-usage-with-gleam)
    - [command](#command)
    - [command structure](#command-structure)

## Description

File tree watcher built using Golang

With Kojin you can watch a file tree for changes and execute a command when a change is detected

## How does it work

Given the folder to watch, Kojin builds a tree structure of the entire folder with its files and subfolders and then, on each iteration which can happen every chosen interval, it compares the two structures.
If a change is detected, the command specified will run

## Why I built Kojin

When trying out a new programming language, I usually start by solving some Leetcodes or Advent of code problems to get a grip of the basic syntax.

Then, when trying out [Gleam](https://gleam.run/), I couldn't stand having to run the same commmand every time I made a change, so I decided to build a general purpose tool that I could use to have a **_live reload_**

## Example of usage (with Gleam)

### command

```bash
./kojin ./src "gleam run" --frequency 200 # run every 200 milliseconds and ignore the build folder/s
```

### command structure

```bash
./kojin [folder_to_watch] [command_to_run] ...[additional flags] # use ./kojin --help for more info
```
