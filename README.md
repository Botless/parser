# parser

Parse chat message events into botless commands


```

Known Event Message Types --> [Parser]  --> Botless Commands
 
```

A command looks like this in Message form:

```
/example arguments to the command
```

And this will be converted to a `botless.bot.command` event payload:

```json
{
    "cmd": "example",
    "args": "arguments to the command",
    "author": "README.md",
    "channel": "github.com"
}
```
