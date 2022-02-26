package anubis

import "testing"

func TestHeaderOpt_SetOpt(t *testing.T) {
	type fields struct {
		Key   string
		Value string
	}
	type args struct {
		anubis *Anubis
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
			opt := HeaderOpt{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			opt.SetOpt(tt.args.anubis)
		})
	}
}

func TestNWorkerOpt_SetOpt(t *testing.T) {
	type args struct {
		anubis *Anubis
	}
	tests := []struct {
		name string
		opt  NWorkerOpt
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opt.SetOpt(tt.args.anubis)
		})
	}
}

func TestOutputOpt_SetOpt(t *testing.T) {
	type args struct {
		anubis *Anubis
	}
	tests := []struct {
		name string
		opt  OutputOpt
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opt.SetOpt(tt.args.anubis)
		})
	}
}

func TestProxyOpt_SetOpt(t *testing.T) {
	type args struct {
		anubis *Anubis
	}
	tests := []struct {
		name string
		opt  ProxyOpt
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opt.SetOpt(tt.args.anubis)
		})
	}
}

func TestResponseHandlerOpt_SetOpt(t *testing.T) {
	type fields struct {
		Handler ResponseHandler
	}
	type args struct {
		anubis *Anubis
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
			opt := ResponseHandlerOpt{
				Handler: tt.fields.Handler,
			}
			opt.SetOpt(tt.args.anubis)
		})
	}
}

func TestWebDriverOpt_SetOpt(t *testing.T) {
	type fields struct {
		Driver WebDriver
	}
	type args struct {
		anubis *Anubis
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
			opt := WebDriverOpt{
				Driver: tt.fields.Driver,
			}
			opt.SetOpt(tt.args.anubis)
		})
	}
}
