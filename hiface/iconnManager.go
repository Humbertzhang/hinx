package hiface

// 连接管理模块
// 作用：管理 Hinx Server的连接

// 连接管理模块
type IConnManager interface {
	// 添加连接
	AddConn(conn IConnection)
	// 删除
	DeleteConn(conn IConnection)
	// 根据ID获取
	GetConnByID(connID uint32) (IConnection, error)
	// 获取连接个数
	Len() int
	// 清除并终止所有连接
	ClearConn()
}