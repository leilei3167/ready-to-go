package main

import (
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
重点是通过sarama的mocks包来模拟2种生产者,产生某些特定的错误,从而测试返回结果
结合httptest包模拟res 和req来进行

*/

//成功场景
func TestCollectSuccessfully(t *testing.T) {
	//模拟两个生产者
	dataCollectorMock := mocks.NewSyncProducer(t, nil)
	dataCollectorMock.ExpectSendMessageAndSucceed()
	accessLogProducerMock := mocks.NewAsyncProducer(t, nil)
	accessLogProducerMock.ExpectInputAndSucceed()

	s := &Server{DataCollector: dataCollectorMock, AccessLogProducer: accessLogProducerMock}

	//安全的关闭非常重要,此处会调用Server实现的close方法
	defer safeClose(t, s)

	req, err := http.NewRequest("GET", "http://localhost:8080/?data", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()   //模拟的writer
	s.Handler().ServeHTTP(res, req) //模拟产生一次处理请求的过程
	if res.Code != 200 {            //检查结果
		t.Errorf("Expected HTTP status 200, found %d", res.Code)
	}
	if res.Body.String() != "Your data is stored with unique identifier query/0/1" {
		t.Error("Unexpected response body", res.Body)
	}

}

//500错误
func TestCollectionFailure(t *testing.T) {
	dataCollectorMock := mocks.NewSyncProducer(t, nil)
	dataCollectorMock.ExpectSendMessageAndFail(sarama.ErrRequestTimedOut) //模拟产生一个错误

	accessLogProducerMock := mocks.NewAsyncProducer(t, nil)
	accessLogProducerMock.ExpectInputAndSucceed() //异步的日志无论如何都应该成功

	s := &Server{
		DataCollector:     dataCollectorMock,
		AccessLogProducer: accessLogProducerMock,
	}
	defer safeClose(t, s)

	req, err := http.NewRequest("GET", "http://example.com/?data", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	s.Handler().ServeHTTP(res, req) //会在SendMessage出产生错误

	if res.Code != 500 {
		t.Errorf("Expected HTTP status 500, found %d", res.Code)
	}
}

//404错误
func TestWrongPath(t *testing.T) {
	dataCollectorMock := mocks.NewSyncProducer(t, nil) //路径错误时datacllector不产生数据

	accessLogProducerMock := mocks.NewAsyncProducer(t, nil)
	accessLogProducerMock.ExpectInputAndSucceed() //但日志仍然要产生

	s := &Server{
		DataCollector:     dataCollectorMock,
		AccessLogProducer: accessLogProducerMock,
	}
	defer safeClose(t, s)

	req, err := http.NewRequest("GET", "http://example.com/wrong?data", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	s.Handler().ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("Expected HTTP status 404, found %d", res.Code)
	}
}
func safeClose(t *testing.T, o io.Closer) {
	if err := o.Close(); err != nil {
		t.Error(err)
	}
}
