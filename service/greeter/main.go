package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-api/proto"
	"github.com/micro/go-micro"
	"go-micro-example/hystrix"
	. "go-micro-example/service/config"
	"go-micro-example/service/constant/micro_c"
	"go-micro-example/service/greeter/dto"
	greeterApi "go-micro-example/service/greeter/proto"
	"go-micro-example/service/greeter/service"
	"go-micro-example/service/user/proto"
	"go-micro-example/service/util"
	"log"
)

type Greeter struct {
	userClient user.UserService
}

func (this *Greeter) Hello(ctx context.Context, req *go_api.Request, rsp *go_api.Response) error {
	log.Println("Received Greeter.Hello API request")
	var helloRequest *dto.HelloRequest
	json.Unmarshal([]byte(req.Body), &helloRequest)
	response, code, err := service.NewGreeterService().Greeter(ctx, this.userClient, helloRequest)
	return util.Resp(code, err, rsp, response)
}

func main() {
	// 配置熔断器
	hystrix.Configure([]string{"go.micro.api.user.User.GetUserInfo"})
	greeterService := micro.NewService(
		micro.Name(micro_c.MicroNameGreeter),
		// 给服务加上装饰器， 用熔断器来装饰服务
		micro.WrapClient(hystrix.NewClientWrapper()),
		micro.Flags(
			cli.StringFlag{
				Name: "profile",
			},
		),
	)

	greeterService.Init(
		micro.Action(func(c *cli.Context) {
			profile := c.String("profile")
			if len(profile) > 0 {
				fmt.Println("profile set to", profile)
			}
			// http://config-server:8081/greeter-prod.yml
			LocalConfig = GetConfig(micro_c.MicroConfigService, "greeter", profile)
			fmt.Printf("config loaded from config-server is: %s\n", LocalConfig)
		}))

	greeterApi.RegisterGreeterHandler(greeterService.Server(), &Greeter{
		userClient: user.NewUserService(micro_c.MicroNameUser, greeterService.Client())})

	if err := greeterService.Run(); err != nil {
		log.Fatal(err)
	}
}
