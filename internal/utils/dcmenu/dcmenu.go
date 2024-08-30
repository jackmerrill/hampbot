package dcmenu

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly/v2"
)

type Menu []Meal

type Meal struct {
	Date      string   `json:"Date"`
	Breakfast []string `json:"Breakfast"`
	Lunch     []string `json:"Lunch"`
	Dinner    []string `json:"Dinner"`
}

func Run(session *discordgo.Session) {
	// Parse the website for the URL
	url, err := ParseWebsite()
	if err != nil {
		panic(err)
	}

	// Parse the URL for the CSV
	parsedURL, err := ParseURL(*url)
	if err != nil {
		panic(err)
	}

	// Parse the CSV for the Menu
	menu, err := ParseCSV(*parsedURL)
	if err != nil {
		panic(err)
	}

	if menu == nil {
		panic("menu is nil")
	}

	refMenu := *menu

	var todaysMenu Meal
	todaysDate := time.Now().Format("1/2/2006")

	for i, v := range refMenu {
		if v.Date == todaysDate {
			todaysMenu = refMenu[i]
			break
		}
	}

	isBrunch := false

	// Check if today is a brunch day (check if breakfast is null)
	if len(todaysMenu.Breakfast) == 0 || todaysMenu.Breakfast == nil {
		isBrunch = true
	}

	fields := []*discordgo.MessageEmbedField{}

	if !isBrunch {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "🧇 **Breakfast**",
			Value:  strings.Join(todaysMenu.Breakfast, "\n"),
			Inline: false,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "🥪 **Lunch**",
		Value:  strings.Join(todaysMenu.Lunch, "\n"),
		Inline: false,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "🍝 **Dinner**",
		Value:  strings.Join(todaysMenu.Dinner, "\n"),
		Inline: false,
	})

	// Send the Menu to the channel
	_, err = session.ChannelMessageSendEmbed("1016558809667350528", &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("🍕 **Menu for %s**", todaysMenu.Date),
		Fields: fields,
		Color:  0x00ff00,
	})

	if err != nil {
		panic(err)
	}
}

func ParseWebsite() (*string, error) {
	c := colly.NewCollector()

	var href string

	c.OnHTML("a[href].button--secondary", func(e *colly.HTMLElement) {
		if e.Text == "MENU" {
			href = e.Attr("href")
		}
	})

	err := c.Visit("https://www.hampshire.edu/student-life/campus-dining-services")
	if err != nil {
		return nil, err
	}

	// Wait for href to be populated
	for href == "" {
	}

	return &href, nil
}

func ParseURL(u string) (*string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	parsedURL = parsedURL.JoinPath("../", "gviz", "tq")

	query := parsedURL.Query()
	query.Del("usp")
	query.Add("tqx", "out:csv")

	parsedURL.RawQuery = query.Encode()

	str := parsedURL.String()
	return &str, nil
}

func ParseCSV(u string) (*Menu, error) {
	// Download the CSV from the URL
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the CSV file
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Remove the first row, its a title
	records = records[1:]

	// Initialize Menu
	var menu Menu

	var week2Row int

	// Find the row where the second week starts
	tmp := records
	for i, row := range tmp[1:] {
		if isDate(row[0]) {
			week2Row = i + 1
			break
		}
	}

	// Initialize separate weeks
	var week1 Menu = make(Menu, 7)
	var week2 Menu = make(Menu, 7)

	// Iterate through each column for Week 1 and Week 2
	for day := 0; day <= 6; day++ {
		// Process Week 1
		week1[day] = parseMeal(records, day, 1, week2Row)

		// Process Week 2
		week2[day] = parseMeal(records, day, week2Row+1, len(records))
	}

	menu = append(menu, week1...)
	menu = append(menu, week2...)

	return &menu, nil
}

func parseMeal(records [][]string, col int, startRow int, endRow int) Meal {
	meal := Meal{Date: records[startRow-1][col]}
	currentMealType := "Breakfast"

	for rowIndex := startRow; rowIndex < endRow; rowIndex++ {
		cell := strings.TrimSpace(records[rowIndex][col])

		// Detect if the cell is a new date for the next week
		if isDate(cell) && rowIndex != startRow {
			// We've hit the next date, so we return the meal as is.
			break
		}

		// Detect and set meal type
		if strings.Contains(strings.ToLower(cell), "breakfast") {
			currentMealType = "Breakfast"
		} else if strings.Contains(strings.ToLower(cell), "brunch") {
			currentMealType = "Brunch"
		} else if strings.Contains(strings.ToLower(cell), "lunch") {
			currentMealType = "Lunch"
		} else if strings.Contains(strings.ToLower(cell), "dinner") {
			currentMealType = "Dinner"
		} else if cell != "" {
			// Append meal item to the current meal type
			switch currentMealType {
			case "Breakfast":
				meal.Breakfast = append(meal.Breakfast, cell)
			case "Brunch":
				meal.Lunch = append(meal.Lunch, cell)
			case "Lunch":
				meal.Lunch = append(meal.Lunch, cell)
			case "Dinner":
				meal.Dinner = append(meal.Dinner, cell)
			}
		}
	}

	return meal
}

// Helper function to check if a string is a date
func isDate(str string) bool {
	// Simple check for date pattern (e.g., "8/19/2024")
	return regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`).MatchString(str)
}
