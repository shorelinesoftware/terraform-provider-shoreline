# Debugging with VS Code

1. Start the Backend app on a local box.
2. Set the `{{ .ENV_VARS_NAME_PREFIX }}_TOKEN` in the `private.env` file.
3. Open `Run and Debug` panel (`shift+cmd+d`).
4. Press the green arrow (`Start Debugging`) or press `F5`.
5. Open `DEBUG CONSOLE` (`shift+cmd+y`).
6. Copy the `TF_REATTACH_PROVIDERS` variable and export it in a terminal.
7. Set breakpoints in the code.
8. In the same terminal, you can now `apply` your `.tf` files and it should stop to your breakpoints (or to the line where TF provider crashes).


## Notes

- If you stop the debugger, when starting it again you should export the `TF_REATTACH_PROVIDERS` again.
- In the `DEBUG CONSOLE`(`shift+cmd+y`) you can also write expressions that will get evaluated in the current debugging context/scope.
- If the TF provider crashes, then the debug server will also crash. You need to start it again (then export `TF_REATTACH_PROVIDERS` again).
- If you make changes in the TF provider code, you need to restart the debug server.