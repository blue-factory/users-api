# Users-API

Microservice implemented in Golang that stores user information into Postgres database.

## Table

```
   Column   |           Type           | Collation | Nullable |      Default
------------+--------------------------+-----------+----------+-------------------
 id         | uuid                     |           | not null | gen_random_uuid()
 email      | character varying(255)   |           | not null |
 name       | character varying(255)   |           | not null |
 password   | character varying(255)   |           | not null |
 created_at | timestamp with time zone |           |          | now()
 updated_at | timestamp with time zone |           |          | now()
 deleted_at | timestamp with time zone |           |          |
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
Triggers:
    update_users_update_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column()
```

## gRPC Service

```go
message User {
  string id = 1;
  string email = 2;
  string name = 3;
  string password = 4;

  int64 created_at = 5;
  int64 updated_at = 6;
}

service UsersService {
  rpc Get(GetRequest) returns (GetResponse) {}
  rpc GetByEmail(GetByEmailRequest) returns (GetByEmailResponse) {}
  rpc Create(CreateRequest) returns (CreateResponse) {}
  rpc VerifyPassword(VerifyPasswordRequest) returns (VerifyPasswordResponse)  {}
  rpc List(ListRequest) returns (ListResponse) {}
  rpc Update(UpdateRequest) returns (UpdateResponse) {}
  rpc Delete(DeleteRequest) returns (DeleteResponse) {}
}
```

## Environments Values

`PORT`: define users service port.

`HOST`: define users service host.

`POSTGRES_DSN`: define postgres database connection DSN.

## Commands (Development)

`make build`: build users service for osx.

`make linux`: build users service for linux os.

`make docker`: build docker.

`make compose`: start docker-docker.

`make stop`: stop docker-docker.

`make run`: run users service.

`docker run -it -p 5020:5020 users-api`: run docker.
