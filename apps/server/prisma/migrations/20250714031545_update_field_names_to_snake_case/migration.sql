/*
  Warnings:

  - You are about to drop the `essay` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropTable
DROP TABLE `essay`;

-- CreateTable
CREATE TABLE `essays` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `content` TEXT NOT NULL,
    `essay_type` VARCHAR(30) NOT NULL DEFAULT '',
    `file` VARCHAR(255) NOT NULL DEFAULT '',
    `file_id` VARCHAR(255) NOT NULL DEFAULT '',
    `sentences` JSON NULL,
    `task_id` VARCHAR(255) NOT NULL DEFAULT '',
    `telegraph` VARCHAR(255) NOT NULL DEFAULT '',
    `thumb` VARCHAR(255) NOT NULL DEFAULT '',
    `title` VARCHAR(36) NOT NULL DEFAULT '',
    `video_link` VARCHAR(255) NOT NULL DEFAULT '',
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,
    `deleted_at` DATETIME(3) NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
