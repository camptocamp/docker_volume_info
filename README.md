volume_info
===========

This container will return volume file access time infos.
It's mainly used to retrieve the last file access time or file modification time of a docker volume. Per default, you get a summary in JSON format.

# Usage

The docker run command must be used with the target volume mounted as /volume in the container.

example:
``` bash
docker run --rm -v $myvolume:**/volume** msutter/volume_info
```

Should you need to get unique values, you can use the 'jq' binary to filter the output.

example
``` bash
docker run --rm -v $myvolume:/volume msutter/volume_info | jq .lastMtimeSince
```


# Options

You can set an option to output **EVERY** file infos of the volume with the OUTPUT_FILE_INFOS=true environment variable. But this could generate a big output in case of many files in the target volume.

example:
``` bash
docker run --rm -e OUTPUT_FILE_INFOS=true -v $myvolume:**/volume** msutter/volume_info
```

# Build

You can use the centurylink/golang-builder to build this utility

``` bash
docker run --rm -v $(pwd):/src centurylink/golang-builder
```