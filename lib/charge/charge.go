package charge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Nekodigi/gpt-playground-backend/config"
	"github.com/Nekodigi/gpt-playground-backend/consts"
	"github.com/gin-gonic/gin"
)

var c *Charge

type (
	Charge struct {
		url       string
		serviceId string
	}

	StatusRes struct {
		Status string `json:"status"`
	}

	CheckQuotaRes struct {
		AllocQuota  float64 `json:"allocQuota"`
		RemainQuota float64 `json:"remainQuota"`
		Status      string  `json:"status"`
	}

	UrlRes struct {
		Url string `json:"url"`
	}
)

func InitCharge(conf *config.Config) *Charge {
	c = &Charge{
		url:       conf.ChargeApiUrl,
		serviceId: conf.ServiceId,
	}
	return c
}

func (c *Charge) UseQuota(ctx *gin.Context, userId string, amount float64) bool {
	fmt.Println(c, userId, amount)

	fmt.Println(fmt.Sprintf("%s/quota/use/%s/%s?amount=%f", c.url, c.serviceId, userId, amount))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/quota/use/%s/%s?amount=%f", c.url, c.serviceId, userId, amount), nil)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	client := &http.Client{}
	resp_, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp_.Body.Close()
	body, _ := io.ReadAll(resp_.Body)
	var resp StatusRes
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("Error unmarshal %+v\n", err)
	}
	if resp_.Status != "200 OK" {
		fmt.Printf("Error using quota: %s\n", resp_.Status)
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}
	if resp.Status != consts.OK {
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}
	return true
}

func (c *Charge) EnsureQuota(ctx *gin.Context, userId string, amount float64) bool {
	fmt.Println(c, userId, amount)

	fmt.Println(fmt.Sprintf("%s/quota/%s/%s", c.url, c.serviceId, userId))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/quota/%s/%s", c.url, c.serviceId, userId), nil)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	client := &http.Client{}
	resp_, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp_.Body.Close()
	body, _ := io.ReadAll(resp_.Body)
	var resp CheckQuotaRes
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Printf("Error unmarshal %+v\n", err)
	}
	if resp_.Status != "200 OK" {
		fmt.Printf("Error checking quota: %s\n", resp_.Status)
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}
	if resp.Status != consts.OK {
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}
	if resp.RemainQuota < amount {
		ctx.JSON(http.StatusBadRequest, resp)
		return false
	}
	return true
}
