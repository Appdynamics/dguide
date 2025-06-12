## 🤝 Steps to run the tool 

Run dguide help command to see list of supported commands.
```
~$ dguide help
```

## Available command type and command names

| Command_type | Command_name     | Example                     |
|--------------|------------------|---------------------------- |
| run          | java, dotnet     | `dguide run java -w -z`     |
| collect      | node, py, websrv | `dguide collect node -w -z` |


To collect the agent log from /tmp/appd, use -z:


```dguide <command_type> <command name> -w -z```


To collect the agent log from a custom path, use -l:
 
```dguide <command_type> <command_name> -w -z -l /opt/<agent_logging_path>```