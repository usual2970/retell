/*
  Warnings:

  - Added the required column `updated_at` to the `essay` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE `essay` ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `deleted_at` DATETIME(3) NULL,
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL;
