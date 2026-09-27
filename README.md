# Kitchen — Go gRPC Microservices Demo

A small example of two Go services talking to each other over **gRPC**, with a public **HTTP/JSON** API in front.

- **orders**: owns order data. It serves a gRPC API (`:9000`) and a small HTTP API (`:8000`).
- **kitchen**: an HTTP gateway (`:8080`). It turns REST/JSON requests into gRPC calls to the orders service.

```
            HTTP/JSON                 gRPC
 client ───────────────▶ kitchen ─────────────▶ orders
          :8080                      :9000       (in-memory store)

 client ───────────────────────────────────────▶ orders
          HTTP/JSON :8000 (POST /orders only)
```

## Project layout

```
.
├── protobuf/orders.proto             # gRPC service + message definitions (source of truth)
├── services/
│   ├── common/
│   │   ├── genproto/orders/          # Go code generated from orders.proto (do not edit)
│   │   └── util/util.go              # JSON read/write helpers shared by HTTP handlers
│   ├── kitchen/                      # HTTP gateway → gRPC client
│   │   ├── main.go
│   │   └── http.go
│   └── orders/                       # Orders service (gRPC + HTTP servers)
│       ├── main.go                   # starts HTTP (:8000) and gRPC (:9000)
│       ├── grpc.go / http.go         # server bootstrapping
│       ├── handlers/orders/          # transport layer (gRPC + HTTP handlers)
│       ├── service/orders.go         # business logic + in-memory "database"
│       └── types/types.go            # OrderService interface
├── Makefile
└── go.mod
```

The orders service has three layers:
**handlers** (transport) → **types.OrderService** (interface) → **service** (implementation).
The gRPC and HTTP handlers both depend only on the interface, so you can swap the storage without changing them.

## gRPC API

Defined in [protobuf/orders.proto](protobuf/orders.proto):

```proto
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrders(GetOrdersRequest) returns (GetOrderResponse);
}
```

## Prerequisites

- Go (the version in `go.mod`)
- These are only needed if you change the `.proto` file:
  - [`protoc`](https://grpc.io/docs/protoc-installation/)
  - Go plugins:
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
    Make sure `$(go env GOPATH)/bin` is on your `PATH`.

## Running

Install dependencies:

```bash
go mod download
```

Start the **orders** service first, because kitchen connects to it on `:9000`:

```bash
make run-orders
# Starting server on:  :8000
# Starting gRPC server on  :9000
```

Then start the **kitchen** gateway in a second terminal:

```bash
make run-kitchen
# Connected to GRPC Server on  :9000
# Starting server on:  :8080
```

## Trying it out

Create an order through the kitchen gateway (HTTP → gRPC):

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"customerId": 1, "productId": 42, "quantity": 2}'
# {"status":"Success"}
```

List orders:

```bash
curl http://localhost:8080/orders/1
# {"orders":[{"OrderId":4,"CustomerId":1,"ProductId":42,"Quantity":2}]}
```

Create an order directly on the orders service's HTTP API:

```bash
curl -X POST http://localhost:8000/orders \
  -H "Content-Type: application/json" \
  -d '{"customerId": 2, "productId": 7, "quantity": 1}'
```

You can also call the gRPC API directly with a tool like [grpcurl](https://github.com/fullstorydev/grpcurl). Server reflection is not enabled, so pass the proto file:

```bash
grpcurl -plaintext -import-path protobuf -proto orders.proto \
  -d '{"customerId": 1}' localhost:9000 common.OrderService/GetOrders
```

## Regenerating protobuf code

After you edit `protobuf/orders.proto`, run:

```bash
make gen
```

This rewrites the files in `services/common/genproto/orders/`.

## Make targets

| Target             | Description                                   |
| ------------------ | --------------------------------------------- |
| `make gen`         | Regenerate Go code from `protobuf/orders.proto` |
| `make run-orders`  | Run the orders service (HTTP :8000, gRPC :9000) |
| `make run-kitchen` | Run the kitchen HTTP gateway (:8080)          |

## Current limitations

This is a learning project. Keep these in mind:

- Orders are stored in a slice in memory. They are lost on restart and are not protected by a mutex.
- The orders HTTP server (`:8000`) and gRPC server (`:9000`) each create their own `OrderService`. They share data only because the store is a package-level global.
- `GetOrders` returns **all** orders and ignores `customerId`.
- Order IDs are hardcoded (`4` via gRPC, `5` via HTTP), and `totalPrice` is never set.
- Ports and addresses are hardcoded. There is no config or env var support.
- gRPC uses insecure (plaintext) transport.
- There are no tests yet.
