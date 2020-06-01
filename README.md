# gotop NVidia extension

Provides NVidia GPU data to gotop

## Configuration

The refresh rate of NVidia data is controlled by the `nvidia-refresh` parameter in the configuration file.  This is a Go `time.Duration` format, for example `2s`, `500ms`, `1m`, etc.
