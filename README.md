# INOUT TESTER

### What is it?

TODO

Screenshot

### Running - native

You can download binaries for your platforms from ...

After that it's just simple as:
```
./inout_tester
```

### Running - docker

```
docker build . -t dist
docker run -p 8080:8080 -it dist /dist/inout_tester
```
or with persistent storage mounted at `/storage`
```
docker run -it -v $(pwd):/storage -p 8080:8080 --memory=1024m dist /dist/inout_tester -problems-dir /storage/problems -submissions-dir /storage/submissions
```
