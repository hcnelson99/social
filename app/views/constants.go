package views

// Query keys: these are the keys in the URI
// Example: `/search?q=term&size=100`, `q` and `size` are query keys
const (
	// URL query key indicating a URI to redirect to
	CONTINUE_QUERY_KEY = "continue"

	// URL query key indicating an error message to display after a redirect
	// to the error page
	ERROR_QUERY_KEY = "error"
)

// Values for the URL query `ERROR_QUERY_KEY`. Note the actual error message
// that is displayed to the user is different, this is just the value that is
// shown in the URL.
const (
	ERROR_INVALIDATE_USER_SESSIONS = "error-invalidate-user-sessions"
)
