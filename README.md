# go-subfinder

A tiny CLI utility to find and download subtitle files, written in Go.

The program will look for a video file in the current working directory, and automatically download a subtitle file from [OpenSubtitles](https://opensubtitles.com) when one is found.

Note: very much a work in progress and by no means complete. Don't use this repository as a reference for anything.

## Usage

```shell
$ go install github.com/Kwintenvdb/go-subfinder
$ cd directory_with_video_file
$ go-subfinder download
```

### Configuration and arguments

All required config properties may also be stored in a YAML config file. To do so, create a `config/config.yml` file in the same directory as the `go-subfinder` executable.


| Config property       | Description                                                                                                                                                                |
| --------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `username` (required) | OpenSubtitles username                                                                                                                                                     |
| `password` (required) | OpenSubtitles password                                                                                                                                                     |
| `api-key` (required)  | OpenSubtitles API key. See the [docs](https://opensubtitles.stoplight.io/docs/opensubtitles-api/ZG9jOjI3NTQ2OTAy-getting-started#api-key) for how to retrieve yours.       |
| `language` (optional) | Specify the language of subtitles to search for using a [ISO-639-1 language code](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `--language=fr` for French) |

All config properties can be passed as command line arguments, e.g.:

```shell
$ go-subfinder download --username=my_username --password=my_password --api-key=my_api_key
```
