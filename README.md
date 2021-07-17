# INOUT TESTER

### What is it?

Simple tester of C++ solutions which take some input from STDIN and write some output to STDOUT.
You can add your own problems with testcases just by copying .in/.out files to subdirectory of 'problems' directory.

![](homepage_screenshot.png?raw=true)

### Running - native

You can download binaries for your platforms from ...

Only single binary is required. Works on Windows, Linux, MacOS. Running is as simple as:
```
./inout_tester
```

It will create "multiply by 2" problem, so you can start testing right away.
The testcase contains single integer on its input, and expects integer multiplied by 2.


### Running - docker

```
docker build . -t dist
docker run -p 8080:8080 -it dist /dist/inout_tester
```
or with persistent storage mounted at `/storage`
```
docker run -it -v $(pwd):/storage -p 8080:8080 --memory=1024m dist /dist/inout_tester -problems-dir /storage/problems -submissions-dir /storage/submissions
```



### Development

```
go build -i .
```

```
go test ./...
```

```
go test ./... -coverprofile=cp.out && go tool cover -html=cp.out
```

test b1
