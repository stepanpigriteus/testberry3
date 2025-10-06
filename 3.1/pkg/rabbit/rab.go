package rabbit

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
)

func InitRab() (*rabbitmq.Connection, *rabbitmq.Channel, *rabbitmq.Publisher) {
	time.Sleep(15 * time.Second)
	conn, err := rabbitmq.Connect("amqp://guest:guest@rabbitmq:5672/", 5, 2*time.Second)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %v", err)
	}

	ex := rabbitmq.NewExchange("my_exchange", "direct")
	if err := ex.BindToChannel(ch); err != nil {
		log.Fatalf("Ошибка при создании обменника: %v", err)
	}

	qm := rabbitmq.NewQueueManager(ch)

	mainQueue, err := qm.DeclareQueue("my_queue")
	if err != nil {
		log.Fatalf("Ошибка при создании основной очереди: %v", err)
	}

	delayQueue, err := qm.DeclareQueue("delay_queue", rabbitmq.QueueConfig{
		Durable: true,
		Args: amqp091.Table{
			"x-dead-letter-exchange":    "my_exchange",
			"x-dead-letter-routing-key": "main_key",
		},
	})
	if err != nil {
		log.Fatalf("Ошибка при создании delay очереди: %v", err)
	}

	if err := ch.QueueBind(mainQueue.Name, "main_key", ex.Name(), false, nil); err != nil {
		log.Fatalf("Ошибка привязки основной очереди: %v", err)
	}

	if err := ch.QueueBind(delayQueue.Name, "delay_key", ex.Name(), false, nil); err != nil {
		log.Fatalf("Ошибка привязки delay очереди: %v", err)
	}

	publisher := rabbitmq.NewPublisher(ch, ex.Name())

	return conn, ch, publisher
}
