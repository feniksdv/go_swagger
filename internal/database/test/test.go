package test

import (
	"fmt"
	"swagger/internal/database"
)

type test struct {
	Id    int
	Name  string
	Value string
}

func getTest() {
	db := database.Connect()

	rows, err := db.Query("select * from test")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	products := []test{}

	for rows.Next() {
		p := test{}
		err := rows.Scan(&p.Id, &p.Name, &p.Value)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	for _, p := range products {
		fmt.Println(p.Id, p.Name, p.Value)
	}

}
