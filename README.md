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
