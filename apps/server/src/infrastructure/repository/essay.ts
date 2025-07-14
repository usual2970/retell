import { Essay } from "../../domain/essay";
import { EssayRepository } from "../../essay/service";
import { getPrismaClient } from "./base";
import { Prisma } from "../../generated/prisma";
import type { Essay as PrismaEssay } from "../../generated/prisma";

export class Repository implements EssayRepository {
  async create(essay: Essay): Promise<Essay> {
    const prisma = getPrismaClient();

    // 将领域模型转换为 Prisma 数据模型
    // 注意：不包含 id，让数据库自动生成
    const prismaData: Prisma.EssayCreateInput = {
      content: essay.content,
      essayType: essay.essayType,
      file: essay.file,
      fileId: essay.fileId,
      sentences:
        essay.sentences && essay.sentences.length > 0
          ? (essay.sentences as Prisma.InputJsonValue)
          : undefined,
      taskId: essay.taskId,
      telegraph: essay.telegraph,
      thumb: essay.thumb,
      title: essay.title,
      videoLink: essay.videoLink,
    };

    const createdEssay = await prisma.essay.create({
      data: prismaData,
    });

    // 将 Prisma 模型转换回领域模型
    return this.mapPrismaToEssay(createdEssay);
  }

  async findById(id: number): Promise<Essay | null> {
    const prisma = getPrismaClient();

    const essay = await prisma.essay.findUnique({
      where: { id },
    });

    return essay ? this.mapPrismaToEssay(essay) : null;
  }

  async findAll(params?: {
    skip?: number;
    take?: number;
    where?: Prisma.EssayWhereInput;
  }): Promise<Essay[]> {
    const prisma = getPrismaClient();

    const essays = await prisma.essay.findMany({
      skip: params?.skip,
      take: params?.take,
      where: params?.where,
      orderBy: { createdAt: "desc" },
    });

    return essays.map((essay) => this.mapPrismaToEssay(essay));
  }

  async update(id: number, data: Partial<Essay>): Promise<Essay | null> {
    const prisma = getPrismaClient();

    // 构建更新数据，只包含已提供的字段
    const updateData: Prisma.EssayUpdateInput = {};

    if (data.content !== undefined) updateData.content = data.content;
    if (data.essayType !== undefined) updateData.essayType = data.essayType;
    if (data.file !== undefined) updateData.file = data.file;
    if (data.fileId !== undefined) updateData.fileId = data.fileId;
    if (data.sentences !== undefined) {
      updateData.sentences =
        data.sentences && data.sentences.length > 0
          ? (data.sentences as Prisma.InputJsonValue)
          : undefined;
    }
    if (data.taskId !== undefined) updateData.taskId = data.taskId;
    if (data.telegraph !== undefined) updateData.telegraph = data.telegraph;
    if (data.thumb !== undefined) updateData.thumb = data.thumb;
    if (data.title !== undefined) updateData.title = data.title;
    if (data.videoLink !== undefined) updateData.videoLink = data.videoLink;

    try {
      const updatedEssay = await prisma.essay.update({
        where: { id },
        data: updateData,
      });

      return this.mapPrismaToEssay(updatedEssay);
    } catch (error) {
      // 如果记录不存在，Prisma 会抛出错误
      return null;
    }
  }

  async delete(id: number): Promise<boolean> {
    const prisma = getPrismaClient();

    try {
      await prisma.essay.delete({
        where: { id },
      });
      return true;
    } catch (error) {
      return false;
    }
  }

  async count(where?: Prisma.EssayWhereInput): Promise<number> {
    const prisma = getPrismaClient();

    return await prisma.essay.count({
      where,
    });
  }

  private mapPrismaToEssay(prismaEssay: PrismaEssay): Essay {
    return {
      id: prismaEssay.id,
      content: prismaEssay.content,
      essayType: prismaEssay.essayType,
      file: prismaEssay.file,
      fileId: prismaEssay.fileId,
      sentences: this.parseSentences(prismaEssay.sentences),
      taskId: prismaEssay.taskId,
      telegraph: prismaEssay.telegraph,
      thumb: prismaEssay.thumb,
      title: prismaEssay.title,
      videoLink: prismaEssay.videoLink,
      createdAt: prismaEssay.createdAt,
      updatedAt: prismaEssay.updatedAt,
      deletedAt: prismaEssay.deletedAt,
    };
  }

  private parseSentences(sentences: Prisma.JsonValue): string[] | undefined {
    if (sentences === null || sentences === undefined) {
      return undefined;
    }

    if (Array.isArray(sentences)) {
      return sentences as string[];
    }

    if (typeof sentences === "string") {
      try {
        const parsed = JSON.parse(sentences);
        return Array.isArray(parsed) ? parsed : undefined;
      } catch {
        return undefined;
      }
    }

    return undefined;
  }
}
