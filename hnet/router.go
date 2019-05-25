package hnet

import "hinx/hiface"

// 在实现Router时，先嵌入这个BaseRouter基类
// 然后根据需要对这个基类的三个方法进行重写即可
type BaseRouter struct {

}

/*
这里每个方法都为空，是因为有的Router可能不需要Pre和Post这两个方法.
这样其他实现IRouter的接口只需要把BaseRouter写进Struct声明中，即实现了IRouter接口
不需要重复实现空方法了。
*/

// 处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request hiface.IRequest) {

}

// 处理conn业务的主方法
func (br *BaseRouter) Handle(request hiface.IRequest) {

}

// 处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request hiface.IRequest) {

}