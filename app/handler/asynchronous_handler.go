package handler

import (
	"go-api/app/shared/response"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AsynchronousHandler struct {
//
}

type WeatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Error       string  `json:"error"`
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func NewAsynchronousHandler() *AsynchronousHandler {
	return &AsynchronousHandler{
		// Initialize any dependencies if needed
	}
}

// getWeatherForCity fetches weather data for a single city using goroutine
func (h *AsynchronousHandler) getWeatherForCity(city string, wg *sync.WaitGroup, results chan<- WeatherData) {
	defer func ()  {
		wg.Done()
	}()

	// Mock weather API call - replace with real API
	// For demo purposes, we'll simulate API delay and generate mock data
	time.Sleep(time.Duration(100+len(city)*10) * time.Millisecond) // Simulate network delay

	// Simulate API failures for some cities (for demonstration)
	var weather WeatherData

	// Simulate different scenarios including failures
	switch city {
	case "Jakarta":
		weather = WeatherData{
			City:        city,
			Temperature: 32.5,
			Description: "Hot and humid",
		}
	case "Singapore":
		weather = WeatherData{
			City:        city,
			Temperature: 28.0,
			Description: "Partly cloudy",
		}
	case "Tokyo":
		weather = WeatherData{
			City:        city,
			Temperature: 22.3,
			Description: "Clear sky",
		}
	case "London":
		weather = WeatherData{
			City:        city,
			Temperature: 15.8,
			Description: "Rainy",
		}
	case "New York":
		weather = WeatherData{
			City:        city,
			Temperature: 18.2,
			Description: "Overcast",
		}
	case "Sydney":
		// Simulate API failure for Sydney
		weather = WeatherData{
			City:  city,
			Error: "API service temporarily unavailable",
		}
	case "Paris":
		// Simulate timeout error for Paris
		weather = WeatherData{
			City:  city,
			Error: "Request timeout - API took too long to respond",
		}
	case "Bangkok":
		// Simulate invalid API response for Bangkok
		weather = WeatherData{
			City:  city,
			Error: "Invalid API response format",
		}
	default:
		weather = WeatherData{
			City:        city,
			Temperature: 25.0,
			Description: "Unknown weather",
		}
	}

	// Send result to channel (always send, even if there's an error)
	results <- weather
}

func (h *AsynchronousHandler) RunAsyncTask(c *fiber.Ctx) error {
	// Define cities to get weather for
	cities := []string{"Jakarta", "Singapore", "Tokyo", "London", "New York", "Sydney", "Paris", "Bangkok"}

	// Create channels and wait group
	channelResults := make(chan WeatherData, len(cities))
	var wg sync.WaitGroup

	// Start goroutines for each city
	for _, city := range cities {
		wg.Add(1)
		go h.getWeatherForCity(city, &wg, channelResults)
	}

	// Close results channel when all goroutines finish
	go func() {
		wg.Wait()
		close(channelResults)
	}()

	// Collect all results
	// var weatherData []WeatherData
	// var successCount int
	// var errorCount int

	resultMap := make(map[string]WeatherData)

	for item := range channelResults {
		resultMap[item.City] = item
	}

	data := make([]WeatherData, 0, len(cities))
	for _, city := range cities {
		if _, exists := resultMap[city]; exists {
			record := resultMap[city]

			weather := WeatherData{
				City:  city,
				Temperature: record.Temperature,
			}
			data = append(data, weather)
		}
	}

	return response.Success(c, map[string]any{
		"data":  data,
	}, "Weather data fetched successfully!")
}
