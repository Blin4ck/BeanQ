package http

import (
	"coffe/internal/menu/delivery/http/dto"
	"coffe/internal/menu/entity"
	"coffe/internal/menu/usecase"
	"coffe/internal/middleware"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuHandler struct {
	middleware  *middleware.JWTMiddleware
	menuUsecase *usecase.MenuUsecase
}

func NewMenuHandler(middleware *middleware.JWTMiddleware, menuUsecase *usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{
		middleware:  middleware,
		menuUsecase: menuUsecase,
	}
}

func (h *MenuHandler) GetAllMenuItems(ctx *gin.Context) {

	menus, err := h.menuUsecase.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"menus": menus,
		"total": len(menus),
	})
}

func (h *MenuHandler) GetMenuItemByID(ctx *gin.Context) {

	idParam := ctx.Param("id")

	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID меню обязателен"})
		return
	}

	munuItemID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID меню"})
		return
	}

	menuItem, err := h.menuUsecase.GetMenuItemByID(ctx, munuItemID)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Позиция меню не найдена"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"menu_item": menuItem,
	})

}

func (h *MenuHandler) GetMenuItemsByCategory(ctx *gin.Context) {

	category := ctx.Param("category_id")

	if category == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Категория обязательна"})
		return
	}

	categoryId, err := uuid.Parse(category)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID категории"})
		return
	}

	menuIteams, err := h.menuUsecase.GetItemsByCategory(ctx, categoryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"menu_items": menuIteams,
		"total":      len(menuIteams),
	})
}

func (h *MenuHandler) SearchMenuItems(ctx *gin.Context) {
	var params struct {
		Query       string  `form:"query"`
		MenuID      string  `form:"menu_id"`
		CategoryID  string  `form:"category_id"`
		ProductID   string  `form:"product_id"`
		IsActive    *bool   `form:"is_active"`
		MinPrice    float64 `form:"min_price"`
		MaxPrice    float64 `form:"max_price"`
		ValidAfter  string  `form:"valid_after"`
		ValidBefore string  `form:"valid_before"`
		Page        int     `form:"page" binding:"min=1"`
		PageSize    int     `form:"page_size" binding:"min=5,max=100"`
		SortBy      string  `form:"sort_by" binding:"oneof=name price created_at sort_order"`
		SortOrder   string  `form:"sort_order" binding:"oneof=asc desc"`
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query params: " + err.Error()})
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}
	if params.SortBy == "" {
		params.SortBy = "sort_order"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	var menuID, categoryID, productID uuid.UUID
	var err error

	if params.MenuID != "" {
		menuID, err = uuid.Parse(params.MenuID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu_id format"})
			return
		}
	}

	if params.CategoryID != "" {
		categoryID, err = uuid.Parse(params.CategoryID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id format"})
			return
		}
	}

	if params.ProductID != "" {
		productID, err = uuid.Parse(params.ProductID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product_id format"})
			return
		}
	}

	var validAfter, validBefore *time.Time
	if params.ValidAfter != "" {
		parsed, err := time.Parse("2006-01-02", params.ValidAfter)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid valid_after format, use YYYY-MM-DD"})
			return
		}
		validAfter = &parsed
	}

	if params.ValidBefore != "" {
		parsed, err := time.Parse("2006-01-02", params.ValidBefore)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid valid_before format, use YYYY-MM-DD"})
			return
		}
		validBefore = &parsed
	}

	searchDTO := dto.MenuSearchDTO{
		Query:       params.Query,
		MenuID:      menuID,
		CategoryID:  categoryID,
		ProductID:   productID,
		IsActive:    params.IsActive,
		PriceRange:  [2]float64{params.MinPrice, params.MaxPrice},
		ValidAfter:  validAfter,
		ValidBefore: validBefore,
		Pagination: dto.Pagination{
			Page:     params.Page,
			PageSize: params.PageSize,
		},
		Sorting: dto.Sorting{
			Field: params.SortBy,
			Order: params.SortOrder,
		},
	}

	results, total, err := h.menuUsecase.SearchMenuItems(ctx, searchDTO)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed: " + err.Error()})
		return
	}

	// Формируем ответ
	ctx.JSON(http.StatusOK, gin.H{
		"data":        results,
		"total":       total,
		"page":        params.Page,
		"page_size":   params.PageSize,
		"total_pages": int(math.Ceil(float64(total) / float64(params.PageSize))),
	})
}

