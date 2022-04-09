package jwt

import "testing"

const ExpireSecond = 3600 * 7
const SecretKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDFiErVWa8++SQefR0CiLC4VCtj\n7gMY4ylQ5WR2gvI/nJ10hsyooNhaqxD1pn5ovXsnIjM3dSU0AJPFSV14abh34HXN\nZ4TktMMs8oAPOeq5nrXyG8g2zjtYOiu6e43WAQnfYNGQ+SFSkZiYB2V1e6YRuk5C\nAh7XxHb5VQbnvEaiFQIDAQAB"

func TestJwt(t *testing.T) {
	jwt := New(SecretKey, ExpireSecond)
	token, _, _ := jwt.GenerateToken(map[string]interface{}{"user_id": "1", "roles": []string{"admin"}})
	t.Log(token)
	if claim, err := jwt.ParseToken(token); err != nil {
		t.Fatal(err)
	} else {
		t.Log(claim.Payload)
	}
}
