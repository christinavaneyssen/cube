# cube
Build an Orchestrator in Go (From Scratch) Timothy Boring
https://github.com/gogococo/orchestrator-in-go

## References

* [Docker GoLang SDK](https://pkg.go.dev/github.com/docker/docker)
* Github [Moby v27.4.1](https://github.com/moby/moby/tree/v27.4.1)

## Go Concepts

### Method Receivers

When you write the following code, you are declaring a method that belongs to the Docker type.
The `(d *Docker)` part is the receiver. It says "this method is attached to the Docker struct, and inside the method, we'll refer to the Docker instance as `d`".
```go
func (d *Docker) PullImage(ctx context.Context) error {
	// implementation
}
```

So, when you have a Docker struct instance and call:
```go
docker := &Docker{}
docker.PullImage(ctx)
```
Go knows to call the `PullImage` method that has Docker as its receiver.

## Interfaces

Any type that implements all the methods specified in the interface automatically implements that interface.

