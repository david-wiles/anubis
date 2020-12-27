package internal

//
//import (
//	"net/http"
//	"testing"
//)
//
//func Test_worker_Start(t *testing.T) {
//	type fields struct {
//		id       int
//		queue    chan string
//		errors   chan error
//		state    chan WorkerState
//		sent     chan CompletedRequest
//		found    chan string
//		pipeline []PipelineFunc
//		config   *Config
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			w := &worker{
//				id:       tt.fields.id,
//				queue:    tt.fields.queue,
//				errors:   tt.fields.errors,
//				state:    tt.fields.state,
//				sent:     tt.fields.sent,
//				found:    tt.fields.found,
//				pipeline: tt.fields.pipeline,
//				config:   tt.fields.config,
//			}
//		})
//	}
//}
//
//func Test_worker_runPipeline(t *testing.T) {
//	type fields struct {
//		id       int
//		queue    chan string
//		errors   chan error
//		state    chan WorkerState
//		sent     chan CompletedRequest
//		found    chan string
//		pipeline []PipelineFunc
//		config   *Config
//	}
//	type args struct {
//		r        *http.Response
//		path     string
//		notifier chan bool
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
//			w := &worker{
//				id:       tt.fields.id,
//				queue:    tt.fields.queue,
//				errors:   tt.fields.errors,
//				state:    tt.fields.state,
//				sent:     tt.fields.sent,
//				found:    tt.fields.found,
//				pipeline: tt.fields.pipeline,
//				config:   tt.fields.config,
//			}
//		})
//	}
//}
//
//func Test_worker_writeBytes(t *testing.T) {
//	type fields struct {
//		id       int
//		queue    chan string
//		errors   chan error
//		state    chan WorkerState
//		sent     chan CompletedRequest
//		found    chan string
//		pipeline []PipelineFunc
//		config   *Config
//	}
//	type args struct {
//		b        []byte
//		filename string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			w := &worker{
//				id:       tt.fields.id,
//				queue:    tt.fields.queue,
//				errors:   tt.fields.errors,
//				state:    tt.fields.state,
//				sent:     tt.fields.sent,
//				found:    tt.fields.found,
//				pipeline: tt.fields.pipeline,
//				config:   tt.fields.config,
//			}
//			if err := w.writeBytes(tt.args.b, tt.args.filename); (err != nil) != tt.wantErr {
//				t.Errorf("writeBytes() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
