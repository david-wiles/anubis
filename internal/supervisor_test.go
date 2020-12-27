package internal

//
//import (
//	"net/http"
//	"sync"
//	"testing"
//)
//
//func TestSupervisor_Start(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//			if err := s.Start(); (err != nil) != tt.wantErr {
//				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestSupervisor_Terminate(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
//
//func TestSupervisor_buildSeed(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//			if err := s.buildSeed(); (err != nil) != tt.wantErr {
//				t.Errorf("buildSeed() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestSupervisor_checkProgramState(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	type args struct {
//		id int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
//
//func TestSupervisor_manageUrls(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
//
//func TestSupervisor_monitorWorker(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	type args struct {
//		w *worker
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
//
//func TestSupervisor_sendWork(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	type args struct {
//		w *worker
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
//
//func TestSupervisor_shiftQueue(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//		want1  bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//			got, got1 := s.shiftQueue()
//			if got != tt.want {
//				t.Errorf("shiftQueue() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("shiftQueue() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func TestSupervisor_startWorkers(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
//
//func TestSupervisor_writeErrors(t *testing.T) {
//	type fields struct {
//		Pipeline    []PipelineFunc
//		client      *http.Client
//		urlQueue    []string
//		qMutex      *sync.Mutex
//		sent        chan CompletedRequest
//		sentUrls    map[string]int
//		found       chan string
//		errors      chan error
//		workers     []*worker
//		workerState []WorkerState
//		config      *Config
//		done        chan bool
//		logger      *Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Supervisor{
//				Pipeline:    tt.fields.Pipeline,
//				client:      tt.fields.client,
//				urlQueue:    tt.fields.urlQueue,
//				qMutex:      tt.fields.qMutex,
//				sent:        tt.fields.sent,
//				sentUrls:    tt.fields.sentUrls,
//				found:       tt.fields.found,
//				errors:      tt.fields.errors,
//				workers:     tt.fields.workers,
//				workerState: tt.fields.workerState,
//				config:      tt.fields.config,
//				done:        tt.fields.done,
//				logger:      tt.fields.logger,
//			}
//		})
//	}
//}
