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

type TimeInfo struct {
    Path      string    `json:"path"`
    FileName  string    `json:"fileName"`
    Time      time.Time `json:"time"`
    TimeSince int       `json:"since"`
}

type Volume struct {
    MountPoint string   `json:"mountPoint"`
    IsEmpty    bool     `json:"isEmpty"`
    LastAccess TimeInfo `json:"lastAccess"`
    LastModify TimeInfo `json:"lastModify"`
    LastChange TimeInfo `json:"lastChange"`
    LastBirth  TimeInfo `json:"lastBirth"`
}

// function to round seconds (be aware this is not a correct rounding)
func LastTimeSinceInSeconds(lastTime time.Time) int {
    return int(time.Since(lastTime).Seconds())
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
    vol := Volume{
        MountPoint: MandatoryMountPoint,
        IsEmpty:    volEmpty,
    }

    if volEmpty == false {
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
            if aTime.After(vol.LastAccess.Time) {
                vol.LastAccess = TimeInfo{
                    Path:      path,
                    FileName:  info.Name(),
                    Time:      aTime,
                    TimeSince: LastTimeSinceInSeconds(aTime),
                }
            }

            mTime := t.ModTime()
            if mTime.After(vol.LastModify.Time) {
                vol.LastModify = TimeInfo{
                    Path:      path,
                    FileName:  info.Name(),
                    Time:      mTime,
                    TimeSince: LastTimeSinceInSeconds(mTime),
                }
            }

            cTime := t.ChangeTime()
            if mTime.After(vol.LastChange.Time) {
                vol.LastChange = TimeInfo{
                    Path:      path,
                    FileName:  info.Name(),
                    Time:      cTime,
                    TimeSince: LastTimeSinceInSeconds(cTime),
                }
            }

            if t.HasBirthTime() {
                btime := t.BirthTime()
                if btime.After(vol.LastBirth.Time) {
                    vol.LastBirth = TimeInfo{
                        Path:      path,
                        FileName:  info.Name(),
                        Time:      btime,
                        TimeSince: LastTimeSinceInSeconds(btime),
                    }
                }
            }

            return err
        })

        if err != nil {
            log.Fatal(err)
        }
    }

    defaultOutputs := []string{
        "mountPoint",
        "isEmpty",
    }

    if volEmpty == false {
        defaultOutputs = append(defaultOutputs,
            "lastAccess",
            "lastModify",
        )
        all_times := os.Getenv("ALL_TIMES")
        if all_times == "true" {
            defaultOutputs = append(defaultOutputs,
                "lastChange",
                "lastBirth",
            )
        }
    }

    outputJson, _ := json.MarshalIndent(vol.SelectFields(defaultOutputs...), "", "  ")
    fmt.Println(string(outputJson))
}
