package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type PubSub struct {
	redis *redis.Client
}

func NewPubSub(redis *redis.Client) *PubSub {
    return &PubSub{redis: redis}
}


func (p *PubSub)Publish(ctx context.Context,roomID int,payload any)error {

	bytes,err:=json.Marshal(payload)

	if err!=nil{
		return err
	}

	channel := fmt.Sprintf("room:%d", roomID)

	return p.redis.Publish(ctx,channel,bytes).Err()
}


func (p *PubSub)Subscribe(ctx context.Context,roomID int,out chan []byte){

	channel := fmt.Sprintf("room:%d", roomID)

	pb:=p.redis.Subscribe(ctx,channel)
	
	msg:=pb.Channel()
	
	go p.ReadSubscription(ctx,msg,pb,out)

}


func (p *PubSub)ReadSubscription(ctx context.Context,ms <-chan *redis.Message,pb *redis.PubSub,out chan []byte){

	defer pb.Close()
	for {
		select{
		case <-ctx.Done():
			return 
		case message,ok:=<-ms:
			if !ok{
				return
			}
			select{
			case out <- []byte(message.Payload):
			default:
			fmt.Println("warning: hub channel full, dropping message")
		}

		}
	}

}