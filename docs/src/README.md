# Introduction

Commander solves one of the biggest headaches with using bash as the shell for the script block in a CI/CD environment.

YAML makes it easy to read and customize the pipeline and execution flow but gets messy once you're in the script 
block.
The current solution is to go fetch a Bash wizard if you want to customize how the command is run.

Commander takes advantage of the same templating language used in tools like Helm for Helm charts.
It's incredibly flexible and expressive which gives developers the ability to define commands in a more declarative 
fashion than the traditional script block.
It's also a bit more readable than bash, and you can test and validate your assumptions locally without having to run
a job and read through the pipeline log.

