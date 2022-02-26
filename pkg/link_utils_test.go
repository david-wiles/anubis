package anubis

import (
	"reflect"
	"testing"
)

func TestGetImageURLs(t *testing.T) {
	type args struct {
		parent string
		html   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageURLs(tt.args.parent, tt.args.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetImageURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLinkURLs(t *testing.T) {
	type args struct {
		parent string
		html   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLinkURLs(tt.args.parent, tt.args.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinkURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetScriptURLs(t *testing.T) {
	type args struct {
		parent string
		html   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetScriptURLs(tt.args.parent, tt.args.html); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScriptURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFullURL(t *testing.T) {
	type args struct {
		parent string
		link   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFullURL(tt.args.parent, tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFullURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFullURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
