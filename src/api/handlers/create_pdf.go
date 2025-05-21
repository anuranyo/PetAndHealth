package handlers

import (
	"fmt"
	"github.com/NureTymofiienkoSnizhana/arkpz-pzpi-22-9-tymofiienko-snizhana/Pract1/arkpz-pzpi-22-9-tymofiienko-snizhana-task2/src/data"
	"github.com/jung-kurt/gofpdf"
	"time"
)

func GeneratePetReportPDF(pet *data.Pet, healthData []*data.HealthData, ownerName string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("DejaVu", "", "fonts/DejaVuSans.ttf")

	pdf.SetFont("DejaVu", "", 14)
	pdf.AddPage()

	pdf.Cell(40, 10, "Звіт")
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Ім'я: %s", pet.Name))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Вид: %s", pet.Species))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Порода: %s", pet.Breed))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Вік: %d", pet.Age))
	pdf.Ln(8)
	pdf.Cell(0, 10, fmt.Sprintf("Власник: %s", ownerName))
	pdf.Ln(8)

	pdf.Ln(10)
	pdf.Cell(0, 10, "Дані стану здоров'я:")
	pdf.Ln(10)
	for _, health := range healthData {
		pdf.Cell(0, 10, fmt.Sprintf("Активність: %.1f", health.Activity))
		pdf.Ln(6)
		pdf.Cell(0, 10, fmt.Sprintf("Сон: %.1f", health.SleepHours))
		pdf.Ln(6)
		pdf.Cell(0, 10, fmt.Sprintf("Температура: %.1f", health.Temperature))
		pdf.Ln(6)
		pdf.Cell(0, 10, fmt.Sprintf("Час: %s", time.Unix(int64(health.Time.T), 0).Format("2006-01-02 15:04:05")))
		pdf.Ln(10)
	}

	return pdf, nil
}
