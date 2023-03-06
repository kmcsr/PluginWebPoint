
USE pluginDatabase;

CREATE TABLE IF NOT EXISTS plugins (
	`id` VARCHAR(64) NOT NULL,
	PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE plugins ADD `name`       VARCHAR(64) NOT NULL;
ALTER TABLE plugins ADD `enabled`    BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE plugins ADD `version`    VARCHAR(32) NOT NULL;
ALTER TABLE plugins ADD `authors`    VARCHAR(64) NOT NULL;
ALTER TABLE plugins ADD `desc`       VARCHAR(256) DEFAULT '' NOT NULL;
ALTER TABLE plugins ADD `desc_zhCN`  VARCHAR(256) DEFAULT '' NOT NULL;
ALTER TABLE plugins ADD `link`       VARCHAR(128) DEFAULT '' NOT NULL;
ALTER TABLE plugins ADD `lastUpdate` DATETIME DEFAULT NULL;
ALTER TABLE plugins ADD `labels`     VARCHAR(64) NOT NULL;

INSERT INTO plugins (`enabled`, `id`, `name`, `version`, `authors`, `desc`, `link`, `labels`) VALUES (
	TRUE, "kpi", "KPI", "0.0.1", "zyxkad", "A MCDR plugins codes share library", "https://github.com/kmcsr/kpi_mcdr", "api"
);
