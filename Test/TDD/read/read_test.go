//在同一包名添加_test后缀,可以使得两个文件共存,可以模拟包外调用的情况(外部测试,需要手动导入包)
package blogposts_test

import (
	"errors"
	"io/fs"
	blogposts "ready-to-go/Test/TDD/read"
	"reflect"
	"testing"
	"testing/fstest"
)

/*
想象需要编写一个解析文件夹中md文件并放到http服务器上的程序。


开始工作时,不应该相信过度活跃的想象力,而去做某种抽象,这不是迭代性的,错过TDD给我们的紧密反馈循环

从消费者的角度考虑如何使用我们要写的代码(编写我们想看到的测试)
专注于"什么"和"为什么"而不是被如何 分心

我们的包需要提供一个可以指向一个文件夹的函数,并返回一些帖子

*/

/*
要编写此类函数的测试,我们是否需要创建一些用于的测试文件夹?里面包含大量帖子?请考虑如下问题:
	-每个测试都需要创建新的文件来测试某种行为
	-有些行为是非常难以测试的,如 读取文件错误
	-这类测试运行会比较慢

将测试和文件系统耦合是个非常不必要的做法,在go 1.16引入了文件系统的抽象,io/fs包,这就使得我们的
测试能够与真实的文件系统解耦,在测试中根据需要来注入不同的实现

testing/fstest提供了一个io/FS的实现(类似于httptest一类的工具)

TDD的工作流程就是:
先写测试
->能够使测试运行的最少代码通过编译
->多写一些代码使得测试能够通过
->重构代码(消除重复,拆分函数,或者抽象,这一步因为已经通过了测试,可以自由的深入思考如何重构解耦,可以大胆的尝试任何构建方法!)
->写新功能测试(重复以上步骤)

*/
const (
	firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
	secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
)

func TestBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}
	posts, _ := blogposts.NewPostsFromFS(fs)
	if len(posts) != len(fs) {
		t.Errorf("expected %d posts, got %d", len(fs), len(posts))
	}

	_, err := blogposts.NewPostsFromFS(&StubFailingFS{})
	if err == nil {
		t.Error("expected error, got nil")
	}

	assertPost(t, posts[0], blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
		Body: `Hello
World`,
	})
}

type StubFailingFS struct{} //用于测试失败情况

func (s *StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("always fail")
}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}

}
