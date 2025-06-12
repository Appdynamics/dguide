## 🤝 Steps to run the tool 

Run dguide help command to see list of supported commands.
```
~$ dguide help
```

Select he sub-command you want to run and collect information for a specific agent installed in your env.

For example, the below command prints the top 10 CPU intensive processes on your system.
``` 
~$ dguide run ptop
```

The below demo shows how to install and run the dguide CLI. 

![dguide-demo](img/dguide_quick_demo1742210637241.gif)



## Available command type and command names

| Command_type | Command_name     | Example                     |
|--------------|------------------|---------------------------- |
| run          | java, dotnet     | `dguide run java -w -z`     |
| collect      | node, py, websrv | `dguide collect node -w -z` |


To collect the agent log from /tmp/appd, use -z:


```dguide <command_type> <command name> -w -z```


To collect the agent log from a custom path, use -l:
 
```dguide <command_type> <command_name> -w -z -l /opt/<agent_logging_path>```