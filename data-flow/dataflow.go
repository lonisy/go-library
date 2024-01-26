package dataflow

import (
	"context"
	"go-library/app"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	INPUT_STAGE   = "input"
	PROCESS_STAGE = "process"
)

type DataProcessorFunc func(ctx context.Context, sourceChan chan interface{}, outputChan chan interface{})

type DataSourceFunc func(ctx context.Context, outputChan chan<- interface{}, ticker *time.Ticker)

type Stage struct {
	Stages            []*Stage
	StageType         string
	DataProcessorFunc DataProcessorFunc
	DataSourceFunc    DataSourceFunc
	Cancel            context.CancelFunc
	Ctx               context.Context
	DataChannel       chan interface{}
	Gc                *GoroutineCounter
	Wg                sync.WaitGroup
	RateLimitPerSec   int
	Workers           int
	OutputChanSize    int
	DataSourceTicker  *time.Ticker
}

func NewDataFlow(rateLimitPerSec, bufferSize int) *Stage {
	ctx, cancel := context.WithCancel(context.Background())
	return &Stage{
		Gc:               &GoroutineCounter{},
		Ctx:              ctx,
		Cancel:           cancel,
		RateLimitPerSec:  rateLimitPerSec,
		DataChannel:      make(chan interface{}, bufferSize),
		DataSourceTicker: time.NewTicker(time.Second / time.Duration(rateLimitPerSec)),
	}
}

func (s *Stage) RegisterDataSource(callback DataSourceFunc, workers int) {
	e := s.initializeStage(INPUT_STAGE, workers)
	e.DataSourceFunc = callback
	s.Stages = append(s.Stages, e)
}

func (s *Stage) RegisterDataProcessor(callback DataProcessorFunc, workers int, chanSize int) {
	e := s.initializeStage(PROCESS_STAGE, workers)
	e.DataProcessorFunc = callback
	e.OutputChanSize = chanSize
	e.DataChannel = make(chan interface{}, chanSize)
	s.Stages = append(s.Stages, e)
}

func (s *Stage) initializeStage(stageType string, workers int) *Stage {
	e := new(Stage)
	e.StageType = stageType
	e.Workers = workers
	e.Gc = &GoroutineCounter{}
	e.Gc.Add(workers)
	e.Wg.Add(workers)
	s.Wg.Add(workers)
	return e
}

func (s *Stage) Run() {
	lastStage := s // 初始化 lastStage 为当前 Stage
	for _, stage := range s.Stages {
		switch stage.StageType {
		case INPUT_STAGE:
			s.runInputStage(stage)
		case PROCESS_STAGE:
			s.runProcessStage(stage, lastStage)
			lastStage = stage // 更新 lastStage 为当前处理阶段
		}
	}
}

func (s *Stage) runInputStage(stage *Stage) {
	for i := 0; i < stage.Workers; i++ {
		go func(st *Stage) {
			defer s.Wg.Done()
			defer st.Wg.Done()
			defer st.Gc.Done()
			st.DataSourceFunc(s.Ctx, s.DataChannel, s.DataSourceTicker)
		}(stage)
	}
}

func (s *Stage) runProcessStage(stage, lastStage *Stage) {
	for i := 0; i < stage.Workers; i++ {
		go func(st, lastSt *Stage) {
			defer s.Wg.Done()
			defer st.Wg.Done()
			defer st.closeChannel()
			st.DataProcessorFunc(s.Ctx, lastSt.DataChannel, st.DataChannel)
		}(stage, lastStage)
	}
}

func (s *Stage) Listen() {
	s.Wg.Add(1)
	go func() {
		defer s.Wg.Done()
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
		for si := range c {
			switch si {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Service Stoping By signal:", si)
				s.Stop()
				return
			case syscall.SIGUSR1:
				log.Println("usr1", s)
			case syscall.SIGUSR2:
				log.Println("usr2", s)
			default:
				log.Println("other", s)
			}
		}
	}()
	s.Wg.Wait()
}

func (s *Stage) closeChannel() {
	s.Gc.Done()
	if s.Gc.Count() == 0 {
		close(s.DataChannel)
	}
}

func (s *Stage) Stop() {
	app.Log.Info("Stopping...")
	s.Cancel()
	s.DataSourceTicker.Stop()
	for idx, stage := range s.Stages {
		if stage.StageType == INPUT_STAGE {
			app.Log.Info("stage", stage.StageType, "close", idx)
			stage.Wg.Wait()
			app.Log.Info("stage", stage.StageType, "closed", idx)
		}
	}
	app.Log.Info("Stopped...")
	close(s.DataChannel)
}
