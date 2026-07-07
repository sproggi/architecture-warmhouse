package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TemperatureData struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

type CreateSensorRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Location string `json:"location" binding:"required"`
	Unit     string `json:"unit" binding:"required"`
}

// Sensor - структура датчика
type Sensor struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Location   string    `json:"location"`
	Unit       string    `json:"unit"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

var (
	sensors      = make(map[string]Sensor)
	sensorsMutex sync.RWMutex
	sensorIDCounter = 0
)

func main() {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Get temperature by location
	router.GET("/temperature", func(c *gin.Context) {
		location := c.Query("location")
		if location == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
			return
		}

		// Generate random temperature data based on location
		data := generateTemperatureData(location, "")

		c.JSON(http.StatusOK, data)
	})

	// Get temperature by sensor ID
	router.GET("/temperature/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sensor ID is required"})
			return
		}

		// Generate random temperature data based on sensor ID
		data := generateTemperatureData("", id)

		c.JSON(http.StatusOK, data)
	})

	// Создание датчика
	router.POST("/sensors", func(c *gin.Context) {
		var req CreateSensorRequest
		
		// Валидация запроса
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
				"details": err.Error(),
			})
			return
		}

		// Генерируем ID для датчика
		sensorsMutex.Lock()
		sensorIDCounter++
		sensorID := string(rune('0' + sensorIDCounter)) // Простая генерация ID
		// Если датчиков больше 9, используем более сложную логику
		if sensorIDCounter > 9 {
			sensorID = string(rune('A' + (sensorIDCounter - 10)))
		}
		sensorsMutex.Unlock()

		// Создаем датчик
		now := time.Now()
		sensor := Sensor{
			ID:         sensorID,
			Name:       req.Name,
			Type:       req.Type,
			Location:   req.Location,
			Unit:       req.Unit,
			Status:     "active",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Сохраняем датчик
		sensorsMutex.Lock()
		sensors[sensorID] = sensor
		sensorsMutex.Unlock()

		log.Printf("Sensor created: ID=%s, Name=%s, Location=%s\n", sensor.ID, sensor.Name, sensor.Location)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Sensor created successfully",
			"sensor":  sensor,
		})
	})

	// Получение всех датчиков
	router.GET("/sensors", func(c *gin.Context) {
		sensorsMutex.RLock()
		defer sensorsMutex.RUnlock()
		
		sensorList := make([]Sensor, 0, len(sensors))
		for _, sensor := range sensors {
			sensorList = append(sensorList, sensor)
		}
		
		c.JSON(http.StatusOK, gin.H{
			"sensors": sensorList,
			"count":   len(sensorList),
		})
	})

	// Получение датчика по ID
	router.GET("/sensors/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		sensorsMutex.RLock()
		sensor, exists := sensors[id]
		sensorsMutex.RUnlock()
		
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sensor not found",
			})
			return
		}
		
		c.JSON(http.StatusOK, sensor)
	})

	// Удаление датчика
	router.DELETE("/sensors/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		sensorsMutex.Lock()
		_, exists := sensors[id]
		if exists {
			delete(sensors, id)
		}
		sensorsMutex.Unlock()
		
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sensor not found",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Sensor deleted successfully",
		})
	})

	// Start server
	log.Println("Temperature API starting on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func generateTemperatureData(location, sensorID string) TemperatureData {
	// Generate a random temperature between 18 and 28 degrees Celsius
	value := 18.0 + float64(time.Now().UnixNano()%10) + float64(time.Now().UnixNano()%100)/100.0

	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	return TemperatureData{
		Value:       value,
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Temperature sensor in " + location,
	}
}
