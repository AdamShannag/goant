# Goant

Goant is a lightweight tool for scanning Go projects for annotated types or interfaces and running commands dynamically
based on those annotations. It works with any tool that can be triggered from the command line.

<p>
  <a href="https://github.com/AdamShannag/goant/actions/workflows/build.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/AdamShannag/goant/build.yml?branch=master&label=CI%2FCD&logo=github" alt="CI/CD">
  </a>
  <a href="https://github.com/AdamShannag/goant/blob/master/go.mod">
    <img src="https://img.shields.io/github/go-mod/go-version/AdamShannag/goant?logo=go" alt="Go Version">
  </a>
  <a href="https://github.com/AdamShannag/goant/releases/latest">
    <img src="https://img.shields.io/github/v/release/AdamShannag/goant?sort=semver" alt="Latest Release">
  </a>
  <a href="https://github.com/AdamShannag/goant/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT">
  </a>
</p>

## Installation

You can install and run Goant in two ways:

### 1. Install via `go install`

```bash
go install github.com/AdamShannag/goant/cmd/goant@v0.1.0
```

Then run it directly:

```bash
goant -keyword=myAnnotation -cmd="your command here"
```

### 2. Run directly with `go run`

```bash
go run github.com/AdamShannag/goant/cmd/goant@v0.1.0 -keyword=myAnnotation -cmd="your command here"
```

## Flags

* `-root` – Root directory to start scanning Go files (default: `.`).
* `-keyword` – Comment keyword to search for, e.g. `gomock`.
* `-cmd` – Command template with placeholders.
* `-dry` – Print the commands that would be executed without actually running them.
* `-s` – Suppress all output.
* `-skip` – Comma-separated list of directories to skip.
* `-version` – Print Goant version and exit.

## Example

Suppose you have an interface annotated like this:

```go
// gomock:package=user_mock out=.gen/mocks/user_mock/mock_repo.go
type Repository interface {
    Create(ctx context.Context, user domain.User) (domain.User, error)
}

```

Run Goant:

```bash
goant -keyword=gomock -cmd="go run go.uber.org/mock/mockgen@v0.1.0 -destination=@out -package=@package -source=@path @type"
```

Goant will substitute placeholders:

* `@out` → `.gen/mocks/user_mock/mock_repo.go`
* `@package` → `user_mock`
* `@path` → `./internal/repository/user/repository.go`
* `@type` → `Repository`

And then execute the generated command.

## Placeholders

### Automatic placeholders

These are resolved automatically by Goant:

- `@path` – Relative path to the Go file.
- `@type` – Name of the type.

### Custom placeholders

Any additional parameters provided in the annotation are also available as placeholders.
For example, if your annotation includes out=... and package=..., then @out and @package become available in the command
template.

You are free to define any parameter names in the annotation, and Goant will substitute
them accordingly.

## Multiple Annotations

Goant supports having **multiple annotations on a single type or interface**. Each annotation is independent and can be
associated with different tools or commands. You can also repeat the **same annotation keyword** with different
parameters. In this case, Goant will run each
annotation separately using its respective parameters.

Example:

```go
// gomock:package=user_mock out=.gen/mocks/user_mock/mock_repo.go
// myTool:echo=hello
// myTool:echo=world
type Repository interface {
    Create(ctx context.Context, user domain.User) (domain.User, error)
}

```

* The first annotation uses the `gomock` keyword with parameters for mock generation.
* The second and third annotations both use the `myTool` keyword but with different `echo` values. Goant will run **each
  annotation separately**, substituting the appropriate parameters.

This allows you to define multiple commands for the same type, either using **different keywords** or the **same keyword
with different parameter sets**.

## License

[MIT](LICENSE)

## Contribution

Contributions are welcome! If you have ideas, bug reports, or improvements, please open an issue or submit a pull
request.