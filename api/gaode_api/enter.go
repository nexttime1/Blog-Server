package gaode_api

import (
	"Blog_server/common/res"
	"Blog_server/global"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type GaodeApi struct{}

type IPResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

type WeatherResponse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Lives    []struct {
		Province         string `json:"province"`
		City             string `json:"city"`
		Adcode           string `json:"adcode"`
		Weather          string `json:"weather"`
		Temperature      string `json:"temperature"`
		Winddirection    string `json:"winddirection"`
		Windpower        string `json:"windpower"`
		Humidity         string `json:"humidity"`
		Reporttime       string `json:"reporttime"`
		TemperatureFloat string `json:"temperature_float"`
		HumidityFloat    string `json:"humidity_float"`
	} `json:"lives"`
}

type WeatherInfoResponse struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	Adcode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Temperature      string `json:"temperature"`
	Winddirection    string `json:"winddirection"`
	Windpower        string `json:"windpower"`
	Humidity         string `json:"humidity"`
	Reporttime       string `json:"reporttime"`
	TemperatureFloat string `json:"temperature_float"`
	HumidityFloat    string `json:"humidity_float"`
}

// WeatherInfoView 获取实时天气
// @Tags 第三方api管理
// @Summary 获取实时天气
// @Description 获取实时天气
// @Param token header string  true  "token"
// @Router /api/gaode/weather [get]
// @Produce json
// @Success 200 {object} res.Response{data=WeatherInfoResponse}
func (GaodeApi) WeatherInfoView(c *gin.Context) {

	var data WeatherInfoResponse
	if !global.Config.Gaode.Enable {
		res.OkWithData(c, data)
		return
	}

	key := global.Config.Gaode.Key

	response, err := http.Get("https://restapi.amap.com/v3/ip?key=" + key)
	if err != nil {
		logrus.Errorf(err.Error())
		res.FailWithMsg(c, "获取定位失败")
		return
	}

	byteData, _ := io.ReadAll(response.Body)

	var ipResponse IPResponse
	err = json.Unmarshal(byteData, &ipResponse)
	if err != nil {
		res.FailWithMsg(c, "解析定位失败")
		return
	}
	if ipResponse.Status != "1" {
		res.FailWithMsg(c, "请求定位失败")
		return
	}

	response, err = http.Get(fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s", key, ipResponse.Adcode))
	if err != nil {
		logrus.Errorf(err.Error())
		res.FailWithMsg(c, "获取天气信息失败")
		return
	}
	byteData, _ = io.ReadAll(response.Body)

	var weatherResponse WeatherResponse
	err = json.Unmarshal(byteData, &weatherResponse)
	if err != nil {
		res.FailWithMsg(c, "解析天气信息失败")
		return
	}

	if weatherResponse.Status != "1" {
		logrus.Errorf(string(byteData))
		res.FailWithMsg(c, "请求天气信息失败")
		return
	}

	if len(weatherResponse.Lives) > 0 {
		res.OkWithData(c, weatherResponse.Lives[0])
		return
	}
	res.FailWithMsg(c, "获取天气信息异常")
	return
}
