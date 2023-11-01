# Commander

Unlike the "system" library call from C and other languages, the os/exec package intentionally does not invoke the 
system shell and does not expand any glob patterns or handle other expansions, pipelines, or redirections typically 
done by shells. 

Environment variable expansion happens before the command is executed.

