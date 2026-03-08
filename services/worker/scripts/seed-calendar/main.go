package main

import (
	"fmt"
	"time"
)

// BankingHoliday represents a Brazilian banking holiday.
type BankingHoliday struct {
	Date       time.Time
	Name       string
	IsMoveable bool // e.g., Easter, Carnival
}

// 2026 Brazilian Banking Holidays (Anbima)
var holidays2026 = []BankingHoliday{
	// Fixed holidays
	{Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Name: "New Year's Day"},
	{Date: time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC), Name: "Tiradentes' Day"},
	{Date: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC), Name: "Labour Day"},
	{Date: time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC), Name: "Independence Day"},
	{Date: time.Date(2026, 10, 12, 0, 0, 0, 0, time.UTC), Name: "Our Lady Aparecida"},
	{Date: time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC), Name: "All Souls' Day"},
	{Date: time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC), Name: "Black Consciousness Day"},
	{Date: time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC), Name: "Christmas"},

	// Moveable holidays for 2026
	// Easter Sunday: April 5, 2026
	// Good Friday: April 3, 2026
	{Date: time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC), Name: "Good Friday", IsMoveable: true},
	{Date: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC), Name: "Easter Sunday", IsMoveable: true},
	// Corpus Christi (60 days after Easter): May 28, 2026
	{Date: time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC), Name: "Corpus Christi", IsMoveable: true},
}

// CalendarEntry represents a single day in the banking calendar.
type CalendarEntry struct {
	Date        time.Time
	IsHoliday   bool
	IsWeekend   bool
	HolidayName string
}

func main() {
	fmt.Println("Brazilian Banking Holiday Calendar Seeder")
	fmt.Println("=========================================")
	fmt.Println()

	// Generate calendar for 2026
	generateCalendar(2026)
}

// generateCalendar creates calendar entries for a given year.
func generateCalendar(year int) {
	// Create a map of holidays for quick lookup
	holidayMap := make(map[string]BankingHoliday)
	for _, h := range holidays2026 {
		key := h.Date.Format("2006-01-02")
		holidayMap[key] = h
	}

	fmt.Printf("Banking Holidays for %d:\n", year)
	fmt.Println()

	// Iterate through the year
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)

	businessDays := 0
	holidayCount := 0
	weekendDays := 0

	for currentDate := startDate; currentDate.Before(endDate.AddDate(0, 0, 1)); currentDate = currentDate.AddDate(0, 0, 1) {
		dayOfWeek := currentDate.Weekday()

		entry := CalendarEntry{
			Date: currentDate,
		}

		isWeekend := dayOfWeek == time.Saturday || dayOfWeek == time.Sunday
		if isWeekend {
			entry.IsWeekend = true
			weekendDays++
		} else {
			businessDays++
		}

		// Check for holiday
		dateKey := currentDate.Format("2006-01-02")
		if holiday, exists := holidayMap[dateKey]; exists {
			entry.IsHoliday = true
			entry.HolidayName = holiday.Name
			holidayCount++
			// Only subtract from businessDays if the holiday falls on a weekday
			if !isWeekend {
				businessDays--
			}
		}

		// Print holidays only
		if entry.IsHoliday {
			fmt.Printf("%s | %s | %s\n",
				entry.Date.Format("2006-01-02"),
				entry.Date.Format("Monday"),
				entry.HolidayName,
			)
		}
	}

	fmt.Println()
	fmt.Printf("Summary for %d:\n", year)
	fmt.Printf("Total days: 365\n")
	fmt.Printf("Business days: %d\n", businessDays)
	fmt.Printf("Weekend days: %d\n", weekendDays)
	fmt.Printf("Banking holidays: %d\n", holidayCount)
	fmt.Println()

	// TODO: Persist to database or file
	fmt.Println("Calendar data ready for persistence to database...")
}
