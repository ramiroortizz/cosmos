package main

import (
	"fmt"
	"time"

	"astro-web/internal/astronomy"
)

func main() {
	loc, _ := time.LoadLocation("America/Argentina/Buenos_Aires")
	date := time.Date(2025, 10, 28, 0, 0, 0, 0, loc)
	rs, err := astronomy.SunRiseSet(-34.6037, -58.3816, date)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tiene:", rs.Has)
	fmt.Println("Amanecer:", rs.Rise.Format("15:04"))
	fmt.Println("Atardecer:", rs.Set.Format("15:04"))
}
