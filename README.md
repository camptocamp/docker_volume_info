volume_info
===========

This container will return file access/modification time infos.
It's mainly used to retrieve the last file access time or file modification time of a docker volume. Per default, you get a summary in JSON format.

# Usage

The docker run command must be used with the target volume mounted as /volume in the container.

example:
``` bash
docker run --rm -v $myvolume:/volume msutter/volume_info
```

this will output something like:
```json
{
  "isEmpty": false,
  "lastAccess": {
    "path": "/volume/test1",
    "fileName": "test1",
    "time": "2017-01-14T16:52:53+01:00",
    "since": 547
  },
  "lastModify": {
    "path": "/volume/test",
    "fileName": "test",
    "time": "2017-01-14T16:52:20+01:00",
    "since": 580
  },
  "mountPoint": "/volume"
}
```


Should you need to get unique values, you can use the 'jq' binary to filter the output.

### Example 1: finding out the number of seconds since last access
``` bash
docker run --rm -v $myvolume:/volume msutter/volume_info | jq .lastAccess.since
```

will return the nuber of seconds since last access. Here it will be 547

### Example 2: finding out the time of the last modification
``` bash
docker run --rm -v $myvolume:/volume msutter/volume_info | jq .lastModify.time
```

will return the date/time of the last modification. Here it will be 2017-01-14T16:52:20+01:00

# Options
Should you need also the ctime and btime, you can enable it with the ALL_TIMES environment variable.

example:
``` bash
docker run --rm -e ALL_TIMES=true -v $myvolume:/volume msutter/volume_info
```

this will output something like:
```json
{
  "isEmpty": false,
  "lastAccess": {
    "path": "/volume/test1",
    "fileName": "test1",
    "time": "2017-01-14T16:52:53+01:00",
    "since": 547
  },
  "lastBirth": {
    "path": "/volume/test1",
    "fileName": "test1",
    "time": "2017-01-14T16:42:56+01:00",
    "since": 1144
  },
  "lastChange": {
    "path": "/volume/test",
    "fileName": "test",
    "time": "2017-01-14T16:52:20+01:00",
    "since": 580
  },
  "lastModify": {
    "path": "/volume/test",
    "fileName": "test",
    "time": "2017-01-14T16:52:20+01:00",
    "since": 580
  },
  "mountPoint": "/volume"
}
```

# Build

You can use the centurylink/golang-builder to build this utility

``` bash
docker run --rm -v $(pwd):/src centurylink/golang-builder
```