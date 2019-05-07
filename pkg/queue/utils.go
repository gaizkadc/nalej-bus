/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package queue

import (
    "github.com/golang/protobuf/proto"
    "github.com/nalej/derrors"
)

// Marshall a given protobuffer message.
// params:
//  pb the protobuffer object
// return:
//  generated message
//  error if any
func MarshallPbMsg(pb proto.Message) ([]byte, derrors.Error){
    msg, err := proto.Marshal(pb)
    if err != nil {
        return nil, derrors.NewInternalError("error when marshalling", err)
    }

    return msg, nil
}

func UnmarshallPbMsg(msg []byte, class proto.Message) derrors.Error {
    err := proto.Unmarshal(msg, class)
    if err != nil {
        return derrors.NewInternalError("error when unmarshalling")
    }
    return nil
}
