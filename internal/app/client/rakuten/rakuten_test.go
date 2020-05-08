package rakuten

import (
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	mock_client "github.com/YumaUeno123/linebot_go/internal/testsupport/mock/client"
	"github.com/golang/mock/gomock"

	"github.com/YumaUeno123/linebot_go/internal/app/client"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const keyword = "mock"

var mockUrl = client.UrlScheme + "://" + urlHost + "/" + urlPath + "?applicationId=" + os.Getenv(applicationID) + "&format=" + format + "&keyword=" + keyword

var noDataResp = strings.NewReader(`
	{"Items": []}
`)

var under10itemsResp = strings.NewReader(`
	{
		"Items": [
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_mock"}],"itemName": "mock","itemPrice": 1000,"itemUrl": "url_mock"}}
		]
	}
`)

var over10itemsResp = strings.NewReader(`
	{
		"Items": [
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a1"}],"itemName": "a1","itemPrice": 1000,"itemUrl": "url_a1"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a2"}],"itemName": "a2","itemPrice": 2000,"itemUrl": "url_a2"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a3"}],"itemName": "a3","itemPrice": 3000,"itemUrl": "url_a3"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a4"}],"itemName": "a4","itemPrice": 4000,"itemUrl": "url_a4"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a5"}],"itemName": "a5","itemPrice": 5000,"itemUrl": "url_a5"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a6"}],"itemName": "a6","itemPrice": 6000,"itemUrl": "url_a6"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a7"}],"itemName": "a7","itemPrice": 7000,"itemUrl": "url_a7"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a8"}],"itemName": "a8","itemPrice": 8000,"itemUrl": "url_a8"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a9"}],"itemName": "a9","itemPrice": 9000,"itemUrl": "url_a9"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a10"}],"itemName": "a10","itemPrice": 10000,"itemUrl": "url_a10"}},
			{"Item": {"mediumImageUrls": [{"imageUrl": "img_a11"}],"itemName": "a11","itemPrice": 11000,"itemUrl": "url_a11"}}
		]
	}
`)

func Test_rakuten_Fetch(t *testing.T) {
	type fields struct {
		kind      string
		apiClient client.Api
	}
	type args struct {
		keyword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		expect  func(*mock_client.MockApi)
		want    *[]linebot.Response
		wantErr bool
	}{
		{
			name: "正常系 検索結果が 0 件",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(noDataResp)}, nil)
			},
			want:    &[]linebot.Response{},
			wantErr: false,
		},
		{
			name: "正常系 検索結果が 10 件以下",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(under10itemsResp)}, nil)
			},
			want: &[]linebot.Response{
				{Title: "mock", Image: "img_mock", Price: "1,000円", LinkURL: "url_mock"},
			},
			wantErr: false,
		},
		{
			name: "正常系 検索結果が 11 件以上",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(over10itemsResp)}, nil)
			},
			want: &[]linebot.Response{
				{Title: "a1", Image: "img_a1", Price: "1,000円", LinkURL: "url_a1"},
				{Title: "a2", Image: "img_a2", Price: "2,000円", LinkURL: "url_a2"},
				{Title: "a3", Image: "img_a3", Price: "3,000円", LinkURL: "url_a3"},
				{Title: "a4", Image: "img_a4", Price: "4,000円", LinkURL: "url_a4"},
				{Title: "a5", Image: "img_a5", Price: "5,000円", LinkURL: "url_a5"},
				{Title: "a6", Image: "img_a6", Price: "6,000円", LinkURL: "url_a6"},
				{Title: "a7", Image: "img_a7", Price: "7,000円", LinkURL: "url_a7"},
				{Title: "a8", Image: "img_a8", Price: "8,000円", LinkURL: "url_a8"},
				{Title: "a9", Image: "img_a9", Price: "9,000円", LinkURL: "url_a9"},
				{Title: "a10", Image: "img_a10", Price: "10,000円", LinkURL: "url_a10"},
			},
			wantErr: false,
		},
		{
			name: "異常系 ステータスコードが 200 ではない場合",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusNotFound}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系 json newDecoder が失敗する場合",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(``))}, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_client.NewMockApi(ctrl)
			tt.expect(mock)

			r := &rakuten{
				kind:      tt.fields.kind,
				apiClient: mock,
			}
			got, err := r.Fetch(tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
