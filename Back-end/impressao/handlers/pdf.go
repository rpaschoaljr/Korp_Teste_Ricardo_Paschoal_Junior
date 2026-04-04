package handlers

import (
	"bytes"
	"fmt"
	"impressao_service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(c *gin.Context) {
	var inv models.Invoice
	if err := c.ShouldBindJSON(&inv); err != nil {
		fmt.Printf("[DEBUG-PDF] Erro ao bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados invalidos: " + err.Error()})
		return
	}

	fmt.Printf("[DEBUG-PDF] Recebido pedido para Nota: %d\n", inv.ID)
	fmt.Printf("[DEBUG-PDF] Cliente recebido: %+v\n", inv.Cliente)
	fmt.Printf("[DEBUG-PDF] Itens recebidos: %d\n", len(inv.Itens))

	pdf := gofpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()

	pdf.SetFillColor(44, 62, 80)
	pdf.Rect(0, 0, 210, 40, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(190, 15, tr("KORP ERP - NOTA FISCAL"), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 5, tr(fmt.Sprintf("Nº DA NOTA: %06d", inv.ID)), "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 5, tr(fmt.Sprintf("DATA DE EMISSÃO: %s", inv.Data)), "", 1, "C", false, 0, "")
	pdf.Ln(15)

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(236, 240, 241)
	pdf.CellFormat(190, 8, tr(" DADOS DO CLIENTE"), "B", 1, "L", true, 0, "")
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 6, tr(fmt.Sprintf("NOME / RAZÃO SOCIAL: %s", inv.Cliente.Nome)))
	pdf.Ln(6)

	docLabel := "CNPJ"
	docValue := inv.Cliente.CNPJ
	if inv.Cliente.CPF != "" {
		docLabel = "CPF"
		docValue = inv.Cliente.CPF
	}
	pdf.Cell(190, 6, tr(fmt.Sprintf("%s: %s", docLabel, docValue)))
	pdf.Ln(6)
	pdf.Cell(190, 6, tr(fmt.Sprintf("ENDEREÇO: %s", inv.Cliente.Endereco)))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 8, tr(" ITENS DA NOTA FISCAL"), "B", 1, "L", true, 0, "")
	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(90, 8, tr(" DESCRIÇÃO"), "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "QTD", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, tr("UNITÁRIO"), "1", 0, "R", true, 0, "")
	pdf.CellFormat(35, 8, "SUBTOTAL", "1", 1, "R", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, item := range inv.Itens {
		pdf.CellFormat(90, 8, " "+tr(item.Descricao), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", item.Quantidade), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("R$ %.2f", item.PrecoUnitario), "1", 0, "R", false, 0, "")
		pdf.CellFormat(35, 8, fmt.Sprintf("R$ %.2f ", item.Subtotal), "1", 1, "R", false, 0, "")
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(39, 174, 96)
	pdf.CellFormat(155, 12, tr("VALOR TOTAL DA NOTA: "), "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 12, tr(fmt.Sprintf("R$ %.2f ", inv.ValorTotal)), "1", 1, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		fmt.Printf("[DEBUG-PDF] Erro ao gerar PDF: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar PDF: " + err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

func HealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
