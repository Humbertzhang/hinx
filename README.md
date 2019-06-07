# hinx: A Toy TCP Server

hinx's architecture 

![hinx](https://picturesbed.oss-cn-hangzhou.aliyuncs.com/img/20190607151119.png)

TODO:
* Complete Connection manage part.
* Change MsgHandler to RequestHandler. It should dispatch Request rather than Connections, it will be more fair for workers. 
* Make logs clear and tidy.
