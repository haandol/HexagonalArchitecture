// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/google/wire"
	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/primary/poller"
	"github.com/haandol/hexagonal/pkg/adapter/primary/router"
	producer2 "github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/connector/producer"
	"github.com/haandol/hexagonal/pkg/port"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/service"
	"gorm.io/gorm"
	"net/http"
)

// Injectors from wire.go:

func InitTripApp(cfg *config.Config) port.App {
	ginRouter := router.NewGinRouter(cfg)
	server := router.NewServerForce(cfg, ginRouter)
	db := provideDB(cfg)
	tripRepository := repository.NewTripRepository(db)
	tripService := service.NewTripService(tripRepository)
	tripRouter := router.NewTripRouter(tripService)
	tripConsumer := provideTripConsumer(cfg, tripService)
	tripApp := NewTripApp(server, ginRouter, tripRouter, tripConsumer)
	return tripApp
}

func InitSagaApp(cfg *config.Config) port.App {
	ginRouter := router.NewGinRouter(cfg)
	server := router.NewServer(cfg, ginRouter)
	sagaProducer := provideSagaProducer(cfg)
	db := provideDB(cfg)
	sagaRepository := repository.NewSagaRepository(db)
	sagaService := service.NewSagaService(sagaProducer, sagaRepository)
	sagaConsumer := provideSagaConsumer(cfg, sagaService)
	sagaApp := NewSagaApp(server, sagaConsumer)
	return sagaApp
}

func InitCarApp(cfg *config.Config) port.App {
	ginRouter := router.NewGinRouter(cfg)
	server := router.NewServer(cfg, ginRouter)
	db := provideDB(cfg)
	carRepository := repository.NewCarRepository(db)
	carService := service.NewCarService(carRepository)
	carConsumer := provideCarConsumer(cfg, carService)
	carApp := NewCarApp(server, carConsumer)
	return carApp
}

func InitMessageRelayApp(cfg *config.Config) port.App {
	kafkaProducer := provideProducer(cfg)
	db := provideDB(cfg)
	outboxRepository := repository.NewOutboxRepository(db)
	messageRelayService := service.NewMessageRelayService(kafkaProducer, outboxRepository)
	outboxPoller := poller.NewOutboxPoller(cfg, messageRelayService)
	messageRelayApp := NewMessageRelayApp(outboxPoller)
	return messageRelayApp
}

func InitHotelApp(cfg *config.Config) port.App {
	ginRouter := router.NewGinRouter(cfg)
	server := router.NewServer(cfg, ginRouter)
	db := provideDB(cfg)
	hotelRepository := repository.NewHotelRepository(db)
	hotelService := service.NewHotelService(hotelRepository)
	hotelConsumer := provideHotelConsumer(cfg, hotelService)
	hotelApp := NewHotelApp(server, hotelConsumer)
	return hotelApp
}

func InitFlightApp(cfg *config.Config) port.App {
	ginRouter := router.NewGinRouter(cfg)
	server := router.NewServer(cfg, ginRouter)
	db := provideDB(cfg)
	flightRepository := repository.NewFlightRepository(db)
	flightService := service.NewFlightService(flightRepository)
	flightConsumer := provideFlightConsumer(cfg, flightService)
	flightApp := NewFlightApp(server, flightConsumer)
	return flightApp
}

// wire.go:

// Common
func provideDB(cfg *config.Config) *gorm.DB {
	db, err := database.Connect(&cfg.TripDB)
	if err != nil {
		panic(err)
	}
	return db
}

func provideProducer(cfg *config.Config) *producer.KafkaProducer {
	kafkaProducer, err := producer.Connect(&cfg.Kafka)
	if err != nil {
		panic(err)
	}
	return kafkaProducer
}

// TripApp
func provideTripConsumer(
	cfg *config.Config,
	tripService *service.TripService,
) *consumer.TripConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "trip", "trip-service")
	return consumer.NewTripConsumer(kafkaConsumer, tripService)
}

var provideTripRouters = wire.NewSet(router.NewGinRouter, wire.Bind(new(http.Handler), new(*router.GinRouter)), router.NewServerForce, wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)), router.NewTripRouter)

var provideRouters = wire.NewSet(router.NewGinRouter, wire.Bind(new(http.Handler), new(*router.GinRouter)), router.NewServer, wire.Bind(new(routerport.RouterGroup), new(*router.GinRouter)))

// SagaApp
func provideSagaConsumer(
	cfg *config.Config,
	sagaService *service.SagaService,
) *consumer.SagaConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "saga", "saga-service")
	return consumer.NewSagaConsumer(kafkaConsumer, sagaService)
}

func provideSagaProducer(cfg *config.Config) *producer2.SagaProducer {
	kafkaProducer := provideProducer(cfg)
	return producer2.NewSagaProducer(kafkaProducer)
}

// CarApp
func provideCarConsumer(
	cfg *config.Config,
	carService *service.CarService,
) *consumer.CarConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "car", "car-service")
	return consumer.NewCarConsumer(kafkaConsumer, carService)
}

// HotelApp
func provideHotelConsumer(
	cfg *config.Config,
	hotelService *service.HotelService,
) *consumer.HotelConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "hotel", "hotel-service")
	return consumer.NewHotelConsumer(kafkaConsumer, hotelService)
}

// FlightApp
func provideFlightConsumer(
	cfg *config.Config,
	flightService *service.FlightService,
) *consumer.FlightConsumer {
	kafkaConsumer := consumer.NewKafkaConsumer(&cfg.Kafka, "flight", "flight-service")
	return consumer.NewFlightConsumer(kafkaConsumer, flightService)
}
