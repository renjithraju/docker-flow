package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type ConsulTestSuite struct {
	suite.Suite
	Server           *httptest.Server
	ConsulScale      int
	ServiceName      string
	ServiceColor     string
	PutScaleResponse string
	PutColorResponse string
}

func (s *ConsulTestSuite) SetupTest() {
	s.ConsulScale = 4
	s.ServiceName = "myService"
	s.ServiceColor = BlueColor
	s.PutScaleResponse = "PUT_SCALE"
	s.PutColorResponse = "PUT_COLOR"
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scaleGetUrl := fmt.Sprintf("/v1/kv/docker-flow/%s/scale?raw", s.ServiceName)
		colorGetUrl := fmt.Sprintf("/v1/kv/docker-flow/%s/color?raw", s.ServiceName)
		scalePutUrl := fmt.Sprintf("/v1/kv/docker-flow/%s/scale?", s.ServiceName)
		colorPutUrl := fmt.Sprintf("/v1/kv/docker-flow/%s/color?", s.ServiceName)
		actualUrl := fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery)
		if r.Method == "GET" {
			if actualUrl == scaleGetUrl {
				fmt.Fprint(w, s.ConsulScale)
			} else if actualUrl == colorGetUrl {
				fmt.Fprint(w, s.ServiceColor)
			} else {
				fmt.Fprint(w, "")
			}
		} else if r.Method == "PUT" {
			if actualUrl == scalePutUrl {
				fmt.Fprint(w, s.PutScaleResponse)
			}
			if actualUrl == colorPutUrl {
				fmt.Fprint(w, s.PutColorResponse)
			}
		}
	}))
}

func (s ConsulTestSuite) Test_GetScaleCalc_Returns1() {
	actual, _ := Consul{}.GetScaleCalc(s.Server.URL, "SERVICE_NEVER_DEPLOYED_BEFORE", "")

	s.Equal(1, actual)
}

func (s ConsulTestSuite) Test_GetScaleCalc_ReturnsNumberFromConsul() {
	actual, _ := Consul{}.GetScaleCalc(s.Server.URL, s.ServiceName, "")

	s.Equal(s.ConsulScale, actual)
}

func (s ConsulTestSuite) Test_GetScaleCalc_ReturnsErrorFromHttpGet() {
	_, err := Consul{}.GetScaleCalc("WRONG_URL", s.ServiceName, "")

	s.Error(err)
}

func (s ConsulTestSuite) Test_GetScaleCalc_ReturnScaleFuncArg() {
	expected := 7

	actual, _ := Consul{}.GetScaleCalc(s.Server.URL, s.ServiceName, strconv.Itoa(expected))

	s.Equal(expected, actual)
}

func (suite ConsulTestSuite) Test_GetScaleCalc_IncrementsScale() {
	actual, _ := Consul{}.GetScaleCalc(suite.Server.URL, suite.ServiceName, "+2")

	suite.Equal(suite.ConsulScale+2, actual)
}

func (suite ConsulTestSuite) Test_GetScaleCalc_DecrementsScale() {
	actual, _ := Consul{}.GetScaleCalc(suite.Server.URL, suite.ServiceName, "-2")

	suite.Equal(suite.ConsulScale-2, actual)
}

func (suite ConsulTestSuite) Test_GetScaleCalc_Returns1_WhenScaleIsNegativeOrZero() {
	actual, _ := Consul{}.GetScaleCalc(suite.Server.URL, suite.ServiceName, "-100")

	suite.Equal(1, actual)
}

func (suite ConsulTestSuite) Test_GetColor_ReturnsGreen() {
	actual, _ := Consul{}.GetColor(suite.Server.URL, "SERVICE_NEVER_DEPLOYED_BEFORE")

	suite.Equal(GreenColor, actual)
}

func (suite ConsulTestSuite) Test_GetColor_ReturnServiceColor() {
	actual, _ := Consul{}.GetColor(suite.Server.URL, suite.ServiceName)

	suite.Equal(suite.ServiceColor, actual)
}

func (suite ConsulTestSuite) Test_GetColor_ReturnsErrorFromHttpGet() {
	_, err := Consul{}.GetColor("WRONG_URL", suite.ServiceName)

	suite.Error(err)
}

func (suite ConsulTestSuite) Test_GetNextColor_ReturnsBlueWhenGreen() {
	actual := Consul{}.GetNextColor(GreenColor)

	suite.Equal(BlueColor, actual)
}

func (suite ConsulTestSuite) Test_GetNextColor_ReturnsGreenWhenBlue() {
	actual := Consul{}.GetNextColor(BlueColor)

	suite.Equal(GreenColor, actual)
}

func (suite ConsulTestSuite) Test_PutScale_PutsToConsul() {
	actual, _ := Consul{}.PutScale(suite.Server.URL, suite.ServiceName, 34)

	suite.Equal(suite.PutScaleResponse, actual)
}

func (suite ConsulTestSuite) Test_GetScaleCalc_ReturnsErrorFromHttpPut() {
	_, err := Consul{}.PutScale("WRONG_URL", suite.ServiceName, 45)

	suite.Error(err)
}

func (suite ConsulTestSuite) Test_PutColor_PutsToConsul() {
	actual, _ := Consul{}.PutColor(suite.Server.URL, suite.ServiceName, "orange")

	suite.Equal(suite.PutColorResponse, actual)
}

func (suite ConsulTestSuite) Test_GetColorCalc_ReturnsErrorFromHttpPut() {
	_, err := Consul{}.PutColor("WRONG_URL", suite.ServiceName, "purple")

	suite.Error(err)
}

func TestConsulTestSuite(t *testing.T) {
	dockerHost := os.Getenv("DOCKER_HOST")
	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	defer func() {
		os.Setenv("DOCKER_HOST", dockerHost)
		os.Setenv("DOCKER_CERT_PATH", dockerCertPath)
	}()
	suite.Run(t, new(ConsulTestSuite))
}
