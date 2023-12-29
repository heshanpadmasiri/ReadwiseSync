package main

import (
	"reflect"
	"testing"
)

func Test_readKeys(t *testing.T) {
	if testing.Short() {
		// our CI can't decrypt gpg keys
		t.Skip("skipping readkeys in short mode")
	}
	type args struct {
		keyFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Read an existing valid key file",
			args:    args{keyFilePath: "./testData/key.gpg"},
			want:    "testKey",
			wantErr: false,
		},
		{
			name:    "Throw an error if key file don't exists",
			args:    args{keyFilePath: "./unexisting.gpg"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readKeys(tt.args.keyFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHightlightRes(t *testing.T) {
	validResponseData := `{"count": 2, "results": [
			{ "title": "title",
			  "readable_title": "readable_title",
			  "source_url": "https://test.com/source",
			  "cover_image_url": "https://test.com/img.png",
			  "category": "articles",
			  "highlights": [
				{
					"text" : "highlight 1",
					"readwise_url": "https:test.com/link1"
				}
			  ]},
			{ "title": "title",
			  "readable_title": "readable_title",
			  "source_url": null,
			  "cover_image_url": "https://test.com/img.png",
			  "category": "books",
			  "highlights": [
				  {
					  "text" : "highlight 1",
					  "readwise_url": "https:test.com/link1"
				  }
			  ]}
	]}`
	type args struct {
		response *[]byte
	}
	tests := []struct {
		name    string
		args    args
		want    *highlightRes
		wantErr bool
	}{
		{
			name: "parse valid response",
			args: args{response: bytePtr(validResponseData)},
			want: &highlightRes{
				Count: 2,
				Sources: []source{
					{
						Readable_title: "readable_title",
						SourceUrl:      stringPtr("https://test.com/source"),
						ImgUrl:         "https://test.com/img.png",
						Category:       article,
						Highlights: []highlight{
							{
								Text: "highlight 1",
								Url:  "https:test.com/link1",
							},
						},
					},
					{
						Readable_title: "readable_title",
						SourceUrl:      nil,
						ImgUrl:         "https://test.com/img.png",
						Category:       book,
						Highlights: []highlight{
							{
								Text: "highlight 1",
								Url:  "https:test.com/link1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHightlightRes(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHightlightRes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHightlightRes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func bytePtr(s string) *[]byte {
	data := []byte(s)
	return &data
}
