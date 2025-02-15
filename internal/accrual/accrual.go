package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/JohnRobertFord/go-market/internal/config"
	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/go-resty/resty/v2"
)

type Storage interface {
	GetUnprocessedOrders(ctx context.Context) ([]model.Order, error)
	UpdateOrder(ctx context.Context, order model.Order) error
}

type Service struct {
	config   *config.Config
	storage  Storage
	register chan model.Order
	results  chan model.Accrual
	client   *resty.Client
}

func NewWorker(accrual string, storage Storage) *Service {
	return &Service{
		storage:  storage,
		register: make(chan model.Order),
		results:  make(chan model.Accrual),
		client:   resty.New().SetBaseURL(accrual),
	}
}

func (s *Service) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go s.accrualWorker(ctx, &wg)

	wg.Add(1)
	go s.updateAccrualWorker(ctx, &wg)

	if err := s.handleUnprocessedOrders(ctx); err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (s *Service) accrualWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-s.register:
			var finalState bool
			for !finalState {
				accrual, err := s.handleOrder(ctx, order)
				if err != nil {
					break
				}
				s.results <- *accrual
				finalState = accrual.Status == "PROCESSED" || accrual.Status == "INVALID"
			}
		}
	}
}

func (s *Service) handleOrder(ctx context.Context, order model.Order) (*model.Accrual, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			accrual, err := s.getAccrual("/api/orders/" + order.OrderID)
			if err == nil {
				accrual.Username = order.Username
				return accrual, nil
			}
			switch {
			case errors.Is(err, model.ErrRequestToAccrualService):
				log.Println(err)
				return &model.Accrual{
					Order:    order.OrderID,
					Status:   "INVALID",
					Accrual:  0,
					Username: order.Username,
				}, err
			case errors.Is(err, model.ErrTooManyRequests):
				log.Println(err)
			}
		}
	}
}

func (s *Service) getAccrual(uri string) (accrual *model.Accrual, err error) {
	resp, err := s.client.R().Get(uri)
	if err != nil {
		return nil, model.ErrRequestToAccrualService
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		makeRetryPause(resp)
		return nil, model.ErrTooManyRequests
	}

	if err := json.NewDecoder(resp.RawBody()).Decode(&accrual); err != nil {
		return nil, model.ErrAccrualServiceDecode
	}

	return
}

func makeRetryPause(resp *resty.Response) {
	retryHeader := resp.Header().Get("Retry-After")
	retry, err := strconv.Atoi(retryHeader)
	if err != nil {
		return
	}
	durationStr := fmt.Sprintf("%ds", retry)
	durationRetry, err := time.ParseDuration(durationStr)
	if err != nil {
		return
	}
	time.Sleep(durationRetry)
}

func (s *Service) updateAccrualWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			close(s.results)
			return
		case a := <-s.results:
			if err := s.storage.UpdateOrder(ctx, model.Order{
				OrderID:  a.Order,
				Status:   a.Status,
				Accrual:  a.Accrual,
				Username: a.Username,
			}); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (s *Service) handleUnprocessedOrders(ctx context.Context) error {
	orders, err := s.storage.GetUnprocessedOrders(ctx)
	if err != nil {
		return err
	}
	for _, order := range orders {
		select {
		case <-ctx.Done():
			close(s.register)
			return ctx.Err()
		default:
			s.SendOrder(ctx, order)
		}
	}

	return nil
}

func (s *Service) SendOrder(ctx context.Context, order model.Order) {
	select {
	case <-ctx.Done():
		return
	default:
		s.register <- order
	}
}
