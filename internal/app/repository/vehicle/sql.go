package vehicle

const (
	QueryGetVehicleNoFilter = `
SELECT
    CASE
        WHEN v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission) THEN ''
        ELSE v.motor_brand
    END                 AS motor_brand
    , CASE
        WHEN
            v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            AND v.motor_transmission = LAG(v.motor_transmission) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            THEN ''
        ELSE v.motor_transmission
    END                 AS motor_transmission
    , v.motor_type
    , v.price
FROM vehicle v
WHERE
    TRUE
ORDER BY v.motor_brand, v.motor_transmission;
`

	QueryGetVehicleWithMotorBrand = `
SELECT
    CASE
        WHEN v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission) THEN ''
        ELSE v.motor_brand
    END                 AS motor_brand
    , CASE
        WHEN
            v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            AND v.motor_transmission = LAG(v.motor_transmission) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            THEN ''
        ELSE v.motor_transmission
    END                 AS motor_transmission
    , v.motor_type
    , v.price
FROM vehicle v
WHERE
    TRUE
	AND v.motor_brand = ?
ORDER BY v.motor_brand, v.motor_transmission;
`

	QueryGetVehicleWithMotorType = `
SELECT
    CASE
        WHEN v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission) THEN ''
        ELSE v.motor_brand
    END                 AS motor_brand
    , CASE
        WHEN
            v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            AND v.motor_transmission = LAG(v.motor_transmission) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            THEN ''
        ELSE v.motor_transmission
    END                 AS motor_transmission
    , v.motor_type
    , v.price
FROM vehicle v
WHERE
    TRUE
	AND v.motor_type = ?
ORDER BY v.motor_brand, v.motor_transmission;
`

	QueryGetVehicleWithMotorTransmission = `
SELECT
    CASE
        WHEN v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission) THEN ''
        ELSE v.motor_brand
    END                 AS motor_brand
    , CASE
        WHEN
            v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            AND v.motor_transmission = LAG(v.motor_transmission) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            THEN ''
        ELSE v.motor_transmission
    END                 AS motor_transmission
    , v.motor_type
    , v.price
FROM vehicle v
WHERE
    TRUE
	AND v.motor_transmission = ?
ORDER BY v.motor_brand, v.motor_transmission;
`

	QueryGetVehicleWithMotorAll = `
SELECT
    CASE
        WHEN v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission) THEN ''
        ELSE v.motor_brand
    END                 AS motor_brand
    , CASE
        WHEN
            v.motor_brand = LAG(v.motor_brand) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            AND v.motor_transmission = LAG(v.motor_transmission) OVER (ORDER BY v.motor_brand, v.motor_transmission)
            THEN ''
        ELSE v.motor_transmission
    END                 AS motor_transmission
    , v.motor_type
    , v.price
FROM vehicle v
WHERE
    TRUE
	AND v.motor_brand = ?
	AND v.motor_type = ?
	AND v.motor_transmission = ?
ORDER BY v.motor_brand, v.motor_transmission;
`
)
