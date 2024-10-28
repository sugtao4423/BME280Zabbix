# BME280Zabbix
BME280で取得した温度・湿度・気圧データをZabbixに送信する

## Usage
short | long | description | default
--- | --- | --- | ---
`-d` | `--i2c-dev` | I2Cデバイス | `/dev/i2c-1`
`-a` | `--bme280-addr` | BME280アドレス | `0x76`
`-n` | `--dry-run` | Dry run | `false`
`-z` | `--zbx-host` | Zabbixサーバー | `localhost:10051`
`-h` | `--item-hostname` | Zabbixアイテムのhostname | `localhost`
`-T` | `--item-temp-key` | Zabbixアイテムの温度Key | `bme280_temp`
`-P` | `--item-press-key` | Zabbixアイテムの気圧Key | `bme280_press`
`-H` | `--item-hum-key` | Zabbixアイテムの湿度Key | `bme280_hum`
