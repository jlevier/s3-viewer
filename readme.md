# s3-viewer

## Debugging

<br />

First, build the executable with proper flags:
```bash
go build -gcflags=all="-N -l"
```

Run the executable:
```bash
./s3-viewer
```

Create at the root of the project a `.vscode` folder with the following `launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Attach to Process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": 0
    }
  ]
}
```

Now hit f5 or click the Start Debugging button.  This will drop down the list of processes.  Find the `s3-viewer` process you started before and now vscode will attach the debugger.