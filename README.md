## 包含service层和dao层的gRPC服务架构图

```mermaid
graph TD
    subgraph "Client 客户端"
        Client[gRPC 客户端]
    end

    subgraph "API 定义层 (PB)"
        PB[protobuf 定义]
        PB --> |包含| ServiceDef[服务接口定义]
        PB --> |包含| MsgDef[消息类型定义]
    end

    subgraph "Server 层 (gRPC 服务实现)"
        GrpcServer[gRPC 服务实现类]
        GrpcServer --> |实现| ServiceDef
        GrpcServer --> |处理| Request[客户端请求]
        GrpcServer --> |返回| Response[响应结果]
        GrpcServer --> |调用| ServiceLayer[Service 层]
    end
    subgraph "Service 层 (业务逻辑)"
        ServiceImpl[业务逻辑实现]
        ServiceImpl --> |包含| Transaction[事务控制]
        ServiceImpl --> |包含| Validation[参数校验]
        ServiceImpl --> |包含| BusinessRule[业务规则]
        ServiceImpl --> |调用| DaoLayer[DAO 层]
    end

    subgraph "DAO 层 (数据访问)"
        DaoImpl[数据访问实现]
        DaoImpl --> |使用| ORM[Xorm/GORM]
        DaoImpl --> |操作| Model["数据模型(Models)"]
        ORM --> |交互| DB[(数据库)]
    end
    Client --> |gRPC 调用| GrpcServer
    GrpcServer --> |解析请求| MsgDef
    MsgDef --> |映射为| Model
    ServiceImpl --> |处理后返回| GrpcServer
    DaoImpl --> |返回数据| ServiceImpl
```

# 自动生成gRPC文件(./pb)
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user_growth.proto
```

# protobuf文件来源
    google/api/*.proto 来自 https://github.com/googleapis/googleapis
    google/protobuf/*.proto 来自 https://github.com/protocolbuffers/protobuf

# 生成protoset文件
```
protoc --proto_path=. --descriptor_set_out=myservice.protoset --include_imports ./user_growth.proto
```

# grpcurl调用
```
# 使用gRPC服务
grpcurl -plaintext localhost:80 list
grpcurl -plaintext localhost:80 list UserGrowth.UserCoin
grpcurl -plaintext localhost:80 describe
grpcurl -plaintext localhost:80 describe UserGrowth.UserCoin
grpcurl -plaintext localhost:80 describe UserGrowth.UserCoin.ListTasks
# 使用proto文件
grpcurl -import-path ./ -proto user_growth.proto list
# 使用protoset文件
grpcurl -protoset myservice.protoset list UserGrowth.UserCoin
# 调用gRPC服务
grpcurl -plaintext localhost:80 UserGrowth.UserCoin/ListTasks
grpcurl -plaintext -d '{"uid":1}' localhost:80 UserGrowth.UserCoin/UserCoinInfo
```

# 生成grpc-gateway代码
````
protoc -I . --grpc-gateway_out ./ \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    user_growth.proto
````
