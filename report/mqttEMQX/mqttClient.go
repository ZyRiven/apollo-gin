package mqttemqx

import (
	"apollo/setting"
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConf struct {
	Addr     string
	ClientId string
	UserName string
	Password string
}

type MqttWrapperClient struct {
	c mqtt.Client
}

func InitMqtt(c *MqttConf) (*MqttWrapperClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.Addr)
	opts.SetClientID(c.ClientId)
	opts.SetUsername(c.UserName)
	opts.SetPassword(c.Password)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	} else {
		return &MqttWrapperClient{mqttClient}, nil
	}
}

func (m *MqttWrapperClient) Close() {
	m.c.Disconnect(uint(250))
}

func (m *MqttWrapperClient) Publish(topic string, payload []byte) error {
	pubToken := m.c.Publish(topic, 2, false, payload)
	return pubToken.Error()
}

func (m *MqttWrapperClient) Subscribe(ctx context.Context, topic string, h mqtt.MessageHandler) (err error) {
	if subErr := m.c.Subscribe(topic, 2, h); subErr.Error() != nil {
		setting.ZAPS.Errorf("subscribe error: %s", subErr.Error())
		return subErr.Error()
	} else {
		return nil
	}
}
