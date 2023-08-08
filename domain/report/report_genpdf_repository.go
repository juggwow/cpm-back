package report

import (
	"context"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/utils/pdf"
)

func GenPdf(db *connection.DBConnection) genPdfFunc {
	return func(ctx context.Context, id uint) (FileResponse, error) {
		var result ReportPdfDB
		var fileResponse FileResponse
		cpm := db.CPM.Model(&result)
		err := cpm.Table("CPM.fGetReportDetial(?)", id).
			Scan(&result).
			Error
		if err != nil {
			return fileResponse, err
		}

		var result2 FileAttachmentDBs
		cpm = db.CPM.Model(&result2)
		err = cpm.Table("CPM.fGetFileAttachment(?)", id).
			Scan(&result2).
			Error
		if err != nil {
			return fileResponse, err
		}

		report := ReportPdfdata{
			Inspection:   result.Inspection,
			Project:      result.Project,
			Station:      result.Station,
			Company:      result.Company,
			Arrival:      result.Arrival,
			Invoice:      result.Invoice,
			Manufacturer: result.Manufacturer,
			Serial:       result.Serial,
			Receive:      result.Receive,
			Doc:          result2.ToResponse(),
			BoqItem: map[string]string{
				"name": result.ItemName,
				"qty":  result.ItemQuantity,
				"unit": result.ItemUnit,
			},
		}

		obj, err := pdf.GetReport("sample/rad-dev/ReceiveAndDamage-v1.0.1.dito", report)
		fileResponse = FileResponse{
			Obj:  obj,
			Ext:  "application/pdf",
			Name: "testfile.pdf",
		}
		return fileResponse, err
	}
}
