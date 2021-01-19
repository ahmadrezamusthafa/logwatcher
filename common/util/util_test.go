package util

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	type args struct {
		plainText string
		key       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Normal",
			args: args{
				plainText: "test",
				key:       "Tiqweku6Casdlasska26ajAS2sSs2cE3",
			},
			want: "test",
		},
		{
			name: "Normal",
			args: args{
				plainText: "data 1625 halo",
				key:       "Miqweku6Casdlasska26ajAS2sSs2cE3",
			},
			want: "data 1625 halo",
		},
		{
			name: "Normal",
			args: args{
				plainText: "tes123@yahuui.com",
				key:       "Biaqr32qrdlassska26ajAsd2sSs2cE3",
			},
			want: "tes123@yahuui.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedText := Encrypt(tt.args.plainText, tt.args.key)
			decryptedText := Decrypt(encryptedText, tt.args.key)
			if decryptedText != tt.want {
				t.Errorf("EncryptDecrypt() = %v, want %v", decryptedText, tt.want)
			}
		})
	}
}
