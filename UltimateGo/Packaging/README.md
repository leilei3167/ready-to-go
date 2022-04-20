关于Package的管理规范

    可按照本项目结构来进行工程管理,要遵循以下几点原则:
        1.永远不要向上导入,即底层的包导入上层的包(如internal中的导入cmd的包)
        2.永远不要出现common utils src这种包,每一个包都要有明确的功能和目的,不能出现涵盖范围过广的包

├── api  对外提供的api接口文件(如以后的protobuf文件,swagger等)
├── build 持续集成相关,如github action,或docker文件
│   └── ci
├── cmd 各个组件的main函数所在地,不存放过多代码,其编译好的文件和main.go放在一起
│   ├── master
│   │   └── main.go
│   └── node
│       └── main.go
├── config 配置文件所在处
│   └── app.yaml
├── internal 私有库代码(项目私有业务逻辑代码所在处)
│   ├── mid
│   │   └── logger.go
│   └── platform 业务代码所在处
│       ├── db
│       │   ├── account.go
│       │   └── db.go
│       └── docker
└── web  web相关的静态资源
    ├── app
    ├── static
    └── template
