package http

import (
	"net/http"
	"time"

	"github.com/Erwanph/be-wan-central-lab/internal/model"
	"github.com/Erwanph/be-wan-central-lab/internal/usecase"
	"github.com/Erwanph/be-wan-central-lab/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)