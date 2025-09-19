# Unique-ip-addr

## Summary
Test assignment for Lightspeed.

The program receives one argument: the name of a valid file containing IP addresses as described in the [task](https://github.com/Ecwid/new-job/blob/master/IP-Addr-Counter-GO.md).

Result of program run is following statistics:

```go
Working thread count: 14
Handled address count: 8000000000
Unique address count: 1000000000
Average velocity: 183095.690385 ip / ms
Spent time: 43.6936231s
Memory usage: 542.885704MB
Mallocs: 186083
```

## Run
1. Clone a repository
```bash
git clone https://github.com/Reasunta/unique-ip-addr.git
```
2. Go to project directory
```bash
cd unique-ip-addr/v1
```
3. Build and run the app
```bash
go run . <path_to_file_with_ip>
```

## Features
- Uses ~540 MB of memory; 512 MB is constant, the rest is multithreading overhead
- The app reads the file with several goroutines, splitting it into chunks. The number of threads is set to `runtime.NumCPU() - 2`, with a minimum of 1 thread
- Constants for file reading, process management, and testing