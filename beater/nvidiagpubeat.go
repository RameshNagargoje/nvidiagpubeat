package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/deepujain/nvidiagpubeat/config"
)

type Nvidiagpubeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Nvidiagpubeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Nvidiagpubeat) Run(b *beat.Beat) error {
	logp.Info("nvidiagpubeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		cmd := NvidiaCommand{query: bt.config.Query, env: bt.config.Env}
		events := Run(cmd, cmd.query)
		for _, event := range events {
			bt.client.PublishEvent(event)
		}
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Nvidiagpubeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
