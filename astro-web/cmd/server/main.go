package main

import (
	"fmt"
	"time"

	"astro-web/internal/astronomy"
)

func main() {
	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	date := time.Date(2026, 10, 28, 0, 0, 0, 0, loc)
	rs, err := astronomy.MoonRiseSet(-34.6037, -58.3816, date)
	if err != nil {
		panic(err)
	}
	fmt.Println("Has:", rs.Has)
	if rs.Has {
		if !rs.Rise.IsZero() {
			fmt.Println("Moonrise:", rs.Rise.Format("15:04"))
		}
		if !rs.Set.IsZero() {
			fmt.Println("Moonset:", rs.Set.Format("15:04"))
		}
	}
}
