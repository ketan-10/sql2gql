SET FOREIGN_KEY_CHECKS=0;

DROP TABLE IF EXISTS `abilities`;
CREATE TABLE `abilities` (
    `id` int NOT NULL AUTO_INCREMENT,
    `ability_name` varchar(255) NOT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `ability_name_index` (`ability_name`),
    FULLTEXT KEY `ability_name_fulltext` (ability_name)
) ENGINE=INNODB;


DROP TABLE IF EXISTS `base_stats`;
CREATE TABLE `base_stats` (
  `id` int NOT NULL AUTO_INCREMENT,
  `fk_pokemon` int NOT NULL,
  `b_hp` int(11) DEFAULT NULL,
  `b_atk` int(11) DEFAULT NULL,
  `b_def` int(11) DEFAULT NULL,
  `b_sp_atk` int(11) DEFAULT NULL,
  `b_sp_def` int(11) DEFAULT NULL,
  `b_speed` int(11) DEFAULT NULL,
  `active` BOOLEAN NOT NULL DEFAULT true,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  
  PRIMARY KEY (`id`),
  FOREIGN KEY (`fk_pokemon`) REFERENCES `pokemon` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=INNODB;

DROP TABLE IF EXISTS `pokemon`;
CREATE TABLE `pokemon` (
  `id` int NOT NULL AUTO_INCREMENT,
  `pokemon_name` varchar(79) NOT NULL,
  `pokemon_height` int(11) DEFAULT NULL,
  `pokemon_weight` int(11) DEFAULT NULL,
  `pokemon_base_experience` int(11) DEFAULT NULL,
  `active` BOOLEAN NOT NULL DEFAULT true,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `pokemon_name_index` (`pokemon_name`),
  FULLTEXT KEY `pokemon_name_fulltext` (pokemon_name)
) ENGINE=INNODB;


DROP TABLE IF EXISTS `pokemon_abilities`;
CREATE TABLE `pokemon_abilities` (
    `id` int NOT NULL AUTO_INCREMENT,
    `fk_pokemon` int NOT NULL,
    `fk_ability` int NOT NULL,
    `is_hidden` tinyint(1) NOT NULL,
    `slot` int(11) NOT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`fk_pokemon`) REFERENCES `pokemon` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (`fk_ability`) REFERENCES `abilities` (`id`)  ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=INNODB;

DROP TABLE IF EXISTS `pokemon_evolution`;
CREATE TABLE `pokemon_evolution` (
    `id` int NOT NULL AUTO_INCREMENT,
    `evolved_species_id` int NOT NULL,
    `evol_minimum_level` int(11) DEFAULT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`evolved_species_id`) REFERENCES `pokemon_evolution_matchup` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=INNODB;


DROP TABLE IF EXISTS `pokemon_evolution_matchup`;
CREATE TABLE `pokemon_evolution_matchup` (
    `id` int NOT NULL AUTO_INCREMENT,
    `pokemon_id` int NOT NULL,
    `evolves_from_species_id` int(11) DEFAULT NULL,
    `habitat` ENUM('cave','forest','grassland','mountain','rare','rough_terrain','sea','urban','waters_edge','Unknown') NULL,
    `gender_rate` int(11) NOT NULL,
    `capture_rate` int(11) NOT NULL,
    `base_happiness` int(11) NOT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`pokemon_id`) REFERENCES `pokemon` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=INNODB;



DROP TABLE IF EXISTS `pokemon_types`;
CREATE TABLE `pokemon_types` (
    `id` int NOT NULL AUTO_INCREMENT,
    `pokemon_id` int(11) NOT NULL,
    `type_id` int NOT NULL,
    `slot` int(11) NOT NULL,
    `active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`pokemon_id`) REFERENCES `pokemon` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (`type_id`) REFERENCES `types` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=INNODB;

DROP TABLE IF EXISTS `types`;
CREATE TABLE `types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type_name` varchar(79) NOT NULL,
  `active` BOOLEAN NOT NULL DEFAULT true,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=INNODB;
