## 🤝 Steps to run the tool 

Run dguide help command to see list of supported commands.
```sh
~$ dguide help
```

## Available command type and command names
General syntax to run the dguide cli is:
```sh
dguide <command_type> <command_name> <options>
```

| Command_type | Command_name     | Example                     |
|--------------|------------------|---------------------------- |
| run          | java, dotnet     | `dguide run java -w -z`     |
| collect      | node, py, websrv | `dguide collect node -w -z` |


To collect the agent log from /tmp/appd, use -z:


```sh
dguide <command_type> <command_name> -w -z
```


To collect the agent log from a custom path, use -l:
 
```sh
dguide <command_type> <command_name> -w -z -l /opt/<agent_logging_path>
```