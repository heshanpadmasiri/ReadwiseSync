package main

import "testing"

func Test_readKeys(t *testing.T) {
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
			name: "Read an existing valid key file",
			args: args { keyFilePath: "./testData/key.gpg" },
			want: "testKey",
			wantErr: false,
		},
		{
			name: "Throw an error if key file don't exists",
			args: args { keyFilePath: "./unexisting.gpg" },
			want: "",
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
