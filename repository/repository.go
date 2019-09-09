package repository

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	stan "github.com/nats-io/go-nats-streaming"
	"template/models"
	"time"

	"template/provider"
)

// Repository interface
type Repository interface {
	CreateNews(news []byte) (response []byte, err error)
	GetNews(newsID []byte) (news []byte, err error)
	getNatsStreamingConn() (stan.Conn, error)
	publishMessage(dataToSend []byte, subjectTheme string) error
	subscribeSimple(subject string, handler func(m *stan.Msg)) error
	subscribeQueue(subject string, handler func(m *stan.Msg)) error
	GetNewsHandler()
	CreateNewsHandler()
	SubscribeNewsGet()
	SubscribeNewsCreate()
}

type repository struct {
	provider   provider.Provider
	timeout    time.Duration
	newsCreate chan []byte
	newsGet    chan []byte
	waitCreate MapWait
	waitGet    MapWait
	respCreate MapResp
	respGet    MapResp
}


// NewNats
func New(pr provider.Provider) Repository {
	return &repository{
		provider:   pr,
		newsCreate: make(chan []byte, 100),
		newsGet:    make(chan []byte, 100),
		waitCreate: MapWait{wait: make(map[string]int)},
		waitGet:    MapWait{wait: make(map[string]int)},
		respCreate: MapResp{resp: make(map[string][]byte)},
		respGet:   	MapResp{resp: make(map[string][]byte)},
	}
}

func (r *repository) getNatsStreamingConn() (stan.Conn, error) {
	return r.provider.GetNatsConnectionStreaming()
}

func (r *repository) CreateNews(news []byte) (response []byte, err error) {
	err = r.publishMessage(news, "create_news")
	if err != nil {
		return nil, err
	}
	r.waitCreate.addWait(string(news))

	waiting := make(chan []byte, 1)

	go func() {
		for {
			if value, ok := r.respCreate.checkSing(string(news)); ok {
				waiting <- value
			}
		}
	}()

	select {
	case res := <-waiting:
		return res, nil
	case <-time.After(time.Minute):
		return nil, errors.New("timeout")
	}

	return nil, nil
}

func (r *repository) GetNews(newsID []byte) (news []byte, err error) {
	err = r.publishMessage(newsID, "get_news")
	if err != nil {
		return nil, err
	}
	r.waitGet.addWait(string(newsID))

	waiting := make(chan []byte, 1)

	go func() {
		for {
			if value, ok := r.respGet.checkSing(string(newsID)); ok {
				waiting <- value
			}
		}
	}()

	select {
	case res := <-waiting:
		return res, nil
	case <-time.After(time.Minute):
		return nil, errors.New("timeout")
	}

	return nil, nil
}

func (r *repository) publishMessage(dataToSend []byte, subjectTheme string) error {
	nc, err := r.getNatsStreamingConn()
	if err != nil {
		return err
	}
	if err := nc.Publish(subjectTheme, dataToSend); err != nil {
		return err
	}
	fmt.Println("send", subjectTheme)
	return nil
}

func (r *repository) subscribeSimple(subject string, handler func(m *stan.Msg)) error {
	nats, err := r.getNatsStreamingConn()
	if err != nil {
		return err
	}
	optErrorsWorker := []stan.SubscriptionOption{
		stan.DurableName("remember" + subject),
		stan.MaxInflight(1),
		stan.SetManualAckMode()}

	_, err = nats.Subscribe(subject, handler, optErrorsWorker...)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) subscribeQueue(subject string, handler func(m *stan.Msg)) error {
	nats, err := r.getNatsStreamingConn()
	if err != nil {
		return err
	}
	optErrorsWorker := []stan.SubscriptionOption{
		stan.DurableName("remember" + subject),
		stan.MaxInflight(1),
		stan.SetManualAckMode()}

	queue := fmt.Sprintf("%s_queue", subject)
	_, err = nats.QueueSubscribe(subject, queue, handler, optErrorsWorker...)

	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SubscribeNews(subject string, channel chan []byte) {
	handlerErrLog := func(m *stan.Msg) {
		defer func() {
			err := m.Ack()
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Println("get", subject)
			select {
			case channel <- m.Data:
			default:
			}
		}()
	}
	r.subscribeQueue(subject, handlerErrLog)
}

func (r *repository) SubscribeNewsGet() {
r.SubscribeNews("news_get_resp", r.newsGet)
}

func (r *repository) SubscribeNewsCreate() {
	r.SubscribeNews("news_create_resp", r.newsCreate)
}

func (r *repository) GetNewsHandler() {
	for news := range r.newsGet {
		var resp models.Response
		err := proto.Unmarshal(news, &resp)
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		if value, ok := r.waitGet.checkSing(resp.Request); value > 0 && ok {
			r.respGet.addWait(resp.Request, resp.Response)
		} else if value == 0 {
			r.waitGet.deleteWait(resp.Request)
		}
	}
}

func (r *repository) CreateNewsHandler() {
	for news := range r.newsCreate {
		var resp models.Response
		err := proto.Unmarshal(news, &resp)
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		if value, ok := r.waitCreate.checkSing(resp.Request); value > 0 && ok {
			r.respCreate.addWait(resp.Request, resp.Response)
		} else if value == 0 {
			r.waitGet.deleteWait(resp.Request)
		}
	}
}
