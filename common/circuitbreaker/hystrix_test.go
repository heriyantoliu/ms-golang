package circuitbreaker

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/h2non/gock"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCallUsingResilienceSucceeds(t *testing.T) {
	defer gock.Off()

	buildGockMatcherTimes(500, 2)

	body := []byte("Some response")
	buildGockMatcherWithBody(200, string(body))

	hystrix.Flush()

	Convey("Given a Call request",t, func() {
		Convey("When", func() {
			bytes, err := CallUsingCircuitBreaker("TEST", "http://quotes-service", "GET")
			Convey("Then", func() {
				So(err, ShouldBeNil)
				So(bytes, ShouldNotBeNil)
				So(string(bytes), ShouldEqual, string(body))
			})
		})
	})
}

func buildGockMatcherTimes(status int, times int){
	for a := 0; a < times; a++ {
		buildGockMatcher(status)
	}
}


func buildGockMatcher(status int) {
buildGockMatcherWithBody(status, "")
}

func buildGockMatcherWithBody(status int, body string) {
	gock.New("http://quotes-service").Reply(status).BodyString(body)
}
