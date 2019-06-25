package token_test

import (
	. "gin_project_starter/src/utils/token"
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("Token", func() {
	Context("test create token", func() {
		It("should success", func() {
			token, err := Create("itsneo1990@gmail.com")
			Expect(err).To(BeNil())
			Expect(Verify(token)).To(BeTrue())
			Expect(GetAccountEmail(token)).To(Equal("itsneo1990@gmail.com"))
		})
	})

	Context("test refresh token", func() {
		It("should success", func() {
			token, _ := Create("itsneo1990@gmail.com")
			token, err := Refresh(token)
			Expect(err).To(BeNil())
			Expect(Verify(token)).To(BeTrue())
		})
		It("should fail while invalid token", func() {
			_, err := Refresh("fake")
			Expect(err).NotTo(BeNil())
		})
		It("should fail while token verified failed", func() {
			claims := JWTClaims{
				Email: "itsneo1990@gmail.com",
			}
			claims.ExpiresAt = time.Now().Add(-time.Second).Unix()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			t, _ := token.SignedString([]byte(viper.GetString("JWT.SECRET")))
			_, err := Refresh(t)
			Expect(err).ToNot(BeNil())
		})
	})
	Context("test verify token", func() {
		It("should success", func() {
			token, _ := Create("itsneo1990@gmail.com")
			Expect(Verify(token)).To(BeTrue())
		})
		It("should fail", func() {
			Expect(Verify("fake")).To(BeFalse())
		})
	})
	Context("test get user email", func() {
		It("should success", func() {
			token, _ := Create("itsneo1990@gmail.com")
			email, err := GetAccountEmail(token)
			Expect(email).To(Equal("itsneo1990@gmail.com"))
			Expect(err).To(BeNil())
		})
		It("should get failed", func() {
			_, err := GetAccountEmail("fake")
			Expect(err).ToNot(BeNil())
		})
	})
})
