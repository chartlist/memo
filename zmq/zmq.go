package zmq

import (
	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/database"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/processor"

	"github.com/pebbe/zmq4"
)

const rawBlock string = "rawblock"
const hashBlock string = "hashblock"
const rawTx string = "rawtx"
const rawTx2 string = "rawtxfee"
const hashTx string = "hashtx"

func SetupZMQ(pool *database.RedisPool) {

	zmqHost := config.GetString("zmq.host")
	zmqPort := config.GetString("zmq.port")
	connectionString := "tcp://" + zmqHost + ":" + zmqPort

	subscriber, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		logger.Error.Println("zmq4.NewSocket err:", err)
		return
	}
	err = subscriber.Connect(connectionString)
	if err != nil {
		logger.Error.Println("Subscriber connect server error:", err, ",ip:", connectionString)
		return
	}
	if config.GetBool("zmq.subscribeTo.rawTx") {
		if err := subscriber.SetSubscribe(rawTx); err != nil {
			logger.Error.Println("Subscriber rawTx error:", err)
			return
		}
	}
	if config.GetBool("zmq.subscribeTo.rawTx2") {
		if err := subscriber.SetSubscribe(rawTx2); err != nil {
			logger.Error.Println("Subscriber rawTx2 error:", err)
			return
		}
	}
	if config.GetBool("zmq.subscribeTo.hashTx") {
		if err := subscriber.SetSubscribe(hashTx); err != nil {
			logger.Error.Println("Subscriber hashTx error:", err)
			return
		}
	}
	if config.GetBool("zmq.subscribeTo.rawBlock") {
		if err := subscriber.SetSubscribe(rawBlock); err != nil {
			logger.Error.Println("Subscriber rawBlock error:", err)
			return
		}
	}
	if config.GetBool("zmq.subscribeTo.hashBlock") {
		if err := subscriber.SetSubscribe(hashBlock); err != nil {
			logger.Error.Println("Subscriber hashBlock error:", err)
			return
		}
	}

	defer func() {
		if err := subscriber.Close(); err != nil {
			logger.Error.Println("Subscriber close error:", err)
		}
	}()

	loopZMQ(subscriber, pool)
}

func loopZMQ(subscriber *zmq4.Socket, pool *database.RedisPool) {
	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			logger.Error.Println(err)
		}
		handleZMQMessage(msg, pool)
	}
}

func handleZMQMessage(zmqMessage []string, pool *database.RedisPool) {
	topic := zmqMessage[0]
	payload := zmqMessage[1]
	switch topic {
	case rawBlock:
		go processor.HandleRawBlock(payload, pool)
	case hashBlock:
		go processor.HandleHashBlock(payload)
	case rawTx:
		go processor.HandleRawTx(payload)
	case rawTx2:
		go processor.HandleRawTxWithSizeAndFee(payload, pool)
	case hashTx:
		go processor.HandleHashTx(payload)
	default:
		logger.Warning.Println("Unhandled ZMQ topic", topic)
	}
}
