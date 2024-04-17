package vehicle

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blockloop/scan/v2"
	"simple-api-go/internal/app/entity"
	"simple-api-go/internal/app/error/core"
	"simple-api-go/internal/pkg/utils/log"
	"time"

	"github.com/aaronjan/hunch"
	"github.com/go-gorp/gorp/v3"
)

type db struct {
	TimeLocation               *time.Location
	WriteDB                    *sql.DB
	ReadDB                     *sql.DB
	TimeZone                   string
	TimeFormat                 string
	TimePostgresFriendlyFormat string
}

var _ Repository = (*db)(nil)

func NewDB(
	writeDBMap *gorp.DbMap,
	readDBMap *gorp.DbMap,
	timeZone string,
	timeFormat string,
	timePostgresFriendlyFormat string,
	timeLocation *time.Location,
) (returnData Repository) {
	if writeDBMap != nil || readDBMap != nil {
		return &db{
			WriteDB:                    writeDBMap.Db,
			ReadDB:                     readDBMap.Db,
			TimeZone:                   timeZone,
			TimeFormat:                 timeFormat,
			TimePostgresFriendlyFormat: timePostgresFriendlyFormat,
			TimeLocation:               timeLocation,
		}
	}

	return
}

func (d *db) GetVehicle(ctx context.Context, motorBrand string, motorType string, motorTransmission string) (returnData []entity.Vehicle, err error) {
	_, err = hunch.Waterfall(
		ctx,
		func(c context.Context, _ interface{}) (interface{}, error) {
			if motorBrand == "" && motorType == "" && motorTransmission == "" {
				return d.ReadDB.QueryContext(c, QueryGetVehicleNoFilter)
			}

			if motorBrand != "" && motorType == "" && motorTransmission == "" {
				return d.ReadDB.QueryContext(
					c,
					QueryGetVehicleWithMotorBrand,
					motorBrand,
				)
			}

			if motorBrand == "" && motorType != "" && motorTransmission == "" {
				return d.ReadDB.QueryContext(
					c,
					QueryGetVehicleWithMotorType,
					motorType,
				)
			}

			if motorBrand == "" && motorType == "" && motorTransmission != "" {
				return d.ReadDB.QueryContext(
					c,
					QueryGetVehicleWithMotorTransmission,
					motorTransmission,
				)
			}

			return d.ReadDB.QueryContext(
				c,
				QueryGetVehicleWithMotorAll,
				motorBrand,
				motorType,
				motorTransmission,
			)
		},
		func(c context.Context, i interface{}) (interface{}, error) {
			rows := i.(*sql.Rows)
			return nil, scan.RowsStrict(&returnData, rows)
		},
	)

	log.AddStrOrAddErr(
		ctx,
		err,
		"[failed] GetVehicle from db",
		"[success] GetVehicle from db",
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		err = core.ErrDBConn
	}

	return
}

func (d *db) DelCache(_ context.Context, _ bool, _ ...string) (err error) {
	return
}
