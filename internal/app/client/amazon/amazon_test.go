package amazon

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/YumaUeno123/linebot_go/internal/app/client"
	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
	mock_client "github.com/YumaUeno123/linebot_go/internal/testsupport/mock/client"
	"github.com/golang/mock/gomock"
)

const keyword = "mock"

const mockUrl = client.UrlScheme + "://" + urlHost + "/s?k=" + keyword

var under10ItemsResp = strings.NewReader(`
	<div class="s-expand-height">
		<div class="a-section">
			<img src="img_mock"><a class="a-link-normal" href="/url_mock"><h2><a><span>mock</span></a></h2><div class="a-price-whole">1,000</div></a>
		</div>
	</div>
`)

var over10itemsResp = strings.NewReader(`
	<div class="s-expand-height">
		<div class="a-section"><img src="img_a1"><a class="a-link-normal" href="/url_a1"><h2><a><span>a1</span></a></h2><div class="a-price-whole">1,000</div></a></div>
		<div class="a-section"><img src="img_a2"><a class="a-link-normal" href="/url_a2"><h2><a><span>a2</span></a></h2><div class="a-price-whole">2,000</div></a></div>
		<div class="a-section"><img src="img_a3"><a class="a-link-normal" href="/url_a3"><h2><a><span>a3</span></a></h2><div class="a-price-whole">3,000</div></a></div>
		<div class="a-section"><img src="img_a4"><a class="a-link-normal" href="/url_a4"><h2><a><span>a4</span></a></h2><div class="a-price-whole">4,000</div></a></div>
		<div class="a-section"><img src="img_a5"><a class="a-link-normal" href="/url_a5"><h2><a><span>a5</span></a></h2><div class="a-price-whole">5,000</div></a></div>
		<div class="a-section"><img src="img_a6"><a class="a-link-normal" href="/url_a6"><h2><a><span>a6</span></a></h2><div class="a-price-whole">6,000</div></a></div>
		<div class="a-section"><img src="img_a7"><a class="a-link-normal" href="/url_a7"><h2><a><span>a7</span></a></h2><div class="a-price-whole">7,000</div></a></div>
		<div class="a-section"><img src="img_a8"><a class="a-link-normal" href="/url_a8"><h2><a><span>a8</span></a></h2><div class="a-price-whole">8,000</div></a></div>
		<div class="a-section"><img src="img_a9"><a class="a-link-normal" href="/url_a9"><h2><a><span>a9</span></a></h2><div class="a-price-whole">9,000</div></a></div>
		<div class="a-section"><img src="img_a10"><a class="a-link-normal" href="/url_a10"><h2><a><span>a10</span></a></h2><div class="a-price-whole">10,000</div></a></div>
		<div class="a-section"><img src="img_a11"><a class="a-link-normal" href="/url_a11"><h2><a><span>a11</span></a></h2><div class="a-price-whole">11,000</div></a></div>
	</div>
`)

var notPriceItemResp = strings.NewReader(`
	<div class="s-expand-height">
		<div class="a-section">
			<img src="img_mock"><a class="a-link-normal" href="/url_mock"><h2><a><span>mock</span></a></h2><div class="a-price-whole"></div></a>
		</div>
	</div>
`)

var noImgSrcResp = strings.NewReader(`
	<div class="s-expand-height">
		<div class="a-section">
			<img><a class="a-link-normal" href="/url_mock"><h2><a><span>mock</span></a></h2><div class="a-price-whole">1,000</div></a>
		</div>
	</div>
`)
var noHrefResp = strings.NewReader(`
	<div class="s-expand-height">
		<div class="a-section">
			<img src="img_mock"><a class="a-link-normal"><h2><a><span>mock</span></a></h2><div class="a-price-whole">1,000</div></a>
		</div>
	</div>
`)

func Test_amazon_Fetch(t *testing.T) {
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
			name: "正常系　検索結果が 10 件以下",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(under10ItemsResp)}, nil)
			},
			want: &[]linebot.Response{
				{Title: "mock", Image: "img_mock", Price: "1,000円", LinkURL: "https://www.amazon.co.jp/url_mock"},
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
				{Title: "a1", Image: "img_a1", Price: "1,000円", LinkURL: "https://www.amazon.co.jp/url_a1"},
				{Title: "a2", Image: "img_a2", Price: "2,000円", LinkURL: "https://www.amazon.co.jp/url_a2"},
				{Title: "a3", Image: "img_a3", Price: "3,000円", LinkURL: "https://www.amazon.co.jp/url_a3"},
				{Title: "a4", Image: "img_a4", Price: "4,000円", LinkURL: "https://www.amazon.co.jp/url_a4"},
				{Title: "a5", Image: "img_a5", Price: "5,000円", LinkURL: "https://www.amazon.co.jp/url_a5"},
				{Title: "a6", Image: "img_a6", Price: "6,000円", LinkURL: "https://www.amazon.co.jp/url_a6"},
				{Title: "a7", Image: "img_a7", Price: "7,000円", LinkURL: "https://www.amazon.co.jp/url_a7"},
				{Title: "a8", Image: "img_a8", Price: "8,000円", LinkURL: "https://www.amazon.co.jp/url_a8"},
				{Title: "a9", Image: "img_a9", Price: "9,000円", LinkURL: "https://www.amazon.co.jp/url_a9"},
				{Title: "a10", Image: "img_a10", Price: "10,000円", LinkURL: "https://www.amazon.co.jp/url_a10"},
			},
			wantErr: false,
		},
		{
			name: "正常系　検索結果の price が空",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(notPriceItemResp)}, nil)
			},
			want: &[]linebot.Response{
				{Title: "mock", Image: "img_mock", Price: "商品ページからご確認ください", LinkURL: "https://www.amazon.co.jp/url_mock"},
			},
			wantErr: false,
		},
		{
			name: "正常系　検索結果の img タグの src が空",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(noImgSrcResp)}, nil)
			},
			want: &[]linebot.Response{
				{Title: "mock", Image: defaultImg, Price: "1,000円", LinkURL: "https://www.amazon.co.jp/url_mock"},
			},
			wantErr: false,
		},
		{
			name: "正常系　検索結果の商品リンクの href が空",
			args: args{keyword},
			expect: func(mock *mock_client.MockApi) {
				mock.EXPECT().Get(mockUrl).Return(&http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(noHrefResp)}, nil)
			},
			want: &[]linebot.Response{
				{Title: "mock", Image: "img_mock", Price: "1,000円", LinkURL: ""},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_client.NewMockApi(ctrl)
			tt.expect(mock)
			a := &amazon{
				kind:      tt.fields.kind,
				apiClient: mock,
			}
			got, err := a.Fetch(tt.args.keyword)
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
