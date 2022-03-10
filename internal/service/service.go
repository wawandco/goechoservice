package service

import (
	"context"
	"fmt"
	"motus/goechoservice/gen-go/echoservice"
)

type Echo struct{}

func (e Echo) Echo(ctx context.Context, input *echoservice.TEchoServiceInputDTO) (*echoservice.TEchoServiceOutputDTO, error) {
	fmt.Printf("[INFO] Received: %s\n", input.GetMessage())

	result := echoservice.NewTEchoServiceOutputDTO()
	result.EchoMessage = fmt.Sprintf("ACK from server: '%s'", input.Message)

	return result, nil
}
