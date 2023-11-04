# Commander

[![Go Reference](https://pkg.go.dev/badge/github.com/BacchusJackson/commander.svg)](https://pkg.go.dev/github.com/BacchusJackson/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/BacchusJackson/commander)](https://goreportcard.com/report/github.com/BacchusJackson/commander)


![Commander Logo](assets/commander-logo-dark-theme-splash.png)


Unlike the "system" library call from C and other languages, the os/exec package intentionally does not invoke the 
system shell and does not expand any glob patterns or handle other expansions, pipelines, or redirections typically 
done by shells. 

Environment variable expansion happens before the command is executed.

