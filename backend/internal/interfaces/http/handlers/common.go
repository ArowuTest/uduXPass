package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/uduxpass/backend/internal/domain/entities"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// validateStruct validates a struct using the validator package
func validateStruct(s interface{}) error {
	return validate.Struct(s)
}

// handleError handles different types of errors and returns appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *entities.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation error",
			"field": e.Field,
			"message": e.Message,
		})
	case *entities.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Resource not found",
			"resource": e.Resource,
		})
	case *entities.ConflictError:
		c.JSON(http.StatusConflict, gin.H{
			"error": "Conflict",
			"message": e.Message,
		})
	case *entities.BusinessRuleError:
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Business rule violation",
			"message": e.Message,
		})
	default:
		// Check for specific error types
		if errors.Is(err, entities.ErrInsufficientInventory) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Insufficient inventory",
				"message": "Not enough tickets available",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"message": err.Error(),
		})
	}
}

// parseUUID parses a UUID from a string parameter
func parseUUID(c *gin.Context, param string) (uuid.UUID, bool) {
	idStr := c.Param(param)
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameter",
			"parameter": param,
		})
		return uuid.Nil, false
	}
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid UUID format",
			"parameter": param,
		})
		return uuid.Nil, false
	}
	
	return id, true
}

// parseQueryUUID parses a UUID from a query parameter
func parseQueryUUID(c *gin.Context, param string) (*uuid.UUID, error) {
	idStr := c.Query(param)
	if idStr == "" {
		return nil, nil
	}
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	
	return &id, nil
}

// parseQueryInt parses an integer from a query parameter
func parseQueryInt(c *gin.Context, param string, defaultValue int) int {
	valueStr := c.Query(param)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}

// parseQueryFloat parses a float from a query parameter
func parseQueryFloat(c *gin.Context, param string) (*float64, error) {
	valueStr := c.Query(param)
	if valueStr == "" {
		return nil, nil
	}
	
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, err
	}
	
	return &value, nil
}

// parseQueryBool parses a boolean from a query parameter
func parseQueryBool(c *gin.Context, param string) (*bool, error) {
	valueStr := c.Query(param)
	if valueStr == "" {
		return nil, nil
	}
	
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return nil, err
	}
	
	return &value, nil
}

// parseQueryTime parses a time from a query parameter
func parseQueryTime(c *gin.Context, param string) (*time.Time, error) {
	valueStr := c.Query(param)
	if valueStr == "" {
		return nil, nil
	}
	
	// Try different time formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, valueStr); err == nil {
			return &t, nil
		}
	}
	
	return nil, fmt.Errorf("invalid time format")
}

// getPaginationParams extracts pagination parameters from query
func getPaginationParams(c *gin.Context) (page, limit int, sortBy, sortOrder string) {
	page = parseQueryInt(c, "page", 1)
	limit = parseQueryInt(c, "limit", 20)
	sortBy = c.DefaultQuery("sort_by", "created_at")
	sortOrder = c.DefaultQuery("sort_order", "desc")
	
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	
	return
}

// getSearchParam extracts and cleans search parameter
func getSearchParam(c *gin.Context) string {
	search := strings.TrimSpace(c.Query("search"))
	if len(search) > 100 {
		search = search[:100]
	}
	return search
}

// bindAndValidate binds JSON request and validates it
func bindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return false
	}
	
	if err := validateStruct(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return false
	}
	
	return true
}

// successResponse returns a success response with data
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": data,
	})
}

// createdResponse returns a created response with data
func createdResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": data,
	})
}

// noContentResponse returns a no content response
func noContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// errorResponse returns an error response
func errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error": message,
	})
}

// validationErrorResponse returns a validation error response
func validationErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error": "Validation error",
		"field": field,
		"message": message,
	})
}


// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

