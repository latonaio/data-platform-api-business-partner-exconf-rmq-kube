package main

import (
	"context"
	"data-platform-business-partner-exconf/database"
	"data-platform-business-partner-exconf/database/models"

	"golang.org/x/xerrors"
)

func ExistenceCheck(ctx context.Context, db *database.Mysql, partnerId string) (bool, error) {
	d, err := models.FindDataPlatformBusinessPartnerGeneralDatum(
		ctx, db, partnerId, models.DataPlatformBusinessPartnerGeneralDatumColumns.BusinessPartner,
	)
	if err != nil {
		return false, xerrors.Errorf("cannot get data: %w", err)
	}
	if d == nil {
		return false, nil
	}
	if d.BusinessPartner != partnerId {
		return false, nil
	}
	return true, nil
}
