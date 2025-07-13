/*
  Warnings:

  - The primary key for the `words` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to alter the column `id` on the `words` table. The data in that column could be lost. The data in that column will be cast from `VarChar(30)` to `Int`.

*/
-- AlterTable
ALTER TABLE `words` DROP PRIMARY KEY,
    MODIFY `id` INTEGER NOT NULL AUTO_INCREMENT,
    ADD PRIMARY KEY (`id`);
