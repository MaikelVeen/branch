# Debugging in Visual Code

Sadly the integrated debug console from the vs code Go plugin does not allow input from stdin. Debugging in Visual Code can be done using the following workaround:

1. Place the following config under the `configurations` key in the `launch.json`:

```json
{
  "name": "Remote debug",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "remotePath": "${workspaceFolder}",
  "port": 2345,
  "host": "127.0.0.1",
  "args": ["login"]
}
```

2. Run the following command, please note that `<command>` is the command you want to test such as login

```
dlv debug --headless --listen=:2345 --log --api-version=2 -- <command>
```

3. Press F5 to attach to the remote debugger, happy programming!
