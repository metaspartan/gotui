# Contributing to gotui

Thank you for your interest in contributing to `gotui`! We strive for **A+ Code Quality** and have strict guidelines to maintain a clean, performant, and maintainable codebase.

## 📝 Core Philosophy

1.  **Modern Go**: Use Go 1.20+ idioms.
2.  **No Comments**: Code must be self-explanatory. Use clear variable/function names and logic instead of comments. Comments are treated as a failure to write clear code.
3.  **Strict Typing**: Always use specific types. Avoid `interface{}`/`any` unless absolutely necessary.
4.  **Low Complexity**: Keep functions simple. Cyclomatic complexity must remain **under 15**. Split large functions.
5.  **Backward Compatibility**: changes **MUST NOT** break existing global APIS (`ui.Init()`, `ui.Render()`, etc.).

## 📂 Code Structure

-   **`types.go`**: All shared/core types (`InitConfig`, `Event`, `Drawable`, etc.) belong here.
-   **`backend.go`**: Encapsulates `tcell` specifics.
-   **`widgets/`**: All widgets go in this package.
-   **`_examples/`**: Examples live here (underscored to avoid import bloat).

## 🛠️ Development Workflow

1.  **Fork & Clone**:
    ```bash
    git clone https://github.com/YOUR_USERNAME/gotui.git
    cd gotui
    ```

2.  **Make Changes**: Implement your feature or fix.

3.  **Verify Quality**:
    *   Format code: `gofmt -s -w .`
    *   Check complexity: `gocyclo -over 15 .` (Install via `go install github.com/fzipp/gocyclo/cmd/gocyclo@latest`)
    *   Ensure no linter warnings.

4.  **Test**:
    *   Run the dashboard example: `go run _examples/dashboard/main.go`
    *   Run the SSH example (if touching backend): `go run _examples/ssh-dashboard/main.go`

5.  **Submit PR**: Open a Pull Request with a clear description of your changes.

## ⚠️ Issues & Bugs

Please report bugs via GitHub Issues. Include:
-   `gotui` version
-   Terminal emulator used
-   Minimal reproduction snippet

Thank you for helping make `gotui` the best TUI library for Go! 🚀
