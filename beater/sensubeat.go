package beater

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/fylie/sensubeat/config"
)

type Sensubeat struct {
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

	bt := &Sensubeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Sensubeat) Run(b *beat.Beat) error {
	logp.Info("sensubeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		resp, err := http.Get("http://" + bt.config.Host + ":" + bt.config.Port + "/events")
	    if err != nil {
	        return err
	    }
	    defer resp.Body.Close()

	    body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var parsed interface{}
		_ = json.Unmarshal(body, &parsed)

		for _, r := range parsed.([]interface{}) {
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"counter":    counter,
				"resp":		  r,
			}
			bt.client.PublishEvent(event)
			logp.Info("Event sent")
		}

		counter++
	}
}

func (bt *Sensubeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
