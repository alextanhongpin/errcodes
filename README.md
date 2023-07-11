# errcodes

Application errors in golang that can map to gRPC and HTTP error codes.

A more lightweight version of [errors](https://github.com/alextanhongpin/errors).


## Why are there no stacktrace?

Mainly because the suggested approach is to use _sentinel error_ in golang.

Declaring a sentinel error shouldn't create a stack trace at that location.

Stack trace should point to the location in the code where that error happens, and that is why `panic` generates a stack trace.

The idea of having stack trace seems to collide with idiomatic go since we should not `panic` but return the error explicitly.


## Why are there no way to set data in error?


Again, if we declare an error as sentinel, it means we have to be careful when setting a data to the pointer of an error.


## How to avoid duplicating stacktrace?

We do not want to expose the stacktrace everytime we wrap an error. This will cause duplication in error stack whenever the stacktrace is extracted from every error chain.

Ideally we want to have only a single stacktrace originating from the source of the error. However, in order to add annotations at specific part of the code, we need to also recover the stacktrace at the line to annotate, and merge and deduplicate both stacktraces.

The issue with having only a single root error holding the stacktrace is, when the error chain is long, we may not be able to recover the error.

Also, if we set a low limit on the PCs (program counter) to recover, we might not obtain enough information for the stacktrace.

To solve the issue of overexposing the stacktrace for every error we wrap, we use the following pseudo code:


```markdown
1. first, create a root error with the stacktrace of size n. The method `StackTrace` will return the []uintptr
2. next, when we wrap an existing error with stacktrace, check the size of the last stacktrace
3. if the size is equal n, then allow exposing the `StackTrace` method
4. otherwise, the `StackTrace` method will return nothing
5. for each error we wrap, we still create the stacktrace
```
