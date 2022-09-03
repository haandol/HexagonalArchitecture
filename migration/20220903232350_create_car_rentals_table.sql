-- +goose Up
-- +goose StatementBegin
CREATE TABLE car_rentals (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  trip_id BIGINT NOT NULL,
  car_id BIGINT NOT NULL,
  quantity INT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_trip_id ON car_rentals (trip_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX idx_car_id ON car_rentals (car_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE car_rentals
-- +goose StatementEnd
