package services

import (
	"encoding/json"
	"fmt"
	"gitlab.com/scalablespace/billing/lib/components"
	"log"
	"time"
)

type Scheduler struct {
	Ticker     *time.Ticker
	Repository components.UsageRepository
	Publisher  *Publisher
}

func (s *Scheduler) Start(done chan struct{}) {
	projects, err := s.Repository.CalculateUsagePerProject()
	if err != nil {
		panic(err)
	}
	for _, ps := range projects {
		projectUsagePerHour, err := s.Repository.ProjectUsagePerHour(ps.Id)
		if err != nil {
			panic(err)
		}
		if projectUsagePerHour > 0 && (ps.Credit+ps.Usage+projectUsagePerHour) > 0 && ps.AccountType != "postpaid" {
			log.Println("project", ps.Id, "will be suspended soon!")
			data, err := json.Marshal(ps)
			if err != nil {
				panic(err)
			}
			if err := s.Publisher.Publish(data); err != nil {
				panic(err)
			}
		}
	}
	usagesCountChanged := make(chan struct{})
	go func() {
		previous := -1
		for {
			time.Sleep(time.Minute)
			current, err := s.Repository.UsagesCount()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if previous == -1 {
				previous = current
				continue
			}
			if previous != current {
				usagesCountChanged <- struct{}{}
			}
			previous = current
		}
	}()
	for {
		select {
		case <-s.Ticker.C:
			s.run()
		case <-usagesCountChanged:
			s.run()
		case <-done:
			return
		}
	}
}

func (s *Scheduler) run() {
	projects, err := s.Repository.CalculateUsagePerProject()
	if err != nil {
		panic(err)
	}
	for _, ps := range projects {
		projectUsagePerHour, err := s.Repository.ProjectUsagePerHour(ps.Id)
		if err != nil {
			panic(err)
		}
		if projectUsagePerHour > 0 && (ps.Credit+ps.Usage+projectUsagePerHour) > 0 && ps.AccountType != "postpaid" {
			log.Println("project", ps.Id, "will be suspended soon!")
			data, err := json.Marshal(ps)
			if err != nil {
				panic(err)
			}
			if err := s.Publisher.Publish(data); err != nil {
				panic(err)
			}
		}
	}
}
