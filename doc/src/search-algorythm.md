# Unique value search algorithm:

By default, the value is saved to the minimally loaded server and returns its number to the Client. The client takes the dirty work of storing and using this number for subsequent accesses to this key. The search for instances occurs only when the -2 flag is specified and is performed in parallel goroutines.
