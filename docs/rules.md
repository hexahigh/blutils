# General Rules

### Parameters

- Parameters must be put in a struct.

- To prevent name collision, parameter should be prefixed with the command name, for example, `report.go`'s struct should be named `ReportParams`.

- Parameter variable names should start with lowercase while struct named should start with uppercase, for example:
```go
 type ReportParams struct {
    /* Parameters here */
}

var reportParams ReportParams
```

- Parameters should be passed as pointers.