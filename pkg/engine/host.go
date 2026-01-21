package engine

import (
	"context"
	"io"
)

// HostInterface defines the capabilities that the Host System (Go, Rust, Node, etc.)
// must provide to the ZenoEngine core.
//
// By programming against this interface, ZenoEngine becomes a pure logic kernel
// that can run in any environment that implements these methods.
type HostInterface interface {
	// --- System ---
	Log(level string, message string)

	// --- Database ---
	// Executes a query that returns rows (SELECT)
	DBQuery(ctx context.Context, dbName string, query string, args []interface{}) (Rows, error)
	// Executes a query that returns modifications (INSERT, UPDATE, DELETE)
	DBExecute(ctx context.Context, dbName string, query string, args []interface{}) (Result, error)

	// --- HTTP ---
	// Sends a response to the current HTTP context
	HTTPSendResponse(ctx context.Context, status int, contentType string, body []byte) error
	// Gets a header from the current HTTP request
	HTTPGetHeader(ctx context.Context, key string) string
	// Gets a query param from the current HTTP request
	HTTPGetQuery(ctx context.Context, key string) string
	// Gets the request body
	HTTPGetBody(ctx context.Context) ([]byte, error)
}

// Rows abstraction (iterated like sql.Rows)
type Rows interface {
	io.Closer
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
}

// Result abstraction (like sql.Result)
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
