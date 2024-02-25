package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

// Переменная для кэширования
var Cache = map[string]string{}

// Вызывает все нижележащие функции для валидации данных, работы с кэшем, генерации и исполнения sql-запроса на вставку
func Store(jsonData string) string {
	isValid := DataValidate(jsonData)
	if !isValid {
		return "Invalid data"
	} else {
		data := ParseJson(jsonData)

		if len(Cache) == 0 {
			RestoreCache(Cache)
		}

		if Cache[data.OrderUID] != "" {
			return "Order already exists"
		}
		Cache[data.OrderUID] = jsonData

		q0, q1, q2, q3, qSlise := MakeQuery(data, jsonData)
		err0 := Query(q0)
		err1 := Query(q1)
		err2 := Query(q2)
		err3 := Query(q3)
		var err4 error
		for _, q := range qSlise {
			err := Query(q)
			if err != nil {
				err4 = err
			}
		}
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return "Query error"
		} else {
			return "Data saved"
		}
	}
}

// Проверка на правильность формата данных
func DataValidate(jsonData string) bool {

	var data DataType
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Парсит json-строку в структуру GO
func ParseJson(jsonData string) DataType {

	var data DataType
	json.Unmarshal([]byte(jsonData), &data)
	return data
}

// Генерирует sql-запрос
func MakeQuery(data DataType, json string) (string, string, string, string, []string) {
	cacheInsert := "insert into cache (order_uid, json)"
	cacheInsert += "values ('" + data.OrderUID + "', '" + json + "')"

	ordersInsert := "insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, "
	ordersInsert += "delivery_service, shardkey, sm_id, date_created, oof_shard)\n"
	ordersInsert += "values "
	ordersInsert += "('" + data.OrderUID + "', '" + data.TrackNumber + "', '" + data.Entry + "', '" + data.Locale + "', '"
	ordersInsert += data.InternalSignature + "', '" + data.CustomerID + "', '" + data.DeliveryService + "', '"
	ordersInsert += data.Shardkey + "', " + fmt.Sprint(data.SmID) + ", '" + strings.Replace(strings.Replace(data.DateCreated, "T", " ", 1), "Z", "", 1)
	ordersInsert += "', '" + data.OofShard + "')"

	deliveryInsert := "insert into delivery (order_uid, name, phone, zip, city, address, region, email)\n"
	deliveryInsert += "values"
	deliveryInsert += "('" + data.OrderUID + "', '" + data.Delivery.Name + "', '" + data.Delivery.Phone + "', '"
	deliveryInsert += data.Delivery.Zip + "', '" + data.Delivery.City + "', '" + data.Delivery.Address + "', '"
	deliveryInsert += data.Delivery.Region + "', '" + data.Delivery.Email + "')"

	paymentInsert := "insert into payments (order_uid, transaction, request_id, currency, provider, amount, "
	paymentInsert += "payment_dt, bank, delivery_cost, goods_total, custom_fee)\n"
	paymentInsert += "values"
	paymentInsert += "('" + data.OrderUID + "', '" + data.Payment.Transaction + "', '" + data.Payment.RequestID + "', '"
	paymentInsert += data.Payment.Currency + "', '" + data.Payment.Provider + "', " + fmt.Sprint(data.Payment.Amount) + ", "
	paymentInsert += fmt.Sprint(data.Payment.PaymentDt) + ", '" + data.Payment.Bank + "', " + fmt.Sprint(data.Payment.DeliveryCost) + ", "
	paymentInsert += fmt.Sprint(data.Payment.GoodsTotal) + ", " + fmt.Sprint(data.Payment.CustomFee) + ")"

	itemsInserts := []string{}

	for i := range data.Items {
		itemInsert := "insert into items (chrt_id, order_uid, track_number, price, rid, "
		itemInsert += "name, sale, size, total_price, nm_id, brand, status)\n"
		itemInsert += "values"
		itemInsert += "(" + fmt.Sprint(data.Items[i].ChrtID) + ", '" + data.OrderUID + "', '"
		itemInsert += data.Items[i].TrackNumber + "', " + fmt.Sprint(data.Items[i].Price) + ", '"
		itemInsert += data.Items[i].Rid + "', '" + data.Items[i].Name + "', " + fmt.Sprint(data.Items[i].Sale) + ", '"
		itemInsert += data.Items[i].Size + "', " + fmt.Sprint(data.Items[i].TotalPrice) + ", " + fmt.Sprint(data.Items[i].NmID) + ", '"
		itemInsert += data.Items[i].Brand + "', " + fmt.Sprint(data.Items[i].Status) + ")"
		itemsInserts = append(itemsInserts, itemInsert)
	}

	return cacheInsert, ordersInsert, deliveryInsert, paymentInsert, itemsInserts
}

// Исполняет sql-запрос
func Query(query string) error {

	db, err := sql.Open("postgres", auth)
	if err != nil {
		log.Println(err)
	} else {

		_, err = db.Exec(query)
		if err != nil {
			log.Println(err)
		}
	}

	defer db.Close()
	return err
}

// Восстанавливает кэш из БД
func RestoreCache(cache map[string]string) {
	db, err := sql.Open("postgres", auth)
	if err != nil {
		log.Println(err)
	} else {
		rows, err := db.Query("select * from cache")
		if err != nil {
			log.Println(err)
		}
		for rows.Next() {
			var uid, json string
			err := rows.Scan(&uid, &json)
			if err != nil {
				log.Println(err)
				continue
			}
			cache[uid] = json
		}
	}
}
