package bridgev2

import (
    "app"
    "errors"
    "libcommon"
    "libservicev2"
    "regexp"
    "strings"
    "util/json"
    "util/logger"
    "validate"
)

var NULL_FRAME_ERR = errors.New("frame is null")


// validate connection
func ValidateConnectionHandler(manager *ConnectionManager, frame *Frame) error {
    if frame == nil {
        return NULL_FRAME_ERR
    }

    var meta = &ConnectMeta{}
    e1 := json.Unmarshal(frame.FrameMeta, meta)
    if e1 != nil {
        return e1
    }

    response := &ConnectResponseMeta{
        UUID: app.UUID,
        New4Tracker: false,
    }

    responseFrame := &Frame{}

    if meta.Secret == app.SECRET {
        responseFrame.SetStatus(STATUS_SUCCESS)
        exist, e2 :=libservicev2.ExistsStorage(meta.UUID)
        if e2 != nil {
            responseFrame.SetStatus(STATUS_INTERNAL_ERROR)
        } else {
            if exist {
                response.New4Tracker = false
            } else {
                response.New4Tracker = true
            }
        }
        // only valid client uuid (means storage client) will log into db.
        if meta.UUID != "" && len(meta.UUID) == 30 {
            storage := &app.StorageDO{
                Uuid: meta.UUID,
                Host: "",
                Port: 0,
                Status: app.STATUS_ENABLED,
                TotalFiles: 0,
                Group: "",
                InstanceId: "",
                HttpPort: 0,
                HttpEnable: false,
                StartTime: 0,
                Download: 0,
                Upload: 0,
                Disk: 0,
                ReadOnly: false,
                Finish: 0,
                IOin: 0,
                IOout: 0,
            }
            e3 := libservicev2.SaveStorage("", storage, nil)
            if e3 != nil {
                responseFrame.SetStatus(STATUS_INTERNAL_ERROR)
            }
        }
        responseFrame.SetMeta(response)
    } else {
        responseFrame.SetStatus(STATUS_INTERNAL_ERROR)
    }
    return writeFrame(manager, responseFrame)
}







