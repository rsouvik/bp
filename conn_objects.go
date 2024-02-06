package bp

// Shared External Connections. Useful for avoiding multiple
// 'go' clients to re-establish the connection everytime

type SharedExtConn struct {
	Msql *SqlContext // SQL Context
}
