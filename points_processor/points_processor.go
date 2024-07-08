package points_processor

import (
	"math"
	"poberhauser/receipt-processor-challenge/receipt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// One point for every alphanumeric character in the retailer name.
func RetailerNamePoints(retailer string) int64 {
	var count int64 = 0
	for _, r := range retailer {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			count++
		}
	}
	return count

}

// 50 points if the total is a round dollar amount with no cents.
func RoundDollarPoints(total string) int64 {
	var points int64 = 0
	if strings.Contains(total, ".00") {
		points = 50
	}

	return points
}

// 25 points if the total is a multiple of .25.
func QuarterMultiplePoints(total string) int64 {
	var points int64 = 0
	floatTotal, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return points
	}

	newReceiptQuarterMultiple := int(floatTotal * 100)

	if newReceiptQuarterMultiple%25 == 0 {
		points = 25
	}
	return points
}

// 5 points for every two items on the receipt.
func CountByTwoItemsPoints(items receipt.Items) int64 {
	lenItems := int64(len(items))
	itemTotalPoints := (lenItems / 2) * 5
	return itemTotalPoints
}

// If the trimmed length description is a multiple of 3 then multiply by .2 and round up.
// Result is the amount of points.
func TrimmedLengthPoints(items receipt.Items) int64 {
	var points int64 = 0
	for i := range items {
		trimmedDescription := strings.TrimSpace(items[i].ShortDescription)
		trimmedDescriptionLength := utf8.RuneCountInString(trimmedDescription)
		if trimmedDescriptionLength%3 == 0 {
			price, err := strconv.ParseFloat(items[i].Price, 64)
			points += int64(math.Ceil(price * 0.2))
			if err != nil {
				return points
			}
		}
	}
	return points
}

// 6 points if the day in the purchase date is odd.
func OddPurchaseDatePoints(purchaseDate string) int64 {
	var points int64 = 0
	purchasedate, err := strconv.ParseInt(purchaseDate[len(purchaseDate)-2:], 0, 64)
	if err != nil {
		return points
	}
	if purchasedate%2 != 0 {
		points = 6
	}
	return points

}

// 10 points if the time of purchase is after 2 pm and before 4 pm.
func PurchaseTimePoints(purchaseTime string) int64 {
	var points int64 = 0
	purchaseHour, err := strconv.ParseFloat(purchaseTime[0:2], 64)
	if err != nil {
		return points
	}
	if purchaseHour >= 14 && purchaseHour <= 16 {
		points = 10
	}
	return points
}
