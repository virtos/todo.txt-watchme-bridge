# todo.txt - watchme bridge

[todo.txt](http://todotxt.org/) is a format for tracking todos. Unfortunately all implementations lack time tracking capabilities. And to tell the truth, it doesn't belong there as the idea behind it is to be provide a very simple solution.

To add time tracking capabilities I have chosen [Watchme](http://www.flamebrain.com/download-watchme/) utility.

**todo-watchme_bridge** parses a todo.txt file, selects tasks assigned for today and adds them to the Watchme config file as timers. There a couple of rules helping with usability:

- No existing timers are destroyed or updated
- If an item (todo.txt task) already exist as a timer in Watchme, it will remain there keeping the current time.

## Usage of todo-watchme_bridge:

```
todo-watchme_bridge [flags]  
  -td string  
        A source todo.txt file name (default "todo.txt")  
  -tt string
        Todo item template file name (default "template.xml")
  -w string
        WatchMe configuration file name (default "WatchMeConfig.xml")
  -h help
If no path specified, the program will look in the current directory
```
HINT: it is easy to create a simple batch file that calls this utility and then runs Watchme. Running this batch file will always make sure you always add current tasks for tracking.

# Compiling

```
go build
```

