# service

 - Business logic layer, between the `handler` layer and the `repository` layer
 - `service` can only fetch data through the `repository` layer
 - Interface-oriented programming
 - Depend on the interface, not on the implementation
 - If there is a transaction, process it at this layer
 - If it is a third-party service called, please do not add `cache` to avoid cache inconsistency (the other party 
 - updates the data, which cannot be known here)
 - Since `service` will be called by `http` or `rpc`, `http` calls are provided by default, for example: `GetUserInfo()`,
   If `rpc` needs to be called, you can encapsulate `GetUserInfo()`, for example: `GetUser()`.
 
 ## Reference
 
 - https://github.com/qiangxue/go-rest-api
 - https://github.com/irahardianto/service-pattern-go
 - https://github.com/golang-standards/project-layout
 - https://www.youtube.com/watch?v=oL6JBUk6tj0
 - https://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html
 - https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
 - [If `rpc` needs to be called, you can encapsulate `GetUserInfo()`, for example: `GetUser()`. Write service layer logic on demand](https://www.5-wow.com/article/detail/89)
 - [Go programming in practice: how to organize code and write tests?](https://www.infoq.cn/article/4TAWp8YNYcVD4t046EGd)
 - https://github.com/sdgmf/go-project-sample
 - [Golang Microservices Best Practices](https://sdgmf.github.io/goproject/)
 - [layout common large web project layering](https://chai2010.cn/advanced-go-programming-book/ch5-web/ch5-07-layout-of-web-project.html)
 - [Try Clean Architecture in Golang](https://studygolang.com/articles/12909)
