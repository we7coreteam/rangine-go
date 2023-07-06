package logger

type Config struct {
	Driver     string         `mapstructure:"driver" json:"driver" yaml:"driver" binding:"required"`
	Path       string         `mapstructure:"path" json:"path" yaml:"path" binding:"required"`
	Level      string         `mapstructure:"level" json:"level" yaml:"level" binding:"required"`
	MaxDays    int            `mapstructure:"days" json:"days" yaml:"days"`
	MaxSize    int            `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups int            `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	Options    map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
