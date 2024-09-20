package loggerer

import "github.com/sirupsen/logrus"

//Эта функция создаёт логер
func NewLogger() *logrus.Logger {

	//Создаю стандартный логер потому что тонкой настройки тут не требуется
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Устанавливаю уровень логирования
	logger.SetLevel(logrus.InfoLevel)

	return logger

}
