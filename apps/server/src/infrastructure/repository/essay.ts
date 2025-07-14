import { Essay } from "../../domain/essay";
import { EssayRepository } from "../../essay/service";
import { getPrismaClient } from "./base";
import { Prisma } from "@prisma/client";
import type { Essay as PrismaEssay } from "@prisma/client";

export class Repository implements EssayRepository {
  async create(essay: Essay): Promise<Essay> {
    const prisma = getPrismaClient();

    // 将领域模型转换为 Prisma 创建数据
    const createData = this.mapEssayToPrismaCreate(essay);

    const createdEssay = await prisma.essay.create({
      data: createData,
    });

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

    // 将领域模型转换为 Prisma 更新数据
    const updateData = this.mapEssayToPrismaUpdate(data);

    try {
      const updatedEssay = await prisma.essay.update({
        where: { id },
        data: updateData,
      });

      return this.mapPrismaToEssay(updatedEssay);
    } catch (error) {
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

  private mapEssayToPrismaCreate(essay: Essay): Prisma.EssayCreateInput {
    return this.buildPrismaData(essay) as Prisma.EssayCreateInput;
  }

  private mapEssayToPrismaUpdate(
    data: Partial<Essay>
  ): Prisma.EssayUpdateInput {
    // 过滤掉不应该更新的字段
    const excludeFields = ["id", "createdAt", "updatedAt", "deletedAt"];
    const filteredData = Object.fromEntries(
      Object.entries(data).filter(([key]) => !excludeFields.includes(key))
    ) as Partial<Essay>;

    return this.buildPrismaData(filteredData);
  }

  private buildPrismaData(data: Partial<Essay>): Record<string, any> {
    const result: Record<string, any> = {};

    Object.entries(data).forEach(([key, value]) => {
      if (value !== undefined) {
        if (key === "sentences") {
          // 特殊处理 sentences 字段
          result.sentences =
            value && (value as string[]).length > 0
              ? (value as Prisma.InputJsonValue)
              : undefined;
        } else if (key !== "id") {
          // 其他字段直接赋值（排除 id）
          result[key] = value;
        }
      }
    });

    return result;
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
