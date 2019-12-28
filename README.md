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

## GRPC Service

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

  // TODO(ca): below methods are not implemented.
  // rpc Update(UpdateRequest) returns (UpdateResponse) {}
  // rpc Delete(DeleteRequest) returns (DeleteResponse) {}
}
```

## Commands (Development)

`make build`: build user service for osx.

`make linux`: build user service for linux os.

`make docker`: build docker.

`docker run -it -p 5020:5020 users-api`: run docker.

`PORT=<port> POSTGRES_DSN=<postgres_dsn> ./bin/users-api`: run user service.
