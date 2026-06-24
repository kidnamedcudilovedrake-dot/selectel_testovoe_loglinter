package a

import (
	"context"
	"log"
	"log/slog"

	"go.uber.org/zap"
)

func TestLogs(ctx context.Context, logger *slog.Logger, zapLogger *zap.Logger, sugar *zap.SugaredLogger) {
	slog.Info("Starting server") // want "log message should start with a lowercase letter"
	slog.Info("starting server")

	slog.Info("запуск сервера") // want "log message should be in English only"

	slog.Info("server started! 🚀")       // want "log message contains emoji, contains forbidden character '!'"
	slog.Error("connection failed!!!")   // want "log message contains forbidden character '!'"
	slog.Warn("warning: something went wrong...") // want "log message contains ellipsis or consecutive dots, contains redundant level prefix 'warning:'"

	password := "1234"
	token := "abc"
	apiKey := "key123"
	config := map[string]string{"password": "123"}

	slog.Info("user password: " + password)          // want "log call contains potentially sensitive variable \"password\""
	slog.Info("key: " + apiKey)                      // want "log call contains potentially sensitive variable \"apiKey\""
	slog.Info("token: " + token)                      // want "log call contains potentially sensitive variable \"token\""
	slog.Info("config token: " + config["password"]) // want "log call contains potentially sensitive map key \"password\""

	log.Println("Starting standard server") // want "log message should start with a lowercase letter"
	log.Printf("Starting: %d", 123)         // want "log message should start with a lowercase letter"

	zapLogger.Info("Failed to run")      // want "log message should start with a lowercase letter"
	sugar.Infof("Starting %s", "server") // want "log message should start with a lowercase letter"
}
