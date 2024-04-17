package entity

type Vehicle struct {
	MotorBrand        string `db:"motor_brand" json:"motor_brand" mapstructure:"motor_brand"`
	MotorType         string `db:"motor_type" json:"motor_type" mapstructure:"motor_type"`
	MotorTransmission string `db:"motor_transmission" json:"motor_transmission" mapstructure:"motor_transmission"`
	Price             int64  `db:"price" json:"price" mapstructure:"price"`
}