func (h *MenuHandler) GetAvailableItems(ctx *gin.Context) {

	items, err := h.menuUsecase.GetActiveItems(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": len(items),
	})
}

func (h *MenuHandler) CreateMenuItem(ctx *gin.Context) {
	var items entity.MenuItem

	if err := ctx.ShouldBindJSON(&items); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных "})
		return
	}

	if err := h.menuUsecase.CreateMenuItem(ctx, &items); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Позиция успешно добавлена"})

}

// вынести валидацию в отдельную вспомогательную функциию
func (h *MenuHandler) UpdateMenuItem(ctx *gin.Context) {
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать product_id"})
		return
	}

	itemID, err := uuid.Parse(productID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	var item entity.MenuItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса: " + err.Error()})
		return
	}

	item.ID = itemID

	if item.MenuID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать menu_id"})
		return
	}

	if item.ProductID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать product_id"})
		return
	}

	if item.CategoryID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать category_id"})
		return
	}

	if item.SortOrder < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "sort_order не может быть отрицательным"})
		return
	}

	if err := h.menuUsecase.UpdateMenuItem(ctx, &item); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Элемент меню не найден"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (h *MenuHandler) DeleteMenuItem(ctx *gin.Context) {

	item := ctx.Param("product_id")

	if item == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Категория обязательна"})
		return
	}

	itemId, err := uuid.Parse(item)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID позиции"})
		return
	}

	if err := h.menuUsecase.DeleteMenuItem(ctx, itemId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Позиция удалена из меню"})

}

func (h *MenuHandler) UpdateAvailability(ctx *gin.Context) {
	itemID := ctx.Param("menu_id")
	if itemID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Необходимо указать ID меню",
		})
		return
	}

	id, ok := uuid.Parse(itemID)
	if ok != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат ID",
		})
		return
	}

	var request struct {
		IsActive bool `json:"is_active" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных: " + err.Error(),
		})
		return
	}

	var err error
	if request.IsActive {
		err = h.menuUsecase.Activate(ctx, id)
	} else {
		err = h.menuUsecase.Deactivate(ctx, id)
	}

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		ctx.AbortWithStatusJSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Доступность меню обновлена",
	})
}

func (h *MenuHandler) GetCategories(ctx *gin.Context) {

	category, err := h.menuUsecase.GetCategories(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"categories": category,
		"total":      len(category),
	})

}

func (h *MenuHandler) CreateCategory(ctx *gin.Context) {
	var category entity.MenuCategory

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных категории"})
		return
	}
	err := h.menuUsecase.CreateCategory(ctx, &category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Категория успешно создана", "category": category})
}

func (h *MenuHandler) UpdateCategory(ctx *gin.Context) {

	iParam, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID категории"})
		return
	}

	var category entity.MenuCategory

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных категории"})
		return
	}

	category.ID = iParam

	err = h.menuUsecase.UpdateCategory(ctx, &category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Категория успешно обновлена"})

}

func (h *MenuHandler) DeleteCategory(ctx *gin.Context) {

	idParam, err := uuid.Parse(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID категории"})
		return
	}
	err = h.menuUsecase.DeleteCategory(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Категория успешно удалена"})
}

func (h *MenuHandler) MenuItemActivation(ctx *gin.Context) {
	itemID := ctx.Param("id")
	if itemID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Необходимо указать ID позиции",
		})
		return
	}

	id, err := uuid.Parse(itemID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат ID",
		})
		return
	}

	var request struct {
		IsActive bool `json:"is_active" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Неверный формат данных: " + err.Error(),
		})
		return
	}

	if request.IsActive {
		err = h.menuUsecase.ActivateMenuItem(ctx, id)
	} else {
		err = h.menuUsecase.DeactivateMenuItem(ctx, id)
	}

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		ctx.AbortWithStatusJSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Статус позиции обновлен",
	})
}
