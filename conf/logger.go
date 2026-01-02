package conf

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func InitLogger() *zap.Logger {
	/*===========================================================================================*/
	//設置要寫到日誌的錯誤級別  依照環境變數設定 輸出對應層級之上的日誌
	logMode := zapcore.InfoLevel
	if viper.GetString("GIN_MODE") == "debug" {
		logMode = zapcore.DebugLevel
	}
	core := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(getWriteSyncer(), zapcore.AddSync(os.Stdout)), logMode)
	return zap.New(core)
}

func getEncoder() zapcore.Encoder {
	/*===========================================================================================*/
	//設定日誌的輸出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "utcTime"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02 15:04:05"))
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriteSyncer() zapcore.WriteSyncer {
	/*===========================================================================================*/
	//設定日誌的輸出位置
	stSeparator := string(filepath.Separator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly) + ".txt"
	/*===========================================================================================*/
	//設定日誌的檔案大小上限 和存活時間
	lumberJackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    viper.GetInt("logger.maxSize_MB"),
		MaxBackups: viper.GetInt("logger.maxBackups"),
		MaxAge:     viper.GetInt("logger.maxAge"),
		Compress:   false, // disabled by default
	}

	return zapcore.AddSync(lumberJackSyncer)
}
