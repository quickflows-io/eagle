# Repository

Repository，is the data access layer，responsible for access DB、MC、external HTTP and other interfaces to shield data 
access details from the upper layer。  
Subsequent replacement、Upgrade ORM engine，does not affect business logic. It can improve the test efficiency. 
During unit testing, Mock objects are used to replace the actual database access, which can double the running speed of test cases.    
The benefits of applying the Repository pattern far outweigh the added code to implement this pattern. 
This pattern should be used whenever projects are layered.
Repository is a concept in DDD, emphasizing that Repository is driven by Domain (this project is mainly Service).
right Model Layers can only operate on a single table，Every method has a parameter `db *grom.DB` Instances to facilitate transaction operations。

Specific responsibilities include:
 - SQL concatenation and DB access logic
 - DB The splitting table logic
 - DB Cache read and write logic
 - HTTP Interface call logic

> Tips: If it is the returned list, try to return the map, which is convenient for the service to use。

Suggest：
 - Writing native SQL is recommended
 - The use of linked table queries is prohibited. The advantage is that it is easy to expand, such as sub-database and sub-table
 - The logic part is processed in the program

One business has one directory, and each repo go file corresponds to one table operation. For example, if the user is in the user directory, everything related to the user can be put here.
Separate into different files according to different modules, while avoiding the problem of too many funcs in a single file. for example:
  - User basic services - user_base_repo.go
  - User concern - user_follow_repo.go
  - users like   - user_like_repo.go
  - user comment - user_comment_repo.go

## unit test

There are several libraries that can be used for unit testing of databases:
 - go-sqlmock https://github.com/DATA-DOG/go-sqlmock Mainly used for interactive operations with the database: additions, deletions and changes
 - GoMock https://github.com/golang/mock

## Reference
 - https://github.com/realsangil/apimonitor/blob/fe1e9ef75dfbf021822d57ee242089167582934a/pkg/rsdb/repository.go
 - https://youtu.be/twcDf_Y2gXY?t=636
 - [Unit testing GORM with go-sqlmock in Go](https://medium.com/@rosaniline/unit-testing-gorm-with-go-sqlmock-in-go-93cbce1f6b5b)
 - [How to unit test GORM application using Sqlmock](https://1024casts.com/topics/R9re7QDaq8MnJoaXRZxdljbNA5BwoK)
