package main

import (
	"net/http"
	"poberhauser/receipt-processor-challenge/points_processor"
	"poberhauser/receipt-processor-challenge/receipt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Receipts sync.Map

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processPoints)
	router.GET("/receipts/:id/points", getPoints)
	router.Run("localhost:8080")
}

func getPoints(c *gin.Context) {
	id := c.Param("id")
	loadedReceipt, ok := Receipts.Load(id)
	if ok {
		typedReceipt := loadedReceipt.(receipt.Receipt)
		c.JSON(http.StatusOK, gin.H{"Points": typedReceipt.Points})
		return
	} else {
		c.String(http.StatusNotFound, "No receipt found for that ID")
		return
	}
}

func processPoints(c *gin.Context) {
	var newReceipt receipt.Receipt

	if err := c.BindJSON(&newReceipt); err != nil {
		c.String(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	var points int64
	points += points_processor.RetailerNamePoints(newReceipt.Retailer)
	points += points_processor.RoundDollarPoints(newReceipt.Total)
	points += points_processor.QuarterMultiplePoints(newReceipt.Total)
	points += points_processor.CountByTwoItemsPoints(newReceipt.Items)
	points += points_processor.TrimmedLengthPoints(newReceipt.Items)
	points += points_processor.OddPurchaseDatePoints(newReceipt.PurchaseDate)
	points += points_processor.PurchaseTimePoints(newReceipt.PurchaseTime)

	newReceipt.Points = points
	ID := uuid.New().String()
	Receipts.Store(ID, newReceipt)
	c.JSON(http.StatusOK, gin.H{"ID": ID})

}
