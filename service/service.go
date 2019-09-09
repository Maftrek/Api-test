package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/pborman/uuid"
	"template/models"
	"template/repository"
)

// Service interface
type Service interface {
	Close()
	CreateNews(news models.Data) (response []byte, err error)
	GetNews(newsID int32) (news []byte, err error)
	Subscribe()
}

type service struct {
	repNatsStreaming repository.Repository
}

// New return new service
func New(rNatsStreaming repository.Repository) Service {
	return &service{
		repNatsStreaming: rNatsStreaming,
	}
}

func (s *service) Close() {}

func (s *service) CreateNews(news models.Data) (response []byte, err error) {
	news.Uuid = uuid.New()
	dataPb, err := proto.Marshal(&news)
	if err != nil {
		return nil, err
	}
	return s.repNatsStreaming.CreateNews(dataPb)
}

func (s *service) GetNews(newsID int32) (news []byte, err error) {
	modelsID := models.Id{
		Id: newsID,
	}
	dataPb, err := proto.Marshal(&modelsID)
	if err != nil {
		return nil, err
	}
	return s.repNatsStreaming.GetNews(dataPb)
}

func (s *service) Subscribe(){
	s.repNatsStreaming.SubscribeNewsCreate()
	s.repNatsStreaming.SubscribeNewsGet()
	go s.repNatsStreaming.GetNewsHandler()
	go s.repNatsStreaming.CreateNewsHandler()
}