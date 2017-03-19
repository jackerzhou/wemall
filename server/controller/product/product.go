package product

import (
	"../../config"
	"../../model"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/jinzhu/gorm"
	"gopkg.in/kataras/iris.v6"
)

// List 产品列表
func List(ctx *iris.Context) {
	var products []model.Product
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)

	if err != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : "error",
			"data"  : iris.Map{},
		})
		return
	}

	defer db.Close()

	pageNo, err := strconv.Atoi(ctx.FormValue("pageNo"))
 
	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	offset   := (pageNo - 1) * config.ServerConfig.PageSize
	queryErr := db.Offset(offset).Limit(config.ServerConfig.PageSize).Find(&products).Error

	if queryErr != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : "error.",
			"data"  : iris.Map{},
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"errNo" : model.ErrorCode.SUCCESS,
		"msg"   : "success",
		"data"  : iris.Map{
			"products": products,
		},
	})
}

// Create 创建产品
func Create(ctx *iris.Context) {
	var product model.Product
	err := ctx.ReadJSON(&product)	

	if err != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : "参数无效",
			"data"  : iris.Map{},
		})
		return
	}

	db, connErr := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)

	if connErr != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : "error",
			"data"  : iris.Map{},
		})
		return
	}

	defer db.Close()

	var isError bool
	var msg = ""

	product.BrowseCount = 0
	product.BuyCount    = 0
	product.TotalSale   = 0
	product.Status      = model.ProductPending

	product.Name = strings.TrimSpace(product.Name)
	if (product.Name == "") {
		isError = true
		msg     = "商品名称不能为空"
	} else if utf8.RuneCountInString(product.Name) > config.ServerConfig.MaxNameLen {
		isError = true
		msg     = "商品名称不能超过" + strconv.Itoa(config.ServerConfig.MaxNameLen) + "个字符"
	} else if product.Status != model.ProductPending {
		isError = true
		msg     = "status无效"
	} else if product.Remark != "" && utf8.RuneCountInString(product.Remark) > config.ServerConfig.MaxRemarkLen {
		isError = true
		msg     = "备注不能超过" + strconv.Itoa(config.ServerConfig.MaxRemarkLen) + "个字符"	
	} else if product.Detail == "" || utf8.RuneCountInString(product.Detail) <= 0 {
		isError = true	
		msg     = "商品详情不能为空"
	}  else if utf8.RuneCountInString(product.Detail) > config.ServerConfig.MaxContentLen {
		isError = true	
		msg     = "商品详情不能超过" + strconv.Itoa(config.ServerConfig.MaxContentLen) + "个字符"
	} else if product.Categories == nil || len(product.Categories) <= 0  {
		msg     = "至少要选择一个商品分类"
	} else if len(product.Categories) > config.ServerConfig.MaxProductCateCount  {
		msg     = "最多只能选择" + strconv.Itoa(config.ServerConfig.MaxProductCateCount) + "个商品分类"
	}

	if isError {
		ctx.JSON(iris.StatusOK, iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : msg,
			"data"  : iris.Map{},
		})
		return
	}

	for i := 0; i < len(product.Categories); i++ {
		var category model.Category
		queryErr := db.First(&category, product.Categories[i].ID).Error
		if queryErr != nil {
			ctx.JSON(iris.StatusOK, iris.Map{
				"errNo" : model.ErrorCode.ERROR,
				"msg"   : "无效的分类id",
				"data"  : iris.Map{},
			})
			return	
		}
		product.Categories[i] = category
	}

	db.LogMode(true)

	saveErr := db.Create(&product).Error

	if saveErr != nil {
		ctx.JSON(iris.StatusOK, iris.Map{
			"errNo" : model.ErrorCode.ERROR,
			"msg"   : "error.",
			"data"  : iris.Map{},
		})
		return	
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"errNo" : model.ErrorCode.SUCCESS,
		"msg"   : "success",
		"data"  : product,
	})
}

// Update 更新产品
func Update(ctx *iris.Context) {
	
}
