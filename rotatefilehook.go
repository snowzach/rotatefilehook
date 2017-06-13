package rotatefilehook

import (
    "io"
    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"
)

type RotateFileConfig struct {
    Filename    string
    MaxSize     int
    MaxBackups  int
    MaxAge      int
    Level       logrus.Level
    Formatter   logrus.Formatter
}

type RotateFileHook struct {
    Config      RotateFileConfig
    logWriter   io.Writer
}

func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {

    hook := RotateFileHook{
        Config: config,
    }

    var zeroLevel logrus.Level
    if hook.Config.Level == zeroLevel {
        hook.Config.Level = logrus.ErrorLevel
    }
    var zeroFormatter logrus.Formatter
    if hook.Config.Formatter == zeroFormatter {
        hook.Config.Formatter = new(logrus.TextFormatter)
    }

    hook.logWriter = &lumberjack.Logger{
        Filename:   config.Filename,
        MaxSize:    config.MaxSize,
        MaxBackups: config.MaxBackups,
        MaxAge:     config.MaxAge,
    }

    return &hook, nil
}

func (hook *RotateFileHook) Levels() []logrus.Level {
    return []logrus.Level{
        logrus.PanicLevel,
        logrus.FatalLevel,
        logrus.ErrorLevel,
        logrus.WarnLevel,
        logrus.InfoLevel,
        logrus.DebugLevel,
    }
}

func (hook *RotateFileHook) Fire(entry *logrus.Entry) (err error) {
    if hook.Config.Level < entry.Level {
        return nil
    }
    b, err := hook.Config.Formatter.Format(entry)
    if err != nil {
        return err
    }
    hook.logWriter.Write(b)
    return nil
}
