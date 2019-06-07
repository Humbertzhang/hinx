package utils

import (
	"encoding/json"
	"hinx/hiface"
	"io/ioutil"
)

/*
存储有关Hinx框架的全局参数，来供其他模块使用
某些配置可以让用户写一个hinx.json配置文件来定义
*/

type GlobalObj struct {
	// 以下是有关Server的配置
	TCPServer  			hiface.IServer 	// 当前Hinx全局的Server对象
	Host 				string 			// 当前服务器监听的IP
	Port 				int 			// 当前服务器监听的端口
	Name 				string 			// 当前服务器的名字

	// 以下是有关Hinx本身的配置
	Version 			string 			// hinx版本号
	MaxConn 			int				// 服务器允许的最大连接数
	MaxPackageSize 		uint32 			// 数据包最大值的大小
	WorkerPoolSize  	uint32 			// Worker工作池大小
	MaxWorkerTaskLen 	uint32			// 每个worker对应的任务数的长度
}

// 定义一个全局的GlobalObj对象，其他包可以直接访问
var GlobalObject *GlobalObj


// 从/conf/hinx.json中加载配置
func (g *GlobalObj) Load() {
	data, err := ioutil.ReadFile("conf/hinx.json")
	if err != nil {
		//配置文件没有，直接Panic
		panic(err)
	}

	err = json.Unmarshal(data,GlobalObject)
	if err != nil {
		panic(err)
	}
}


// 在导入包时，会先调用这个包的init方法
// 在这里我们初始化GlobalObject对象
func init() {
	// 这里我们先默认给一些值
	GlobalObject  =  &GlobalObj{
		Host: "0.0.0.0",
		Port: 8999,
		Name : "HinxServer",

		Version: "v0.4",
		MaxConn: 1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}

	// 尝试从配置文件(conf/hinx.json)中拿相应配置，如果拿到了，则覆盖之前的默认值
	GlobalObject.Load()
}