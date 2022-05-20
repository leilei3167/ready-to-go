// Package v1 暂时提供master解析用户输入ip,并验证解析后包装为任务发送至kafka的功能,节点作为消费者从kafka中获取任务,并执行扫描
//将扫描结果写入数据库(mongo或者mysql)
//master只提供任务的发送,以及结果查询接口
//后续增加按ip查询,分页查询端口开放情况等api
package v1
