/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
func MarshallPbMsg(pb proto.Message) ([]byte, derrors.Error) {
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
