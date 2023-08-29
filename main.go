package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ROWS = 9
	COLUMNS
)

type Seat struct {
	Available bool `json:"Available"`
	Id        int  `json:"Id"`
}

var seating = make([][]Seat, ROWS)
var count = 1
var counter = 0

func init() {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLUMNS; j++ {
			seating[i] = append(seating[i], Seat{true, count})
			count++
		}
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", indexPage)
	router.GET("/seating", showSeating)
	router.GET("/options", optionFields)
	router.POST("/purchase", purchaseSeat)
	router.POST("/return/:row/:column", undoPurchase)

	router.Run(":8080")
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Movie Theater",
	})
}

func showSeating(c *gin.Context) {
	c.JSON(http.StatusOK, seating)
}

func optionFields(c *gin.Context) {
	switch c.Query("options") {
	case "purchase":
		c.HTML(http.StatusOK, "purchase.html", gin.H{})
	case "return":
		c.HTML(http.StatusOK, "return.html", gin.H{})
	}
}

func purchaseSeat(c *gin.Context) {
	row, err := strconv.Atoi(c.PostForm("row"))
	if err != nil {
		c.HTML(http.StatusOK, "invalid.html", gin.H{})
		return
	}
	column, err := strconv.Atoi(c.PostForm("column"))
	if err != nil {
		c.HTML(http.StatusOK, "invalid.html", gin.H{})
		return
	}

	if row >= ROWS || column >= COLUMNS || row <= 0 || column <= 0 {
		c.HTML(http.StatusOK, "invalid.html", gin.H{})
	} else {
		seat := &seating[row-1][column-1]
		c.HTML(http.StatusOK, "seat.html", gin.H{
			"row":       row,
			"column":    column,
			"available": seating[row-1][column-1].Available,
		})
		seat.Available = false
	}
}

func undoPurchase(c *gin.Context) {
	id := strings.Trim(c.Param("row"), "/")
	row, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	id = strings.Trim(c.Param("column"), "/")
	column, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	info := struct{ Id int }{-1}
	if err = c.BindJSON(&info); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if row >= ROWS || column >= COLUMNS {
		c.JSON(http.StatusBadRequest, "invalid seat option")
		return
	}

	if seat := &seating[row-1][column-1]; !seat.Available && seat.Id == info.Id {
		c.JSON(http.StatusOK, "return processed")
		seat.Available = true
	} else {
		c.JSON(http.StatusOK, "seat invalid for return")
	}
}
