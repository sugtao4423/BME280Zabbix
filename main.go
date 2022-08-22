package main

import (
	"fmt"
	"os"
	"time"

	"github.com/quhar/bme280"
	"github.com/spetr/go-zabbix-sender"
	flag "github.com/spf13/pflag"
	"golang.org/x/exp/io/i2c"
)

var (
	i2cDev       string
	bme280Addr   int
	zbxHost      string
	itemHostname string
	itemTempKey  string
	itemPressKey string
	itemHumKey   string
)

func main() {
	flag.StringVarP(&i2cDev, "i2c-dev", "d", "/dev/i2c-1", "I2C device")
	flag.IntVarP(&bme280Addr, "bme280-addr", "a", 0x76, "BME280 address")
	flag.StringVarP(&zbxHost, "zbx-host", "z", "localhost:10051", "Zabbix host")
	flag.StringVarP(&itemHostname, "item-hostname", "h", "localhost", "Zabbix item hostname")
	flag.StringVarP(&itemTempKey, "item-temp-key", "T", "bme280_temp", "Zabbix item temp key")
	flag.StringVarP(&itemPressKey, "item-press-key", "P", "bme280_press", "Zabbix item press key")
	flag.StringVarP(&itemHumKey, "item-hum-key", "H", "bme280_hum", "Zabbix item hum key")
	flag.Parse()

	temp, press, hum := retryGetBme280()
	t := fmt.Sprintf("%.2f", temp)
	p := fmt.Sprintf("%.2f", press)
	h := fmt.Sprintf("%.2f", hum)

	res := retrySendZbx(t, p, h)
	fmt.Println("send data to Zabbix")
	fmt.Println("response:", res.Response)
	fmt.Println("info:", res.Info)
}

func retryGetBme280() (temp, press, hum float64) {
	for i := 0; i < 5; i++ {
		t, p, h, err := getBme280()
		if err != nil {
			fmt.Println("failed to get BME280 data:", err)
			fmt.Printf("retrying... (attempt %d)\n", i+1)
			time.Sleep(time.Second)
			continue
		}
		return t, p, h
	}
	fmt.Println("failed to get BME280 data")
	fmt.Println("aborting...")
	os.Exit(1)
	return
}

func retrySendZbx(temp, press, hum string) *zabbix.Response {
	for i := 0; i < 5; i++ {
		res, err := sendZbx(temp, press, hum)
		if err != nil {
			fmt.Println("failed to send data to Zabbix:", err)
			fmt.Printf("retrying... (attempt %d)\n", i+1)
			time.Sleep(time.Second)
			continue
		}
		return res
	}
	fmt.Println("failed to send data to Zabbix")
	fmt.Println("aborting...")
	os.Exit(1)
	return nil
}

func getBme280() (temp, press, hum float64, err error) {
	d, err := i2c.Open(&i2c.Devfs{Dev: i2cDev}, bme280Addr)
	if err != nil {
		return 0, 0, 0, err
	}
	defer d.Close()

	b := bme280.New(d)
	err = b.Init()
	if err != nil {
		return 0, 0, 0, err
	}

	t, p, h, err := b.EnvData()
	if err != nil {
		return 0, 0, 0, err
	}
	return t, p, h, nil
}

func sendZbx(temp, press, hum string) (*zabbix.Response, error) {
	metrics := []*zabbix.Metric{
		zabbix.NewMetric(itemHostname, itemTempKey, temp, false),
		zabbix.NewMetric(itemHostname, itemPressKey, press, false),
		zabbix.NewMetric(itemHostname, itemHumKey, hum, false),
	}
	sender := zabbix.NewSender(zbxHost)

	_, _, resTrapper, errTrapper := sender.SendMetrics(metrics)
	if errTrapper != nil {
		return nil, fmt.Errorf("zabbix send error: %v", errTrapper)
	}
	return &resTrapper, nil
}
