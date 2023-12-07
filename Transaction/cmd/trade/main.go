package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/BielLopes/FullCicle/Transaction/internal/infra/kafka"
	"github.com/BielLopes/FullCicle/Transaction/internal/market/dto"
	"github.com/BielLopes/FullCicle/Transaction/internal/market/entity"
	"github.com/BielLopes/FullCicle/Transaction/internal/market/transformer"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() { // Primeira Thread T1
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)

	wg := &sync.WaitGroup{}
	defer wg.Wait() // ultima linha de código dessa função

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers":	"host.docker.internal:9094",
		"group.id":				"myGroup",
		"auto.offset.reset":	"latest",
	}

	producer := kafka.NewProducer(configMap)
	consumer := kafka.NewConsumer(configMap, []string{"input"})

	go consumer.Consume(kafkaMsgChan)  // Criada uma segunda Thread T2

	// Recebe do canal do kafka no Input (ordersIn), Processa para achar matchs entre Orders e Joga as ordens no canal de saída (ordersout) para serem publicadas no kafka
	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // Cria a tercerira Thread T3

	go func() { // Cria a quarta Thread T4 - Isso é uma função anônima
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))

			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}

			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "   ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}

		producer.Publish(outputJson, []byte("orders"), "output")
	}
}