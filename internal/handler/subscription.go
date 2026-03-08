package handler

import (
	"net/http"
	"strconv"
	"lighthouse/internal/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

// GetPlans godoc
// @Summary      Get subscription plans
// @Description  Get all available subscription plans
// @Tags         subscription
// @Produce      json
// @Success      200     {object}  map[string]interface{}
// @Router       /subscription/plans [get]
func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans, err := h.subscriptionService.GetPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// GetMySubscription godoc
// @Summary      Get current user subscription
// @Description  Get active subscription for current user
// @Tags         subscription
// @Security     BearerAuth
// @Produce      json
// @Success      200     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]string
// @Router       /subscription/me [get]
func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	subscription, err := h.subscriptionService.GetUserActiveSubscription(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if subscription == nil {
		c.JSON(http.StatusOK, gin.H{"subscription": nil, "has_subscription": false})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"subscription": subscription, "has_subscription": true})
}

// GetMySubscriptions godoc
// @Summary      Get all user subscriptions
// @Description  Get subscription history for current user
// @Tags         subscription
// @Security     BearerAuth
// @Produce      json
// @Success      200     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]string
// @Router       /subscription/history [get]
func (h *SubscriptionHandler) GetMySubscriptions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	subscriptions, err := h.subscriptionService.GetUserSubscriptions(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// CreateSubscription godoc
// @Summary      Create subscription
// @Description  Subscribe to a plan (payment integration needed)
// @Tags         subscription
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Subscription data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Router       /subscription/subscribe [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		PlanID          int    `json:"plan_id" binding:"required"`
		PaymentProvider string `json:"payment_provider"`
		PaymentID       string `json:"payment_id"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	subscription, err := h.subscriptionService.CreateSubscription(
		userID.(string), req.PlanID, req.PaymentProvider, req.PaymentID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"subscription": subscription})
}

// CancelSubscription godoc
// @Summary      Cancel subscription
// @Description  Cancel active subscription
// @Tags         subscription
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /subscription/cancel/:id [post]
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	
	if err := h.subscriptionService.CancelSubscription(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled"})
}
