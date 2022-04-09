package db

//模拟一个包含数据库交互的接口
type DB interface {
	Get(key string) (int, error)
}

//模拟某个api与数据库交互
func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1

}

//mockgen -source=db.go -destination=mock/db_mock.go 生成mock
