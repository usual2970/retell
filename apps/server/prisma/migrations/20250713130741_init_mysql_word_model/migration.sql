-- CreateTable
CREATE TABLE `words` (
    `id` VARCHAR(30) NOT NULL,
    `word` VARCHAR(255) NOT NULL DEFAULT '',
    `means` JSON NULL,
    `essays` TEXT NULL,
    `labels` JSON NULL,
    `deleted` BOOLEAN NOT NULL DEFAULT false,
    `easiness` FLOAT NOT NULL DEFAULT 0,
    `interval` FLOAT NOT NULL DEFAULT 0,
    `need_review_at` DATETIME(3) NULL,
    `proficiency` VARCHAR(50) NOT NULL DEFAULT '',
    `repetitions` INTEGER NOT NULL DEFAULT 0,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    INDEX `words_word_idx`(`word`),
    INDEX `words_need_review_at_idx`(`need_review_at`),
    INDEX `words_deleted_idx`(`deleted`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
