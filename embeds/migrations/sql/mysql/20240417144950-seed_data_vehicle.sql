
-- +migrate Up
INSERT INTO `vehicle` (`motor_brand`, `motor_type`, `motor_transmission`, `price`)
SELECT
    *
FROM (
    WITH input_data AS (
        SELECT CAST('
            [
              {
                "motor_brand": "Honda",
                "motor_type": "Beat",
                "motor_transmission": "Automatic",
                "price": 18500000
              },
              {
                "motor_brand": "Honda",
                "motor_type": "Genio",
                "motor_transmission": "Automatic",
                "price": 18500000
              },
              {
                "motor_brand": "Honda",
                "motor_type": "Scoopy",
                "motor_transmission": "Automatic",
                "price": 22500000
              },
              {
                "motor_brand": "Honda",
                "motor_type": "SuperCub",
                "motor_transmission": "Manual",
                "price": 77160000
              },
              {
                "motor_brand": "Honda",
                "motor_type": "GTR",
                "motor_transmission": "Manual",
                "price": 25180000
              },
              {
                "motor_brand": "Yamaha",
                "motor_type": "XMax",
                "motor_transmission": "Automatic",
                "price": 66450000
              },
              {
                "motor_brand": "Yamaha",
                "motor_type": "NMax",
                "motor_transmission": "Automatic",
                "price": 31925000
              },
              {
                "motor_brand": "Yamaha",
                "motor_type": "XSR",
                "motor_transmission": "Manual",
                "price": 3807500
              }]
                ' AS JSON
                               ) AS json_data
    )
    SELECT
        vehicles_data.motor_brand,
        vehicles_data.motor_type,
        vehicles_data.motor_transmission,
        vehicles_data.price
    FROM
        input_data id,
        JSON_TABLE(
                id.json_data,
                '$[*]' COLUMNS (
                    motor_brand VARCHAR(100) PATH '$.motor_brand',
                    motor_type VARCHAR(100) PATH '$.motor_type',
                    motor_transmission VARCHAR(100) PATH '$.motor_transmission',
                    price BIGINT PATH '$.price'
                    )
        ) vehicles_data
) cte_data;

-- +migrate Down
DELETE FROM `vehicle`;
