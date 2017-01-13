package main // import "github.com/msutter/docker_volume_info"

import (
    "encoding/json"
    "fmt"
    "github.com/djherbis/times"
    "io"
    "log"
    "os"
    "path/filepath"
    "reflect"
    "time"
)

const MandatoryMountPoint string = "/volume"

type FileInfos struct {
    Path    string      // path of the file
    Name    string      // base name of the file
    Size    int64       // length in bytes for regular files; system-dependent for others
    Mode    os.FileMode // file mode bits
    ModTime time.Time   // modification time
    IsDir   bool        // abbreviation for Mode().IsDir()
    Sys     interface{} // underlying data source (can return nil)

    Atime time.Time
    Mtime time.Time
    Ctime time.Time
    Btime time.Time
}

type Volume struct {
    IsEmpty        bool         `json:"isEmpty"`
    FilesInfos     []*FileInfos `json:"fileInfos"`
    LastAtime      time.Time    `json:"lastAtime"`
    LastMtime      time.Time    `json:"lastMtime"`
    LastCtime      time.Time    `json:"lastCtime"`
    LastBtime      time.Time    `json:"lastBtime"`
    LastAtimeSince string       `json:"lastAtimeSince"`
    LastMtimeSince string       `json:"lastMtimeSince"`
    LastCtimeSince string       `json:"lastCtimeSince"`
    LastBtimeSince string       `json:"lastBtimeSince"`
}

type Output struct{}

// function to round seconds (be aware this is not a correct rounding)
func LastTimeSinceInSeconds(lastTime time.Time) string {
    return fmt.Sprintf("%.0f", time.Since(lastTime).Seconds())
}

// function fieldSet for json outputs selection
func fieldSet(fields ...string) map[string]bool {
    set := make(map[string]bool, len(fields))
    for _, v := range fields {
        set[v] = true
    }
    return set
}

// helper function SelectFields for json outputs selection
func (v *Volume) SelectFields(fields ...string) map[string]interface{} {
    fs := fieldSet(fields...)
    rt, rv := reflect.TypeOf(*v), reflect.ValueOf(*v)
    out := make(map[string]interface{}, rt.NumField())
    for i := 0; i < rt.NumField(); i++ {
        field := rt.Field(i)
        jsonKey := field.Tag.Get("json")
        if fs[jsonKey] {
            out[jsonKey] = rv.Field(i).Interface()
        }
    }
    return out
}

// helper function IsEmpty to check if volume is empty (no files or not mounted)
func IsEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

func main() {

    // check if mounted or files within volume
    volEmpty, _ := IsEmpty(MandatoryMountPoint)
    vol := Volume{IsEmpty: volEmpty}

    err := filepath.Walk(MandatoryMountPoint, func(path string, info os.FileInfo, err error) error {

        if err != nil {
            log.Print(err)
            return nil
        }

        if path == MandatoryMountPoint { // Will skip walking of directory pictures and its contents.
            return err
        }

        t, err := times.Stat(path)
        if err != nil {
            log.Fatal(err.Error())
        }

        aTime := t.AccessTime()
        if aTime.After(vol.LastAtime) {
            vol.LastAtime = aTime
        }

        mTime := t.ModTime()
        if mTime.After(vol.LastMtime) {
            vol.LastMtime = mTime
        }

        cTime := t.ChangeTime()
        if cTime.After(vol.LastCtime) {
            vol.LastCtime = cTime
        }

        fileInfos := &FileInfos{
            Path:  path,
            Name:  info.Name(),
            Size:  info.Size(),
            Mode:  info.Mode(),
            IsDir: info.IsDir(),
            Sys:   info.Sys(),
            Atime: aTime,
            Mtime: mTime,
            Ctime: cTime,
        }

        if t.HasBirthTime() {
            fileInfos.Btime = t.BirthTime()
            if fileInfos.Btime.After(vol.LastBtime) {
                vol.LastBtime = fileInfos.Btime
            }
        }

        vol.FilesInfos = append(vol.FilesInfos, fileInfos)
        return err

    })

    // Get time since now in seconds
    vol.LastAtimeSince = LastTimeSinceInSeconds(vol.LastAtime)
    vol.LastMtimeSince = LastTimeSinceInSeconds(vol.LastMtime)
    vol.LastCtimeSince = LastTimeSinceInSeconds(vol.LastCtime)
    vol.LastBtimeSince = LastTimeSinceInSeconds(vol.LastBtime)

    defaultOutputs := []string{
        "isEmpty",
        "lastAtimeSince",
        "lastMtimeSince",
        "lastCtimeSince",
        "lastBtimeSince",
        "lastAtime",
        "lastMtime",
        "lastCtime",
        "lastBtime",
    }

    if os.Getenv("OUTPUT_FILE_INFOS") == "true" {
        defaultOutputs = append(defaultOutputs, "fileInfos")
    }

    outputJson, _ := json.MarshalIndent(vol.SelectFields(defaultOutputs...), "", "  ")
    fmt.Println(string(outputJson))

    if err != nil {
        log.Fatal(err)
    }
}
