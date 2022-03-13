package initialize

import "go.uber.org/zap"

func InitLogger() {
	// S()可以获取一个全局的sugar, 可以让我们自己设置一个全局的logger
	// 日志锁分级别的，debug， info， warn， error， fetal
	// S函数和L函数很有用，线程安全
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
