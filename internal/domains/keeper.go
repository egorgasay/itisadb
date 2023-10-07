package domains

// Keeper is a middleware between storage and service.
// It reimplements some methods from storage to do transaction logging.
type Keeper Storage
