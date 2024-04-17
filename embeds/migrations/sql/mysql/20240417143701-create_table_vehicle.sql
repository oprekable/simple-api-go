
-- +migrate Up
CREATE TABLE IF NOT EXISTS vehicle (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "vehicle ID",
    motor_brand VARCHAR(100) NOT NULL COMMENT "vehicle motor brand",
    motor_type VARCHAR(100) NOT NULL COMMENT "vehicle motor type",
    motor_transmission VARCHAR(100) NOT NULL COMMENT "vehicle motor transmission",
    price BIGINT NOT NULL COMMENT "vehicle price",
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "The date and time of creation",
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "The date and time of update");

CREATE INDEX idx_vehicle_motor_brand ON vehicle (motor_brand);
CREATE INDEX idx_vehicle_motor_type ON vehicle (motor_type);
CREATE INDEX idx_vehicle_motor_transmission ON vehicle (motor_transmission);
CREATE INDEX idx_vehicle_combine ON vehicle (motor_brand, motor_type, motor_transmission);

-- +migrate Down
DROP TABLE IF EXISTS vehicle;