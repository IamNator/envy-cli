package encryption

import "testing"

func TestEncryptionAndDecryptionAdapter(t *testing.T) {

	tt := []struct {
		Name   string
		Secret string
		Value  string
	}{
		{
			Name:   "simple",
			Secret: "secret",
			Value:  "value",
		},
		{
			Name:   "empty secret",
			Secret: "",
			Value:  "3454$%$&%63234",
		},
		{
			Name:   "empty value",
			Secret: "heavy_secret", //
			Value:  "",
		},
		{
			Name:   "empty secret and value",
			Secret: "",
			Value:  "",
		},
		{
			Name:   "long secret",
			Secret: "https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/09.6.html",
			Value:  "lgorithm in golang - the caesar cipher. Is there something i could do more efficiently? I am quite new to go and any improveme",
		},
		{
			Name:   "Url with special characters",
			Secret: "d15jidNhuBYKGHJyh0E",
			Value:  "https://www.google.com/search?q=go+url+encode&oq=go+url+encode&aqs=chrome..69i57j0l5.1001j0j7&sourceid=chrome&ie=UTF-8",
		},
	}

	for _, tc := range tt {

		enc := MakeEncrypter(tc.Secret)
		dec := MakeDecrypter(tc.Secret)

		encrypted, err := enc(tc.Value)
		if err != nil {
			t.Fatal(err)
		}

		decrypted, err := dec(encrypted)
		if err != nil {
			t.Fatal(err)
		}

		if decrypted != tc.Value {
			t.Fatalf("expected %s, got %s", tc.Value, decrypted)
		}

	}

}
