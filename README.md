# gotop NVidia extension

Provides NVidia GPU data to gotop

## Compiling

The easiest way is to not try to compile from this repository, but use the [gotop-builder](github.com/xxxserxxx/gotop-builder):

```
mkdir gotop-nvidia
cd gotop-nvidia
go run github.com/xxxserxxx/gotop-builder -r v4.1.0 github.com/xxxserxxx/gotop-nvidia
go build -o gotop-nvidia ./gotop.go
```

Check the documentation for gotop-builder, and the source -- it's not a very long file.  What this program does is downloads the `gotop.go` command file, modifies it to import the nvidia code (this project), and creates an appropropriate `go.mod` to import the main gotop project. You end up with `gotop.go` and a `go.mod` files in the directory, which you can then build into an executable with the nvidia plugin.

If you need or want to compile locally (such as, for development), there are a couple of variations. It ain't trivial.

### If you want to use the gotop code from github directly

From this repository, run:

```
mkdir tmp
cd tmp
go run github.com/xxxserxxx/gotop-builder -r v4.1.0 github.com/xxxserxxx/gotop-nvidia
echo 'replace github.com/xxxserxxx/gotop-nvidia => ../' >> go.mod
go build -o gotop .
```

This builds with the checked-out nvidia code.

### If you want to use the gotop code from a local path

then do everything in the previous example, but before compiling do an additional replacement:

```
echo 'replace github.com/xxxserxxx/gotop/v4 => ../../gotop' >> go.mod
```

## Dependencies

- [nvidia-smi](https://wiki.archlinux.org/index.php/NVIDIA/Tips_and_tricks#nvidia-smi)

## Configuration

The refresh rate of NVidia data is controlled by the `nvidia-refresh` parameter in the configuration file.  This is a Go `time.Duration` format, for example `2s`, `500ms`, `1m`, etc.

## Alternatives to test

https://github.com/mindprince/gonvml
