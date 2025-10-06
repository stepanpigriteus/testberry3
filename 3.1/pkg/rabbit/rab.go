package rabbit

import (
	"log"
	"time"

	"github.com/wb-go/wbf/rabbitmq"
)

func InitRab() (*rabbitmq.Connection, *rabbitmq.Channel, *rabbitmq.Publisher) {
	time.Sleep(10 * time.Second)
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
	queue, err := qm.DeclareQueue("my_queue")
	if err != nil {
		log.Fatalf("Ошибка при создании очереди: %v", err)
	}

	if err := ch.QueueBind(queue.Name, "my_key", ex.Name(), false, nil); err != nil {
		log.Fatalf("Ошибка привязки очереди: %v", err)
	}

	publisher := rabbitmq.NewPublisher(ch, ex.Name())
	return conn, ch, publisher
}
