#一,模型定义

可在结构体中嵌入gorm.Model,实现快速添加ID主键等

// gorm.Model 的定义
type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}

使用`gorm`:""来添加字段标签,实现权限控制,设置主键,索引等,甚至外键
约束等


#二,链接数据库
一般业务场景会自定义数据库的选项

dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	//基本模式,使用默认配置,但实际基本都是自己配置
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(mysql.New(mysql.Config{
		//数据库相关配置
		DSN: dsn,
          DefaultStringSize: 256, // string 类型字段的默认长度
  DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
  DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
  DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
  SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库链接出错!", err)

	}

要记得设置连接池
    先获取系统DB
    
	//设置连接池,获取底层db
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)// 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxOpenConns(100)//设置打开数据库连接的最大数量。
	sqlDB.SetConnMaxLifetime(time.Hour)//设置了连接可复用的最大时间。

    ##性能优化
    默认设置适用于大部分场景,但是仍可以根据自身实际情况来设置以下选项进一步提高性能:
    ###禁用默认事务