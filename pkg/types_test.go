package anubis

import (
	"net/http"
	"reflect"
	"testing"
)

func TestDefaultResponseHandler_Handle(t *testing.T) {
	type fields struct {
		Anubis      *Anubis
		NeededLinks map[string]bool
	}
	type args struct {
		req  *http.Request
		resp *http.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := DefaultResponseHandler{
				Anubis:      tt.fields.Anubis,
				NeededLinks: tt.fields.NeededLinks,
			}
			if err := handler.Handle(tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultWebDriver_DoRequest(t *testing.T) {
	type fields struct {
		client http.Client
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver := DefaultWebDriver{
				client: tt.fields.client,
			}
			got, err := driver.DoRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
