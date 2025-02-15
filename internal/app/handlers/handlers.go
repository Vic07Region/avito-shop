package handlers

import (
	"context"
	"errors"
	"github.com/Vic07Region/avito-shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type ServiceInterface interface {
	GetWalletInfo(ctx context.Context, userID uuid.UUID) (*service.FullInfo, error)
	SendCoins(ctx context.Context, userID uuid.UUID, toUsername string, amount int) error
	PurchaseMerch(ctx context.Context, userID uuid.UUID, merchName string, quantity int) error
	LoginUser(ctx context.Context, userdata service.UserData) (string, error)
}

type Handlers struct {
	Service ServiceInterface
	log     *zap.Logger
}

func New(srv ServiceInterface, zapLogger *zap.Logger) *Handlers {
	return &Handlers{
		Service: srv,
		log:     zapLogger,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required,alphanum"`
	Amount int    `json:"amount" binding:"required,min=1"`
}

func (h *Handlers) AuthUser(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "request body is empty",
			})
			return
		}

		if vErr, ok := err.(validator.ValidationErrors); ok {
			for _, e := range vErr {
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": e.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()

	token, err := h.Service.LoginUser(ctx, service.UserData{
		Username: strings.ToLower(req.Username),
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvelidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": "login service filed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func (h *Handlers) WalletInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "unauthorized",
		})
		return
	}
	ctx := c.Request.Context()
	walletInfo, err := h.Service.GetWalletInfo(ctx, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "wallet service failed",
		})
		return
	}

	var inventoryList = []Inventory{}
	var sentList = []Sent{}
	var receivedList = []Received{}

	if walletInfo.Inventory != nil {
		for _, i := range walletInfo.Inventory {
			inventoryList = append(inventoryList, Inventory{
				Type:     i.Type,
				Quantity: i.Quantity,
			})
		}
	}
	if walletInfo.CoinHistory.Sent == nil {
		for _, s := range walletInfo.CoinHistory.Sent {
			sentList = append(sentList, Sent{
				ToUser: s.ToUser,
				Amount: s.Amount,
			})
		}
	}
	if walletInfo.CoinHistory.Received == nil {
		for _, r := range walletInfo.CoinHistory.Received {
			receivedList = append(receivedList, Received{
				FromUser: r.FromUser,
				Amount:   r.Amount,
			})
		}
	}

	c.JSON(http.StatusOK, FullInfo{
		Coins:     walletInfo.Coins,
		Inventory: inventoryList,
		CoinHistory: CoinHistory{
			Received: receivedList,
			Sent:     sentList,
		},
	})
}

func (h *Handlers) SendCoin(c *gin.Context) {
	var req SendCoinRequest

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "unauthorized",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "request body is empty",
			})
			return
		}

		if vErr, ok := err.(validator.ValidationErrors); ok {
			for _, e := range vErr {
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": e.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	err := h.Service.SendCoins(ctx, userID.(uuid.UUID), req.ToUser, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrNotEnoughCoins) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": "wallet service failed",
		})
		return
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "success",
	//})
	c.Status(http.StatusOK)

}

func (h *Handlers) BuyMerch(c *gin.Context) {
	merchName := c.Param("merchName")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": "unauthorized",
		})
		return
	}

	ctx := c.Request.Context()

	err := h.Service.PurchaseMerch(ctx, userID.(uuid.UUID), merchName, 1)
	if err != nil {
		if errors.Is(err, service.ErrMerchNotFound) || errors.Is(err, service.ErrNotEnoughCoins) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": "wallet service failed",
		})
		return
	}

	c.Status(http.StatusOK)

}
