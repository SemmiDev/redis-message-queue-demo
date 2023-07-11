package rmq_test

import (
	rmq "redis-message-queue"
	"testing"
)

func TestGenerateUTMURL(t *testing.T) {
	type args struct {
		baseURL  string
		source   string
		medium   string
		campaign string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should generate a valid URL",
			args: args{
				baseURL:  "https://www.example.com",
				source:   "newsletter",
				medium:   "email",
				campaign: "new post",
			},
			want:    "https://www.example.com?utm_campaign=new+post&utm_medium=email&utm_source=newsletter",
			wantErr: false,
		},
		{
			name: "should return an error if the base URL is invalid",
			args: args{
				baseURL:  "://www.example.com",
				source:   "newsletter",
				medium:   "email",
				campaign: "new post",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rmq.GenerateUTMURL(tt.args.baseURL, tt.args.source, tt.args.medium, tt.args.campaign)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUTMURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateUTMURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEventFromParams(t *testing.T) {
	type args struct {
		queryParams string
	}
	tests := []struct {
		name string
		args args
		want rmq.EventType
	}{
		{
			name: "should return LinkClick if the utm_medium is newsletter",
			args: args{
				queryParams: "utm_medium=newsletter",
			},
			want: rmq.LinkClick,
		},
		{
			name: "should return Read if the utm_medium is not newsletter",
			args: args{
				queryParams: "utm_medium=not-newsletter",
			},
			want: rmq.Read,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rmq.GetEventFromParams(tt.args.queryParams); got != tt.want {
				t.Errorf("GetEventFromParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
