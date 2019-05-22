package siface

type IConnManager interface {
	Add(conn IConnection)

	Remove(conn IConnection)

	Get(connID uint32) (IConnection, error)

	Len() int

	ClearConn()
}
