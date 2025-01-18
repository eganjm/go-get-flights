package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getAllFlights(w http.ResponseWriter, r *http.Request) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	svc := dynamodb.New(sess)

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("FlightSchedules"),
	})
	if err != nil {
		http.Error(w, "Error scanning DynamoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var flights []map[string]string
	for _, item := range result.Items {
		flight := map[string]string{
			"flight_id":      *item["flight_id"].S,
			"departure_time": *item["departure_time"].S,
			"destination":    *item["destination"].S,
			"status":         *item["status"].S,
		}
		flights = append(flights, flight)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Simple CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin for testing; in production, specify your domain
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's an OPTIONS request (preflight), return now
		if r.Method == http.MethodOptions {
			return
		}

		// Otherwise, continue to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/flights", getAllFlights)

	// Wrap routes with CORS middleware
	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
