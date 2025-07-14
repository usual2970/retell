-- CreateTable
CREATE TABLE `essay` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `essayType` VARCHAR(30) NOT NULL DEFAULT '',
    `file` VARCHAR(255) NOT NULL DEFAULT '',
    `fileId` VARCHAR(255) NOT NULL DEFAULT '',
    `sentences` JSON NULL,
    `taskId` VARCHAR(255) NOT NULL DEFAULT '',
    `telegraph` VARCHAR(255) NOT NULL DEFAULT '',
    `thumb` VARCHAR(255) NOT NULL DEFAULT '',
    `title` VARCHAR(36) NOT NULL DEFAULT '',
    `videoLink` VARCHAR(255) NOT NULL DEFAULT '',

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
