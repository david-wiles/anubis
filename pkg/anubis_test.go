package anubis

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestAnubis_AddURL(t *testing.T) {
	type fields struct {
		Output  string
		Workers int
		Headers map[string]string
		Driver  WebDriver
		Handler ResponseHandler
		wg      *sync.WaitGroup
		queue   chan string
		Cancel  context.CancelFunc
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Anubis{
				Output:  tt.fields.Output,
				Workers: tt.fields.Workers,
				Headers: tt.fields.Headers,
				Driver:  tt.fields.Driver,
				Handler: tt.fields.Handler,
				wg:      tt.fields.wg,
				queue:   tt.fields.queue,
				Cancel:  tt.fields.Cancel,
			}
			if got := a.AddURL(tt.args.url); got != tt.want {
				t.Errorf("AddURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnubis_Commit(t *testing.T) {
	type fields struct {
		Output  string
		Workers int
		Headers map[string]string
		Driver  WebDriver
		Handler ResponseHandler
		wg      *sync.WaitGroup
		queue   chan string
		Cancel  context.CancelFunc
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Anubis{
				Output:  tt.fields.Output,
				Workers: tt.fields.Workers,
				Headers: tt.fields.Headers,
				Driver:  tt.fields.Driver,
				Handler: tt.fields.Handler,
				wg:      tt.fields.wg,
				queue:   tt.fields.queue,
				Cancel:  tt.fields.Cancel,
			}
			if err := a.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnubis_Start(t *testing.T) {
	type fields struct {
		Output  string
		Workers int
		Headers map[string]string
		Driver  WebDriver
		Handler ResponseHandler
		wg      *sync.WaitGroup
		queue   chan string
		Cancel  context.CancelFunc
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Anubis{
				Output:  tt.fields.Output,
				Workers: tt.fields.Workers,
				Headers: tt.fields.Headers,
				Driver:  tt.fields.Driver,
				Handler: tt.fields.Handler,
				wg:      tt.fields.wg,
				queue:   tt.fields.queue,
				Cancel:  tt.fields.Cancel,
			}
			if err := a.Start(); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAnubis_worker(t *testing.T) {
	type fields struct {
		Output  string
		Workers int
		Headers map[string]string
		Driver  WebDriver
		Handler ResponseHandler
		wg      *sync.WaitGroup
		queue   chan string
		Cancel  context.CancelFunc
	}
	type args struct {
		ctx   context.Context
		queue chan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Anubis{
				Output:  tt.fields.Output,
				Workers: tt.fields.Workers,
				Headers: tt.fields.Headers,
				Driver:  tt.fields.Driver,
				Handler: tt.fields.Handler,
				wg:      tt.fields.wg,
				queue:   tt.fields.queue,
				Cancel:  tt.fields.Cancel,
			}
			a.worker(tt.args.ctx, tt.args.queue)
		})
	}
}

func TestNewAnubis(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *Anubis
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAnubis(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnubis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processURL(t *testing.T) {
	type args struct {
		url       string
		headers   map[string]string
		webdriver WebDriver
		handler   ResponseHandler
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processURL(tt.args.url, tt.args.headers, tt.args.webdriver, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("processURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
