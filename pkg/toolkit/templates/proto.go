package templates

const ModelProtoTemplate = `syntax = "proto3";

package exmsg.models;

option go_package = "github.com/nhdms/base-go/proto/exmsg/models;models";
import "google/protobuf/timestamp.proto";
import "google/protobuf/struct.proto";

message {{.Handler}} {
    int64 id = 1;
	string type = 2;
    google.protobuf.Timestamp created_at = 3;
    google.protobuf.Struct metadata = 4;
}
`

const ServiceProtoTemplate = `syntax = "proto3";

package exmsg.services;

option go_package = "github.com/nhdms/base-go/proto/exmsg/services;services";

import "models/{{.ServiceName}}.proto";
import "models/common.proto";

service {{.Handler}}Service {
    rpc Get{{.Handler}}s ({{.Handler}}Request) returns ({{.Handler}}Response);
}

message {{.Handler}}Request {
  exmsg.models.Query query = 1;
  int64 id = 2;
  repeated int64 ids = 3;
}

message {{.Handler}}Response {
  repeated exmsg.models.{{.Handler}} {{ .HandlerLower }} = 1;
  exmsg.models.SQLResult exec_result = 2;
}

`
